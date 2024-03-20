 /*global process */

import { fileURLToPath, URL } from 'url'
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
   define: {
      // enable hydration mismatch details in production build
      __VUE_PROD_HYDRATION_MISMATCH_DETAILS__: 'true'
   },
   plugins: [vue()],
   resolve: {
      alias: {
         '@': fileURLToPath(new URL('./src', import.meta.url))
      }
   },
   server: { // this is used in dev mode only
      port: 8080,
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
          '/healthcheck': {
            target: process.env.CURIO_SRV,
            changeOrigin: true,
            logLevel: 'debug'
          },
          '/version': {
            target: process.env.CURIO_SRV,
            changeOrigin: true,
            logLevel: 'debug'
          },
      }
   },
   //  configureWebpack: {
   //    performance: {
   //      // bump max sizes to 1024
   //      maxEntrypointSize: 1024000,
   //      maxAssetSize: 1024000
   //    }
   //  },
})


