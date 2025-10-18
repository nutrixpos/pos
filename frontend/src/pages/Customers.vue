<template>
    <div class="w-full">
       <div class="grid mx-2">
           <div class="col-12 flex">
               <div class="gird w-full">
                   <div class="col-12">
                       <h3>{{$t('customer',3)}}</h3>
                   </div>
                   <div class="col-12">
                       <DataTable @page="updatCustomersTableRowsPerPage" :lazy="true" :totalRecords="customersTableTotalRecords" :loading="isCustomersTableLoading"  paginatorPosition="both"  paginator :rows="customersTableRowsPerPage" :rowsPerPageOptions="[50, 100, 500]" :value="customers" stripedRows tableStyle="min-width: 50rem;max-height:50vh;" class="w-full pr-2">
                           <template #header>
                               <div class="flex justify-start">
                                   <Button icon="pi pi-plus" :label="$t('add_customer')"  rounded raised @click="customerAddDialog=true" />
                               </div>
                           </template>
                           <Column sortable field="name" :header="$t('name')"></Column>
                           <Column field="address" :header="$t('address')"></Column>
                           <Column field="phone" :header="$t('phone',1)"></Column>
                           <Column :header="$t('actions')">
                               <template #body="slotProps">
                                   <ConfirmPopup></ConfirmPopup>
                                   <ButtonGroup>
                                       <Button icon="pi pi-pencil" severity="secondary" aria-label="Edit" @click="prepareCustomerToEdit(slotProps.data)"/>
                                       <Button icon="pi pi-trash" severity="danger" aria-label="Remove" @click="confirmDeleteCustomer($event,slotProps.data.id)"/>
                                   </ButtonGroup>
                               </template>
                           </Column>
                       </DataTable>
                       <Dialog v-model:visible="customerAddDialog" modal :header="$t('add_customer')" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
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
                       <Dialog v-model:visible="customerEditDialog" modal :header="`${$t('edit_customer')} ${customerToEdit.name}`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
                           <div class="flex flex-column gap-2 w-5">
                               <label for="name">{{$t('name')}}</label>
                               <InputText id="name" v-model="customerToEdit.name" aria-describedby="name" />
                           </div>
                           <div class="flex flex-column gap-2 w-5 mt-2">
                               <label for="address">{{$t('address')}}</label>
                               <InputText id="address" v-model.number="customerToEdit.address" aria-describedby="address" />
                           </div>
                           <div class="flex flex-column gap-2 w-5 mt-2">
                               <label for="phone">{{$t('phone')}}</label>
                               <InputText id="phone" v-model="customerToEdit.phone" aria-describedby="phone" />
                           </div>
                           <template #footer>
                               <ButtonGroup>
                                   <Button :label="$t('cancel')"  severity="secondary" aria-label="Cancel"  />
                                   <Button class="ml-2" severity="primary" @click="updateCustomer" :label="$t('save')" aria-label="Save" />
                               </ButtonGroup>
                           </template>
                       </Dialog>
                   </div>
               </div>
           </div>
       </div>
   </div>
</template>

<script setup lang="ts">
import ConfirmPopup from 'primevue/confirmpopup';
import DataTable from 'primevue/datatable';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Column from 'primevue/column';
import Button from 'primevue/button';
import ButtonGroup from 'primevue/buttongroup';
import axios from 'axios'
import { ref,getCurrentInstance } from "vue";
import { useToast } from "primevue/usetoast";
import { useConfirm } from "primevue/useconfirm";
// import { Material } from '@/classes/OrderItem';

const { proxy } = getCurrentInstance();
const confirm = useConfirm();

const customerToEdit = ref<any>({})
const customerEditDialog =  ref(false)

const customerAddDialog = ref(false)
const new_customer_name = ref("")
const new_customer_phone = ref("")
const new_customer_id = ref("")
const new_customer_address = ref("")


const customersTableTotalRecords  = ref(0)
const customersTableRowsPerPage = ref(50)
const isCustomersTableLoading = ref(true)
const customers = ref([])
const toast = useToast()




const deleteCustomer = (customer_id: string) => {
   axios.delete(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/customers/${customer_id}`, {
       headers: {
           Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
       }
   }).then(() => {
       toast.add({ severity: 'success', summary: 'Success', detail: 'Customer deleted successfully',group:'br',life:3000 });
       loadCustomers();
   }).catch(() => {
       toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to delete customer',group:'br',life:3000 });
   });
}


const confirmDeleteCustomer = (event,customer_id) => {
   confirm.require({
       target: event.currentTarget,
       message: 'Are you sure you want to delete this custmoer ?',
       icon: 'pi pi-exclamation-triangle',
       rejectProps: {
           label: 'Cancel',
           severity: 'secondary',
           outlined: true
       },
       acceptProps: {
           label: 'Yes'
       },
       accept: () => {


           deleteCustomer(customer_id)

        
       },
       reject: () => {
       }
   });
 }


const prepareCustomerToEdit = (customer: any) => {
   customerToEdit.value = JSON.parse(JSON.stringify(customer))
   customerEditDialog.value = true
}


const updateCustomer = () => {

   axios.patch(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/customers/${customerToEdit.value.id}`, 
   {
       data : customerToEdit.value
   },{
       headers: {
           Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
       }
   })
   .then(() => {
       toast.add({ severity: 'success', summary: 'Customer Updated', detail: proxy.$t('done'),group:'br',life:3000 });
       customerEditDialog.value = false;
       customerToEdit.value = {}
       loadCustomers()
   })
   .catch(error => {
       toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.message || 'An error occurred',group:'br' });
   });
}


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
           customerAddDialog.value = false;
           loadCustomers()
       })
       .catch(error => {
           toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.data || 'An error occurred',group:'br' });
       });
}



const updatCustomersTableRowsPerPage = (event: any) => {
   loadCustomers(event.first,event.rows)
}


const loadCustomers = (first=0,rows=customersTableRowsPerPage.value) => {
   isCustomersTableLoading.value = true;

   if (first == 0){
       first = 1
   }

   const page_number = Math.ceil((first/rows))
   
   axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/customers`, {
       params: {
           "page[number]": page_number,
           "page[size]": rows
       },
       headers: {
           Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
       }
   })
   .then(response => {
       customers.value = response.data.data;
       customersTableTotalRecords.value = response.data.meta.total_records;
   })
   .catch(() => {
       toast.add({severity:'error', summary: 'Error', detail: 'Failed to load customers'});
   })
   .finally(() => {
       isCustomersTableLoading.value = false;
   });
}

loadCustomers()


</script>