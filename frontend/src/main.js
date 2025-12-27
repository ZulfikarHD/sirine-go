import { createApp } from 'vue'
import './style.css'
import App from './App.vue'

// Register Service Worker for PWA
if ('serviceWorker' in navigator) {
  window.addEventListener('load', () => {
    navigator.serviceWorker.register('/sw.js').then(
      (registration) => {
        console.log('Service Worker registered:', registration)
      },
      (error) => {
        console.log('Service Worker registration failed:', error)
      }
    )
  })
}

createApp(App).mount('#app')
