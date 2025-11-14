<!-- App.vue -->
<script setup>
import AppHeader from '@components/AppHeader.vue';
import sidebar from '@components/sidebar.vue';
import emitter from '@core/eventBus.js';
import LayoutPanel from '@layouts/LayoutPanel.vue';
import { defineAsyncComponent, nextTick, ref } from 'vue';

import { apiCaller } from '@core/apiCaller';
import { loadComponent } from '@core/dynamic-import-helper.js';
var test="abc";

function handleToggleSidebar() {
  debugger;
  sidebarOpen.value = !sidebarOpen.value
  test="bvbvbbv";
}
const currentComponent = ref(null)
const doLogin=()=> {

}
let currentPagrPath="";

const loadPage=async (pagePath)=>{
  currentPagrPath=pagePath;
  debugger;
  if (apiCaller.getToken()==''){
    if(pagePath.split('?')[0]!='/auth/login'){
      window.history.pushState({}, "", `/auth/login?ret=${encodeURI(pagePath)}`);
    }
  
    await nextTick() 
    
    currentComponent.value = await defineAsyncComponent(() => loadComponent('/auth/login'));
    
    return
  }
  const parts = pagePath.split('/').filter(Boolean) // ["system","users"]
  await nextTick() 
  currentComponent.value =  defineAsyncComponent(() => loadComponent(pagePath))
  emitter.emit('route-change', { speed: 500 }) // that ra bat dau tu day
}
emitter.on('login-success',()=>{
  currentComponent.value =  defineAsyncComponent(() => loadComponent(currentPagrPath))
  emitter.emit('route-change', { speed: 500 }) // that ra bat dau tu day
  window.history.pushState({}, "", currentPagrPath);
})
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

