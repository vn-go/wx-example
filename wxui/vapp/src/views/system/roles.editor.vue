<template>
    <div style="min-width: 500px;min-height: 400px;" class="flex flex-col debug">
        {{ roleId }}
    </div>
</template>
<script lang="ts" setup>
import libs from '@core/lib';

const { roleId } = defineProps({
  roleId: {
    type: String,
    required: true
  }
});

class Role extends libs.BaseUI {
    async onPreInit(){
        this.apiDiscovery(
            [
               
                "roles/get-item",
                
            ]
        )
    }
    async onInit() {
        let res= await this.remoteCaller.post("roles/get-item",roleId);
        console.log(res);
    }

}
const role=new Role();
</script>