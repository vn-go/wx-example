<template>
    <div class="editor-form">

        <div class="body grid grid-cols-4 grid-cols-[1fr_120px_1fr_3fr]">
            <!-- Row 1: Code | Name -->
            <div class="col-span-1">
                <label>Code</label>
            </div>
            <div class="col-span-1">
                <input v-model="role.dataItem.data.code" type="text" />
            </div>

            <div class="col-span-1">
                <label>Name</label>
            </div>
            <div class="col-span-1">
                <input v-model="role.dataItem.data.name" type="text" />
            </div>

            <!-- Row 2: CreatedOn | CreatedBy -->
            <div class="col-span-1">
                <label>Created On</label>
            </div>
            <div class="col-span-1">
                <input v-model="role.dataItem.data.createdOn" type="text" readonly />
            </div>

            <div class="col-span-1">
                <label>Created By</label>
            </div>
            <div class="col-span-1">
                <input v-model="role.dataItem.data.createdBy" type="text" readonly />
            </div>

            <!-- Row 3: ModifiedOn | ModifiedBy -->
            <div class="col-span-1">
                <label>Modified On</label>
            </div>
            <div class="col-span-1">
                <input v-model="role.dataItem.data.modifiedOn" type="text" />
            </div>

            <div class="col-span-1">
                <label>Modified By</label>
            </div>
            <div class="col-span-1">
                <input v-model="role.dataItem.data.modifiedBy" type="text" />
            </div>

            <!-- Row 4: Description -->
            <div class="col-span-1 label-top">
                <label>Description</label>
            </div>
            <div class="col-span-3">
                <textarea v-model="role.dataItem.data.description" rows="6"></textarea>
            </div>
        </div>
        <div class="footer">

            <button class="btn-update" v-if="role.dataItem.data.id" @click="() => {
                role.doUpdate()
            }">Update</button>
            <button class="btn-insert" v-if="role.dataItem.data.id == null" @click="() => {
                role.doAddNew()
            }">Add</button>
            <button class="btn-cancel" @click="() => {
                role.doClose()
            }">Cancel</button>
        </div>
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
    dataItem = libs.newRef({
        data: {}
    });
    token = ""
    async onPreInit() {
        this.apiDiscovery(
            [
                "system/roles/new-item",
                "system/roles/item",
                "system/roles/update",
                "system/roles/insert",
            ]
        )
    }
    async onInit() {
        if (roleId) {
            let res = await this.remoteCaller.post("system/roles/item", roleId);
            if (res.ok) {
                this.dataItem.value = res.data;
                this.token = res.data.token;
            }
        } else {
            let res = await this.remoteCaller.post("system/roles/new-item", roleId);
            if (res.ok) {
                this.dataItem.value = res.data;
                this.token = res.data.token;
            }
        }

    }
    async doUpdate() {
        let res = await this.remoteCaller.post("system/roles/update", this.dataItem);
        if (res.ok) {
            this.dataItem.value = res.data;
            this.token = res.data.token;
        }
        console.log(res);
    }
    async doAddNew() {
        let res = await this.remoteCaller.post("system/roles/insert", this.dataItem);
        if (res.ok) {
            this.dataItem.value = res.data;
            this.token = res.data.token;
        }
        console.log(res);
    }

}
const role = libs.newReactive(new Role());
</script>