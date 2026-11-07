package model

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/apimgr/jokes/src/data"
)

type Joke struct {
	ID         int      `json:"id"`
	Type       string   `json:"type"`
	Joke       string   `json:"joke,omitempty"`
	Setup      string   `json:"setup,omitempty"`
	Delivery   string   `json:"delivery,omitempty"`
	Categories []string `json:"categories"`
}

type JokesData struct {
	Jokes      []Joke   `json:"jokes"`
	Categories []string `json:"categories"`
}

var (
	jokesMu   sync.RWMutex
	jokesData *JokesData
)

// LoadEmbeddedJokes loads the curated dataset compiled into the application.
func LoadEmbeddedJokes() error {
	contents, err := data.ReadFile("jokes.json")
	if err != nil {
		return fmt.Errorf("failed to read embedded jokes: %w", err)
	}
	return loadJokes(contents)
}

// LoadJokes loads jokes from an operator-supplied JSON file.
func LoadJokes(filePath string) error {
	contents, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read jokes file: %w", err)
	}
	return loadJokes(contents)
}

func loadJokes(contents []byte) error {
	var next JokesData
	if err := json.Unmarshal(contents, &next); err != nil {
		return fmt.Errorf("failed to parse jokes file: %w", err)
	}

	seen := make(map[int]struct{}, len(next.Jokes))
	categories := make(map[string]struct{})
	for i := range next.Jokes {
		joke := &next.Jokes[i]
		if joke.ID < 1 {
			return fmt.Errorf("joke at index %d has invalid ID %d", i, joke.ID)
		}
		if _, ok := seen[joke.ID]; ok {
			return fmt.Errorf("duplicate joke ID %d", joke.ID)
		}
		seen[joke.ID] = struct{}{}
		if joke.Type == "" {
			if joke.Setup != "" || joke.Delivery != "" {
				joke.Type = "twopart"
			} else {
				joke.Type = "single"
			}
		}
		switch joke.Type {
		case "single":
			if strings.TrimSpace(joke.Joke) == "" {
				return fmt.Errorf("single joke %d has no text", joke.ID)
			}
		case "twopart":
			if strings.TrimSpace(joke.Setup) == "" || strings.TrimSpace(joke.Delivery) == "" {
				return fmt.Errorf("two-part joke %d must have both setup and delivery", joke.ID)
			}
		case "":
			return fmt.Errorf("joke %d has no type or text", joke.ID)
		default:
			return fmt.Errorf("joke %d has unsupported type %q", joke.ID, joke.Type)
		}
		for _, category := range joke.Categories {
			category = strings.TrimSpace(category)
			if category == "" {
				return fmt.Errorf("joke %d has an empty category", joke.ID)
			}
			categories[category] = struct{}{}
		}
		joke.Categories = append([]string(nil), joke.Categories...)
	}
	if len(next.Categories) == 0 {
		for category := range categories {
			next.Categories = append(next.Categories, category)
		}
	}
	sort.Strings(next.Categories)
	jokesMu.Lock()
	jokesData = &next
	jokesMu.Unlock()
	return nil
}

// Text returns the complete user-facing text for either joke format.
func (j Joke) Text() string {
	if j.Type == "twopart" {
		return j.Setup + "\n" + j.Delivery
	}
	return j.Joke
}

func cloneJoke(joke Joke) Joke {
	joke.Categories = append([]string(nil), joke.Categories...)
	return joke
}

func snapshot() JokesData {
	jokesMu.RLock()
	defer jokesMu.RUnlock()
	if jokesData == nil {
		return JokesData{}
	}
	copyData := JokesData{
		Jokes:      make([]Joke, len(jokesData.Jokes)),
		Categories: append([]string(nil), jokesData.Categories...),
	}
	for i, joke := range jokesData.Jokes {
		copyData.Jokes[i] = cloneJoke(joke)
	}
	return copyData
}

// GetAllJokes returns all jokes
func GetAllJokes() []Joke {
	return snapshot().Jokes
}

// GetJokeByID returns a joke by its ID
func GetJokeByID(id int) *Joke {
	for _, joke := range GetAllJokes() {
		if joke.ID == id {
			return &joke
		}
	}
	return nil
}

