<template>
    <div>
        <h1>Inventory</h1>
        <div class="grid">
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
                <Panel style="width:100%;">
                    <InputText v-model="value1" placeholder="Search" class="mb-4" />
                    <div class="flex">
                        <InventoryItem class="w-3" />
                        <InventoryItem class="w-3" />
                        <InventoryItem class="w-3" />
                    </div>
                </Panel>
            </div>
        </div>
    </div>
</template>

<script setup>
import Listbox from 'primevue/listbox';
import Panel from 'primevue/panel';
import InventoryItem from '@/components/InventoryItem.vue'
import InputText from 'primevue/inputtext'

  
  import { ref } from "vue";
  
  const selectedCategory = ref();
  const categories = ref([
      { name: 'All', icon:'star-of-life' },
      { name: 'Cheese', icon:'cheese' },
      { name: 'Chicken', icon: 'drumstick-bite' },
      { name: 'Meat', icon: 'burger' },
      { name: 'Pastery', icon: 'bread-slice' },
      { name: 'Sauces', icon: 'bowl-food' },
      { name: 'Vegetables', icon: 'apple-whole' },
      { name: 'Liquids', icon: 'water' },
  ]);
  
</script>