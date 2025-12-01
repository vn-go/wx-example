<template>
    <div class="fixed top-0 left-0 w-full h-[6px] z-50 overflow-hidden">
      <div
        class="h-full w-full bg-gradient-to-r from-gray-400 via-blue-500 to-blue-600 animate-gradient"
      ></div>
    </div>
  </template>
  
  <script setup lang="ts">
  // Không cần props hay state, chỉ layout
  </script>
  
  <style>
  @keyframes gradientSlide {
    0% { background-position: 0% 0; }
    100% { background-position: 200% 0; }
  }
  
  .animate-gradient {
    background-size: 200% 100%;
    animation: gradientSlide 1s linear infinite;
  }
  </style>
  