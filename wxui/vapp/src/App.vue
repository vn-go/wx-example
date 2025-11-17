<script setup lang="ts">
import AppLayout from '@app-layouts/app.vue';
import BurgerButton from '@widgets/app-header-buger.vue';
import AppSiderbar from '@widgets/app-siderbar.vue';
</script>

<template>
  <AppLayout>
    <template #burger> <BurgerButton/></template>
    <template #sidebar><AppSiderbar/></template>
    <template #body>
     

    </template>
  </AppLayout>
</template>


