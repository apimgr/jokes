package graphql

import (
	"strings"

	models "github.com/apimgr/jokes/src/model"
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
	jokes := models.GetAllJokes()
	if args.Category != nil && *args.Category != "" {
		jokes = models.GetJokesByCategory(*args.Category)
	}
	if args.LimitTo != nil && *args.LimitTo != "" {
		jokes = models.FilterJokesByCategories(jokes, parseCategoryList(*args.LimitTo), nil)
	}
	if args.Exclude != nil && *args.Exclude != "" {
		jokes = models.FilterJokesByCategories(jokes, nil, splitCategories(*args.Exclude))
	}
	selected := models.RandomJokesFrom(jokes, 1)
	if len(selected) == 0 {
		return nil
	}
	return &JokeResolver{joke: &selected[0]}
}

func (r *Resolver) RandomJokes(args struct {
	Count   int32
	Exclude *string
	LimitTo *string
}) *[]*JokeResolver {
	jokes := models.GetAllJokes()
	if args.LimitTo != nil && *args.LimitTo != "" {
		jokes = models.FilterJokesByCategories(jokes, parseCategoryList(*args.LimitTo), nil)
	}
	if args.Exclude != nil && *args.Exclude != "" {
		jokes = models.FilterJokesByCategories(jokes, nil, splitCategories(*args.Exclude))
	}
	count := int(args.Count)
	if count < 1 {
		count = 1
	}
	if count > 100 {
		count = 100
	}
	jokes = models.RandomJokesFrom(jokes, count)
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

func (r *JokeResolver) Type() string {
	return r.joke.Type
}

func (r *JokeResolver) Setup() *string {
	if r.joke.Setup == "" {
		return nil
	}
	return &r.joke.Setup
}

func (r *JokeResolver) Delivery() *string {
	if r.joke.Delivery == "" {
		return nil
	}
	return &r.joke.Delivery
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
func splitCategories(value string) []string {
	categories := strings.Split(strings.Trim(value, "[]"), ",")
	result := make([]string, 0, len(categories))
	for _, category := range categories {
		if category = strings.TrimSpace(category); category != "" {
			result = append(result, category)
		}
	}
	return result
}

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
