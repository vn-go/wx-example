<template>
    <div class="flex flex-col items-center justify-center h-full bg-red-50 text-red-800">
      <h1 class="text-4xl font-bold mb-4">404</h1>
      <p class="text-lg mb-2">Oops! The view you are looking for was not found.</p>
      <p class="text-sm">Please check the URL or go back to the home page.</p>
      <button
        class="mt-4 px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
        @click="$router.push('/')"
      >
        Go Home
      </button>
    </div>
  </template>
  
  <script setup>
  // Không cần logic gì thêm
  </script>
  
  <style scoped>
  /* Tailwind đã handle style */
  </style>