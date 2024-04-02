
<template>

      <router-view /> 
  


</template>


<script setup lang="ts">
import router from "./router";
import { useAuthStore } from "./stores/auth";

const authStore = useAuthStore()

// @ts-ignore
router.beforeEach((to, from, next) => {
  if (from) {

  }

  if (to.matched.some(record => record.meta.requiresAuth)) {
    if(authStore.loggedIn){
      next()
      return
    }
    next('/login')
    return
  }

  next()
});


//import router from "./router";
//import { useAuthStore } from "./stores/auth";
//const authStore = useAuthStore()

/*

router.beforeEach((to, from, next) => {
  console.log(to)
  console.log(authStore)
  if (from) {

  }
/*
  //console.log(from)
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if(authStore.loggedIn){
      next()
      return
    }
    next('/login')
    return
  }

});

  next()
  */
</script>
