package graphql

const Schema = `
	type Joke {
		id: Int!
		joke: String!
		categories: [String!]!
	}

	type JokeStats {
		total: Int!
		categories: [CategoryStats!]!
	}

	type CategoryStats {
		name: String!
		count: Int!
	}

	type Query {
		joke(id: Int!): Joke
		randomJoke(category: String, exclude: String, limitTo: String): Joke
		randomJokes(count: Int!, exclude: String, limitTo: String): [Joke!]!
		allJokes(limitTo: String, exclude: String): [Joke!]!
		categories: [String!]!
		jokesByCategory(category: String!): [Joke!]!
		stats: JokeStats!
	}
`
