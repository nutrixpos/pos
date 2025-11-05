<template>
     <div class="w-full">
        <div class="grid mx-2">
            <div class="col-12 flex">
                <div class="gird w-full">
                    <div class="col-12">
                        <h3>{{ $t('material',3) }}</h3>
                    </div>
                    <div class="col-12 flex justify-content-center align-items-center w-full">
                        <DataTable :value="inventory_components" stripedRows tableStyle="min-width: 50rem" class="w-full">
                            <template #header>
                                <div class="flex flex-wrap items-center justify-between align-items-center gap-2">
                                    <Button icon="pi pi-plus" :label="$t('add_material')" @click="add_component_dialog = true" rounded raised />
                                </div>
                            </template>
                            <Column :header="$t('actions')" style="width:5rem">
                                <template #body="slotProps">
                                    <Button icon="fa fa-dolly-flatbed" severity="secondary" aria-label="ShowEntries" @click="entries_dialog_material = slotProps.data; entries_dialog=true; loadEntries()"/>
                                </template>
                            </Column>
                            <Column field="name" :header="$t('name')"></Column>
                            <Column field="totalAmount" :header="$t('quantity')"></Column>
                            <Column field="unit" :header="$t('unit')"></Column>
                            <Column :header="$t('status')">
                                <template #body="slotProps">
                                    <Tag :value="slotProps.data.totalAmount > slotProps.data.settings.stock_alert_treshold ? 'INSTOCK' : 'LOWSTOCK'" :severity="slotProps.data.totalAmount > slotProps.data.settings.stock_alert_treshold ? 'success' : 'danger'" />
                                </template>
                            </Column>
                            <Column :header="$t('actions')" style="width:30rem">
                                <template #body="slotProps">
                                    <ButtonGroup>
                                        <Button icon="pi pi-clock" :label="$t('history')" @click="loadComponentLogs(slotProps.data.id)" severity="secondary" aria-label="Save"  />
                                        <Button icon="pi pi-pencil" severity="secondary" @click="edit_material = slotProps.data; edit_material_dialog=true" aria-label="Edit"  />
                                        <Button icon="pi pi-cog" severity="secondary" aria-label="Settings" @click="material_settings = slotProps.data; material_settings_dialog=true"  />
                                        <Button icon="pi pi-times" severity="danger" aria-label="Delete" @click="confirmDeleteMaterial(slotProps.data.id)"/>
                                    </ButtonGroup>
                                </template>
                            </Column>
                        </DataTable>
                    </div>
                </div>
            </div>
            <Dialog v-model:visible="edit_material_dialog" modal :header="`Edit data for  ${edit_material?.name}`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
                <EditMaterial v-if="edit_material != undefined" @returnMaterial="(material) => saveEditedMaterial(material)" :material="edit_material" />
            </Dialog>
            <Dialog v-model:visible="material_settings_dialog" modal :header="`Settings for  ${material_settings?.name}`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
                <div class="flex align-items-center">
                    <h4>stock_alert_treshold</h4>
                    <InputText type="number" class="ml-2" id="stock_alert_treshold" v-model.number="material_settings.settings.stock_alert_treshold" aria-describedby="stock_alert_treshold" />
                </div>
                <template #footer>
                    <ButtonGroup>
                        <Button label="Cancel" @click="material_settings_dialog=false" severity="secondary" aria-label="Save"  />
                        <Button class="ml-2" severity="primary" label="Save" aria-label="Save" @click="saveMaterialSettings"/>
                    </ButtonGroup>
                </template>
            </Dialog>
            <Dialog v-if="entries_dialog" :header="`Entries for ${ entries_dialog_material?.name }`" v-model:visible="entries_dialog" :style="{ width: '90rem' }">
                <div class="p-0">
                    <div class="flex flex-column w-8 xl:w-4">
                        <InputText class="m-1" :placeholder="$t('company')" v-model="new_entry_company" aria-describedby="name-help" />
                        <InputText class="m-1" :placeholder="$t('purchase_quantity')" v-model="new_entry_quantity" aria-describedby="name-help" />
                        <InputText class="m-1" :placeholder="$t('purchase_price')" v-model="new_entry_price" aria-describedby="name-help" />
                        <FloatLabel class="mx-1">
                            <Calendar inputId="new_entry_expiration_date" v-model="new_entry_expiration_date" showIcon />
                            <label for="new_entry_expiration_date">Expiration date</label>
                        </FloatLabel>
                        <Button icon="pi pi-plus" label="Add Entry" class="my-1" @click="addNewEntry(entries_dialog_material.id)" severity="info" raised />
                    </div>
                    <DataTable @page="updatEntriesTableRowsPerPage" :lazy="true" :totalRecords="entriesTableTotalRecords" :loading="isEntriesTableLoading"  paginatorPosition="both"  paginator :rows="entriesTableRowsPerPage" :rowsPerPageOptions="[50, 100, 500]"stripedRows :value="entries_dialog_material.entries" v-model:expandedRows="expandedEntryRows">
                        <Column expander style="width: 5rem" />
                        <Column field="company" :header="$t('company')"></Column>
                        <Column field="quantity" :header="$t('quantity')" sortable></Column>
                        <Column field="expiration_date" :header="$t('expiration_date')" sortable></Column>
                        <Column header="Actions" style="width:30rem">
                            <template #body="slotProps">
                                <ButtonGroup>
                                    <Button icon="pi pi-times" :label="$t('delete')" severity="secondary" aria-label="Delete" @click="confirmDeleteEntry($event,entries_dialog_material?.id,slotProps.data.id)" />
                                </ButtonGroup>
                            </template>
                        </Column>
                        <template #expansion="slotProps">
                            <DataTable :value="slotProps">
                                <Column field="purchase_quantity" :header="$t('purchase_quantity')" sortable></Column>
                                <Column field="purchase_price" :header="$t('purchase_price')" sortable></Column>
                            </DataTable>
                        </template>
                    </DataTable>
                </div>
            </Dialog>
            <Dialog v-model:visible="add_component_dialog" modal :header="`Add new inventory material`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '90vw', '575px': '90vw' }">
               <div class="md:w-full">
                    <div class="flex flex-column gap-2">
                        <label for="name">Name</label>
                        <InputText id="name" v-model="new_component_name" aria-describedby="name-help" />
                        <!-- <small id="name-help">Enter the component name</small> -->
                    </div>
                    <div class="flex flex-column gap-2 mt-3 ">
                        <label for="unit">Unit</label>
                        <InputText id="unit" v-model="new_component_unit" aria-describedby="unit-help" />
                    </div>
                    <Divider />
                    <h4>Entries</h4>
                    <div class="flex flex-column w-full xl:w-3">
                        <InputText class="m-1" placeholder="Company" v-model="new_component_entry_company" aria-describedby="name-help" />
                        <InputText class="m-1" placeholder="Quantity" v-model="new_component_entry_quantity" aria-describedby="name-help" />
                        <InputText class="m-1" placeholder="Total Price" v-model="new_component_entry_price" aria-describedby="name-help" />
                        <div>
                            <Button class="mx-1 my-2" label="Add" @click="new_component_entries.push({company: new_component_entry_company, quantity: new_component_entry_quantity, unit: new_component_unit, purchase_price: new_component_entry_price})" />
                        </div>
                    </div>
                    <DataTable :value="new_component_entries">
                        <Column field="company" header="Company"></Column>
                        <Column field="quantity" header="Quantity"></Column>
                        <Column field="unit" header="Unit">
                            <template #body="slotProps">
                                {{ slotProps.data.unit }}
                            </template>
                        </Column>
                        <Column field="purchase_price" header="Total Price"></Column>
                        <Column :header="$t('actions')" style="width:30rem">
                            <template #body="slotProps">
                                <ButtonGroup>
                                    <Button icon="pi pi-times" severity="secondary" aria-label="Delete" @click="new_component_entries.splice(new_component_entries.findIndex(el => el === slotProps.data), 1)" />
                                </ButtonGroup>
                            </template>
                        </Column>
                    </DataTable>

                    <div class="flex w-full mt-5 justify-content-center align-items-center">
                        <Button label="Submit" class="lg:w-6" @click="submitNewComponent" />
                    </div>
               </div>
            </Dialog>
            <Dialog v-model:visible="component_logs_dialog" modal :header="`Consumption for  ${component_logs_name}`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
                <DataTable @rowExpand="onComponentLogRowExpand" @page="updatLogsTableRowsPerPage" :lazy="true" :totalRecords="logsTableTotalRecords" :loading="isLogsTableLoading" v-model:expandedRows="expandedComponentLogsRows" paginatorPosition="both"  paginator :rows="logsTableRowsPerPage" :rowsPerPageOptions="[7, 14, 30, 90]" :value="component_logs" stripedRows tableStyle="min-width: 50rem;max-height:50vh;" class="w-full pr-2">
                    <Column expander style="width: 5rem" />
                    <Column field="date" header="Date"></Column>
                    <Column  header="Quantity">
                        <template #body="slotProps">
                            <Tag :value="slotProps.data.type == 'component_consume' ? `- ${slotProps.data.quantity}` : `+ ${slotProps.data.quantity}`" :severity="slotProps.data.type == 'component_consume' ? 'danger' : 'success'" />
                        </template>
                    </Column>
                    <Column field="order_id" header="Order Id"></Column>
                    <template #expansion="slotProps">
                        <div class="p-4">
                            <h4>Order Items</h4>
                            <MaterialLogsOrderItemsTable v-if="slotProps.data.order" :items="slotProps.data.order.items" :order_item_index="slotProps.data.order_item_index" />
                            <div v-else>
                                Loading ...
                            </div>
                        </div>
                    </template>
                </DataTable>
            </Dialog>
        </div>
        <ConfirmDialog></ConfirmDialog>
    </div>
