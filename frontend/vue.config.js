// NOTES on this are found here:
//    https://cli.vuejs.org/config/#devserver
//    https://github.com/chimurai/http-proxy-middleware#proxycontext-config
module.exports = {
  devServer: {
    host: '0.0.0.0',
    proxy: {
      '/api': {
        target: process.env.CURIO_SRV, // export CURIO_SRV=http://localhost:8185
        changeOrigin: true,
        logLevel: 'debug'
      },
      '/oembed': {
        target: process.env.CURIO_SRV,
        changeOrigin: true,
        logLevel: 'debug'
      },
      '/version': {
        target: process.env.CURIO_SRV,
        changeOrigin: true,
        logLevel: 'debug'
      },
      '/healthcheck': {
        target: process.env.CURIO_SRV,
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
  },
}
