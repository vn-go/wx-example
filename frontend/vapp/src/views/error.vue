<script setup lang="ts">
import libs from '@core/lib';



function goHome() {
  libs.urlNav.move('/')
}
</script>

<template>
  <div class="flex flex-col items-center justify-center h-full bg-gray-50 p-6">
    <h1 class="text-6xl font-bold text-red-500 mb-4">404</h1>
    <h2 class="text-2xl font-semibold mb-2">Page Not Found</h2>
    <p class="text-gray-600 mb-6 text-center">
      Sorry, the page you are looking for does not exist or has been moved.
    </p>
    <button
      @click="goHome"
      class="px-6 py-3 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition"
    >
      Go to Home
    </button>
  </div>
</template>
