<template>
   
   <div class="w-full px-2 py-2 border-t border-gray-200 bg-gray-50 flex justify-end gap-1">

<button v-if="showNew"
        class="btn btn-new"
        :disabled="loading"
        @click="$emit('onNew')"
        >
    New
</button>

<button v-if="showUpdate"
        class="btn btn-update"
        :disabled="loading"
        @click="$emit('onUpdate')">
    Update
</button>

<button v-if="showDelete"
        class="btn btn-delete"
        :disabled="loading"
        @click="$emit('onDelete')">
    Delete
</button>

<button class="btn btn-close"
        :disabled="loading"
        @click="$emit('onClose')">
    Close
</button>

</div>

  </template>
  
  <script setup lang="ts">
  defineProps({
    loading: { type: Boolean, default: false },
  
    showNew: { type: Boolean, default: true },
    showUpdate: { type: Boolean, default: true },
    showDelete: { type: Boolean, default: true }
  });
  
  defineEmits(["onNew", "onUpdate", "onDelete", "onClose"]);
  </script>
  
  <style scoped>
  /* Bạn có thể tuỳ chỉnh thêm nếu cần */
  </style>
  