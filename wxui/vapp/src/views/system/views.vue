<template>
<div class="h-full debug">
    {{ views.getViewPath() }}
    <button class="btn" @click="()=>{
        views.showViewPath();
    }">Test</button>
</div>
</template>
<script setup>
import libs from '@core/lib';


    class Views extends libs.BaseUI {
        onInit() {
            this.loadAllViews();
        }
        async loadAllViews() {
            const res= await this.remoteCaller.post("/view-manager/get-list-of-views");
            console.log(res);

        }
        showViewPath(){
            alert(this.getViewPath());
        }
    }
    const views= new Views("system/views");
    views.loadAllViews()

</script>