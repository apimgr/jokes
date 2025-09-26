const express = require('express');
const helmet = require('helmet');
const cors = require('cors');
const rateLimit = require('express-rate-limit');

// Import route handlers
const jokesRoutes = require('./routes/jokes');
const healthRoutes = require('./routes/health');
const docsRoutes = require('./routes/docs');

const app = express();

// Security middleware
app.use(helmet());
app.use(cors());
app.use(express.json());

// Rate limiting
const limiter = rateLimit({
  windowMs: 60 * 60 * 1000,
  max: 2000,
  message: {
    type: "error",
    value: "Too many requests from this IP, please try again after an hour."
  },
  standardHeaders: true,
  legacyHeaders: false,
  statusCode: 429
});

app.use(limiter);

// Handle unsupported HTTP methods early
app.use((req, res, next) => {
  if (req.method !== 'GET' && req.method !== 'HEAD' && req.method !== 'OPTIONS') {
    return res.status(405).json({
      type: "error", 
      value: `Method ${req.method} not allowed. Only GET, HEAD, and OPTIONS are supported.`
    });
  }
  next();
});

// Mount routes
app.use('/', healthRoutes);
app.use('/', docsRoutes);
app.use('/api/v1', jokesRoutes);
app.use('/', require('./routes/legacy')); // Legacy redirects

// Handle 404 for unknown endpoints
app.use((req, res) => {
  res.status(404).json({
    type: "error",
    value: "Endpoint not found"
  });
});

// Error handler
app.use((err, req, res, next) => {
  console.error(err.stack);
  res.status(500).json({
    type: "error",
    value: "Internal server error"
  });
});

module.exports = app;