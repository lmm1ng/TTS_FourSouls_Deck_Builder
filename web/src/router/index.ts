import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import GamesView from '@/views/GamesView.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Game',
    component: GamesView
  }
]

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default router
