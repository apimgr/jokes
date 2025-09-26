const express = require('express');
const {
  getRandomJoke,
  getJokeById,
  getJokesByCategory,
  getRandomJokeByCategory,
  getAllCategories,
  getRandomJokes,
  replaceNameInJoke,
  getJokesCount
} = require('../models/jokes');

const router = express.Router();

function createResponse(type, value) {
  return { type, value };
}

function handleNameReplacement(joke, firstName, lastName) {
  if (firstName || lastName) {
    const first = firstName || "Chuck";
    const last = lastName || "Norris";
    return {
      ...joke,
      joke: replaceNameInJoke(joke.joke, first, last)
    };
  }
  return joke;
}

// Get joke count and statistics
router.get('/jokes/count', (req, res) => {
  res.status(200).json(createResponse("success", {
    total: getJokesCount(),
    categories: getAllCategories().map(cat => ({
      name: cat,
      count: getJokesByCategory(cat).length
    }))
  }));
});

// Get all categories
router.get('/jokes/categories', (req, res) => {
  const categories = getAllCategories();
  res.status(200).json(createResponse("success", categories));
});

// Get random joke
router.get('/jokes/random', (req, res) => {
  let { firstName, lastName, category, exclude, limitTo } = req.query;
  
  // Handle invalid category by defaulting to all (excluding explicit)
  let validCategory = category;
  if (category && !getAllCategories().includes(category)) {
    validCategory = null; // Will use all categories
    // Add explicit to exclude list if not already there
    if (!exclude || !exclude.includes('explicit')) {
      exclude = exclude ? `${exclude},explicit` : 'explicit';
    }
  }
  
  // Handle limitTo parameter to restrict joke pool
  let availableJokes = null;
  if (limitTo) {
    try {
      // Parse limitTo parameter - handle both [category] and [category1,category2] formats
      const cleanLimitTo = limitTo.replace(/[\[\]]/g, '').trim();
      if (!cleanLimitTo) {
        return res.status(400).json(createResponse("error", "limitTo cannot be empty"));
      }
      
      const limitToCategories = cleanLimitTo.split(',').map(cat => cat.trim()).filter(cat => cat.length > 0);
      
      if (limitToCategories.length === 0) {
        return res.status(400).json(createResponse("error", "limitTo must contain at least one category"));
      }
      
      // Validate categories
      const invalidCategories = limitToCategories.filter(cat => !getAllCategories().includes(cat));
      if (invalidCategories.length > 0) {
        return res.status(400).json(createResponse("error", `Invalid limitTo categories: ${invalidCategories.join(', ')}. Available categories: ${getAllCategories().join(', ')}`));
      }
      
      // Get jokes from limitTo categories
      const jokesModule = require('../models/jokes');
      availableJokes = jokesModule.jokes.filter(joke => joke.categories.some(cat => limitToCategories.includes(cat)));
      
      if (availableJokes.length === 0) {
        return res.status(404).json(createResponse("error", "No jokes found for the specified limitTo categories"));
      }
      
    } catch (error) {
      return res.status(400).json(createResponse("error", "Invalid limitTo format. Use [category] or [category1,category2]"));
    }
  }
  
  let joke;
  
  if (validCategory) {
    if (availableJokes) {
      // Filter available jokes by specific category
      const categoryJokes = availableJokes.filter(j => j.categories.includes(validCategory));
      if (categoryJokes.length === 0) {
        return res.status(404).json(createResponse("error", "No jokes found matching both category and limitTo criteria"));
      }
      joke = categoryJokes[Math.floor(Math.random() * categoryJokes.length)];
    } else {
      joke = getRandomJokeByCategory(validCategory);
      if (!joke) {
        return res.status(404).json(createResponse("error", "No jokes found for the specified category"));
      }
    }
  } else {
    if (availableJokes) {
      // Pick random from limited set
      joke = availableJokes[Math.floor(Math.random() * availableJokes.length)];
    } else {
      joke = getRandomJoke();
      if (!joke) {
        return res.status(503).json(createResponse("error", "Service temporarily unavailable - no jokes available"));
      }
    }
  }
  
  if (exclude) {
    const excludeCategories = exclude.split(',').map(cat => cat.trim());
    
    // Validate exclude categories
    const invalidCategories = excludeCategories.filter(cat => !getAllCategories().includes(cat));
    if (invalidCategories.length > 0) {
      return res.status(400).json(createResponse("error", `Invalid exclude categories: ${invalidCategories.join(', ')}`));
    }
    
    const hasExcludedCategory = joke.categories.some(cat => excludeCategories.includes(cat));
    
    if (hasExcludedCategory) {
      let attempts = 0;
      const maxAttempts = 50;
      
      do {
        joke = validCategory ? getRandomJokeByCategory(validCategory) : getRandomJoke();
        attempts++;
      } while (
        joke && 
        joke.categories.some(cat => excludeCategories.includes(cat)) && 
        attempts < maxAttempts
      );
      
      if (attempts >= maxAttempts || !joke) {
        return res.status(404).json(createResponse("error", "No jokes found matching the criteria"));
      }
    }
  }
  
  const modifiedJoke = handleNameReplacement(joke, firstName, lastName);
  res.status(200).json(createResponse("success", modifiedJoke));
});

