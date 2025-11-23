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
            <Column   style="width:120">
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
        async onPreInit() {
        // register all api for this UI
            this.apiDiscovery([
                "accounts/get-edit",
                "accounts/update-by-id",
                "accounts/delete",
                "accounts/new",
            ]);
        }
        async onInit(uiEle) {
            
            let retApi=await this.loadData();
          
            this.users.value=retApi.data ;
        }
        async loadData(){
            
            return await libs.api.post(this.getViewPath(),"/accounts/get-list-of-accounts",{
                "first": 0,
                "last": 1000
                })
        }
        
        async doEdit(data) {
           
            await this.newModal("views/system/users.editor").setTitle("test").setData({userId:data.userId}).setSize(700,500).render();
            
        }

    }
    const user=libs.newReactive(new Users("system/users"));
   

   
</script>