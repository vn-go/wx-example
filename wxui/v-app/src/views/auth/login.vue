<template>
    <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div class="max-w-md w-full">
        <!-- Simple Card -->
        <div class="bg-white shadow rounded-2xl p-8">
          <div class="mb-6 text-center">
            <h2 class="text-2xl font-semibold text-gray-900">Đăng nhập</h2>
            <p class="text-sm text-gray-500">Vui lòng nhập username và mật khẩu</p>
          </div>
  
          <form @submit.prevent="onSubmit" class="space-y-4" novalidate>
            <!-- Username -->
            <div>
              <label for="username" class="block text-sm font-medium text-gray-700">Username</label>
              <input
                id="username"
                v-model="form.username"
                type="text"
                autocomplete="username"
                required
                class="mt-1 block w-full px-4 py-2 rounded-lg border border-gray-200 shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:border-transparent"
                :class="{'ring-rose-400 border-rose-200': errors.username}"/>
              
              <p v-if="errors.username" class="mt-1 text-xs text-rose-600">{{ errors.username }}</p>
            </div>
  
            <!-- Password -->
            <div>
              <label for="password" class="block text-sm font-medium text-gray-700">Mật khẩu</label>
              <div class="relative mt-1">
                <input
                  id="password"
                  v-model="form.password"
                  :type="showPassword ? 'text' : 'password'"
                  autocomplete="current-password"
                  required
                  class="block w-full px-4 py-2 rounded-lg border border-gray-200 shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-400 focus:border-transparent"
                  :class="{'ring-rose-400 border-rose-200': errors.password}"
                />
  
                <button t
                ype="button" class="absolute right-2 top-1/2 -translate-y-1/2 p-1.5 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-300" 
                @click="togglePassword" :aria-pressed="showPassword">
                  <span v-if="showPassword" class="sr-only">Ẩn mật khẩu</span>
                  <span v-else class="sr-only">Hiện mật khẩu</span>
                  <svg v-if="!showPassword" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                  <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M3 3l18 18" />
                    <path stroke-linecap="round" stroke-linejoin="round" d="M10.477 10.477A3 3 0 0113.523 13.523" />
                  </svg>
                </button>
              </div>
              <p v-if="errors.password" class="mt-1 text-xs text-rose-600">{{ errors.password }}</p>
            </div>
  
            <div>
              <button
                type="submit"
                :disabled="loading"
                class="w-full flex justify-center items-center gap-2 px-4 py-2 rounded-lg bg-gradient-to-r from-indigo-600 to-violet-600 text-white font-semibold shadow-md hover:opacity-95 focus:outline-none focus:ring-2 focus:ring-indigo-300 disabled:opacity-60"
              >
                <svg v-if="loading" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 animate-spin" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"></path>
                </svg>
                <span>{{ loading ? 'Đang xử lý...' : 'Đăng nhập' }}</span>
              </button>
            </div>
          </form>
  
          <p class="mt-4 text-xs text-gray-500 text-center">Quên mật khẩu? <button @click="$emit('forgot')" class="text-indigo-600 hover:underline">Nhấn vào đây</button></p>
        </div>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
  import { apiCaller } from '@core/apiCaller.js';
import emitter from '@core/eventBus.ts';
import { reactive, ref } from 'vue';
  const emit = defineEmits(['submit', 'forgot'])
  
  const form = reactive({
    username: '',
    password: '',
  })
  
  const errors = reactive<{ username?: string; password?: string }>({})
  const loading = ref(false)
  const showPassword = ref(false)
  
  function validate() {
    errors.username = undefined
    errors.password = undefined
  
    if (!form.username) errors.username = 'Username không được để trống.'
    if (!form.password) errors.password = 'Mật khẩu không được để trống.'
    else if (form.password.length < 4) errors.password = 'Mật khẩu phải ít nhất 4 ký tự.'
  
    return !errors.username && !errors.password
  }
  
  function togglePassword() {
    showPassword.value = !showPassword.value
  }
  
  async function onSubmit() {
    if (!validate()) return
    loading.value = true
    try {
      let loginInfo=await apiCaller.login(form.username,form.password)
      console.log(loginInfo);
      emitter.emit('login-success', {
        user: loginInfo.user,
        token: loginInfo.token,
        timestamp: new Date().toISOString()
      });

      //await emit('submit', { username: form.username, password: form.password })
    } finally {
      loading.value = false
    }
  }

  </script>
  
  <style scoped></style>