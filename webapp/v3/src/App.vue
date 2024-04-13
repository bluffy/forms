
<template>

<div class="d-flex flex-column min-vh-100" >
  <PageHeader/>
  <router-view class="content"/>
  <PageFooter class="mt-auto"/>
</div>    
</template>


<script setup lang="ts">
import PageHeader from "./components/PageHeader.vue";
import PageFooter from "./components/PageFooter.vue";
</script>
<style>
    .content {
        min-height: 200px;
    }
</style>