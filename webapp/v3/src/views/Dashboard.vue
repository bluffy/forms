<template>
<div>
      <p>Dashboard</p>
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


function responseError(err: any) {
  const appError = genResponseError(err)
  dialog.value?.alert(appError?.message);
}

onMounted(() => {

    return ApiService.get("/").then(
        (page: PageIndex) => {
            console.log(page)

        },
        (err: any) => { 
            responseError(err);
          }
      );

});
</script>
