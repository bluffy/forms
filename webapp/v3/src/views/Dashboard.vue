<template>
  <div class="content container">
        <div v-if="loaded">Dashboard</div>
        <AlertDialog ref="dialog"></AlertDialog>
  </div>
</template>
  
<script lang="ts" setup>

import { onMounted , ref} from 'vue'
import ApiService from '../services/api.service'
import {PageIndex} from '../models/page.model'
import AlertDialog from "../components/AlertDialog.vue";
import { genResponseError } from "../utils/errorMessage";

const dialog = ref()
const loaded = ref(false)


function responseError(err: any) {
  const appError = genResponseError(err)
  dialog.value?.alert(appError?.message);
}

onMounted(() => {
     loaded.value = false;
    return ApiService.get("/").then(
        (page: PageIndex) => {
            console.log(page.data)
            loaded.value = true
        },
        (err: any) => { 
            responseError(err);
          }
      );

});
</script>
