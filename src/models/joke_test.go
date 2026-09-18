package models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func writeJokesFile(t *testing.T, data JokesData) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "jokes.json")
	raw, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("failed to marshal fixture: %v", err)
	}
	if err := os.WriteFile(path, raw, 0644); err != nil {
		t.Fatalf("failed to write fixture: %v", err)
	}
	return path
}

func resetJokesData(t *testing.T) {
	t.Helper()
	orig := jokesData
	t.Cleanup(func() { jokesData = orig })
}

func TestLoadJokes(t *testing.T) {
	resetJokesData(t)

	t.Run("valid file with explicit categories", func(t *testing.T) {
		path := writeJokesFile(t, JokesData{
			Jokes: []Joke{
				{ID: 1, Joke: "Chuck Norris joke", Categories: []string{"dev"}},
			},
			Categories: []string{"dev", "explicit"},
		})
		if err := LoadJokes(path); err != nil {
			t.Fatalf("LoadJokes() error = %v", err)
		}
		if GetJokesCount() != 1 {
			t.Errorf("GetJokesCount() = %d, want 1", GetJokesCount())
		}
		cats := GetCategories()
		if len(cats) != 2 {
			t.Errorf("GetCategories() = %v, want 2 categories", cats)
		}
	})

	t.Run("categories derived from jokes when missing", func(t *testing.T) {
		path := writeJokesFile(t, JokesData{
			Jokes: []Joke{
				{ID: 1, Joke: "a", Categories: []string{"dev", "nerdy"}},
				{ID: 2, Joke: "b", Categories: []string{"dev"}},
			},
		})
		if err := LoadJokes(path); err != nil {
			t.Fatalf("LoadJokes() error = %v", err)
		}
		cats := GetCategories()
		sort.Strings(cats)
		want := []string{"dev", "nerdy"}
		if !reflect.DeepEqual(cats, want) {
			t.Errorf("derived categories = %v, want %v", cats, want)
		}
	})

	t.Run("missing file returns error", func(t *testing.T) {
		if err := LoadJokes(filepath.Join(t.TempDir(), "nope.json")); err == nil {
			t.Error("LoadJokes() with missing file should return an error")
		}
	})

	t.Run("malformed JSON returns error", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "bad.json")
		if err := os.WriteFile(path, []byte("{not json"), 0644); err != nil {
			t.Fatalf("setup failed: %v", err)
		}
		if err := LoadJokes(path); err == nil {
			t.Error("LoadJokes() with malformed JSON should return an error")
		}
	})
}

func TestAccessorsBeforeLoad(t *testing.T) {
	resetJokesData(t)
	jokesData = nil

	if got := GetAllJokes(); len(got) != 0 {
		t.Errorf("GetAllJokes() before load = %v, want empty slice", got)
	}
	if got := GetJokeByID(1); got != nil {
		t.Errorf("GetJokeByID() before load = %v, want nil", got)
	}
	if got := GetRandomJoke(); got != nil {
		t.Errorf("GetRandomJoke() before load = %v, want nil", got)
	}
	if got := GetRandomJokes(3); len(got) != 0 {
		t.Errorf("GetRandomJokes() before load = %v, want empty slice", got)
	}
	if got := GetJokesByCategory("dev"); len(got) != 0 {
		t.Errorf("GetJokesByCategory() before load = %v, want empty slice", got)
	}
	if got := GetCategories(); len(got) != 0 {
		t.Errorf("GetCategories() before load = %v, want empty slice", got)
	}
	if got := GetJokesCount(); got != 0 {
		t.Errorf("GetJokesCount() before load = %d, want 0", got)
	}
	if got := IsCategoryValid("dev"); got {
		t.Errorf("IsCategoryValid() before load = %v, want false", got)
	}
}

func testFixture() JokesData {
	return JokesData{
		Jokes: []Joke{
			{ID: 1, Joke: "Chuck Norris counted to infinity twice.", Categories: []string{"dev", "nerdy"}},
			{ID: 2, Joke: "Chuck Norris can divide by zero.", Categories: []string{"dev"}},
			{ID: 3, Joke: "Explicit joke about Chuck Norris.", Categories: []string{"explicit"}},
		},
		Categories: []string{"dev", "nerdy", "explicit"},
	}
}

func TestGetJokeByID(t *testing.T) {
	resetJokesData(t)
	path := writeJokesFile(t, testFixture())
	if err := LoadJokes(path); err != nil {
		t.Fatalf("LoadJokes() error = %v", err)
	}

	tests := []struct {
		name    string
		id      int
		wantNil bool
	}{
		{"existing id", 1, false},
		{"nonexistent id", 999, true},
		{"zero id", 0, true},
		{"negative id", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetJokeByID(tt.id)
			if tt.wantNil && got != nil {
				t.Errorf("GetJokeByID(%d) = %v, want nil", tt.id, got)
			}
			if !tt.wantNil {
				if got == nil {
					t.Fatalf("GetJokeByID(%d) = nil, want a joke", tt.id)
				}
				if got.ID != tt.id {
					t.Errorf("GetJokeByID(%d).ID = %d, want %d", tt.id, got.ID, tt.id)
				}
			}
		})
	}
}

