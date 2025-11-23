<template>
    <div style="min-width: 500px;min-height: 400px;" class="flex flex-col debug">
        <input type="text" v-model="role.dataItem.data.code" />
        <button @click="()=>{
            role.doUpdate()
        }">Update</button>
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
    dataItem=libs.newRef({
        data:{}
    });
    token=""
    async onPreInit(){
        this.apiDiscovery(
            [
               
                "roles/get-item",
                "roles/update-item"
            ]
        )
    }
    async onInit() {
        let res= await this.remoteCaller.post("roles/get-item",roleId);
        if (res.ok){
            this.dataItem.value=res.data;
            this.token= res.data.token;
        }
        console.log(res);
    }
    async doUpdate(){
        let res= await this.remoteCaller.post("roles/update-item",this.dataItem);
        if (res.ok){
            this.dataItem.value=res.data;
            this.token= res.data.token;
        }
        console.log(res);
    }

}
const role=libs.newReactive(new Role()) ;
</script>