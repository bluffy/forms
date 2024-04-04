import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { viteStaticCopy } from 'vite-plugin-static-copy'
import { fileURLToPath, URL } from 'node:url'


// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue(),viteStaticCopy({
    targets: [
      {
        src: 'public/.gitkeep',
        dest: '../../app/client'
      }
    ]
  })],
  build: {
    outDir: '../../app/client'
  },  
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
})
