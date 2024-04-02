<template>

<div class="modal" v-if="dialog">
  <div class="modal-dialog">
    <div class="modal-content">
      <div class="modal-header">
        <h5 class="modal-title" v-if="options.title">{{ options.title }}</h5>
        <!-- <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button> -->
      </div>
      <div class="modal-body" v-html="message" />
      <div class="modal-footer">
        <button type="button" class="btn btn-primary" @click="agree">{{ options.ok }}</button>
        <button type="button" class="btn btn-secondary"  @click="disagree" v-if="isConfirm"> {{ options.cancle }}</button>
      </div>
    </div>
  </div>
</div>


</template>

<script setup lang="ts">
import { ref } from 'vue'


const dialog = ref(false)
const message = ref("")
const options = ref();
const isConfirm = ref(false)

var resolve = (v: boolean) => {console.log(v)}

function disagree() {
  dialog.value = false;
  resolve(false)
}
function agree() {
  resolve(true)
  dialog.value = false;
}

function alert(pMessage: string, params: any) {
  isConfirm.value = false
  dialog.value = true
  message.value = ""
  options.value = {
    title: "",
    ok: "OK",
}

  message.value = pMessage
  options.value = Object.assign(options.value,params)
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

defineExpose({
  alert,
  confirm,
});

</script>