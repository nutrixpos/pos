<template>
    <div>
        <h2 class="m-2">Analytics</h2>
        <div class="grid m-2">
            <div class="col-2">
                <Listbox  v-model="selectedCategory" :options="categories" optionLabel="category" class="w-full mt-2">
                    <template #option="slotProps">
                        <div class="flex align-items-center">
                            <fa :icon="slotProps.option.icon" class="mr-2" />
                            <div>{{ slotProps.option.name }}</div>
                        </div>
                    </template>
                </Listbox>
            </div>
            <div class="col-10 flex pt-3">
                <div class="gird w-full" v-if="selectedCategory.name == 'Inventory'">
                    <div class="col-12">
                        <h3>Inventory</h3>
                    </div>
                    <div class="col-12 flex justify-content-center align-items-center w-full">
                        <DataTable v-model:expandedRows="expandedRows" :value="inventory_components" stripedRows tableStyle="min-width: 50rem" class="w-full pr-5">
                            <Column expander style="width: 5rem" />
                            <Column field="name" header="Name"></Column>
                            <Column field="totalAmount" header="Quantity"></Column>
                            <Column field="unit" header="Unit"></Column>
                            <Column header="Actions" style="width:30rem">
                                <template #body="slotProps">
                                    <ButtonGroup>
                                        <Button icon="pi pi-clock" label="History" @click="loadComponentLogs(slotProps.data.name)" severity="secondary" aria-label="Save"  />
                                        <Button icon="pi pi-plus" label="Add" aria-label="Add" severity="info" @click="console.log(slotProps)" />
                                    </ButtonGroup>
                                </template>
                            </Column>
                            <template #expansion="slotProps">
                                <div class="p-4">
                                    <h4>Entries for {{ slotProps.data.name }}</h4>
                                    <DataTable :value="slotProps.data.entries">
                                        <Column field="company" header="Company"></Column>
                                        <Column field="quantity" header="Quantity" sortable></Column>
                                        <Column field="unit" header="Unit" sortable></Column>
                                        <Column header="Actions" style="width:30rem">
                                            <template #body="slotProps">
                                                <ButtonGroup>
                                                    <Button icon="pi pi-pencil" label="Edit" severity="secondary" aria-label="Edit" @click="console.log(slotProps)" />
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
            <Dialog v-model:visible="component_logs_dialog" modal :header="`Consumption for  ${component_logs_name}`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
                <DataTable @rowExpand="onComponentLogRowExpand"  v-model:expandedRows="expandedComponentLogsRows" :value="component_logs" stripedRows tableStyle="min-width: 50rem" class="w-full pr-5">
                    <Column expander style="width: 5rem" />
                    <Column field="date" header="Date"></Column>
                    <Column field="component_name" header="Unit"></Column>
                    <Column field="quantity" header="Quantity"></Column>
                    <Column field="company" header="Company"></Column>
                    <Column field="order_id" header="Order Id"></Column>
                    <template #expansion="slotProps">
                        <div class="p-4">
                            <h4>Order Items</h4>
                            <DataTable :value="slotProps.data.order.items" v-if="slotProps.data.order">
                                <Column field="name" header="Name"></Column>
                                <Column header="Ingredients">
                                    <template #body="slotProps">
                                        <ul>
                                            <li v-for="(ingredient,index) in slotProps.data.ingredients" :key="index">
                                                {{ ingredient.name }}: {{ ingredient.quantity }}
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
    </div>
</template>

<script setup>
import Listbox from 'primevue/listbox';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import axios from 'axios'
import Button from 'primevue/button'
import ButtonGroup from 'primevue/buttongroup'
import Dialog from 'primevue/dialog'

  
  import { ref } from "vue";
  
  const selectedCategory = ref({ name: 'Inventory', icon:'inbox' });
  const categories = ref([
      { name: 'Inventory', icon:'inbox' },
  ]);

  const expandedRows = ref([]);
  const expandedComponentLogsRows = ref([])
  
  const inventory_components = ref([])

  const component_logs = ref([])
  const component_logs_dialog = ref(false)
  const component_logs_name = ref("")



  const onComponentLogRowExpand = (event) => {
    axios.get('http://localhost:8000/api/order?id='+event.data.order_id)
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


  const loadComponentLogs = (component_name) => {
    axios.get('http://localhost:8000/api/componentlogs?name='+component_name)
    .then((result)=>{
        component_logs.value = result.data
        component_logs_dialog.value = true
        component_logs_name.value = component_name
    })
  }


  const loadInventory = () => {
    axios.get('http://localhost:8000/api/components')
    .then((result)=>{

        result.data.forEach(component => {
            var totalAmount;

            component.entries?.forEach(entry => {
                totalAmount = totalAmount ? totalAmount + entry.quantity : entry.quantity
            });

            component.totalAmount = totalAmount
        });
        
        inventory_components.value = result.data
    })
  }


  loadInventory();
  
</script>