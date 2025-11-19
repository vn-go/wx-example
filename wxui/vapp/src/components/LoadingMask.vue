<template>
  <transition name="fade">
    <div 
      v-if="visible" 
      class="fixed inset-0 z-50 pointer-events-auto"
    >
      <TopLoadingBar/>
    </div>
  </transition>
</template>

<script setup lang="ts">
import TopLoadingBar from './TopLoadingBar.vue';
defineProps<{ visible: boolean }>();
</script>

<style>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.2s;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>
