<script setup lang="ts">
import AppLayout from '@app-layouts/app.vue';
import emitter from '@core/eventBus';
import libs from '@core/lib';
import BurgerButton from '@widgets/app-header-buger.vue';
import AppSiderbar from '@widgets/app-siderbar.vue';

import { ref, shallowRef } from "vue";
import LoadingMask from './components/LoadingMask.vue';


const currentComponent = shallowRef<any>(null);
const currentView = ref();

class Application extends libs.BaseUI {
  currentComponent = libs.newRef<any>(null);
  $currentView: any;
  isShowMask = libs.newRef(false);
  userInfo: any;
  constructor() {
    super();
    this.refKey("$currentView")
    emitter.on("on-api-dial", () => {
      this.isShowMask.value = true;
    });
    emitter.on("on-api-complete", () => {
      this.isShowMask.value = false;
    });

    libs.onAfterLogin(async () => {
      debugger;
      let userInfoResult = await this.getUserInfo();
      if (!userInfoResult.error) {
        this.userInfo = libs.newReactive(userInfoResult.data)
      } else {
        libs.urlNav.move("/auth/login");
      }
    })
  }

  async onInit() {
    let self = this;
    emitter.on('require-login', async () => {
      debugger;
      let currentPatname = libs.urlNav.getPathname();

      if (currentPatname != "auth/login") {
        libs.urlNav.changeUrl("/auth/login", libs.urlNav.makeQuery("ret", currentPatname));
        await self.loadView("/auth/login");
      }
    });
    emitter.on("after-login", async () => {
      debugger;
      let redirectTo = libs.urlNav.getQuery("ret") ?? "/";
      libs.urlNav.move(redirectTo);
    })
    libs.urlNav.onNav(async (pathName, search) => {
      debugger;
      if (!pathName) {
        await self.loadView("home");
      } else {
        await self.loadView(pathName);
      }

    });

    libs.urlNav.init();

  }
  async loadView(pathName?: string) {
    if (!pathName) return;
    currentComponent.value = await libs.loadViews(pathName, 'error')

  }
  async getUserInfo() {
    return await libs.api.post("app", "/accounts/me")
  }

}
const application = new Application();
</script>

<template>
  <LoadingMask :visible="application.isShowMask.value"></LoadingMask>
  <AppLayout>
    <template #burger>
      <BurgerButton />
    </template>
    <template #sidebar>
      <AppSiderbar />
    </template>
    <template #body>
      <Transition mode="out-in" enter-active-class="transition transform duration-300 ease-out"
        enter-from-class="translate-x-full opacity-0" enter-to-class="translate-x-0 opacity-100"
        leave-active-class="transition transform duration-300 ease-in" leave-from-class="translate-x-0 opacity-100"
        leave-to-class="translate-x-full opacity-0">
        <component :is="currentComponent" :key="currentComponent?.__asyncLoader?.toString() || 'component'" class=""
          :ref="ele => application.$currentView.value = ele" />
      </Transition>
    </template>
  </AppLayout>
</template>
