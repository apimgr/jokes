package graphql

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/apimgr/jokes/src/models"
)

func loadFixture(t *testing.T) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "jokes.json")

	data := models.JokesData{
		Jokes: []models.Joke{
			{ID: 1, Joke: "Chuck Norris counted to infinity twice.", Categories: []string{"dev", "nerdy"}},
			{ID: 2, Joke: "Chuck Norris can divide by zero.", Categories: []string{"dev"}},
			{ID: 3, Joke: "Explicit joke about Chuck Norris.", Categories: []string{"explicit"}},
		},
		Categories: []string{"dev", "nerdy", "explicit"},
	}

	raw, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("failed to marshal fixture: %v", err)
	}
	if err := os.WriteFile(path, raw, 0644); err != nil {
		t.Fatalf("failed to write fixture: %v", err)
	}
	if err := models.LoadJokes(path); err != nil {
		t.Fatalf("LoadJokes() error = %v", err)
	}
}

func TestResolverJoke(t *testing.T) {
	loadFixture(t)
	r := &Resolver{}

	t.Run("found", func(t *testing.T) {
		got := r.Joke(struct{ ID int32 }{ID: 1})
		if got == nil {
			t.Fatal("Joke(1) = nil, want a resolver")
		}
		if got.ID() != 1 {
			t.Errorf("ID() = %d, want 1", got.ID())
		}
		if got.Joke() != "Chuck Norris counted to infinity twice." {
			t.Errorf("Joke() = %q, unexpected", got.Joke())
		}
		wantCats := []string{"dev", "nerdy"}
		if !reflect.DeepEqual(*got.Categories(), wantCats) {
			t.Errorf("Categories() = %v, want %v", *got.Categories(), wantCats)
		}
	})

	t.Run("not found", func(t *testing.T) {
		got := r.Joke(struct{ ID int32 }{ID: 999})
		if got != nil {
			t.Errorf("Joke(999) = %v, want nil", got)
		}
	})
}

func TestResolverRandomJoke(t *testing.T) {
	loadFixture(t)
	r := &Resolver{}

	t.Run("no filters returns a joke", func(t *testing.T) {
		got := r.RandomJoke(struct {
			Category *string
			Exclude  *string
			LimitTo  *string
		}{})
		if got == nil {
			t.Fatal("RandomJoke() = nil, want a resolver")
		}
	})

	t.Run("category filters the pool", func(t *testing.T) {
		category := "explicit"
		got := r.RandomJoke(struct {
			Category *string
			Exclude  *string
			LimitTo  *string
		}{Category: &category})
		if got == nil {
			t.Fatal("RandomJoke(category=explicit) = nil, want a resolver")
		}
	})

	t.Run("category with no matches returns nil", func(t *testing.T) {
		category := "nonexistent"
		got := r.RandomJoke(struct {
			Category *string
			Exclude  *string
			LimitTo  *string
		}{Category: &category})
		if got != nil {
			t.Errorf("RandomJoke(category=nonexistent) = %v, want nil", got)
		}
	})
}

func TestResolverRandomJokes(t *testing.T) {
	loadFixture(t)
	r := &Resolver{}

	tests := []struct {
		name  string
		count int32
		want  int
	}{
		{"zero count clamps to one", 0, 1},
		{"count within range", 2, 2},
		{"negative count clamps to one", -5, 1},
		{"count exceeds max clamps to 100", 500, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := r.RandomJokes(struct {
				Count   int32
				Exclude *string
				LimitTo *string
			}{Count: tt.count})
			if len(*got) != tt.want {
				t.Errorf("RandomJokes(%d) returned %d resolvers, want %d", tt.count, len(*got), tt.want)
			}
		})
	}

	t.Run("limitTo filters results", func(t *testing.T) {
		limitTo := "explicit"
		got := r.RandomJokes(struct {
			Count   int32
			Exclude *string
			LimitTo *string
		}{Count: 10, LimitTo: &limitTo})
		if len(*got) != 1 {
			t.Fatalf("RandomJokes(limitTo=explicit) returned %d, want 1", len(*got))
		}
		if (*got)[0].ID() != 3 {
			t.Errorf("RandomJokes(limitTo=explicit)[0].ID() = %d, want 3", (*got)[0].ID())
		}
	})

	t.Run("exclude filters results", func(t *testing.T) {
		exclude := "dev, nerdy"
		got := r.RandomJokes(struct {
			Count   int32
			Exclude *string
			LimitTo *string
		}{Count: 10, Exclude: &exclude})
		if len(*got) != 1 {
			t.Fatalf("RandomJokes(exclude=dev,nerdy) returned %d, want 1", len(*got))
		}
		if (*got)[0].ID() != 3 {
			t.Errorf("RandomJokes(exclude=dev,nerdy)[0].ID() = %d, want 3", (*got)[0].ID())
		}
	})
}

