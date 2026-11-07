package swagger

// SwaggerJSON documents the public, read-only API routes served by this binary.
const SwaggerJSON = `{
  "swagger": "2.0",
  "info": {
    "title": "Jokes API",
    "description": "Public read-only access to the curated joke collection.",
    "version": "1.0.0",
    "license": {"name": "MIT", "url": "https://github.com/apimgr/jokes/blob/main/LICENSE.md"}
  },
  "host": "jokes.apimgr.us",
  "basePath": "/api/v1",
  "schemes": ["https", "http"],
  "consumes": ["application/json"],
  "produces": ["application/json"],
  "paths": {
    "/jokes/random": {"get": {"summary": "Get a random joke", "parameters": [{"name": "category", "in": "query", "type": "string"}, {"name": "limitTo", "in": "query", "type": "string"}, {"name": "exclude", "in": "query", "type": "string"}], "responses": {"200": {"description": "Joke response"}, "400": {"description": "Invalid filter"}, "404": {"description": "No matching joke"}}}},
    "/jokes/random/{count}": {"get": {"summary": "Get up to 100 random jokes", "parameters": [{"name": "count", "in": "path", "required": true, "type": "integer", "minimum": 1, "maximum": 100}], "responses": {"200": {"description": "Joke list"}, "400": {"description": "Invalid count"}}}},
    "/jokes/{id}": {"get": {"summary": "Get a joke by numeric ID", "parameters": [{"name": "id", "in": "path", "required": true, "type": "integer", "minimum": 1}], "responses": {"200": {"description": "Joke response"}, "404": {"description": "Joke not found"}}}},
    "/jokes/all": {"get": {"summary": "Browse jokes", "parameters": [{"name": "page", "in": "query", "type": "integer", "minimum": 1}, {"name": "limit", "in": "query", "type": "integer", "minimum": 1, "maximum": 100}, {"name": "q", "in": "query", "type": "string"}, {"name": "limitTo", "in": "query", "type": "string"}, {"name": "exclude", "in": "query", "type": "string"}], "responses": {"200": {"description": "Paginated joke list"}, "400": {"description": "Invalid query"}}}},
    "/jokes/search": {"get": {"summary": "Search jokes", "parameters": [{"name": "q", "in": "query", "required": true, "type": "string"}, {"name": "category", "in": "query", "type": "string"}, {"name": "page", "in": "query", "type": "integer"}, {"name": "limit", "in": "query", "type": "integer"}], "responses": {"200": {"description": "Paginated search results"}, "400": {"description": "Invalid query"}}}},
    "/jokes/categories": {"get": {"summary": "List categories with counts", "responses": {"200": {"description": "Category names and joke counts"}}}},
    "/jokes/categories/{category}": {"get": {"summary": "Browse one category", "parameters": [{"name": "category", "in": "path", "required": true, "type": "string"}, {"name": "page", "in": "query", "type": "integer"}, {"name": "limit", "in": "query", "type": "integer"}], "responses": {"200": {"description": "Paginated jokes"}, "400": {"description": "Unknown category"}}}},
    "/server/healthz": {"get": {"summary": "Check service health", "responses": {"200": {"description": "Service health"}}}},
    "/server/swagger": {"get": {"summary": "Get this Swagger document", "responses": {"200": {"description": "Swagger 2.0 JSON"}}}},
    "/server/graphql": {"post": {"summary": "Execute a GraphQL query", "parameters": [{"name": "body", "in": "body", "required": true, "schema": {"type": "object", "required": ["query"], "properties": {"query": {"type": "string"}, "variables": {"type": "object"}}}}], "responses": {"200": {"description": "GraphQL response"}, "400": {"description": "Invalid request"}}}}
  }
}`
