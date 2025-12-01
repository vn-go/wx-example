<template>
    <div class="w-full h-full">
        <DataTable :value="roles.listOfRole" :reorderableColumns="true" scrollable sortMode="multiple" resizableColumns
            columnResizeMode="fit" showGridlines>
            <Column header="Code" field="code"></Column>
            <Column header="Name" field="name"></Column>
            <Column header="Users" field="totalUsers"></Column>
            <Column header="description" field="description"></Column>
            <Column>
                <template #body="{ data }">
                    <button v-if="data.id != null" class="btn" @click="() => {
                        roles.doEdit(data.id)
                    }">Edit</button>
                    <button v-if="data.id == null" class="btn" @click="() => {
                        roles.doNewItem()
                    }">New</button>
                </template>
            </Column>
        </DataTable>
    </div>
</template>
<script lang="ts" setup>
import libs from '@core/lib';
import { Column, DataTable } from 'primevue';
class Roles extends libs.BaseUI {
    listOfRole = libs.newRef();
    async onPreInit() {
        this.apiDiscovery(
            [
                "roles/list",
                "roles/item-get",
                "roles/item-update",
                "roles/add-users",
                "roles/remove-users",
                "roles/item-delete",
            ]
        )
    }
    async onInit() {
        this.listOfRole.value = await this.getListOfRoles();

    }
    async getListOfRoles() {
        let res = await this.remoteCaller.post("system/roles/list", {
            "index": 0,
            "size": 20,

        });
        if (res.ok) {
            res.data.push({});
            return res.data;
        }
    }
    async doEdit(roleId) {
        await this.newModal("views/system/roles.editor").setTitle("Edit role").setData({
            roleId: roleId
        }).render();
    }
    async doNewItem() {
        await this.newModal("views/system/roles.editor").setTitle("New role").setData({
            roleId: null
        }).render();
    }
}
const roles = libs.newReactive(new Roles("system/roles"));

</script>