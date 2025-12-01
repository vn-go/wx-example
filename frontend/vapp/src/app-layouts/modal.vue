<template>
    <div id="modal" class="fixed inset-0 flex items-center justify-center bg-black bg-opacity-30 opacity-0 pointer-events-none transition-opacity duration-300">
  <div class="bg-white rounded-lg shadow-lg w-full max-w-lg mx-4 transform scale-90 transition-transform duration-300">
    
    <!-- Header -->
    <div class="flex justify-between items-center border-b px-4 py-2">
      <h3 class="text-lg font-semibold">Modal Title</h3>
      <a href="#" class="text-gray-500 hover:text-gray-700 font-bold text-xl">&times;</a>
    </div>

    <!-- Body -->
    <div class="p-4">
        <slot name="content">
            <p>This is the content of the modal. You can put forms, messages, etc.</p>
        </slot>
     
    </div>

    <!-- Footer -->
    <div class="flex justify-end gap-2 border-t px-4 py-2">
      <a href="#" class="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300 transition">Close</a>
      <a href="#" class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 transition">Save</a>
    </div>
  </div>
</div>
</template>