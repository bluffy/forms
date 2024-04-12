<template>
<div>
      <div v-if="loaded"></div>

      <AlertDialog ref="dialog"></AlertDialog>
</div>
</template>
  
<script lang="ts" setup>

import { onMounted , ref} from 'vue'
import ApiService from '../services/api.service'
import AlertDialog from "../components/AlertDialog.vue";
import { genResponseError } from "../utils/errorMessage";
import { useRoute } from 'vue-router'
import router from '../router'
import { PageMessage } from '../models/page.model';

const route = useRoute();


const dialog = ref()
const loaded = ref(false)


function responseError(err: any) {
  const appError = genResponseError(err)
  dialog.value?.alert(appError?.message);
}

onMounted(() => {
     
    route.params.link
    loaded.value = false;
  
    if (route.name == "register_link") {
      return ApiService.postPage("/user/register/link", {"decoded": route.params.link}).then(
        (page: PageMessage) => {
            router.push("/login")
            if (page.data.message) {
              dialog.value?.alert(page.data.message);
            }
            return
        },
        (err: any) => { 
            router.replace("/login")
            responseError(err);
          }
      );
    }
    router.push("/notfound")
    return


});
</script>
