<template>
    <header class="bg-white border-b border-gray-200 flex items-center px-1 shadow-sm z-30">
      <div class="flex items-center space-x-3">
        <!-- Burger icon -->
        <slot name="buger"></slot>

        <span class="text-sm font-medium hidden sm:inline text-gray-700">My App</span>
      </div>
    </header>
</template>