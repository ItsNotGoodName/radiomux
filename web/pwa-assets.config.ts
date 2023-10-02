import { defaultAssetName, defineConfig, minimalPreset } from '@vite-pwa/assets-generator/config'


export default defineConfig({
  preset: {
    ...minimalPreset,
    apple: {
      sizes: [180],
      padding: 0,
    },
    assetName: (type, size) => {
      if (type != 'apple') {
        return defaultAssetName(type, size)
      }
      return 'apple-touch-icon.png'
    }
  },
  images: [
    'public/favicon.svg',
  ]
})