// Get multiple random jokes
router.get('/jokes/random/:count', (req, res) => {
  const count = parseInt(req.params.count);
  const { firstName, lastName, exclude, limitTo } = req.query;
  
  if (isNaN(count) || count < 1 || count > 100) {
    return res.status(400).json(createResponse("error", "Count must be a number between 1 and 100"));
  }
  
  let jokes;
  
  // Handle limitTo parameter to restrict joke pool
  if (limitTo) {
    try {
      // Parse limitTo parameter - handle both [category] and [category1,category2] formats
      const cleanLimitTo = limitTo.replace(/[\[\]]/g, '').trim();
      if (!cleanLimitTo) {
        return res.status(400).json(createResponse("error", "limitTo cannot be empty"));
      }
      
      const limitToCategories = cleanLimitTo.split(',').map(cat => cat.trim()).filter(cat => cat.length > 0);
      
      if (limitToCategories.length === 0) {
        return res.status(400).json(createResponse("error", "limitTo must contain at least one category"));
      }
      
      // Validate categories
      const invalidCategories = limitToCategories.filter(cat => !getAllCategories().includes(cat));
      if (invalidCategories.length > 0) {
        return res.status(400).json(createResponse("error", `Invalid limitTo categories: ${invalidCategories.join(', ')}. Available categories: ${getAllCategories().join(', ')}`));
      }
      
      // Get jokes from limitTo categories
      const jokesModule = require('../models/jokes');
      const availableJokes = jokesModule.jokes.filter(joke => joke.categories.some(cat => limitToCategories.includes(cat)));
      
      if (availableJokes.length === 0) {
        return res.status(404).json(createResponse("error", "No jokes found for the specified limitTo categories"));
      }
      
      // Get random jokes from limited set
      const result = [];
      const usedIndices = new Set();
      
      while (result.length < count && result.length < availableJokes.length) {
        const randomIndex = Math.floor(Math.random() * availableJokes.length);
        if (!usedIndices.has(randomIndex)) {
          usedIndices.add(randomIndex);
          result.push(availableJokes[randomIndex]);
        }
      }
      
      jokes = result;
      
    } catch (error) {
      return res.status(400).json(createResponse("error", "Invalid limitTo format. Use [category] or [category1,category2]"));
    }
  } else {
    jokes = getRandomJokes(count);
  }
  
  if (exclude) {
    const excludeCategories = exclude.split(',').map(cat => cat.trim());
    jokes = jokes.filter(joke => !joke.categories.some(cat => excludeCategories.includes(cat)));
    
    if (jokes.length === 0) {
      return res.status(404).json(createResponse("error", "No jokes found matching the criteria"));
    }
  }
  
  const modifiedJokes = jokes.map(joke => handleNameReplacement(joke, firstName, lastName));
  res.status(200).json(createResponse("success", modifiedJokes));
});

