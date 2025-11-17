<template>
    <nav class="flex flex-col space-y-2 p-4" style="min-width: 250px;">
      <div v-for="group in sideBar.data" :key="group.title">
        <!-- group title -->
        <button
          @click="sideBar.toggleGroup(group.title)"
          class="w-full flex justify-between items-center px-3 py-2 text-gray-700 font-medium rounded hover:bg-gray-100 transition"
        >
          <span>{{ group.title }}</span>
          <svg
            class="w-4 h-4 transition-transform duration-300"
            :class="{'rotate-90': sideBar.openGroups.value.has(group.title)}"
            fill="none" stroke="currentColor" viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M9 5l7 7-7 7" />
          </svg>
        </button>
  
        <!-- children -->
        <div
          v-show="sideBar.openGroups.value.has(group.title)"
          class="flex flex-col pl-5 mt-1 space-y-1"
        >
          <a
            v-for="child in group.children"
            :key="child.id"
            href="javascript:void(0);"
            @click="sideBar.doLoadView(child.pathname)"
            class="px-2 py-1 rounded hover:bg-gray-200 transition"
          >
            {{ child.title }}
          </a>
        </div>
      </div>
    </nav>
  </template>
<script setup lang="ts">
import libs from '@core/lib';
class Sidebar extends libs.BaseUI {
    
    data=libs.getAppMenuData();
    openGroups = libs.newRef<Set<string>>(new Set());
    doLoadView(viewPath?:string) {
        libs.urlNav.move(viewPath);
    }
    toggleGroup(title: string) {
        if (this.openGroups.value.has(title)) {
            this.openGroups.value.delete(title);
        } else {
            this.openGroups.value.add(title);
        }
    }
}
const sideBar =new Sidebar();
console.log(libs.getViewMap());
</script>