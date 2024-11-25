<template>
    <div>
        <DataTable :value="props.items" stripedRows tableStyle="width: 50rem" class="w-full pr-5">
            <Column field="product.name" header="Name"></Column>
            <Column field="price" header="Price"></Column>
            <Column field="quantity" header="Quantity"></Column>
            <Column header="Actions">
                <template #body="slotProps">
                    <ButtonGroup>
                        <Button icon="pi pi-book" label="Details" severity="secondary" aria-label="Details" @click="itemToDisplay=slotProps.data; item_details_dialog=true" />
                    </ButtonGroup>
                </template>
            </Column>
        </DataTable>
        <Dialog v-model:visible="item_details_dialog" modal header="Item Details" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
            <pre>
                {{ itemToDisplay  }}
            </pre>
        </Dialog>
    </div>
</template>

<script setup lang="ts">
import Dialog from 'primevue/dialog';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import { ref,defineProps } from 'vue';
import { OrderItem } from '@/classes/OrderItem';
import Button from 'primevue/button';
import ButtonGroup from 'primevue/buttongroup';


const itemToDisplay = ref<OrderItem>()
const item_details_dialog = ref(false)

const props = defineProps({
    items: {
        type: Array<OrderItem>,
        required: true
    }
})


</script>