</template>

<script setup lang="ts">
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import axios from 'axios'
import Button from 'primevue/button'
import ButtonGroup from 'primevue/buttongroup'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import { useConfirm } from "primevue/useconfirm";
import Tag from 'primevue/tag'
import { Material, MaterialEntry } from '@/classes/OrderItem';
import Calendar from 'primevue/calendar';
import FloatLabel from 'primevue/floatlabel'
import EditMaterial from '@/components/EditMaterial.vue'
import ConfirmDialog from 'primevue/confirmdialog'
import MaterialLogsOrderItemsTable from '@/components/MaterialLogsOrderItemsTable.vue';
import {Divider} from 'primevue'
// import Message from 'primevue/message'
  
import { ref,getCurrentInstance } from "vue";
import { useToast } from "primevue/usetoast";

const { proxy } = getCurrentInstance();


const entries_dialog = ref(false)
const entries_dialog_material = ref<Material>()
const entriesTableTotalRecords = ref(0)
const isEntriesTableLoading = ref(false)
const entriesTableRowsPerPage = ref(50)

const confirm = useConfirm();

const edit_material_dialog = ref(false)
const edit_material = ref<Material>()


const logsTableRowsPerPage = ref(7)
const logsTableTotalRecords = ref(0)
const isLogsTableLoading = ref(true)



  const toast = useToast();

  const material_settings_dialog = ref(false)
  const material_settings = ref<Material>()
  
  const expandedEntryRows = ref([]);
  const expandedComponentLogsRows = ref([])
  
  const inventory_components = ref([])

  const component_logs = ref([])
  const component_logs_dialog = ref(false)
  const component_logs_name = ref("")


  const add_component_dialog = ref(false)


  const new_component_name = ref("")
  const new_component_unit = ref("")

  const new_component_entry_company = ref("")
  const new_component_entry_quantity = ref("")
  const new_component_entry_price = ref("")
  const new_component_entries = ref([])

  const new_entry_company = ref("")
  const new_entry_quantity = ref("")
  const new_entry_price = ref("")
  const new_entry_expiration_date = ref("")


  const material_logs_id = ref("")


