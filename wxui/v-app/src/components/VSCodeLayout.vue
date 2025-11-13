<script setup>
import { onMounted, onUnmounted, ref } from 'vue'

// Sidebar collapsed state (mobile-first)
const sidebarOpen = ref(true)

const toggleSidebar = () => {
  sidebarOpen.value = !sidebarOpen.value
}

// Auto-close on very small screens
const handleResize = () => {
  if (window.innerWidth < 640) sidebarOpen.value = false
  else sidebarOpen.value = true
}
onMounted(() => {
  handleResize()
  window.addEventListener('resize', handleResize)
})
onUnmounted(() => window.removeEventListener('resize', handleResize))
</script>

<template>
  <div class="h-screen flex flex-col bg-[#1e1e1e] text-[#d4d4d4] overflow-hidden">
    <!-- ==== HEADER ==== -->
    <header class="bg-[#252526] border-b border-[#3c3c3c] flex items-center px-2 h-9 shrink-0">
      <slot name="header">
        <!-- Default VS Code style header -->
        <div class="flex items-center space-x-1 flex-1">
          <button
            @click="toggleSidebar"
            class="p-1.5 hover:bg-[#3c3c3c] rounded"
            aria-label="Toggle sidebar"
          >
            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 16 16">
              <path d="M2 2h12v12H2zM5 3v10" stroke="currentColor" stroke-width="1.5" />
            </svg>
          </button>

          <div class="flex space-x-1 ml-2">
            <span class="px-2 py-0.5 text-xs bg-[#007acc] text-white rounded">
              App.vue
            </span>
            <span class="px-2 py-0.5 text-xs bg-[#3c3c3c] rounded">index.html</span>
          </div>
        </div>

        <div class="flex space-x-2">
          <button class="p-1 hover:bg-[#3c3c3c] rounded">
            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 16 16">
              <path d="M8 2v12M2 8h12" stroke="currentColor" stroke-width="1.5"/>
            </svg>
          </button>
        </div>
      </slot>
    </header>

    <!-- ==== BODY (SIDEBAR + CONTENT) ==== -->
    <div class="flex flex-1 overflow-hidden">
      <!-- ==== SIDEBAR ==== -->
      <aside
        :class="[
          'bg-[#252526] border-r border-[#3c3c3c] transition-all duration-200 overflow-y-auto',
          sidebarOpen ? 'w-64' : 'w-0 sm:w-12'
        ]"
      >
        <div
          :class="[
            'h-full flex flex-col',
            sidebarOpen ? '' : 'sm:flex'
          ]"
        >
          <slot name="sidebar">
            <!-- Default Explorer tree -->
            <div class="p-2 text-xs space-y-1">
              <div class="flex items-center space-x-1 text-[#cccccc] font-medium">
                <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 16 16">
                  <path d="M2 3h3l1 1h7v9H2z" stroke="currentColor" stroke-width="1"/>
                </svg>
                <span v-if="sidebarOpen">EXPLORER</span>
              </div>

              <template v-if="sidebarOpen">
                <div class="ml-4 space-y-0.5">
                  <div class="flex items-center space-x-1 hover:bg-[#3c3c3c] px-1 py-0.5 rounded cursor-pointer">
                    <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 16 16">
                      <path d="M3 3h10v9H3z" stroke="currentColor" stroke-width="1"/>
                    </svg>
                    <span>src</span>
                  </div>
                  <div class="ml-4 flex items-center space-x-1 hover:bg-[#3c3c3c] px-1 py-0.5 rounded cursor-pointer">
                    <svg class="w-3 h-3 text-[#4db8ff]" fill="currentColor" viewBox="0 0 16 16">
                      <path d="M2 2h12v12H2z"/>
                    </svg>
                    <span>App.vue</span>
                  </div>
                </div>
              </template>
            </div>
          </slot>
        </div>
      </aside>

      <!-- ==== MAIN CONTENT ==== -->
      <main class="flex-1 overflow-y-auto bg-[#1e1e1e] p-4">
        <slot>
          <!-- Default editor placeholder -->
          <div class="bg-[#252526] rounded p-4 text-sm">
            <p class="text-[#6c6c6c]">/* Your editor content goes here */</p>
            <pre class="mt-2 text-[#d4d4d4]">
<span class="text-[#569cd6]">const</span> <span class="text-[#9cdcfe]">message</span> = <span class="text-[#ce9178]">"Hello VS Code layout!"</span>;
console.<span class="text-[#dcdcaa]">log</span>(message);</pre>
          </div>
        </slot>
      </main>
    </div>
  </div>
</template>

<style scoped>
/* Tiny tweak for smooth scrollbar on sidebar */
aside::-webkit-scrollbar {
  width: 6px;
}
aside::-webkit-scrollbar-thumb {
  background: #444;
  border-radius: 3px;
}
</style>