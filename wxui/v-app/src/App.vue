<!-- App.vue -->
<script setup>
import AppHeader from '@components/AppHeader.vue';
import sidebar from '@components/sidebar.vue';
import LayoutPanel from '@layouts/LayoutPanel.vue';
import { defineAsyncComponent, ref } from 'vue';
import { loadComponent } from './dynamic-import-helper.js';
import emitter from './event-bus.js';
var test="abc";
function handleToggleSidebar() {
  debugger;
  sidebarOpen.value = !sidebarOpen.value
  test="bvbvbbv";
}
const currentComponent = ref(null)

const loadPage=(pagePath)=>{
  const parts = pagePath.split('/').filter(Boolean) // ["system","users"]
  currentComponent.value = defineAsyncComponent(() => loadComponent(pagePath))
  emitter.emit('route-change', { speed: 500 }) // that ra bat dau tu day
}
const layoutRef = ref(null)
</script>


<template>
  <LayoutPanel>
    <template #header>
      <AppHeader @toggle-sidebar="handleToggleSidebar"/>
    </template>

    <template #sidebar>
      <sidebar  @on-page-change="loadPage"/>
    </template>

    <template #body>
      <component :is="currentComponent" /> <!-- that ra no da render ngay cho nay roi-->
    </template>
  </LayoutPanel>
</template>

