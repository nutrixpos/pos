<template>
    <div>
        <div class="grid m-2">
            <div class="col-2 mt-6">
                <Listbox  v-model="selected_list_item" :options="list_items" optionGroupLabel="name" optionGroupChildren="items" class="w-full mt-2">
                    <template #option="slotProps">
                        <RouterLink style="text-decoration: none;color: inherit" :to="slotProps.option.link" class="flex align-items-center">
                            <fa :icon="slotProps.option.icon" class="mr-2" />
                            <div>{{ slotProps.option.name }}</div>
                        </RouterLink>
                    </template>
                </Listbox>
            </div>
            <div class="col-10 flex pt-3">
                <RouterView />
            </div>
        </div>
    </div>
</template>

<script setup>
import {ref} from "vue";
import Listbox from 'primevue/listbox';

const selected_list_item = ref ({ name: 'Inventory', icon:'inbox', link:'inventory' })

const list_items = ref([
    { 
        name: 'Admin', 
        items: [
            { name: 'Inventory', icon:'inbox', link:'/admin/inventory' }
        ]
        
    },
    { 
        name: 'Observe',
        items: [
            { name: 'Sales', icon:'dollar-sign', link: '/admin/sales'}
        ]
        
    },
]);
</script>