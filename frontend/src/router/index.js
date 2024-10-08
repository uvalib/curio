import { createRouter, createWebHistory } from 'vue-router'

const routes = [
   {
      path: '/',
      name: 'home',
      component: () => import('../views/Home.vue')
   },
   {
      path: '/view/:pid',
      name: 'viewer',
      component: () => import('../views/Viewer.vue')
   },
]

const router = createRouter({
   history: createWebHistory(import.meta.env.BASE_URL ),
   routes
})
export default router