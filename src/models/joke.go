package models

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

type Joke struct {
	ID         int      `json:"id"`
	Joke       string   `json:"joke"`
	Categories []string `json:"categories"`
}

type JokesData struct {
	Jokes      []Joke   `json:"jokes"`
	Categories []string `json:"categories"`
}

var (
	jokesData *JokesData
)

// LoadJokes loads jokes from the JSON file
func LoadJokes(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read jokes file: %w", err)
	}

	jokesData = &JokesData{}
	if err := json.Unmarshal(data, jokesData); err != nil {
		return fmt.Errorf("failed to parse jokes file: %w", err)
	}

	// If categories not in file, extract from jokes
	if len(jokesData.Categories) == 0 {
		categoryMap := make(map[string]bool)
		for _, joke := range jokesData.Jokes {
			for _, cat := range joke.Categories {
				categoryMap[cat] = true
			}
		}
		jokesData.Categories = make([]string, 0, len(categoryMap))
		for cat := range categoryMap {
			jokesData.Categories = append(jokesData.Categories, cat)
		}
	}

	return nil
}

// GetAllJokes returns all jokes
func GetAllJokes() []Joke {
	if jokesData == nil {
		return []Joke{}
	}
	return jokesData.Jokes
}

// GetJokeByID returns a joke by its ID
func GetJokeByID(id int) *Joke {
	if jokesData == nil {
		return nil
	}
	for _, joke := range jokesData.Jokes {
		if joke.ID == id {
			return &joke
		}
	}
	return nil
}

// GetRandomJoke returns a random joke
func GetRandomJoke() *Joke {
	if jokesData == nil || len(jokesData.Jokes) == 0 {
		return nil
	}
	return &jokesData.Jokes[rand.Intn(len(jokesData.Jokes))]
}

// GetRandomJokes returns multiple random jokes
func GetRandomJokes(count int) []Joke {
	if jokesData == nil || len(jokesData.Jokes) == 0 {
		return []Joke{}
	}

	if count > len(jokesData.Jokes) {
		count = len(jokesData.Jokes)
	}

	// Create a copy and shuffle
	jokes := make([]Joke, len(jokesData.Jokes))
	copy(jokes, jokesData.Jokes)

	// Fisher-Yates shuffle
	for i := len(jokes) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		jokes[i], jokes[j] = jokes[j], jokes[i]
	}

	return jokes[:count]
}

// GetJokesByCategory returns all jokes in a category
func GetJokesByCategory(category string) []Joke {
	if jokesData == nil {
		return []Joke{}
	}

	result := []Joke{}
	for _, joke := range jokesData.Jokes {
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
	if jokesData == nil {
		return []string{}
	}
	return jokesData.Categories
}

// GetJokesCount returns the total number of jokes
func GetJokesCount() int {
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
	if jokesData == nil {
		return false
	}
	for _, cat := range jokesData.Categories {
		if cat == category {
			return true
		}
	}
	return false
}
