// next.config.js
const withPWA = require('next-pwa')({
  dest: 'public',
  register: true,
  skipWaiting: true,
  disableDevLogs: true,
  buildExcludes: [/middleware-manifest\.json$/],
  mode: process.env.NODE_ENV === 'development' ? 'development' : 'production',
  disable: process.env.NODE_ENV === 'development',
});

module.exports = withPWA({
  reactStrictMode: true,
  typescript: { ignoreBuildErrors: false },
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://localhost:8181/api/:path*',
      },
      {
        source: '/auth/:path*',
        destination: 'http://localhost:8181/auth/:path*',
      },
    ];
  },
});