import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Viewer from '../views/Viewer.vue'


const routes = [
   {
      path: '/',
      name: 'home',
      component: Home
   },
   {
      path: '/view/:pid',
      name: 'viewer',
      component: Viewer
   },
]

const router = createRouter({
   history: createWebHistory(),
   routes
})
export default router