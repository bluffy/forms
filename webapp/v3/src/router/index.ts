import { createRouter, createWebHistory,  } from 'vue-router'
import { defineAsyncComponent } from 'vue'


import NotFound from '../views/NotFound.vue'


const routes = [  
    {
        name: "login",
        path: "/login",
        meta: {
          requiresAuth: false
        },
        component: defineAsyncComponent(() => import('../views/Login.vue'))
      },   
      {
        name: "dashboard",
        path: "/",
        meta: {
          requiresAuth: true
        },        
        component: defineAsyncComponent(() => import('../views/Dashboard.vue'))
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