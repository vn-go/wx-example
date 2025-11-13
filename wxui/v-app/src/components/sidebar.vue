<script setup>
import { defineEmits, onMounted, ref } from 'vue'

const menuData = [
  {
    title: "System",
    children: [
      { 
        title: "Users",
        link:"/system/users"
     },
      { 
        title: "Roles",
        link:"/system/roles"
     },
      { title: "Permissions" }
    ]
  },
  {
    title: "Settings",
    children: [
      { title: "Language" },
      { title: "Theme" }
    ]
  }
]
const emmit = defineEmits(["on-page-change"])
// danh sách mở/đóng (dựa trên index)
const openIndexes = ref([])

function toggleMenu(i) {
  const idx = openIndexes.value.indexOf(i)
  if (idx >= 0) openIndexes.value.splice(idx, 1)
  else openIndexes.value.push(i)
}

function isOpen(i) {
  return openIndexes.value.includes(i)
}
function doLoadPage(data) {
    window.history.pushState({}, "", data.link)
    emmit("on-page-change",data.link);


}
onMounted(()=>{
    window.addEventListener("popstate", (event) => {
        emmit("on-page-change",window.location.pathname);

    })
    emmit("on-page-change",window.location.pathname);
})

</script>

<template>
    <ul>
      <li v-for="(menu, i) in menuData" :key="i" class="mb-1">
        <!-- Menu cấp 1 -->
        <button
          @click="toggleMenu(i)"
          class="w-full flex justify-between items-center bg-gray-100 hover:bg-gray-200 px-3 py-2 rounded-md font-medium text-gray-800 transition"
        >
          <span>{{ menu.title }}</span>
          <svg
            class="w-4 h-4 transform transition-transform duration-200"
            :class="{ 'rotate-90': isOpen(i) }"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
          </svg>
        </button>

        <!-- Menu cấp 2 -->
        <ul
          v-if="isOpen(i)"
          class="ml-4 mt-1 border-l border-gray-200 pl-2 space-y-1 transition-all"
        >
        <li v-for="child in menu.children" :key="child.link">
            <a @click="doLoadPage(child)" href="javascript:void(0);">{{ child.title }}</a>
            </li>
        </ul>
      </li>
    </ul>
</template>
