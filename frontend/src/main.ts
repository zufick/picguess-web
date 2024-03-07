import './assets/main.css'

import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import { createPinia } from 'pinia'
import App from './App.vue'
import RouterWrap from './Router.vue'
import DrawingCanvas from './components/DrawingCanvas.vue'

const pinia = createPinia()
const app = createApp(RouterWrap)


const routes = [
  { path: '/', component: App },
  { path: '/room/:room', component: App },
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})  

app.use(pinia)
app.use(router)
app.mount('#app')