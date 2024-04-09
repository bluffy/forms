import { createRouter, createWebHistory, } from 'vue-router'
import { defineAsyncComponent } from 'vue'


import NotFound from '../views/NotFound.vue'


const routes = [
  {
    name: "register",
    path: "/register",
    component: defineAsyncComponent(() => import('../views/Register.vue'))
  },
  {
    name: "forgot_password",
    path: "/forgot_password",
    component: defineAsyncComponent(() => import('../views/ForgotPassword.vue'))
  },
  {
    name: "login",
    path: "/login",
    component: defineAsyncComponent(() => import('../views/Login.vue'))
  },
  {
    name: "dashboard",
    path: "/",
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