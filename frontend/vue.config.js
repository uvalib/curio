// NOTES on this are found here:
//    https://cli.vuejs.org/config/#devserver
//    https://github.com/chimurai/http-proxy-middleware#proxycontext-config
module.exports = {
  devServer: {
    // public: process.env.BASE_URL,
    host: '0.0.0.0',
    public: '0.0.0.0:8080',
    disableHostCheck: true,
    proxy: {
      '/api': {
        target: process.env.CURIO_SRV, // or 'http://localhost:8085',
        changeOrigin: true,
        logLevel: 'debug'
      },
      '/oembed': {
        target: process.env.CURIO_SRV, // or 'http://localhost:8085',
        changeOrigin: true,
        logLevel: 'debug'
      },
      '/version': {
        target: process.env.CURIO_SRV, // or 'http://localhost:8085',
        changeOrigin: true,
        logLevel: 'debug'
      },
      '/healthcheck': {
        target: process.env.CURIO_SRV, // or 'http://localhost:8085',
        changeOrigin: true,
        logLevel: 'debug'
      },
    }
  },
  configureWebpack: {
    performance: {
      // bump max sizes to 1024
      maxEntrypointSize: 1024000,
      maxAssetSize: 1024000
    }
  }
}
