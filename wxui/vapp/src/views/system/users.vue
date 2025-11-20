<template>
    <div class="flex flex-col w-full h-full">
        <DataTable :value="user.users" 
               :reorderableColumns="true"
                scrollable 
               sortMode="multiple"
               resizableColumns columnResizeMode="fit" showGridlines
               >
            <Column field="username" header="Username"   style="width: 220px;" sortable ></Column>
            <Column field="email" header="Email"   style="width: 220px;" sortable ></Column>
            <Column field="createdOn" header="Created On"   style="width: 220px;" sortable >
                <template #body="{ data }">
                    {{ user.formatDate(data.createdOn) }}
                </template>
            </Column>
            <Column field="isActive" header="Active" style="width: 220px;" sortable ></Column>
            <Column field="roleName" header="Role" tyle="width: 100%;"></Column>
            <Column   tyle="width:120">
                <template #body="{ data }">
                    <button class="btn" @click="user.doEdit(data)">
                            Edit
                    </button>
                </template>
            </Column>
        </DataTable>
    </div>
</template>
<script setup>
import libs from '@core/lib';
import { Column, DataTable } from 'primevue';
    
    class Users extends libs.BaseUI {
        users=libs.newRef();
        async onInit() {
            
            let retApi=await this.loadData();
          
            this.users.value=retApi.data ;
        }
        async loadData(){
            return await libs.api.post("/accounts/get-list-of-accounts",{
                "first": 0,
                "last": 1000
                })
        }
        async doEdit(data) {
            // let modal=libs.newModal("html/modal.html");
            // console.log(modal);
            // await modal.load("views/system/users");
            await libs.showModal("views/system/users.editor",{userId:data.userId})
        }

    }
    const user=libs.newReactive(new Users());
    
   
</script>