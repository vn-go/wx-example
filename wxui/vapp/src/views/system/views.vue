<template>
<div class="flex flex-row w-full h-full">
    <div class="flex flex-col h-full debug" style="width:400px;">
        <DataTable :value="views.data"
    :reorderableColumns="true"
                scrollable 
               sortMode="multiple"
               resizableColumns 
               columnResizeMode="fit" 
               showGridlines
               selectionMode="single"
               v-model:selection="views.selecteView"
               @rowSelect="()=>{
                    views.onViewSelect();
               }">
        <Column field="viewPath" header="Form"></Column>

    </DataTable>
    </div>
    <div class="flex flex-col h-full w-full debug">
        <DataTable :value="views.apiData"
    :reorderableColumns="true"
                scrollable 
               sortMode="multiple"
               resizableColumns columnResizeMode="fit" showGridlines>
        <Column field="apiPath" header="API"></Column>
        <Column field="title" header="Title"></Column>
        <Column field="description" header="Description"></Column>
        <Column field="createdOn" header="Create On">
            <template #body="{ data }">
                {{ dayjs(data.createdOn).format("DD/MM/YYYY HH:mm") }}
            </template>
        </Column>
        <Column>
            <template #body="{ data }">
                <button class="btn" @click="()=>{
                    views.deleteApi(data.apiPath)
                }">Delete</button>
            </template>
        </Column>

    </DataTable>
    </div>
    
</div>
</template>
<script setup>
import libs from '@core/lib';
import dayjs from "dayjs";
import { Column, DataTable } from 'primevue';

    class Views extends libs.BaseUI {
        data=libs.newRef();
        apiData=libs.newRef();
        selecteView=libs.newRef();
        async onPreInit(){
            this.apiDiscovery([
                "view-manager/get-list-of-views",
                "view-manager/get-list-api-of-view",
                "view-manager/api-delete"
            ]);
        }
        onInit() {
            this.loadAllViews();
        }
        async onViewSelect(){
             await this.loadApiOfView(this.selecteView.viewPath);
        }
        async loadAllViews() {
            const res= await this.remoteCaller.post("view-manager/get-list-of-views");
            if(res.ok){
                this.data=res.data;
                if(this.data.length>0) {
                    await this.loadApiOfView(this.data[0].viewPath);
                }
            }
            console.log(this.data);

        }
        async loadApiOfView(viewPath) {
            const res= await this.remoteCaller.post("view-manager/get-list-api-of-view",{
                viewPath:viewPath
            });
            if(res.ok){
                this.apiData=res.data;
            }
            console.log(this.apiData);
        }
        async deleteApi(apiPath) {
            const res= await this.remoteCaller.post("view-manager/api-delete",{
                apiPath:apiPath
            });
            if(res.ok){
                this.apiData=res.data;
            }
            console.log(this.apiData);
        }
    }
  
   
    const views=libs.newReactive(new Views("system/views"));
    views.loadAllViews()
</script>