// GetRandomJoke returns a random joke
func GetRandomJoke() *Joke {
	jokes := GetAllJokes()
	if len(jokes) == 0 {
		return nil
	}
	return &jokes[rand.Intn(len(jokes))]
}

// GetRandomJokes returns multiple random jokes
func GetRandomJokes(count int) []Joke {
	return RandomJokesFrom(GetAllJokes(), count)
}

// RandomJokesFrom returns up to count unique random jokes from a candidate set.
func RandomJokesFrom(candidates []Joke, count int) []Joke {
	if count <= 0 || len(candidates) == 0 {
		return []Joke{}
	}
	if count > len(candidates) {
		count = len(candidates)
	}
	jokes := make([]Joke, len(candidates))
	for i := range candidates {
		jokes[i] = cloneJoke(candidates[i])
	}
	for i := 0; i < count; i++ {
		j := i + rand.Intn(len(jokes)-i)
		jokes[i], jokes[j] = jokes[j], jokes[i]
	}
	return jokes[:count]
}

// GetJokesByCategory returns all jokes in a category
func GetJokesByCategory(category string) []Joke {
	result := []Joke{}
	for _, joke := range GetAllJokes() {
		for _, cat := range joke.Categories {
			if cat == category {
				result = append(result, joke)
				break
			}
		}
	}
	return result
}

// GetCategories returns all available categories
func GetCategories() []string {
	return snapshot().Categories
}

// GetJokesCount returns the total number of jokes
func GetJokesCount() int {
	jokesMu.RLock()
	defer jokesMu.RUnlock()
	if jokesData == nil {
		return 0
	}
	return len(jokesData.Jokes)
}

// ReplaceNameInJoke replaces "Chuck Norris" with custom names
func ReplaceNameInJoke(joke string, firstName, lastName string) string {
	if firstName == "" {
		firstName = "Chuck"
	}
	if lastName == "" {
		lastName = "Norris"
	}

	// Replace full name first, then just first name
	result := strings.ReplaceAll(joke, "Chuck Norris", firstName+" "+lastName)
	result = strings.ReplaceAll(result, "Chuck", firstName)

	return result
}

// FilterJokesByCategories filters jokes by including/excluding categories
func FilterJokesByCategories(jokes []Joke, limitTo []string, exclude []string) []Joke {
	result := []Joke{}

	for _, joke := range jokes {
		// Check if joke should be limited to certain categories
		if len(limitTo) > 0 {
			hasCategory := false
			for _, cat := range joke.Categories {
				for _, limitCat := range limitTo {
					if cat == limitCat {
						hasCategory = true
						break
					}
				}
				if hasCategory {
					break
				}
			}
			if !hasCategory {
				continue
			}
		}

		// Check if joke should be excluded
		if len(exclude) > 0 {
			shouldExclude := false
			for _, cat := range joke.Categories {
				for _, excludeCat := range exclude {
					if cat == excludeCat {
						shouldExclude = true
						break
					}
				}
				if shouldExclude {
					break
				}
			}
			if shouldExclude {
				continue
			}
		}

		result = append(result, joke)
	}

	return result
}

// IsCategoryValid checks if a category exists
func IsCategoryValid(category string) bool {
	for _, cat := range GetCategories() {
		if cat == category {
			return true
		}
	}
	return false
}

// SearchJokes returns jokes whose content or category contains query,
// case-insensitively. An empty query matches every joke.
func SearchJokes(jokes []Joke, query string) []Joke {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return jokes
	}
	result := make([]Joke, 0)
	for _, joke := range jokes {
		if strings.Contains(strings.ToLower(joke.Text()), query) {
			result = append(result, joke)
			continue
		}
		for _, category := range joke.Categories {
			if strings.Contains(strings.ToLower(category), query) {
				result = append(result, joke)
				break
			}
		}
	}
	return result
}

// Paginate returns a bounded window of records. Invalid or out-of-range
// offsets return a non-nil empty slice.
func Paginate(jokes []Joke, offset, limit int) []Joke {
	if offset < 0 || limit <= 0 || offset >= len(jokes) {
		return []Joke{}
	}
	end := offset + limit
	if end < offset || end > len(jokes) {
		end = len(jokes)
	}
	result := make([]Joke, end-offset)
	for i, joke := range jokes[offset:end] {
		result[i] = cloneJoke(joke)
	}
	return result
}
