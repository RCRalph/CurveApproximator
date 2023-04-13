// Plugins
import vue from '@vitejs/plugin-vue'
import vuetify, { transformAssetUrls } from 'vite-plugin-vuetify'
import { copy, mkdirSync, renameSync } from 'fs-extra'

// Utilities
import { defineConfig } from 'vite'
import { fileURLToPath, URL } from 'node:url'

// https://vitejs.dev/config/
export default defineConfig({
  build: {
    outDir: "dist",
    copyPublicDir: false,
    rollupOptions: {
      input: {
        main: './index.html'
      }
    }
  },
  plugins: [
    {
      name: "move-files-after-build",
      closeBundle() {
        mkdirSync("./dist/public")
        copy("./public", "./dist/public")
        renameSync("./dist/assets", "./dist/public/assets");
      }
    },
    vue({
      template: { transformAssetUrls }
    }),
    // https://github.com/vuetifyjs/vuetify-loader/tree/next/packages/vite-plugin
    vuetify({
      autoImport: true,
      styles: {
        configFile: 'src/styles/settings.scss',
      },
    }),
  ],
  define: { 'process.env': {} },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
    extensions: [
      '.js',
      '.json',
      '.jsx',
      '.mjs',
      '.ts',
      '.tsx',
      '.vue',
    ],
  },
  server: {
    port: 3000,
  },
})
