<template>
    <div class="w-full h-full debug">

        <EditorForm>
            <template #header>
                <span>Edit user</span>
            </template>
            <template #body>
                <div class="debug">{{ userEditor.errMsg }}</div>
            </template>
        </EditorForm>
        
    </div>
</template>
<script lang="ts" setup>
import EditorForm from '@app-layouts/editor-form.vue';
import libs from '@core/lib';
class UserEditor extends libs.BaseUI {
    userId= undefined;
    userData=undefined;
    errMsg= libs.newRef();
    async onInit() {
        let res= await libs.api.post("accounts/get-edit",{userId:this.userId});
        if (!res.ok){
            this.errMsg= res.error.statusText
        }

    }
}
let  userEditor=libs.newReactive(new UserEditor()) ;
const { userId } = defineProps({
  userId: {
    type: Number,
    required: true
  }
});
userEditor.userId=userId;
</script>