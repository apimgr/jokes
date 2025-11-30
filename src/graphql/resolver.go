package graphql

import (
	"strings"

	"github.com/apimgr/jokes/src/models"
)

type Resolver struct{}

// Query resolvers
func (r *Resolver) Joke(args struct{ ID int32 }) *JokeResolver {
	joke := models.GetJokeByID(int(args.ID))
	if joke == nil {
		return nil
	}
	return &JokeResolver{joke: joke}
}

func (r *Resolver) RandomJoke(args struct {
	Category *string
	Exclude  *string
	LimitTo  *string
}) *JokeResolver {
	var jokes []models.Joke

	// Handle limitTo
	if args.LimitTo != nil {
		limitToCategories := parseCategoryList(*args.LimitTo)
		jokes = models.GetAllJokes()
		jokes = models.FilterJokesByCategories(jokes, limitToCategories, nil)
	} else if args.Category != nil && *args.Category != "" {
		jokes = models.GetJokesByCategory(*args.Category)
	} else {
		jokes = models.GetAllJokes()
	}

	// Handle exclude
	if args.Exclude != nil && *args.Exclude != "" {
		excludeCategories := strings.Split(*args.Exclude, ",")
		for i := range excludeCategories {
			excludeCategories[i] = strings.TrimSpace(excludeCategories[i])
		}
		jokes = models.FilterJokesByCategories(jokes, nil, excludeCategories)
	}

	if len(jokes) == 0 {
		return nil
	}

	// Get random joke from filtered list
	randomJokes := models.GetRandomJokes(1)
	if len(randomJokes) > 0 {
		// Filter the random joke as well
		if args.Exclude != nil || args.LimitTo != nil {
			filtered := models.FilterJokesByCategories(randomJokes, nil, nil)
			if len(filtered) > 0 {
				return &JokeResolver{joke: &filtered[0]}
			}
		}
		return &JokeResolver{joke: &randomJokes[0]}
	}

	return nil
}

func (r *Resolver) RandomJokes(args struct {
	Count   int32
	Exclude *string
	LimitTo *string
}) *[]*JokeResolver {
	count := int(args.Count)
	if count < 1 {
		count = 1
	}
	if count > 100 {
		count = 100
	}

	jokes := models.GetRandomJokes(count)

	// Handle limitTo
	if args.LimitTo != nil && *args.LimitTo != "" {
		limitToCategories := parseCategoryList(*args.LimitTo)
		jokes = models.FilterJokesByCategories(jokes, limitToCategories, nil)
	}

	// Handle exclude
	if args.Exclude != nil && *args.Exclude != "" {
		excludeCategories := strings.Split(*args.Exclude, ",")
		for i := range excludeCategories {
			excludeCategories[i] = strings.TrimSpace(excludeCategories[i])
		}
		jokes = models.FilterJokesByCategories(jokes, nil, excludeCategories)
	}

	resolvers := make([]*JokeResolver, len(jokes))
	for i := range jokes {
		resolvers[i] = &JokeResolver{joke: &jokes[i]}
	}
	return &resolvers
}

func (r *Resolver) AllJokes(args struct {
	LimitTo *string
	Exclude *string
}) *[]*JokeResolver {
	jokes := models.GetAllJokes()

	// Handle limitTo
	if args.LimitTo != nil && *args.LimitTo != "" {
		limitToCategories := parseCategoryList(*args.LimitTo)
		jokes = models.FilterJokesByCategories(jokes, limitToCategories, nil)
	}

	// Handle exclude
	if args.Exclude != nil && *args.Exclude != "" {
		excludeCategories := strings.Split(*args.Exclude, ",")
		for i := range excludeCategories {
			excludeCategories[i] = strings.TrimSpace(excludeCategories[i])
		}
		jokes = models.FilterJokesByCategories(jokes, nil, excludeCategories)
	}

	resolvers := make([]*JokeResolver, len(jokes))
	for i := range jokes {
		resolvers[i] = &JokeResolver{joke: &jokes[i]}
	}
	return &resolvers
}

func (r *Resolver) Categories() *[]string {
	categories := models.GetCategories()
	return &categories
}

func (r *Resolver) JokesByCategory(args struct{ Category string }) *[]*JokeResolver {
	jokes := models.GetJokesByCategory(args.Category)
	resolvers := make([]*JokeResolver, len(jokes))
	for i := range jokes {
		resolvers[i] = &JokeResolver{joke: &jokes[i]}
	}
	return &resolvers
}

func (r *Resolver) Stats() *StatsResolver {
	return &StatsResolver{}
}

// Type resolvers
type JokeResolver struct {
	joke *models.Joke
}

func (r *JokeResolver) ID() int32 {
	return int32(r.joke.ID)
}

func (r *JokeResolver) Joke() string {
	return r.joke.Joke
}

func (r *JokeResolver) Categories() *[]string {
	return &r.joke.Categories
}

type StatsResolver struct{}

func (r *StatsResolver) Total() int32 {
	return int32(models.GetJokesCount())
}

func (r *StatsResolver) Categories() *[]*CategoryStatsResolver {
	categories := models.GetCategories()
	resolvers := make([]*CategoryStatsResolver, len(categories))
	for i, cat := range categories {
		jokes := models.GetJokesByCategory(cat)
		resolvers[i] = &CategoryStatsResolver{
			name:  cat,
			count: len(jokes),
		}
	}
	return &resolvers
}

type CategoryStatsResolver struct {
	name  string
	count int
}

func (r *CategoryStatsResolver) Name() string {
	return r.name
}

func (r *CategoryStatsResolver) Count() int32 {
	return int32(r.count)
}

// Helper functions
func parseCategoryList(limitTo string) []string {
	limitTo = strings.Trim(limitTo, "[]")
	limitTo = strings.TrimSpace(limitTo)

	if limitTo == "" {
		return []string{}
	}

	categories := strings.Split(limitTo, ",")
	result := []string{}
	for _, cat := range categories {
		cat = strings.TrimSpace(cat)
		if cat != "" {
			result = append(result, cat)
		}
	}
	return result
}
