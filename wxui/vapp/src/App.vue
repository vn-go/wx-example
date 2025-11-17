<script setup lang="ts">
import AppLayout from '@app-layouts/app.vue';
import libs from '@core/lib';
import BurgerButton from '@widgets/app-header-buger.vue';
import AppSiderbar from '@widgets/app-siderbar.vue';

import { shallowRef } from 'vue';
const currentComponent = shallowRef<any>(null);

class Application extends libs.BaseUI {
  currentComponent = libs.newRef<any>(null);
  onInit(): void {
      let self=this;
      libs.urlNav.onNav(async (pathName,search)=>{
        await self.loadView(pathName);
      });
      libs.urlNav.init();
  }
  async loadView(pathName?:string){
    if(!pathName) return;
    currentComponent.value= await libs.loadViews(pathName,'error')
  }
  
}
const application=new Application();
</script>

<template>
  <AppLayout>
    <template #burger> <BurgerButton/></template>
    <template #sidebar><AppSiderbar/></template>
    <template #body>
    <Transition
      mode="out-in"
      enter-active-class="transition transform duration-300 ease-out"
      enter-from-class="translate-x-full opacity-0"
      enter-to-class="translate-x-0 opacity-100"
      leave-active-class="transition transform duration-300 ease-in"
      leave-from-class="translate-x-0 opacity-100"
      leave-to-class="translate-x-full opacity-0"
    >
      <component :is="currentComponent" :key="currentComponent?.__asyncLoader?.toString() || 'component'" class=""/>
    </Transition>
  </template>
  </AppLayout>
</template>


