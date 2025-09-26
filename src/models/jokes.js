const fs = require('fs');
const path = require('path');

// Load jokes from JSON file
let jokesData;
try {
  const jokesPath = path.join(__dirname, 'jokes.json');
  const rawData = fs.readFileSync(jokesPath, 'utf8');
  jokesData = JSON.parse(rawData);
} catch (error) {
  console.error('Error loading jokes.json:', error);
  // Fallback to empty array if file can't be loaded
  jokesData = { jokes: [], categories: [] };
}

const jokes = jokesData.jokes || [];
const categories = jokesData.categories || ["explicit", "nerdy", "movie", "history", "animal", "food", "sports", "work", "travel", "music", "medical", "lawyer", "school", "science"];

function getRandomJoke() {
  if (jokes.length === 0) return null;
  return jokes[Math.floor(Math.random() * jokes.length)];
}

function getJokeById(id) {
  return jokes.find(joke => joke.id === parseInt(id));
}

function getJokesByCategory(category) {
  return jokes.filter(joke => joke.categories.includes(category));
}

function getRandomJokeByCategory(category) {
  const categoryJokes = getJokesByCategory(category);
  if (categoryJokes.length === 0) return null;
  return categoryJokes[Math.floor(Math.random() * categoryJokes.length)];
}

function getAllCategories() {
  return categories;
}

function getRandomJokes(count) {
  if (jokes.length === 0) return [];
  
  const result = [];
  const usedIndices = new Set();
  
  while (result.length < count && result.length < jokes.length) {
    const randomIndex = Math.floor(Math.random() * jokes.length);
    if (!usedIndices.has(randomIndex)) {
      usedIndices.add(randomIndex);
      result.push(jokes[randomIndex]);
    }
  }
  
  return result;
}

function replaceNameInJoke(joke, firstName, lastName) {
  let modifiedJoke = joke.replace(/Chuck Norris/g, `${firstName} ${lastName}`);
  modifiedJoke = modifiedJoke.replace(/Chuck/g, firstName);
  return modifiedJoke;
}

function getJokesCount() {
  return jokes.length;
}

function getJokesByRange(start, end) {
  return jokes.slice(start - 1, end);
}

// Log the number of jokes loaded
console.log(`Loaded ${jokes.length} jokes from jokes.json`);

module.exports = {
  jokes,
  categories,
  getRandomJoke,
  getJokeById,
  getJokesByCategory,
  getRandomJokeByCategory,
  getAllCategories,
  getRandomJokes,
  replaceNameInJoke,
  getJokesCount,
  getJokesByRange
};