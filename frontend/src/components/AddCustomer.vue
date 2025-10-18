<template>
    <Dialog v-model:visible="isVisible" modal :header="$t('add_customer')" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
        <div class="flex flex-column gap-2 w-5">
            <label for="name">{{$t('name')}}</label>
            <InputText id="name" v-model="new_customer_name" aria-describedby="name" />
        </div>
        <div class="flex flex-column gap-2 w-5 mt-2">
            <label for="address">{{$t('address')}}</label>
            <InputText id="address" v-model="new_customer_address" aria-describedby="address" />
        </div>
        <div class="flex flex-column gap-2 w-5 mt-2">
            <label for="phone">{{$t('phone')}}</label>
            <InputText v-model="new_customer_phone" aria-describedby="phone" />
        </div>
        <template #footer>
            <ButtonGroup>
                <Button :label="$t('cancel')"  severity="secondary" aria-label="Cancel"  />
                <Button class="ml-2" severity="primary" @click="submitCustomer" :label="$t('save')" aria-label="Save" />
            </ButtonGroup>
        </template>
    </Dialog>
</template>

<script setup lang="ts">
import { ref,getCurrentInstance,defineEmits,defineProps,watch } from 'vue';
import { useToast } from "primevue/usetoast";
import axios from 'axios'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import ButtonGroup from 'primevue/buttongroup'
import Button from 'primevue/button'

const emit = defineEmits(['customer-added','update:visible'])

const props = defineProps({
    visible: {
        type: Boolean,
        required: true
    }
})

const isVisible = ref(false)

watch(() => props.visible, (newVal) => {
    isVisible.value = newVal;
});

watch(isVisible, (newVal) => {
    emit('update:visible', newVal);
});


const init = () => {
    isVisible.value = props.visible
}

init()


const new_customer_name = ref("")
const new_customer_phone = ref("")
const new_customer_address = ref("")
const new_customer_id = ref("")

const toast = useToast()
const { proxy } = getCurrentInstance();


const submitCustomer = () => {
   const payload = {
       name: new_customer_name.value,
       phone: new_customer_phone.value,
       address: new_customer_address.value
   };

   axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/customers`,{
       data : payload
   } , {
       headers: {
           Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
       }
   })
       .then(response => {
           toast.add({ severity: 'success', summary: 'Customer Added', detail: proxy.$t('done'),group:'br' });
           new_customer_id.value = response.data.data.id
           emit('customer-added')
       })
       .catch(error => {
           toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.data || 'An error occurred',group:'br' });
       });
}
</script>