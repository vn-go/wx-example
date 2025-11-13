<script setup>
import { defineEmits } from 'vue'
import buger from '@components/buger.vue'
import { ref } from 'vue'
// Props
defineProps({
  title: { type: String, default: 'My App' }
})

// Emits
const emit = defineEmits(['toggle-sidebar'])

// Dropdown state
const dropdownOpen = ref(false)

const toggleDropdown = () => {
  dropdownOpen.value = !dropdownOpen.value
}
function emitToggleSidebar(){
  debugger;
  emit('toggle-sidebar')
}
// Close dropdown when clicking outside
const handleClickOutside = (e) => {
  if (!e.target.closest('.dropdown')) dropdownOpen.value = false
}
</script>

<template>
  <header
    
    @click="handleClickOutside"
  >
    <!-- LEFT: Burger + Title -->
    <div class="flex items-center space-x-3">
      <buger/>
      <span class="text-sm font-medium hidden sm:inline"></span>
    </div>

    <!-- RIGHT: Dropdown -->
    <div class="relative dropdown">
      <button
        @click.stop="toggleDropdown"
        class="flex items-center space-x-1 p-1.5 hover:bg-[#3c3c3c] rounded transition-colors"
        aria-label="User menu"
      >
        <div class="w-6 h-6 rounded-full bg-gradient-to-br from-indigo-500 to-purple-600"></div>
        <svg
          class="w-4 h-4 transition-transform"
          :class="{ 'rotate-180': dropdownOpen }"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
        </svg>
      </button>

      <!-- Dropdown Menu -->
      <transition
        enter-active-class="transition ease-out duration-100"
        enter-from-class="transform opacity-0 scale-95"
        enter-to-class="transform opacity-100 scale-100"
        leave-active-class="transition ease-in duration-75"
        leave-from-class="transform opacity-100 scale-100"
        leave-to-class="transform opacity-0 scale-95"
      >
        <div
          v-if="dropdownOpen"
          class="absolute right-0 mt-2 w-48 bg-[#252526] rounded-md shadow-lg border border-[#3c3c3c] overflow-hidden z-50"
        >
          <div class="py-1">
            <a
              href="#"
              class="block px-4 py-2 text-sm hover:bg-[#3c3c3c] transition-colors"
            >
              Profile
            </a>
            <a
              href="#"
              class="block px-4 py-2 text-sm hover:bg-[#3c3c3c] transition-colors"
            >
              Settings
            </a>
            <hr class="border-[#3c3c3c] my-1" />
            <a
              href="#"
              class="block px-4 py-2 text-sm text-red-400 hover:bg-[#3c3c3c] transition-colors"
            >
              Sign out
            </a>
          </div>
        </div>
      </transition>
    </div>
  </header>
</template>