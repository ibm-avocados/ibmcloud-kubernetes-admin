const { createProxyMiddleware } = require('http-proxy-middleware');
const proxy = createProxyMiddleware({
  target: 'http://localhost:9030',
  changeOrigin: true,
});
module.exports = function (app) {
  app.use('/api', proxy);
  app.use('/auth', proxy);
};