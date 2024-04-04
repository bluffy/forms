<template>

<div  class="modal fade"  tabindex="-1" aria-labelledby=""
    aria-hidden="true"  ref="modalEle" >
  <div class="modal-dialog">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" v-if="options && options.title">{{ options.title }}</h5>
        <!-- <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button> -->
      </div>
      <div class="modal-body" v-html="message" />
      <div class="modal-footer">
        <button v-if="options" type="button" class="btn btn-primary" @click="agree">{{ options.ok }}</button>
        <button v-if="options && isConfirm" type="button" class="btn btn-secondary"  @click="disagree"> {{ options.cancle }}</button>
      </div>
    </div>
  </div>
</div>


</template>

<script setup lang="ts" name="AlertDialog">
import { ref, onMounted} from 'vue'
import { Modal } from "bootstrap";

//import { Modal } from "bootstrap";

let thisModalObj = null;
let modalEle = ref(null);


const dialog = ref(false)
const message = ref("")
const options = ref()
const isConfirm = ref(false)

var resolve = (v: boolean) => {console.log(v)}

function disagree() {
  dialog.value = false;
  thisModalObj.hide();
  resolve(false)
}
function agree() {
  resolve(true)
  dialog.value = false;
  thisModalObj.hide();
}

function alert(pMessage: string, params: any) {
  console.log("alert")



  isConfirm.value = false
  dialog.value = true
  message.value = ""
  options.value = {
    title: "",
    ok: "OK",
}

  message.value = pMessage
  options.value = Object.assign(options.value,params)

  thisModalObj.show();
  return new Promise((res) => {
    resolve = res
  })

}
function confirm(pMessage: string, params: any) {
  isConfirm.value = true
  dialog.value = true
  message.value = ""
  options.value = {
    title: "",
    ok: "OK",
    cancle: "Abbrechen",
}

  message.value = pMessage
  options.value = Object.assign(options.value,params)
  return new Promise((res) => {
    resolve = res
  })
}

onMounted(() => {
  thisModalObj = new Modal(modalEle.value, {backdrop: "static"});
  console.log(thisModalObj)
});


defineExpose({
  alert,
  confirm,
});



</script>