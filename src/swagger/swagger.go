package swagger

const SwaggerJSON = `{
  "swagger": "2.0",
  "info": {
    "title": "Jokes API",
    "description": "A comprehensive REST API serving over 5,000 jokes across 16 categories",
    "version": "1.0.0",
    "contact": {
      "name": "APIMGR",
      "url": "https://jokes.apimgr.us"
    },
    "license": {
      "name": "MIT",
      "url": "https://github.com/apimgr/jokes/blob/main/LICENSE.md"
    }
  },
  "host": "jokes.apimgr.us",
  "basePath": "/api/v1",
  "schemes": ["https", "http"],
  "consumes": ["application/json"],
  "produces": ["application/json"],
  "paths": {
    "/jokes/random": {
      "get": {
        "summary": "Get a random joke",
        "description": "Returns a random joke, optionally filtered by category",
        "parameters": [
          {
            "name": "firstName",
            "in": "query",
            "description": "Replace 'Chuck' with this name",
            "type": "string"
          },
          {
            "name": "lastName",
            "in": "query",
            "description": "Replace 'Norris' with this name",
            "type": "string"
          },
          {
            "name": "category",
            "in": "query",
            "description": "Filter by specific category",
            "type": "string"
          },
          {
            "name": "limitTo",
            "in": "query",
            "description": "Limit to categories. Format: [category] or [cat1,cat2]",
            "type": "string"
          },
          {
            "name": "exclude",
            "in": "query",
            "description": "Exclude categories. Format: cat1,cat2",
            "type": "string"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "$ref": "#/definitions/JokeResponse"
            }
          },
          "400": {
            "description": "Bad request",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          },
          "404": {
            "description": "No jokes found",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/jokes/random/{count}": {
      "get": {
        "summary": "Get multiple random jokes",
        "description": "Returns multiple random jokes (1-100)",
        "parameters": [
          {
            "name": "count",
            "in": "path",
            "required": true,
            "description": "Number of jokes to return (1-100)",
            "type": "integer"
          },
          {
            "name": "firstName",
            "in": "query",
            "description": "Replace 'Chuck' with this name",
            "type": "string"
          },
          {
            "name": "lastName",
            "in": "query",
            "description": "Replace 'Norris' with this name",
            "type": "string"
          },
          {
            "name": "limitTo",
            "in": "query",
            "description": "Limit to categories",
            "type": "string"
          },
          {
            "name": "exclude",
            "in": "query",
            "description": "Exclude categories",
            "type": "string"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "$ref": "#/definitions/JokesResponse"
            }
          }
        }
      }
    },
    "/jokes/{id}": {
      "get": {
        "summary": "Get joke by ID",
        "description": "Returns a specific joke by its ID",
        "parameters": [
          {
            "name": "id",
            "in": "path",
            "required": true,
            "description": "Joke ID (1-5160)",
            "type": "integer"
          },
          {
            "name": "firstName",
            "in": "query",
            "description": "Replace 'Chuck' with this name",
            "type": "string"
          },
          {
            "name": "lastName",
            "in": "query",
            "description": "Replace 'Norris' with this name",
            "type": "string"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "$ref": "#/definitions/JokeResponse"
            }
          },
          "404": {
            "description": "Joke not found",
            "schema": {
              "$ref": "#/definitions/ErrorResponse"
            }
          }
        }
      }
    },
    "/jokes/all": {
      "get": {
        "summary": "Get all jokes",
        "description": "Returns all jokes with optional filtering",
        "parameters": [
          {
            "name": "limitTo",
            "in": "query",
            "description": "Limit to categories",
            "type": "string"
          },
          {
            "name": "exclude",
            "in": "query",
            "description": "Exclude categories",
            "type": "string"
          },
          {
            "name": "firstName",
            "in": "query",
            "description": "Replace 'Chuck' with this name",
            "type": "string"
          },
          {
            "name": "lastName",
            "in": "query",
            "description": "Replace 'Norris' with this name",
            "type": "string"
          }
        ],
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "$ref": "#/definitions/AllJokesResponse"
            }
          }
        }
      }
    },
    "/jokes/categories": {
      "get": {
        "summary": "Get all categories",
        "description": "Returns a list of all joke categories",
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "$ref": "#/definitions/CategoriesResponse"
            }
          }
        }
      }
    },
    "/jokes/count": {
      "get": {
        "summary": "Get joke statistics",
        "description": "Returns total joke count and per-category statistics",
        "responses": {
          "200": {
            "description": "Successful response",
            "schema": {
              "$ref": "#/definitions/CountResponse"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "Joke": {
      "type": "object",
      "properties": {
        "id": {
          "type": "integer",
          "description": "Unique joke ID"
        },
        "joke": {
          "type": "string",
          "description": "The joke text"
        },
        "categories": {
          "type": "array",
          "items": {
            "type": "string"
          },
          "description": "Categories this joke belongs to"
        }
      }
    },
    "JokeResponse": {
      "type": "object",
      "properties": {
        "type": {
          "type": "string",
          "enum": ["success", "error"]
        },
        "value": {
          "$ref": "#/definitions/Joke"
        }
      }
    },
    "JokesResponse": {
      "type": "object",
      "properties": {
        "type": {
          "type": "string",
          "enum": ["success", "error"]
        },
        "value": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/Joke"
          }
        }
      }
    },
    "AllJokesResponse": {
      "type": "object",
      "properties": {
        "type": {
          "type": "string",
          "enum": ["success"]
        },
        "value": {
          "type": "object",
          "properties": {
            "jokes": {
              "type": "array",
              "items": {
                "$ref": "#/definitions/Joke"
              }
            },
            "meta": {
              "type": "object",
              "properties": {
                "total_in_database": {
                  "type": "integer"
                },
                "returned": {
                  "type": "integer"
                },
                "limited_to_categories": {
                  "type": "array",
                  "items": {
                    "type": "string"
                  }
                },
                "excluded_categories": {
                  "type": "array",
                  "items": {
                    "type": "string"
                  }
                }
              }
            }
          }
        }
      }
    },
    "CategoriesResponse": {
      "type": "object",
      "properties": {
        "type": {
          "type": "string",
          "enum": ["success"]
        },
        "value": {
          "type": "array",
          "items": {
            "type": "string"
          }
        }
      }
    },
    "CountResponse": {
      "type": "object",
      "properties": {
        "type": {
          "type": "string",
          "enum": ["success"]
        },
        "value": {
          "type": "object",
          "properties": {
            "total": {
              "type": "integer"
            },
            "categories": {
              "type": "array",
              "items": {
                "type": "object",
                "properties": {
                  "name": {
                    "type": "string"
                  },
                  "count": {
                    "type": "integer"
                  }
                }
              }
            }
          }
        }
      }
    },
    "ErrorResponse": {
      "type": "object",
      "properties": {
        "type": {
          "type": "string",
          "enum": ["error"]
        },
        "value": {
          "type": "string",
          "description": "Error message"
        }
      }
    }
  }
}`
