<template>
    <div class="h-full  flex items-center justify-center">
      <div class="bg-white rounded-2xl shadow-2xl p-10 w-full max-w-md widget-container border-gray-600">
        <!-- Header -->
        <h2 class="text-3xl font-bold text-center text-gray-800 mb-6">Sign In</h2>
        
        <!-- Form -->
        <div class="space-y-4">
          <div>
            <label class="block text-gray-700 mb-1" for="username">Username</label>
            <input
              type="text"
              id="username"
              placeholder="Your username"
              v-model="login.data.username"
              class="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-indigo-400 focus:outline-none"
            />
          </div>
          <div>
            <label class="block text-gray-700 mb-1" for="password">Password</label>
            <input
              type="password"
              id="password"
              placeholder="********"
              v-model="login.data.password"
              class="w-full px-4 py-2 border rounded-lg focus:ring-2 focus:ring-indigo-400 focus:outline-none"
            />
          </div>
          <button
            @click="login.doLogin"
            class="w-full py-2 bg-indigo-600 text-white font-semibold rounded-lg shadow hover:bg-indigo-700 transition"
          >
            Sign In
          </button>
        </div>
        <p class="mt-6 text-center text-gray-500 text-sm">
           {{ login.errMsg }}
        </p>
        <!-- Footer -->
        <p class="mt-6 text-center text-gray-500 text-sm">
          Forgot your password? <a href="#" class="text-indigo-600 hover:underline">Reset</a>
        </p>
      </div>
    </div>
  </template>
  
  <script setup lang="ts">
    import emitter from '@core/eventBus';
import libs from '@core/lib';

   class Login extends libs.BaseUI{
      data ={
        username:"",
        password:""
      }
      errMsg=libs.newRef();
      private _onAfterLogin=undefined;
      async doLogin() {
        let res= await libs.apiPublic.formPost(this.getViewPath(), 'auth/login',{
          grant_type: 'password',
          username:this.data.username,
          password:this.data.password
        });
       
        if(res.error){
          this.errMsg=res.error.statusText;
        }else {
          libs.sessionStore.set("tk",res.data.access_token);
          emitter.emit("after-login",{})
        }
        
         
      }
     async raiseAfterLogin() {
       if(this._onAfterLogin){
          await this._onAfterLogin();
       }
     }
     onAfterLogin(fn:()=>{}){
        this._onAfterLogin=fn;
     }
    
   }
   const instance= new Login("auth/login");

   const login=libs.newReactive(instance) ;
   function getIns() {
    return instance;
   }
   defineExpose({
    instance
   })
  </script>
  