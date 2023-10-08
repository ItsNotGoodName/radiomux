import { defineConfig } from 'vite'
import solid from 'vite-plugin-solid'
import path from "path";
import { macaronVitePlugin } from '@macaron-css/vite';

export default defineConfig({
  plugins: [
    macaronVitePlugin(),
    solid(),
  ],
  server: {
    proxy: {
      "/api": "http://localhost:8080/",
      "/api/ws": {
        target: "http://localhost:8080/",
        ws: true,
      },
      "/rpc": "http://localhost:8080/",
    }
  },
  resolve: {
    alias: {
      "~": path.resolve(__dirname, "./src"),
    },
  },
})
