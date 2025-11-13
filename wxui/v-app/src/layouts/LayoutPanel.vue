<!-- LayoutPanel.vue -->
<script setup>
import Buger from '@components/buger.vue';
import { defineEmits, nextTick, onMounted, ref } from 'vue';
import emitter from './../event-bus.js';

const sidebarRef = ref(null);
const bodyRef = ref(null);
const headerRef =ref(null);
const onMousemove=(event)=>{ 
  
  isOpen.value =true;
}
const onMouseOut=(event) =>{
  isOpen.value =false;
  
}
const isOpen = ref(false) // sidebar status
const emit = defineEmits(["buger-click"])
const toggleSidber=()=>{
    isOpen.value = !isOpen.value
    if(isOpen.value) {
          if (isOpen.value) {
              sidebarRef.value.addEventListener('mouseout',onMouseOut)
              sidebarRef.value.addEventListener('mousemove',onMousemove)
            
            } else {
              sidebarRef.value.removeEventListener ('mouseout',onMouseOut)
              sidebarRef.value.removeEventListener('mousemove',onMousemove)
            }
    }
}
const updateBodyHeight = () => {
  if (!headerRef.value || !bodyRef.value) return;

  const headerHeight = headerRef.value.offsetHeight; // chiều cao header
  const viewportHeight = window.innerHeight; // chiều cao trình duyệt

  bodyRef.value.style.height = `${viewportHeight - headerHeight}px`;
};
onMounted(()=>{
  nextTick(() => {
    updateBodyHeight();
  });

  
  window.addEventListener('resize', updateBodyHeight);
  emitter.on('route-change', async () => {
    // đợi Vue render slot mới xong
    await nextTick();
    slideFromRight(500);
  });
})

const slideFromRight = (duration = 500) => {
  const el = bodyRef.value;
  if (!el) return;

  // reset position
  el.style.transition = 'none';
  el.style.transform = 'translateX(100%)';

  // cần 2 frame để browser nhận transform ban đầu trước khi trigger animation
  requestAnimationFrame(() => {
    requestAnimationFrame(() => {
      el.style.transition = `transform ${duration}ms ease-out`;
      el.style.transform = 'translateX(0)';
    });
  });
};
</script>

<style scoped>
/* sidebar animation*/ 
.sidebar { 
  position: relative;
  width: 250px;
  /* max-width: 80%; */
  /* height: 100vh; */
  /* background-color: #f0f0f0; */
  overflow: hidden;
  transform: translateX(-100%);
  transition: transform 0.3s ease-in-out; /* animation */
  box-shadow: 2px 0 5px rgba(0,0,0,0.2);
}
.sidebar.open {
  transform: translateX(0);
}

</style>

<template>
    <div class="main" >
   
      <header @buger-click="toggleSidber" ref="headerRef">
  
    <div class="flex items-center space-x-3">
      <Buger @buger-click="toggleSidber"/>
      <span class="text-sm font-medium hidden sm:inline"></span>
    </div>

    
      </header>
    
  
      <!-- Main content -->
      <div class="flex flex-1">
        <!-- Left Sidebar -->
        <aside 
          @mouseove="onMOuseMove"
          @mouseout="onMouseOut"
          ref="sidebarRef"
          style="position: fixed;top:40px;height:calc(100vh - 40px);" 
          :class="{ open: isOpen }" 
          class="sidebar h-[200px] bg-blue-200 z-20"
        >
          <slot name="sidebar">Default Sidebar</slot>
        </aside>
        <div class="h-full w-full flex-1" ref="bodyRef">
        <slot name="body">Default Body Content</slot>
        </div>
       
      </div>
    </div>
  </template>
  
