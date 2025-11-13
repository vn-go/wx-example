<script setup>
    import { defineEmits } from 'vue';
    const emit = defineEmits(["buger-click"])
    const onClick=()=>{
        debugger;
        emit("buger-click");
    }
</script>
<template>
    <button
        @click.stop="onClick"
        class="p-1.5 hover:bg-[#3c3c3c] rounded transition-colors"
        aria-label="Toggle sidebar"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M4 6h16M4 12h16M4 18h16"
          />
        </svg>
      </button>
</template>