// Get all jokes endpoint  
router.get('/jokes/all', (req, res) => {
  const { limitTo, firstName, lastName, exclude } = req.query;
  
  // Start with all jokes - get fresh copy from jokes module
  const jokesModule = require('../models/jokes');
  let jokes = [...jokesModule.jokes];
  
  // Apply limitTo category filter if provided
  if (limitTo) {
    let limitToCategories;
    
    // Parse limitTo parameter - handle both [category] and [category1,category2] formats
    try {
      // Remove brackets and split by comma
      const cleanLimitTo = limitTo.replace(/[\[\]]/g, '').trim();
      if (!cleanLimitTo) {
        return res.status(400).json(createResponse("error", "limitTo cannot be empty"));
      }
      
      limitToCategories = cleanLimitTo.split(',').map(cat => cat.trim()).filter(cat => cat.length > 0);
      
      if (limitToCategories.length === 0) {
        return res.status(400).json(createResponse("error", "limitTo must contain at least one category"));
      }
      
      // Validate categories
      const invalidCategories = limitToCategories.filter(cat => !getAllCategories().includes(cat));
      if (invalidCategories.length > 0) {
        return res.status(400).json(createResponse("error", `Invalid limitTo categories: ${invalidCategories.join(', ')}. Available categories: ${getAllCategories().join(', ')}`));
      }
      
      // Filter jokes to only include those with specified categories
      jokes = jokes.filter(joke => joke.categories.some(cat => limitToCategories.includes(cat)));
      
    } catch (error) {
      return res.status(400).json(createResponse("error", "Invalid limitTo format. Use [category] or [category1,category2]"));
    }
  }
  
  // Apply exclude filter if provided
  if (exclude) {
    const excludeCategories = exclude.split(',').map(cat => cat.trim());
    
    // Validate exclude categories
    const invalidCategories = excludeCategories.filter(cat => !getAllCategories().includes(cat));
    if (invalidCategories.length > 0) {
      return res.status(400).json(createResponse("error", `Invalid exclude categories: ${invalidCategories.join(', ')}`));
    }
    
    jokes = jokes.filter(joke => !joke.categories.some(cat => excludeCategories.includes(cat)));
  }
  
  // Apply name replacement if requested
  if (firstName || lastName) {
    jokes = jokes.map(joke => handleNameReplacement(joke, firstName, lastName));
  }
  
  // Return response with metadata
  res.status(200).json(createResponse("success", {
    jokes: jokes,
    meta: {
      total_in_database: getJokesCount(),
      returned: jokes.length,
      limited_to_categories: limitTo ? limitTo.replace(/[\[\]]/g, '').split(',').map(cat => cat.trim()) : null,
      excluded_categories: exclude ? exclude.split(',').map(cat => cat.trim()) : []
    }
  }));
});

// Get specific joke by ID
router.get('/jokes/:id', (req, res) => {
  const id = parseInt(req.params.id);
  const { firstName, lastName } = req.query;
  
  if (isNaN(id)) {
    return res.status(400).json(createResponse("error", "Invalid joke ID - must be a number"));
  }
  
  if (id < 1) {
    return res.status(400).json(createResponse("error", "Invalid joke ID - must be greater than 0"));
  }
  
  if (id > getJokesCount()) {
    return res.status(404).json(createResponse("error", `Joke ID ${id} not found. Valid range: 1-${getJokesCount()}`));
  }
  
  const joke = getJokeById(id);
  
  if (!joke) {
    return res.status(404).json(createResponse("error", "Joke not found"));
  }
  
  const modifiedJoke = handleNameReplacement(joke, firstName, lastName);
  res.status(200).json(createResponse("success", modifiedJoke));
});

// Get jokes by category
router.get('/jokes', (req, res) => {
  let { category, exclude } = req.query;
  
  let jokes;
  
  if (category) {
    // Handle invalid category by defaulting to all (excluding explicit)
    if (!getAllCategories().includes(category)) {
      // Get all jokes instead of specific category
      const jokesModule = require('../models/jokes');
      jokes = [...jokesModule.jokes];
      // Add explicit to exclude list if not already there
      if (!exclude || !exclude.includes('explicit')) {
        exclude = exclude ? `${exclude},explicit` : 'explicit';
      }
    } else {
      jokes = getJokesByCategory(category);
      if (jokes.length === 0) {
        return res.status(404).json(createResponse("error", "No jokes found for the specified category"));
      }
    }
  } else {
    return res.status(400).json(createResponse("error", "Please specify a category or use /jokes/random"));
  }
  
  if (exclude) {
    const excludeCategories = exclude.split(',').map(cat => cat.trim());
    jokes = jokes.filter(joke => !joke.categories.some(cat => excludeCategories.includes(cat)));
    
    if (jokes.length === 0) {
      return res.status(404).json(createResponse("error", "No jokes found matching the criteria"));
    }
  }
  
  res.status(200).json(createResponse("success", jokes));
});

module.exports = router;