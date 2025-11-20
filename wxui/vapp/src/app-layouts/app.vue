<template>
    
    <header :ref="el => app.headerEle.value = el" class="widget-container  flex items-center">
      <div class="flex items-center space-x-3">
        <!-- Burger icon -->
        <button
            @click="app.toggleSidebar"
             class="w-8 h-8 flex items-center justify-center cursor-pointer hover:bg-gray-100 rounded transition"
        >
         <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
           <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                 d="M4 6h16M4 12h16M4 18h16"/>
         </svg>
       </button>

        <span class="text-sm font-medium hidden sm:inline text-gray-700">My App</span>
      </div>
    </header>
    <aside style="z-index:10000;"
  :ref="el => app.sideBarEle.value = el"
  class="sidebar widget-container"
 :class="{ 'translate-x-0': app.sidebarOpen }"
>
    <slot name="sidebar">
            
        </slot>
    </aside>
    
    <div class="flex flex-col widget-container " :ref="ele=>app.bodyEle.value=ele">
        <slot name="body">
            body

        </slot>
    </div>
   
</template>
<script lang="ts" setup>
import libs from '@core/lib';
class AppLayout extends libs.BaseUI {
    headerEle = libs.newDOMRef();
    bodyEle = libs.newDOMRef();
    sideBarEle = libs.newDOMRef();
    sidebarOpen = libs.newRef(true);
    onInit(): void {
       
        this.fixBodySize();
        this.fixSiderBar();
        const self=this;
        debugger;
        window.addEventListener("resize",()=>{
            self.fixBodySize();
        });
        this.hideSideBarWhenMouseMoveOnBody();
    }
    hideSideBarWhenMouseMoveOnBody() {
        const self=this;
        this.bodyEle.value.addEventListener("mousemove",()=>{
            if (self.isSideBarOpen()){
                self.toggleSidebar();
            }
           
        })
    }
    isSideBarOpen() {
        const sideBar = this.sideBarEle?.value;
        return sideBar.style.transform === 'translateX(0px)';
    }
    /*
        this funx will make position of sideBar top is height of header
        bottom is bottom of body
    */
    fixSiderBar() {
        const header = this.headerEle?.value;
        const sideBar = this.sideBarEle?.value;
        const body = this.bodyEle?.value;

        if (!header || !sideBar || !body) return;

        // Lấy chiều cao header
        const headerHeight = header.offsetHeight;

        // Lấy chiều cao body (viewport - header)
        const bodyRect = body.getBoundingClientRect();
        const top = headerHeight; // vị trí top
        const bottom = window.innerHeight - bodyRect.bottom; // khoảng cách tới bottom

        // Gán style cho sidebar
        sideBar.style.position = 'fixed';
        sideBar.style.top = `${top}px`;
        sideBar.style.bottom = `${0}px`;
        sideBar.style.overflowY='auto';
        //sideBar.style.height = `${bottom-top-8}px`;
    }
    /*
        this function will fix this.bodyEle.value: HTMLElement
        by make height of this.bodyEle.value= height of window- this.header.value.height
        and this.bodyEle.value.width= window.innerWidth
    */
    fixBodySize() {
        const header = this.headerEle?.value;
        const body = this.bodyEle?.value;

        if (!header || !body) return;

        const headerHeight = header.getBoundingClientRect().height;

        body.style.height = (window.innerHeight - headerHeight) + "px";
        body.style.width = window.innerWidth + "px";
        
    }
    toggleSidebar() {
        const sideBar = this.sideBarEle?.value;
    if (!sideBar) return;

    const isOpen = sideBar.style.transform === 'translateX(0px)';
    sideBar.style.transform = isOpen ? 'translateX(-100%)' : 'translateX(0)';
        
    }
}
const app=new AppLayout();
</script>
<style scoped>
.sidebar {
  transform: translateX(-100%);
  transition: transform 0.3s ease-in-out;
}
.sidebar.open {
  transform: translateX(0);
}
.debug2 {
    border: solid 4px blue;
}
</style>