func TestGetRandomJoke(t *testing.T) {
	resetJokesData(t)

	t.Run("empty dataset returns nil", func(t *testing.T) {
		path := writeJokesFile(t, JokesData{Jokes: []Joke{}, Categories: []string{}})
		if err := LoadJokes(path); err != nil {
			t.Fatalf("LoadJokes() error = %v", err)
		}
		if got := GetRandomJoke(); got != nil {
			t.Errorf("GetRandomJoke() with empty dataset = %v, want nil", got)
		}
	})

	t.Run("returns a joke from the dataset", func(t *testing.T) {
		path := writeJokesFile(t, testFixture())
		if err := LoadJokes(path); err != nil {
			t.Fatalf("LoadJokes() error = %v", err)
		}
		got := GetRandomJoke()
		if got == nil {
			t.Fatal("GetRandomJoke() = nil, want a joke")
		}
		if GetJokeByID(got.ID) == nil {
			t.Errorf("GetRandomJoke() returned joke with unknown id %d", got.ID)
		}
	})
}

func TestGetRandomJokes(t *testing.T) {
	resetJokesData(t)
	path := writeJokesFile(t, testFixture())
	if err := LoadJokes(path); err != nil {
		t.Fatalf("LoadJokes() error = %v", err)
	}

	tests := []struct {
		name  string
		count int
		want  int
	}{
		{"zero count", 0, 0},
		{"count within range", 2, 2},
		{"count equal to total", 3, 3},
		{"count exceeds total is clamped", 100, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetRandomJokes(tt.count)
			if len(got) != tt.want {
				t.Errorf("GetRandomJokes(%d) returned %d jokes, want %d", tt.count, len(got), tt.want)
			}
			seen := map[int]bool{}
			for _, j := range got {
				if seen[j.ID] {
					t.Errorf("GetRandomJokes(%d) returned duplicate id %d", tt.count, j.ID)
				}
				seen[j.ID] = true
			}
		})
	}
}

func TestGetJokesByCategory(t *testing.T) {
	resetJokesData(t)
	path := writeJokesFile(t, testFixture())
	if err := LoadJokes(path); err != nil {
		t.Fatalf("LoadJokes() error = %v", err)
	}

	tests := []struct {
		name     string
		category string
		want     int
	}{
		{"category with multiple matches", "dev", 2},
		{"category with one match", "explicit", 1},
		{"unknown category", "nope", 0},
		{"empty category", "", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetJokesByCategory(tt.category)
			if len(got) != tt.want {
				t.Errorf("GetJokesByCategory(%q) = %d jokes, want %d", tt.category, len(got), tt.want)
			}
		})
	}
}

func TestIsCategoryValid(t *testing.T) {
	resetJokesData(t)
	path := writeJokesFile(t, testFixture())
	if err := LoadJokes(path); err != nil {
		t.Fatalf("LoadJokes() error = %v", err)
	}

	if !IsCategoryValid("dev") {
		t.Error("IsCategoryValid(\"dev\") = false, want true")
	}
	if IsCategoryValid("nonexistent") {
		t.Error("IsCategoryValid(\"nonexistent\") = true, want false")
	}
}

func TestReplaceNameInJoke(t *testing.T) {
	tests := []struct {
		name      string
		joke      string
		firstName string
		lastName  string
		want      string
	}{
		{
			name:      "full name replaced",
			joke:      "Chuck Norris can divide by zero.",
			firstName: "Bruce",
			lastName:  "Lee",
			want:      "Bruce Lee can divide by zero.",
		},
		{
			name:      "first name only replaces remaining Chuck",
			joke:      "Chuck once punched Chuck Norris.",
			firstName: "Jackie",
			lastName:  "Chan",
			want:      "Jackie once punched Jackie Chan.",
		},
		{
			name:      "empty names default to Chuck Norris",
			joke:      "Chuck Norris can divide by zero.",
			firstName: "",
			lastName:  "",
			want:      "Chuck Norris can divide by zero.",
		},
		{
			name:      "no Chuck in joke leaves it unchanged",
			joke:      "This joke has no target name.",
			firstName: "Bruce",
			lastName:  "Lee",
			want:      "This joke has no target name.",
		},
		{
			name:      "only firstName provided keeps default lastName",
			joke:      "Chuck Norris can divide by zero.",
			firstName: "Bruce",
			lastName:  "",
			want:      "Bruce Norris can divide by zero.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReplaceNameInJoke(tt.joke, tt.firstName, tt.lastName)
			if got != tt.want {
				t.Errorf("ReplaceNameInJoke(%q, %q, %q) = %q, want %q", tt.joke, tt.firstName, tt.lastName, got, tt.want)
			}
		})
	}
}

func TestFilterJokesByCategories(t *testing.T) {
	jokes := []Joke{
		{ID: 1, Joke: "a", Categories: []string{"dev", "nerdy"}},
		{ID: 2, Joke: "b", Categories: []string{"dev"}},
		{ID: 3, Joke: "c", Categories: []string{"explicit"}},
	}

	tests := []struct {
		name    string
		limitTo []string
		exclude []string
		wantIDs []int
	}{
		{"no filters returns all", nil, nil, []int{1, 2, 3}},
		{"limitTo single category", []string{"nerdy"}, nil, []int{1}},
		{"limitTo multiple categories dedupes matches", []string{"dev"}, nil, []int{1, 2}},
		{"exclude removes matches", nil, []string{"explicit"}, []int{1, 2}},
		{"limitTo and exclude combined", []string{"dev"}, []string{"nerdy"}, []int{2}},
		{"limitTo with no matches returns empty", []string{"nonexistent"}, nil, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FilterJokesByCategories(jokes, tt.limitTo, tt.exclude)
			gotIDs := make([]int, 0, len(got))
			for _, j := range got {
				gotIDs = append(gotIDs, j.ID)
			}
			if !reflect.DeepEqual(gotIDs, tt.wantIDs) {
				t.Errorf("FilterJokesByCategories() ids = %v, want %v", gotIDs, tt.wantIDs)
			}
		})
	}
}
