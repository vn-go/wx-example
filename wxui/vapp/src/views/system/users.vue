<template>
    <div class="debug flex flex-col w-full h-full">

    </div>
</template>
<script setup>
import libs from '@core/lib';

    class Users extends libs.BaseUI {
        onInit() {
        
        }
    }
    const user=libs.newReactive(new Users());
</script>