func TestResolverAllJokes(t *testing.T) {
	loadFixture(t)
	r := &Resolver{}

	t.Run("no filters returns all", func(t *testing.T) {
		got := r.AllJokes(struct {
			LimitTo *string
			Exclude *string
		}{})
		if len(*got) != 3 {
			t.Errorf("AllJokes() returned %d, want 3", len(*got))
		}
	})

	t.Run("limitTo filters", func(t *testing.T) {
		limitTo := "dev"
		got := r.AllJokes(struct {
			LimitTo *string
			Exclude *string
		}{LimitTo: &limitTo})
		if len(*got) != 2 {
			t.Errorf("AllJokes(limitTo=dev) returned %d, want 2", len(*got))
		}
	})

	t.Run("exclude filters", func(t *testing.T) {
		exclude := "explicit"
		got := r.AllJokes(struct {
			LimitTo *string
			Exclude *string
		}{Exclude: &exclude})
		if len(*got) != 2 {
			t.Errorf("AllJokes(exclude=explicit) returned %d, want 2", len(*got))
		}
	})
}

func TestResolverCategories(t *testing.T) {
	loadFixture(t)
	r := &Resolver{}

	got := r.Categories()
	if len(*got) != 3 {
		t.Errorf("Categories() returned %d, want 3", len(*got))
	}
}

func TestResolverJokesByCategory(t *testing.T) {
	loadFixture(t)
	r := &Resolver{}

	t.Run("known category", func(t *testing.T) {
		got := r.JokesByCategory(struct{ Category string }{Category: "dev"})
		if len(*got) != 2 {
			t.Errorf("JokesByCategory(dev) returned %d, want 2", len(*got))
		}
	})

	t.Run("unknown category returns empty", func(t *testing.T) {
		got := r.JokesByCategory(struct{ Category string }{Category: "nope"})
		if len(*got) != 0 {
			t.Errorf("JokesByCategory(nope) returned %d, want 0", len(*got))
		}
	})
}

func TestResolverStats(t *testing.T) {
	loadFixture(t)
	r := &Resolver{}

	stats := r.Stats()
	if stats.Total() != 3 {
		t.Errorf("Stats().Total() = %d, want 3", stats.Total())
	}

	catStats := stats.Categories()
	if len(*catStats) != 3 {
		t.Fatalf("Stats().Categories() returned %d, want 3", len(*catStats))
	}

	counts := map[string]int32{}
	for _, cs := range *catStats {
		counts[cs.Name()] = cs.Count()
	}
	if counts["dev"] != 2 {
		t.Errorf("category dev count = %d, want 2", counts["dev"])
	}
	if counts["explicit"] != 1 {
		t.Errorf("category explicit count = %d, want 1", counts["explicit"])
	}
}

func TestParseCategoryList(t *testing.T) {
	tests := []struct {
		name    string
		limitTo string
		want    []string
	}{
		{"empty string", "", []string{}},
		{"single category", "dev", []string{"dev"}},
		{"comma separated", "dev,nerdy", []string{"dev", "nerdy"}},
		{"with brackets", "[dev,nerdy]", []string{"dev", "nerdy"}},
		{"with spaces", " dev , nerdy ", []string{"dev", "nerdy"}},
		{"only brackets", "[]", []string{}},
		{"trailing comma skips empty entries", "dev,,nerdy", []string{"dev", "nerdy"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCategoryList(tt.limitTo)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCategoryList(%q) = %v, want %v", tt.limitTo, got, tt.want)
			}
		})
	}
}
