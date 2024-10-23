<template>
    <div class="w-full">
        <div class="grid m-2">
            <div class="col-12 flex pt-3">
                <div class="gird w-full">
                    <div class="col-12">
                        <h3>Inventory</h3>
                    </div>
                    <div class="col-12 flex justify-content-center align-items-center w-full">
                        <DataTable v-model:expandedRows="expandedRows"  @rowExpand="componentRowExpand" :value="inventory_components" stripedRows tableStyle="min-width: 50rem" class="w-full pr-5">
                            <template #header>
                                <div class="flex flex-wrap items-center justify-between align-items-center gap-2">
                                    <Button icon="pi pi-plus" label="Add Component" @click="add_component_dialog = true" rounded raised />
                                </div>
                            </template>
                            <Column expander style="width: 5rem" />
                            <Column field="name" header="Name"></Column>
                            <Column field="totalAmount" header="Quantity"></Column>
                            <Column field="unit" header="Unit"></Column>
                            <Column header="Status">
                                <template #body="slotProps">
                                    <Tag :value="slotProps.data.totalAmount > slotProps.data.settings.stock_alert_treshold ? 'INSTOCK' : 'LOWSTOCK'" :severity="slotProps.data.totalAmount > slotProps.data.settings.stock_alert_treshold ? 'success' : 'warning'" />
                                </template>
                            </Column>
                            <Column header="Actions" style="width:30rem">
                                <template #body="slotProps">
                                    <ButtonGroup>
                                        <Button icon="pi pi-clock" label="History" @click="loadComponentLogs(slotProps.data._id)" severity="secondary" aria-label="Save"  />
                                        <Button icon="pi pi-cog" class="ml-2" severity="secondary" aria-label="Edit" @click="material_settings = slotProps.data; material_settings_dialog=true"  />
                                    </ButtonGroup>
                                </template>
                            </Column>
                            <template #expansion="slotProps">
                                <div class="p-4">
                                    <h4>Entries for {{ slotProps.data.name }}</h4>
                                    <div class="flex justify-content-center align-items-center">
                                        <InputText class="m-1" placeholder="Company" v-model="new_entry_company" aria-describedby="name-help" />
                                        <InputText class="m-1" placeholder="Purchase Quantity" v-model="new_entry_quantity" aria-describedby="name-help" />
                                        <InputText class="m-1" placeholder="Purchase Price" v-model="new_entry_price" aria-describedby="name-help" />
                                        <Button icon="pi pi-plus" label="Add Entry" @click="addNewEntry(slotProps.data.id)" severity="info" raised />
                                    </div>
                                    <DataTable :value="slotProps.data.entries">
                                        <Column field="company" header="Company"></Column>
                                        <Column field="quantity" header="Quantity" sortable></Column>
                                        <Column field="purchase_quantity" header="Purchase Quantity" sortable></Column>
                                        <Column field="purchase_price" header="Purchase Price" sortable></Column>
                                        <Column header="Actions" style="width:30rem">
                                            <template #body="slotProps">
                                                <ButtonGroup>
                                                    <Button icon="pi pi-times" label="Delete" severity="secondary" aria-label="Delete" @click="confirmDeleteEntry($event,slotProps.data.id)" />
                                                </ButtonGroup>
                                            </template>
                                        </Column>
                                    </DataTable>
                                </div>
                            </template>
                        </DataTable>
                    </div>
                </div>
            </div>
            <Dialog v-model:visible="material_settings_dialog" modal :header="`Settings for  ${material_settings.name}`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
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
            <Dialog v-model:visible="add_component_dialog" modal :header="`Add new inventory component`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
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
                    <h4>Entries</h4>
                    <div class="flex justify-content-center align-items-center">
                        <InputText class="m-1" placeholder="Company" v-model="new_component_entry_company" aria-describedby="name-help" />
                        <InputText class="m-1" placeholder="Quantity" v-model="new_component_entry_quantity" aria-describedby="name-help" />
                        <InputText class="m-1" placeholder="Price" v-model="new_component_entry_price" aria-describedby="name-help" />
                        <Button label="Add" @click="new_component_entries.push({company: new_component_entry_company, quantity: new_component_entry_quantity, unit: new_component_unit, purchase_price: new_component_entry_price})" />
                    </div>
                    <DataTable :value="new_component_entries">
                        <Column field="company" header="Company"></Column>
                        <Column field="quantity" header="Quantity"></Column>
                        <Column field="unit" header="Unit">
                            <template #body="slotProps">
                                {{ slotProps.data.unit }}
                            </template>
                        </Column>
                        <Column field="price" header="Price"></Column>
                    </DataTable>

                    <div class="flex w-full mt-5 justify-content-center align-items-center">
                        <Button label="Submit" class="lg:w-6" @click="submitNewComponent" />
                    </div>
               </div>
            </Dialog>
            <Dialog v-model:visible="component_logs_dialog" modal :header="`Consumption for  ${component_logs_name}`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
                <DataTable @rowExpand="onComponentLogRowExpand"  v-model:expandedRows="expandedComponentLogsRows" :value="component_logs" stripedRows tableStyle="min-width: 50rem" class="w-full pr-5">
                    <Column expander style="width: 5rem" />
                    <Column field="date" header="Date"></Column>
                    <Column field="quantity" header="Quantity"></Column>
                    <Column  header="Order Item">
                        <template #body="slotProps">
                            [{{ slotProps.data.item_order_index }}] {{ slotProps.data.item_name }}
                        </template>
                    </Column>
                    <Column field="order_id" header="Order Id"></Column>
                    <template #expansion="slotProps">
                        <div class="p-4">
                            <h4>Order Items</h4>
                            <DataTable :value="slotProps.data.order.items" v-if="slotProps.data.order">
                                <Column field="product.name" header="Name"></Column>
                                <Column header="Ingredients">
                                    <template #body="slotProps">
                                        <ul>
                                            <li v-for="(material,index) in slotProps.data.materials" :key="index">
                                                {{ material.material.name }}: {{ material.quantity }} {{ material.material.unit }}
                                            </li>
                                        </ul>
                                    </template>
                                </Column>
                            </DataTable>
                            <div v-else>
                                Loading ...
                            </div>
                        </div>
                    </template>
                </DataTable>
            </Dialog>
        </div>
        <ConfirmPopup></ConfirmPopup>
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
import ConfirmPopup from 'primevue/confirmpopup';
import Tag from 'primevue/tag'
import { Material } from '@/classes/OrderItem';
// import Message from 'primevue/message'
  
