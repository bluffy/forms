import { createRouter, createWebHistory,  } from 'vue-router'
import { defineAsyncComponent } from 'vue'


import NotFound from '../views/NotFound.vue'


const routes = [
    {
        name: "login",
        path: "/login",
        meta: {
          analyticsIgnore: true
        },
        component: defineAsyncComponent(() => import('../views/Login.vue'))
      },    
      {
        name: "NotFound",
        path: "/:pathMatch(.*)",
        component: NotFound,
      }   
    ]    

const router = createRouter({
  history: createWebHistory(),
  routes
})
 //  history: createWebHistory(),
 

export default router