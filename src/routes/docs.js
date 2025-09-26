const express = require('express');
const { getAllCategories, getJokesCount } = require('../models/jokes');

const router = express.Router();

// Root endpoint with API documentation
router.get('/', (req, res) => {
  res.status(200).json({
    message: "Welcome to the JOKES API",
    version: "1.0.0",
    stats: {
      total_jokes: getJokesCount(),
      categories: getAllCategories().length
    },
    endpoints: {
      "/healthz": "Health check endpoint",
      "/api/v1/jokes/random": "Get a random joke",
      "/api/v1/jokes/random/:count": "Get multiple random jokes (1-100)",
      "/api/v1/jokes/:id": "Get a specific joke by ID",
      "/api/v1/jokes/all": "Get all jokes with optional limiting",
      "/api/v1/jokes/categories": "Get all available categories",
      "/api/v1/jokes/count": "Get total number of jokes and category stats",
      "/docs": "API documentation"
    },
    legacy_endpoints: {
      "/jokes/*": "Legacy endpoints (deprecated, use /api/v1/jokes/*)"
    }
  });
});

// API Documentation endpoint
router.get('/docs', (req, res) => {
  res.status(200).json({
    title: "JOKES API Documentation",
    version: "1.0.0",
    description: "A comprehensive jokes API with 5000+ jokes across multiple categories",
    base_url: `${req.protocol}://${req.get('host')}`,
    
    health_check: {
      endpoint: "/healthz",
      method: "GET",
      description: "Health check endpoint for monitoring",
      response: {
        status: "healthy",
        timestamp: "2025-09-18T12:15:00.000Z",
        uptime: 3600,
        version: "1.0.0",
        jokes_loaded: getJokesCount()
      }
    },
    
    endpoints: {
      "/api/v1/jokes/random": {
        method: "GET",
        description: "Get a random joke",
        parameters: {
          firstName: "Optional. Replace 'Chuck' with this name",
          lastName: "Optional. Replace 'Norris' with this name", 
          category: "Optional. Filter by category",
          limitTo: "Optional. Array of categories to limit random selection to. Format: [category] or [category1,category2]",
          exclude: "Optional. Comma-separated categories to exclude"
        },
        example: "/api/v1/jokes/random?limitTo=[animal,nerdy]&firstName=John"
      },
      
      "/api/v1/jokes/random/{count}": {
        method: "GET",
        description: "Get multiple random jokes (1-100)",
        parameters: {
          count: "Required. Number of jokes (1-100)",
          firstName: "Optional. Replace 'Chuck' with this name",
          lastName: "Optional. Replace 'Norris' with this name",
          limitTo: "Optional. Array of categories to limit random selection to. Format: [category] or [category1,category2]",
          exclude: "Optional. Comma-separated categories to exclude"
        },
        example: "/api/v1/jokes/random/5?limitTo=[animal,nerdy]&exclude=explicit"
      },
      
      "/api/v1/jokes/{id}": {
        method: "GET", 
        description: "Get a specific joke by ID (1-5160)",
        parameters: {
          id: "Required. Joke ID (1-5160)",
          firstName: "Optional. Replace 'Chuck' with this name",
          lastName: "Optional. Replace 'Norris' with this name"
        },
        example: "/api/v1/jokes/1?firstName=John&lastName=Doe"
      },
      
      "/api/v1/jokes/categories": {
        method: "GET",
        description: "Get all available joke categories",
        example: "/api/v1/jokes/categories"
      },
      
      "/api/v1/jokes/count": {
        method: "GET",
        description: "Get total number of jokes and category statistics",
        example: "/api/v1/jokes/count"
      },
      
      "/api/v1/jokes": {
        method: "GET",
        description: "Get jokes by category",
        parameters: {
          category: "Required. Category name",
          exclude: "Optional. Comma-separated categories to exclude"
        },
        example: "/api/v1/jokes?category=nerdy"
      },
      
      "/api/v1/jokes/all": {
        method: "GET",
        description: "Get all jokes with optional category filtering",
        parameters: {
          limitTo: "Optional. Array of categories to limit results to. Format: [category] or [category1,category2]",
          firstName: "Optional. Replace 'Chuck' with this name",
          lastName: "Optional. Replace 'Norris' with this name",
          exclude: "Optional. Comma-separated categories to exclude"
        },
        example: "/api/v1/jokes/all?limitTo=[animal,nerdy]&exclude=explicit",
        response: {
          type: "success",
          value: {
            jokes: "Array of joke objects",
            meta: {
              total_in_database: "Total jokes in database",
              returned: "Number of jokes returned",
              limited_to_categories: "Array of categories results were limited to (null if no limit)",
              excluded_categories: "Array of excluded categories"
            }
          }
        }
      }
    },
    
    categories: getAllCategories(),
    
    response_format: {
      success: {
        type: "success",
        value: "Joke object or array of jokes"
      },
      error: {
        type: "error", 
        value: "Error message"
      }
    },
    
    http_status_codes: {
      200: "OK - Request successful",
      400: "Bad Request - Invalid parameters",
      404: "Not Found - Resource not found",
      405: "Method Not Allowed - Only GET, HEAD, OPTIONS allowed",
      429: "Too Many Requests - Rate limit exceeded (2000/hour)",
      500: "Internal Server Error - Server error"
    },
    
    rate_limiting: {
      limit: "2000 requests per hour per IP",
      headers: {
        "X-RateLimit-Limit": "Rate limit",
        "X-RateLimit-Remaining": "Remaining requests",
        "X-RateLimit-Reset": "Reset time"
      }
    },
    
    legacy_endpoints: {
      note: "Legacy endpoints redirect to /api/v1/* with 301 status",
      "/jokes/*": "Redirects to /api/v1/jokes/*"
    }
  });
});

module.exports = router;