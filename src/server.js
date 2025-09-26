const app = require('./app');

const PORT = process.env.PORT || 3009;

app.listen(PORT, () => {
  const protocol = process.env.HTTPS === 'true' ? 'https' : 'http';
  const host = process.env.HOST || 'localhost';
  const baseUrl = `${protocol}://${host}:${PORT}`;
  
  console.log(`JOKES API server is running on port ${PORT}`);
  console.log(`Visit ${baseUrl} for available endpoints`);
  console.log(`Health check: ${baseUrl}/healthz`);
  console.log(`API docs: ${baseUrl}/docs`);
});

// Graceful shutdown
process.on('SIGTERM', () => {
  console.log('SIGTERM received, shutting down gracefully');
  process.exit(0);
});

process.on('SIGINT', () => {
  console.log('SIGINT received, shutting down gracefully');  
  process.exit(0);
});