import { ref } from "vue";
import { useToast } from "primevue/usetoast";



const confirm = useConfirm();



  const toast = useToast();

  const material_settings_dialog = ref(false)
  const material_settings = ref<Material>({})
  
  const expandedRows = ref([]);
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

  const expanded_component_id = ref("")


  const saveMaterialSettings = () => {
    axios.post(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/editmaterial`, {
        material: material_settings.value
    }).then(() => {
      toast.add({severity:'success', summary: 'Success', detail: 'Material settings saved', life: 3000,group:'br'});
      material_settings_dialog.value = false
    }).catch(error => {
      toast.add({severity:'error', summary: 'Error', detail: error.message, life: 3000});
    })
  }

  const componentRowExpand = (event) => {
    expanded_component_id.value = event.data.id
  }


  const confirmDeleteEntry = (event,entry_id) => {
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




            axios.get(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/entry?entry_id=${entry_id}&component_id=`+expanded_component_id.value)
            .then(() => {
                    toast.add({ severity: 'success', summary: 'Done', detail: "Entry deleted !",life: 3000,group:'br' });
                    inventory_components.value.forEach((component) => {
                        if (component._id == expanded_component_id.value){
                            component.entries.splice(component.entries.findIndex(el => el._id == entry_id), 1)
                        }
                    })
            })
        },
        reject: () => {
            toast.add({ severity: 'error', summary: 'Transaction failed !', detail: event.data.name, life: 3000,group:'br' });
        }
    });
  }


  const addNewEntry = (component_id) => {

    var newEntry = {
                "quantity": parseFloat(new_entry_quantity.value),
                "purchase_price": parseFloat(new_entry_price.value),
                "company": new_entry_company.value
            }

    axios.post(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/materialentry`, {
        component_id,
        entries: [
           newEntry
        ]
      })
      .then(() => {
        toast.add({ severity: 'success', summary: 'Success', detail: 'Entry saved successfully !', life: 3000,group:'br' });

        inventory_components.value.forEach((component) => {
            if (component._id == component_id)
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

      var entries = []
      new_component_entries.value.forEach((entry) => {
          entries.push({
            company: entry.company,
            quantity: parseInt(entry.quantity),
            purchase_price: parseFloat(entry.purchase_price)
          })
      })

      axios.post(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/material`, {
        name: new_component_name.value,
        unit: new_component_unit.value,
        entries: entries,
      })
      .then(() => {
        toast.add({ severity: 'success', summary: 'Success', detail: 'Component saved successfully !', life: 3000,group:'br' });
        add_component_dialog.value = false
      });
  }


  const onComponentLogRowExpand = (event) => {
    axios.get(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/order?id=`+event.data.order_id)
    .then((result)=>{

        component_logs.value.forEach((log) => {
            if (log._id == event.data._id){
                log.order = result.data
                log.order.items.forEach((_,index) => {
                    log.order.items[index].ingredients = log.order.ingredients[index]
                })
            }
        })

    })
  };


  const loadComponentLogs = (component_id) => {
    axios.get(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/materiallogs?id=`+component_id)
    .then((result)=>{
        component_logs.value = result.data
        component_logs_dialog.value = true
        component_logs_name.value = component_id
    })
  }


  const loadInventory = () => {
    axios.get(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/materials`)
    .then((result)=>{

        result.data.forEach(component => {
            var totalAmount;

            component.entries?.forEach(entry => {
                if (entry.quantity > 0)
                    totalAmount = totalAmount ? totalAmount + entry.quantity : entry.quantity
            });

            component.totalAmount = totalAmount
        });
        
        inventory_components.value = result.data
    })
  }


  loadInventory();
  
</script>