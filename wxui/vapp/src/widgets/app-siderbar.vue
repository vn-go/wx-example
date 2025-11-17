<template>
    <ul>
        <li>
            <a @click="sideBar.doLoadView('/system/user')">Users</a>
        </li>
    </ul>
</template>
<script setup lang="ts">
import libs from '@core/lib';
class Sidebar extends libs.BaseUI {
    doLoadView(viewPath?:string) {
        libs.urlNav.move(viewPath);
    }
}
const sideBar =new Sidebar();
console.log(libs.getViewMap());
</script>