const updatEntriesTableRowsPerPage = (event: any) => {
    loadEntries(event.first,event.rows)
}

const loadEntries = (first=0,rows=entriesTableRowsPerPage.value) => {
    isEntriesTableLoading.value = true;

    const page_number = Math.ceil((first/rows))
    
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/materials/${entries_dialog_material.value.id}/entries`, {
        params: {
            "page[number]": page_number,
            "page[size]": rows
        },
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {
        entries_dialog_material.value.entries = response.data.data;
        entriesTableTotalRecords.value = response.data.meta.total_records;
    })
    .catch(() => {
        toast.add({severity:'error', summary: 'Error', detail: 'Failed to load products'});
    })
    .finally(() => {
        isEntriesTableLoading.value = false;
    });
}


const confirmDeleteMaterial = (material_id: string) => {
    confirm.require({
        message: 'Are you sure you want to delete this material ?',
        header: 'Confirmation',
        icon: 'pi pi-exclamation-triangle',
        rejectProps: {
            label: 'Cancel',
            severity: 'secondary',
            outlined: true
        },
        acceptProps: {
            label: 'Yes',
            severity: 'danger'
        },
        accept: () => {
            deleteMaterial(material_id)
        },
        reject: () => {
            
        }
    });
};

  const deleteMaterial = (material_id: string) => {
    axios.delete(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/materials/${material_id}`, {
    headers: {
        Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
    }
    }).then(() => {
        toast.add({ severity: 'success', summary: 'Success', detail: "Material deleted successfully", life: 3000 });
        loadInventory()
    }).catch(() => {
        toast.add({ severity: 'error', summary: 'Error', detail: "Failed to delete material", life: 3000 });
    });
  }

  const saveEditedMaterial = (material: Material) => {

    axios.patch(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/materials/${edit_material.value?.id}`, {
        data: {
            ...material,
            ...material_settings.value
        }
    },{
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    }).then(() => {
      edit_material_dialog.value = false
      toast.add({severity:'success', summary: 'Success', detail: "successfully edit material data", life: 3000});
      loadInventory()
    }).catch((error) => {
      toast.add({severity:'error', summary: 'Error', detail: error.message, life: 3000});
    })
  }

  const updatLogsTableRowsPerPage = (event) => {
    const { first, rows } = event;
    loadComponentLogs(material_logs_id.value,first,rows)
}


  const saveMaterialSettings = () => {
    axios.patch(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/materials/${material_settings.value?.id}`, {
        data: material_settings.value
    },{
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    }).then(() => {
      toast.add({severity:'success', summary: 'Success', detail: 'Material settings saved', life: 3000,group:'br'});
      material_settings_dialog.value = false
    }).catch(error => {
      toast.add({severity:'error', summary: 'Error', detail: error.message, life: 3000});
    })
  }


  const confirmDeleteEntry = (event,material_id, entry_id) => {
    confirm.require({
        target: event.currentTarget,
        message: 'Are you sure you want to delete this entry ?',
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




            axios.delete(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/materials/${material_id}/entries/${entry_id}`,{
                headers: {
                    Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
                }
            })
            .then(() => {
                    toast.add({ severity: 'success', summary: 'Done', detail: "Entry deleted !",life: 3000,group:'br' });
                    inventory_components.value.forEach((component) => {
                        if (component.id == expanded_component_id.value){
                            component.entries.splice(component.entries.findIndex(el => el.id == data), 1)
                        }
                    })
            })
        },
        reject: () => {
            
        }
    });
  }


  const addNewEntry = (component_id) => {

    var newEntry = {
                "quantity": parseFloat(new_entry_quantity.value),
                "purchase_price": parseFloat(new_entry_price.value),
                "company": new_entry_company.value,
                "expiration_date": new_entry_expiration_date.value
            }

    axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/materials/${component_id}/entries`, {
        data: [
           newEntry
        ]
      },{
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
      })
      .then(() => {
        toast.add({ severity: 'success', summary: 'Success', detail: 'Entry saved successfully !', life: 3000,group:'br' });

        inventory_components.value.forEach((component) => {
            if (component.id == component_id)
            {
                component.entries.push(newEntry)
            }
        })


        new_entry_company.value = ""
        new_entry_quantity.value = ""
        new_entry_price.value = ""


      });
  }



  const submitNewComponent = () => {

      var entries : any = []
      new_component_entries.value.forEach((entry) => {
          entries.push({
            company: entry.company,
            quantity: parseInt(entry.quantity),
            purchase_price: parseFloat(entry.purchase_price)
          })
      })

      axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/materials`, {
        data : {
            name: new_component_name.value,
            unit: new_component_unit.value,
            entries: entries,
        }
      },{
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
      })
      .then(() => {
        toast.add({ severity: 'success', summary: 'Success', detail: 'Component saved successfully !', life: 3000,group:'br' });
        add_component_dialog.value = false
        loadInventory()
        new_component_entries.value = []
        new_component_name.value = ""
        new_component_unit.value = ""
        new_component_entry_company.value = ""
        new_component_entry_quantity.value = ""
        new_component_entry_price.value = ""
      });
  }


  const onComponentLogRowExpand = (event) => {
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders/`+event.data.order_id,{
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then((result)=>{

        component_logs.value.forEach((log,index) => {
            if (log.order_id == result.data.data.id){
                component_logs.value[index].order = result.data.data
            }
        })

    })
  };


  const loadComponentLogs = (component_id:string,first=0,rows=50) => {

    material_logs_id.value = component_id
    isLogsTableLoading.value = true

    let page_number = Math.floor((first/rows)) + 1

    


    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/materials/`+component_id+`/logs?page[number]=${page_number}&page[size]=${rows}`,{
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then((result)=>{
        component_logs.value = result.data.data
        component_logs_dialog.value = true
        component_logs_name.value = component_id
        logsTableTotalRecords.value = result.data.meta.total_records
        isLogsTableLoading.value = false
    })
  }


  const loadInventory = () => {
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/materials`,{
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then((result)=>{

        result.data.data.forEach(component => {
            component.totalAmount = component.quantity
        });
        
        inventory_components.value = result.data.data
    })
  }


  loadInventory();
  
</script>