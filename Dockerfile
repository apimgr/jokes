# Use official Node.js runtime as base image
FROM node:22-alpine

# Set working directory in container
WORKDIR /app

# Copy package files first for better Docker layer caching
COPY package*.json ./

# Install dependencies
RUN npm ci --only=production

# Copy application source code
COPY . .

# Create non-root user for security
RUN addgroup -g 1001 -S nodejs && adduser -S jokes -u 1001

# Change ownership of app directory to non-root user
RUN chown -R jokes:nodejs /app
USER jokes

# Expose port 3009
EXPOSE 3009

# Add health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD node -e "const http = require('http'); \
    const options = { host: 'localhost', port: 3009, path: '/healthz', timeout: 2000 }; \
    const req = http.request(options, (res) => { \
      if (res.statusCode === 200) { process.exit(0); } \
      else { process.exit(1); } \
    }); \
    req.on('error', () => process.exit(1)); \
    req.end();"

# Start the application
CMD ["npm", "start"]
