<template>

<nav class="navbar navbar-expand-lg bg-light">
  <div class="container">
    <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarTogglerDemo03" aria-controls="navbarTogglerDemo03" aria-expanded="false" aria-label="Toggle navigation">
      <span class="navbar-toggler-icon"></span>
    </button>
    <a class="navbar-brand" href="#">Navbar</a>
    <div class="collapse navbar-collapse" id="navbarTogglerDemo03">
      <ul class="navbar-nav ms-auto">
        <li class="nav-item" v-if="route.meta.intern">
          <RouterLink class="nav-link" to="/"><i class="bi bi-house"></i> Home</RouterLink>
        </li>
        <li class="nav-item" v-if="route.meta.intern">
          <a class="nav-link" href="#"><i class="bi bi-gear"></i> Settings</a>
        </li>
        <li class="nav-item" v-if="route.meta.intern">
            <a class="nav-link" href="#" @click="onLogout"><i class="bi bi-box-arrow-left"></i> Logout</a>
        </li>
        <li class="nav-item" v-if="!route.meta.intern && route.name != 'register'">
            <RouterLink class="nav-link" to="/user/register"><i class="bi bi-person-fill-add"></i> Register</RouterLink>
        </li>
        <li class="nav-item" v-if="!route.meta.intern && route.name != 'login'">
            <RouterLink class="nav-link" to="/login"><i class="bi bi-box-arrow-left"></i> Login</RouterLink>
        </li>        
      </ul>


    </div>
  </div>
</nav>



        <AlertDialog ref="dialog"></AlertDialog>
  
</template>
<script lang="ts" setup name="PageHeader">
import {  ref} from 'vue'

import ApiService from '../services/api.service'
import AlertDialog from "../components/AlertDialog.vue";
import { PageNoContent } from '../models/page.model';
import router from '../router';
import { genResponseError } from "../utils/errorMessage";
import { useRoute } from 'vue-router'

const route = useRoute();

const dialog = ref()


function onLogout(){
    console.log("clock")
      return ApiService.get("/user/logout").then(
        (page: PageNoContent) => {
            if (page.status == 204 ) {
                router.push("/login")
            }
       
            return
        },
        (err: any) => { 
            dialog.value?.alert( genResponseError(err)?.message);
          }
      );
 
} 
</script>
<style>
    .header {
        min-height: 120px;;
    }
</style>