<template>
    <div class="w-full h-full debug">

        <EditorForm>
            <template #header>
                <span>Edit user</span>
            </template>
            <template #body>
                <div class="grid grid-cols-2 gap-6">

<!-- Left Column -->
<div class="space-y-4">
  

  <div>
    <label class="font-semibold">Username</label>
    <input disabled value="root" class="input readonly">
  </div>

  <div>
    <label class="font-semibold">Email</label>
    <input type="email" v-model="userEditor.userData.email" class="input">
  </div>

  <div class="flex items-center gap-3">
    <input type="checkbox" v-model="userEditor.userData.isActive">
    <label>Active</label>
  </div>

  <div class="flex items-center gap-3">
    <input type="checkbox" v-model="userEditor.userData.isSysAdmin">
    <label>System Admin</label>
  </div>
</div>

<!-- Right Column -->
<div class="space-y-4">

  <div>
    <label class="font-semibold">Role ID</label>
    <select v-model="userEditor.userData.roleId" class="input">
      <option :value="null">-- No role --</option>
      
    </select>
  </div>

  <div>
    <label class="font-semibold">Role Code</label>
    <input v-model="userEditor.userData.roleCode" class="input">
  </div>

  <div>
    <label class="font-semibold">Created By</label>
    <input disabled value="admin" class="input readonly">
  </div>

  <div>
    <label class="font-semibold">Created On</label>
    <input disabled value="2025-11-10T10:06:48Z" class="input readonly">
  </div>

  <div>
    <label class="font-semibold">Modified On</label>
    <input disabled value="" class="input readonly">
  </div>

</div>
</div>



            </template>
            <template #footer>
                <FormFooter/>
            </template>
        </EditorForm>
    </div>
</template>
<script lang="ts" setup>
import EditorForm from '@app-layouts/editor-form.vue';
import FormFooter from '@components/FormFooter.vue';
import libs from '@core/lib';
class UserEditor extends libs.BaseUI {
    userId= undefined;
    userData=libs.newRef({});
    errMsg= libs.newRef();
    
    async onInit() {
        let res= await libs.api.post("accounts/get-edit",{userId:this.userId});
        if (!res.ok){
            this.errMsg.value= res.error.statusText
        }else {
            this.userData.value= res.data;
        }

    }
}
const instance =new UserEditor()
const  userEditor=libs.newReactive(instance) ;
const { userId } = defineProps({
  userId: {
    type: Number,
    required: true
  }
});
userEditor.userId=userId;
defineExpose({
    instance
   })
</script>