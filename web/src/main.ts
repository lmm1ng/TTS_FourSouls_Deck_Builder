import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import naive from 'naive-ui'
import { createPinia } from 'pinia'

createApp(App).use(router).use(createPinia()).use(naive).mount('#app')
