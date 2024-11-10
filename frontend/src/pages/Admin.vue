<template>
    <div>
        <div class="grid m-2">
            <div class="col-12">
                <Menubar :model="items">
                    <template #item="{ item, props }">
                        <a class="flex align-items-center" v-bind="props.action" :href="item.link">
                            <span :class="item.icon" />
                            <span class="ml-2">{{ item.label }}</span>
                        </a>
                    </template>
                </Menubar>
            </div>
            <div class="col-2">
                <Listbox  v-model="selected_list_item" :options="list_items" optionGroupLabel="name" optionGroupChildren="items" class="w-full mt-2">
                    <template #option="slotProps">
                        <RouterLink style="text-decoration: none;color: inherit" :to="slotProps.option.link" class="flex align-items-center">
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
import Menubar from 'primevue/menubar'
import Listbox from 'primevue/listbox';

const selected_list_item = ref ({ name: 'Inventory', icon:'inbox', link:'inventory' })

const list_items = ref([
    { 
        name: 'Admin', 
        items: [
            { name: 'Inventory', icon:'inbox', link:'/admin/inventory' },
            { name: 'Products', icon:'barcode', link:'/admin/products' },
            { name: 'Categories', icon:'barcode', link:'/admin/categories' },
            { name: 'Orders', icon:'barcode', link:'/admin/orders' },
        ]
        
    },
    { 
        name: 'Observe',
        items: [
            { name: 'Sales', icon:'dollar-sign', link: '/admin/sales'}
        ]
        
    },
]);

const items = ref([
      {
          label: 'Cashier',
          icon: 'pi pi-home',
          link: '/home',
      },
      {
          label: 'Kitchen',
          icon: 'pi pi-star',
          link:'/kitchen'
      },
      {
          label: 'Admin',
          icon: 'pi pi-cog',
          link: '/admin',
      }
  ]);
</script>