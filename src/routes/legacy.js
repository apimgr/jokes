const express = require('express');

const router = express.Router();

// Legacy endpoints (redirect to v1)
router.get('/jokes/count', (req, res) => {
  res.redirect(301, '/api/v1/jokes/count');
});

router.get('/jokes/categories', (req, res) => {
  res.redirect(301, '/api/v1/jokes/categories');
});

router.get('/jokes/random/:count?', (req, res) => {
  const queryString = req.url.split('?')[1] ? `?${req.url.split('?')[1]}` : '';
  if (req.params.count) {
    res.redirect(301, `/api/v1/jokes/random/${req.params.count}${queryString}`);
  } else {
    res.redirect(301, `/api/v1/jokes/random${queryString}`);
  }
});

router.get('/jokes/:id', (req, res) => {
  const queryString = req.url.split('?')[1] ? `?${req.url.split('?')[1]}` : '';
  res.redirect(301, `/api/v1/jokes/${req.params.id}${queryString}`);
});

router.get('/jokes', (req, res) => {
  const queryString = req.url.split('?')[1] ? `?${req.url.split('?')[1]}` : '';
  res.redirect(301, `/api/v1/jokes${queryString}`);
});

router.get('/jokes/all', (req, res) => {
  const queryString = req.url.split('?')[1] ? `?${req.url.split('?')[1]}` : '';
  res.redirect(301, `/api/v1/jokes/all${queryString}`);
});

module.exports = router;