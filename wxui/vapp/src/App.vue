<script setup lang="ts">
import AppLayout from '@app-layouts/app.vue';
import libs from '@core/lib';
import BurgerButton from '@widgets/app-header-buger.vue';
import AppSiderbar from '@widgets/app-siderbar.vue';
import { ref, shallowRef } from "vue";
const currentComponent = shallowRef<any>(null);
const currentView = ref();

class Application extends libs.BaseUI {
  currentComponent = libs.newRef<any>(null);
  $currentView:any;

  userInfo: any;
  constructor() {
    super();
    this.refKey("$currentView")
    // watch(this.$currentView, (val) => {
    //   if (val) {
    //     console.log("Mounted instance:", val);
        
    //   }
    // });
    libs.onAfterLogin(async ()=>{
          debugger;
          let userInfoResult=await this.getUserInfo();
          if(!userInfoResult.error){
            this.userInfo=libs.newReactive(userInfoResult.data)
          } else {
              libs.urlNav.move("/auth/login");
          }
      })
  }
  
  async onInit() {
      
      let userInfoResult=await this.getUserInfo();
      let isRequireLogin=false;
      if (userInfoResult.error && userInfoResult.error.status==401){
         isRequireLogin=true;

      } else {
          this.userInfo=libs.newReactive(userInfoResult.data)
      }
      
      let self=this;
      libs.urlNav.onNav(async (pathName,search)=>{
          if(!pathName) {
            libs.urlNav.move("home");
            return;
          }
          if(!libs.sessionStore.get("tk")||isRequireLogin){
            let currentPatname=libs.urlNav.getPathname();
            if (currentPatname!="auth/login"){
              libs.urlNav.changeUrl("/auth/login",libs.urlNav.makeQuery("ret",currentPatname));
              await self.loadView("/auth/login");
              self.$currentView= await self.getBindComponent("$currentView"); 
              if(self.$currentView && self.$currentView.value && self.$currentView.value.instance){
                self.$currentView.value.instance.onAfterLogin(async ()=>{
                   
                    userInfoResult=await self.getUserInfo();
                  
                    if (userInfoResult.error && userInfoResult.error.status==401){
                      isRequireLogin=true;

                    } else {
                      self.userInfo=libs.newReactive(userInfoResult.data);
                      isRequireLogin=false;
                    }
                    libs.urlNav.move(currentPatname);
                })
              }
            } else {
              await self.loadView("/auth/login");
              self.$currentView= await self.getBindComponent("$currentView"); 
              if(self.$currentView && self.$currentView.value && self.$currentView.value.instance){
                self.$currentView.value.instance.onAfterLogin(async ()=>{
                 
                  userInfoResult=await self.getUserInfo();
                  
                    if (userInfoResult.error && userInfoResult.error.status==401){
                      isRequireLogin=true;

                    } else {
                      self.userInfo=libs.newReactive(userInfoResult.data);
                      isRequireLogin=false;
                    }
                    libs.urlNav.move("home");
                })
              }
              
            }
            
          } else {
            if (pathName=="auth/login"){
                libs.urlNav.move("home");
            } else {
              if(!pathName) {
                pathName="home";
              }
              await self.loadView(pathName);
            }
            
          }
        });
      libs.urlNav.init();
      
  }
  async loadView(pathName?:string){
    if(!pathName) return;
    currentComponent.value= await libs.loadViews(pathName,'error')
   
  }
  async getUserInfo() {
      return await libs.api.post("/accounts/me")
  }
  
}
const application=new Application();
</script>

<template>
  <AppLayout>
    <template #burger> <BurgerButton/></template>
    <template #sidebar><AppSiderbar/></template>
    <template #body>
    <Transition
      mode="out-in"
      enter-active-class="transition transform duration-300 ease-out"
      enter-from-class="translate-x-full opacity-0"
      enter-to-class="translate-x-0 opacity-100"
      leave-active-class="transition transform duration-300 ease-in"
      leave-from-class="translate-x-0 opacity-100"
      leave-to-class="translate-x-full opacity-0"
    >
      <component 
        :is="currentComponent" 
        :key="currentComponent?.__asyncLoader?.toString() || 'component'" 
        class=""
        :ref="ele=>application.$currentView.value=ele"
      />
    </Transition>
  </template>
  </AppLayout>
</template>


