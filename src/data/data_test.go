package data

import "testing"

func TestReadFile(t *testing.T) {
	t.Run("reads an existing embedded json file", func(t *testing.T) {
		got, err := ReadFile("jokes.json")
		if err != nil {
			t.Fatalf("ReadFile(jokes.json) error = %v", err)
		}
		if len(got) == 0 {
			t.Error("ReadFile(jokes.json) returned empty content")
		}
	})

	t.Run("returns error for a missing file", func(t *testing.T) {
		_, err := ReadFile("does-not-exist.json")
		if err == nil {
			t.Error("ReadFile(does-not-exist.json) error = nil, want non-nil")
		}
	})

	t.Run("returns error for an empty name", func(t *testing.T) {
		_, err := ReadFile("")
		if err == nil {
			t.Error("ReadFile(\"\") error = nil, want non-nil")
		}
	})
}
