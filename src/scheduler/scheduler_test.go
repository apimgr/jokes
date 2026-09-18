package scheduler

import (
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestAddRemoveTask(t *testing.T) {
	s := New()
	s.AddTask("task1", time.Hour, func() error { return nil })

	tasks := s.GetTasks()
	if len(tasks) != 1 {
		t.Fatalf("GetTasks() len = %d, want 1", len(tasks))
	}
	if tasks[0].Name != "task1" {
		t.Errorf("task name = %q, want task1", tasks[0].Name)
	}
	if !tasks[0].Enabled {
		t.Error("newly added task should be enabled")
	}

	s.RemoveTask("task1")
	if got := len(s.GetTasks()); got != 0 {
		t.Fatalf("GetTasks() len after remove = %d, want 0", got)
	}
}

func TestEnableDisableTask(t *testing.T) {
	s := New()
	s.AddTask("task1", time.Hour, func() error { return nil })

	s.DisableTask("task1")
	tasks := s.GetTasks()
	if tasks[0].Enabled {
		t.Error("task should be disabled after DisableTask")
	}

	s.EnableTask("task1")
	tasks = s.GetTasks()
	if !tasks[0].Enabled {
		t.Error("task should be enabled after EnableTask")
	}
}

func TestEnableDisableTaskUnknownName(t *testing.T) {
	s := New()
	s.EnableTask("missing")
	s.DisableTask("missing")
	if got := len(s.GetTasks()); got != 0 {
		t.Fatalf("GetTasks() len = %d, want 0", got)
	}
}

func TestRunNow(t *testing.T) {
	s := New()
	var ran int32
	s.AddTask("task1", time.Hour, func() error {
		atomic.AddInt32(&ran, 1)
		return nil
	})

	if err := s.RunNow("task1"); err != nil {
		t.Fatalf("RunNow() error = %v", err)
	}
	if atomic.LoadInt32(&ran) != 1 {
		t.Errorf("task function ran %d times, want 1", ran)
	}

	tasks := s.GetTasks()
	if tasks[0].LastRun.IsZero() {
		t.Error("LastRun should be set after RunNow")
	}
	if !tasks[0].NextRun.After(tasks[0].LastRun) {
		t.Error("NextRun should be after LastRun")
	}
}

func TestRunNowPropagatesError(t *testing.T) {
	s := New()
	wantErr := errors.New("boom")
	s.AddTask("task1", time.Hour, func() error { return wantErr })

	if err := s.RunNow("task1"); !errors.Is(err, wantErr) {
		t.Errorf("RunNow() error = %v, want %v", err, wantErr)
	}
}

func TestRunNowUnknownTask(t *testing.T) {
	s := New()
	if err := s.RunNow("missing"); err != nil {
		t.Errorf("RunNow() on unknown task error = %v, want nil", err)
	}
}

func TestGetTasksEmpty(t *testing.T) {
	s := New()
	tasks := s.GetTasks()
	if tasks == nil {
		t.Error("GetTasks() should return non-nil empty slice")
	}
	if len(tasks) != 0 {
		t.Errorf("GetTasks() len = %d, want 0", len(tasks))
	}
}

func TestStartStop(t *testing.T) {
	s := New()
	ran := make(chan struct{}, 1)
	s.AddTask("task1", time.Millisecond, func() error {
		select {
		case ran <- struct{}{}:
		default:
		}
		return nil
	})

	s.Start()
	defer s.Stop()

	select {
	case <-ran:
	case <-time.After(2 * time.Second):
		t.Fatal("task did not run within timeout")
	}

	s.Stop()

	select {
	case <-ran:
	default:
	}
	select {
	case <-ran:
		t.Error("task ran again after Stop")
	case <-time.After(50 * time.Millisecond):
	}
}

func TestStartTwiceIsNoop(t *testing.T) {
	s := New()
	s.AddTask("task1", time.Hour, func() error { return nil })
	s.Start()
	s.Start()
	s.Stop()
}

func TestStopWithoutStartIsSafe(t *testing.T) {
	s := New()
	s.Stop()
	s.Stop()
}

func TestStopIsIdempotent(t *testing.T) {
	s := New()
	s.Start()
	s.Stop()
	s.Stop()
}

func TestParseInterval(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  time.Duration
	}{
		{"minutely", "minutely", time.Minute},
		{"hourly", "hourly", time.Hour},
		{"daily", "daily", 24 * time.Hour},
		{"weekly", "weekly", 7 * 24 * time.Hour},
		{"monthly", "monthly", 30 * 24 * time.Hour},
		{"duration string", "5m", 5 * time.Minute},
		{"duration hours", "2h", 2 * time.Hour},
		{"duration seconds", "30s", 30 * time.Second},
		{"empty string falls back to daily", "", 24 * time.Hour},
		{"garbage falls back to daily", "not-a-duration", 24 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseInterval(tt.input); got != tt.want {
				t.Errorf("ParseInterval(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
