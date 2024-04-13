import { createRouter, createWebHistory, } from 'vue-router'
import { defineAsyncComponent } from 'vue'


import NotFound from '../views/NotFound.vue'


const routes = [

  {
    name: "register_link",
    path: "/user/register/:link",
    component: defineAsyncComponent(() => import('../views/Redirects.vue'))
  },
  {
    name: "register",
    path: "/user/register",
    component: defineAsyncComponent(() => import('../views/auth/Register.vue'))
  },
  {
    name: "forgot_password_recover_form",
    path: "/user/forgot_password",
    component: defineAsyncComponent(() => import('../views/auth/ForgotPassword.vue'))
  },  
  {
    name: "forgot_password_link",
    path: "/user/forgot_password/:link",
    component: defineAsyncComponent(() => import('../views/auth/ForgotPassword.vue'))
  },  
  {
    name: "login",
    path: "/login",
    component: defineAsyncComponent(() => import('../views/auth/Login.vue'))
  },
  {
    name: "dashboard",
    path: "/",
    meta: {
      intern: true
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