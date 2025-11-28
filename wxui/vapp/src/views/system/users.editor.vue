<template>
    <div class="w-full h-full">

        <EditorForm>
            <template #header>
                <span class="font-semibold text-gray-800">Edit user</span>
            </template>
            <template #body>
                <div class="p-6 form-data">
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-8">

                        <div class="space-y-6">
                            <div>
                                <label>Full name</label>
                                <input v-model="userEditor.userData.displayName">
                            </div>
                            <div>
                                <label>Username</label>
                                <input v-model="userEditor.userData.username">
                            </div>

                            <div>
                                <label>Email</label>
                                <input type="email" v-model="userEditor.userData.email">
                            </div>

                            <div class="flex items-center gap-3 pt-2">
                                <input type="checkbox" v-model="userEditor.userData.isActive">
                                <label>Active</label>
                            </div>

                            <div class="flex items-center gap-3">
                                <input type="checkbox" v-model="userEditor.userData.isSysAdmin">
                                <label>System Admin</label>
                            </div>
                        </div>

                        <div class="space-y-6">

                            <div>
                                <label>Role ID</label>
                                <v-select label="name" v-model="userEditor.userData.roleId" :options="userEditor.roles"
                                    :reduce="(item: any) => item.id"></v-select>
                            </div>



                            <div>
                                <label>Created By</label>
                                <input v-model="userEditor.userData.createdBy">
                            </div>

                            <div>
                                <label class="block text-sm font-medium text-gray-700 mb-1">Created On</label>
                                <!-- <input 
                        disabled 
                        value="2025-11-10T10:06:48Z" 
                    > -->
                                <!-- <vue-datePicker :text-input="true" v-model="userEditor.userData.createdOn"
                                    :enable-tab-navigation="true" /> -->
                                <DatePicker v-model="userEditor.userData.createdOn" />

                            </div>

                            <div>
                                <label class="block text-sm font-medium text-gray-700 mb-1">Modified On</label>
                                <input v-model="userEditor.userData.modifiedOn">
                            </div>

                        </div>
                    </div>
                </div>

            </template>
            <template #footer>
                <div class="">
                    <FormFooter @onClose="() => {
                        userEditor.doClose();
                    }" @onUpdate="() => {
                        userEditor.doUpdate();
                    }" />
                </div>
            </template>
        </EditorForm>
    </div>
</template>
<script lang="ts" setup>
import EditorForm from '@app-layouts/editor-form.vue';
import FormFooter from '@components/FormFooter.vue';
import libs from '@core/lib';
// import '@vuepic/vue-datepicker/dist/main.css';
import DatePicker from 'primevue/datepicker';
import vSelect from 'vue-select';
class UserEditor extends libs.BaseUI {
    userId = undefined;
    userData = libs.newRef({});
    roles = libs.newRef([]);
    errMsg = libs.newRef();
    updateToken: any;

    async onPreInit() {
        // register all api for this UI
        // this.apiDiscovery([
        //     "system/users/item",
        //     "accounts/update",
        //     "accounts/delete",
        //     "accounts/new",
        // ]);
    }
    async onInit() {
        let resOfRoles = await libs.api.post(this.getViewPath(), "system/users/roles", {});
        console.log(resOfRoles);
        if (!resOfRoles.ok) {
            this.errMsg.value = resOfRoles.error.statusText
        } else {
            this.roles.value = resOfRoles.data;
        }
        let res = await libs.api.post(this.getViewPath(), "system/users/get-item", {
            "userId": this.userId
        });
        if (!res.ok) {
            this.errMsg.value = res.error.statusText
        } else {
            this.userData.value = res.data.data;
            this.updateToken = res.data.token;
        }

    }
    async doUpdate() {
        let res = await libs.api.post(this.getViewPath(), "system/users/update", {
            data: this.userData,
            token: this.updateToken,
        });
    }

}
const instance = new UserEditor()
const userEditor = libs.newReactive(instance);
const { userId } = defineProps({
    userId: {
        type: Number,
        required: true
    }
});
userEditor.userId = userId;
defineExpose({
    instance
})
</script>