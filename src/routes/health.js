const express = require('express');
const { getJokesCount } = require('../models/jokes');

const router = express.Router();

// Health check endpoint
router.get('/healthz', (req, res) => {
  res.status(200).json({
    status: "healthy",
    timestamp: new Date().toISOString(),
    uptime: process.uptime(),
    version: "1.0.0",
    jokes_loaded: getJokesCount()
  });
});

module.exports = router;