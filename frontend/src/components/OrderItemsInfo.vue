<template>
    <div>
        <DataTable  v-model:expandedRows="expandedRows" :value="props.order.items" stripedRows tableStyle="width: 50rem" class="w-full pr-5">
            <Column expander style="width: 5rem" />
            <Column field="product.name" :header="$t('name')"></Column>
            <Column field="price" :header="$t('price')"></Column>
            <Column field="quantity" :header="$t('quantity')"></Column>
            <Column field="status" :header="$t('status')">
                <template #body="slotProps">
                    <Badge :value="order_items_status[slotProps.data.id].title" :severity="order_items_status[slotProps.data.id].severity" />
                </template>
            </Column>
            <Column :header="$t('actions')">
                <template #body="slotProps">
                    <ButtonGroup>
                        <Button icon="pi pi-undo" :label="$t('return')" severity="warn" aria-label="refund" @click="itemToRefund=slotProps.data;item_refund_dialog=true;" v-if="order_items_status[slotProps.data.id].title != 'REFUNDED'">
                            <span class="pi pi-undo"></span>
                            <div>return</div>
                        </Button>
                    </ButtonGroup>
                </template>
            </Column>
            <template #expansion="slotProps">
                <h4>{{ $t('material',3) }}</h4>
                <DataTable :value="slotProps.data.materials">
                    <Column field="material.name" :header="$t('name')"></Column>
                    <Column field="quantity" :header="$t('quantity')"></Column>
                    <Column field="material.unit" :header="$t('unit')"></Column>
                </DataTable>
            </template>
        </DataTable>
        <Dialog v-model:visible="item_refund_dialog" modal :header="`Returning ${itemToRefund?.product.name} | order #${props.order.display_id}`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
            <OrderItemRefund :item="itemToRefund" :order="props.order" @finished="emit('updated')" />
        </Dialog>
    </div>
</template>

<script setup lang="ts">
import Dialog from 'primevue/dialog';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import { ref,defineProps,computed,defineEmits } from 'vue';
import { OrderItem } from '@/classes/OrderItem';
import Button from 'primevue/button';
import ButtonGroup from 'primevue/buttongroup';
import OrderItemRefund from './OrderItemRefund.vue';
import Order from '@/classes/Order';
import {Badge} from 'primevue';

const expandedRows = ref([]);
const itemToRefund = ref<OrderItem>()
const item_refund_dialog = ref(false)

const emit = defineEmits(['updated'])

const props = defineProps({
    order: {
        type: Order,
        required: true
    },
})

const order_items_status : any = computed(() => {


let statuses = {}


props.order.items.forEach((item) => {

    let item_status = {}

    
    if (item.status == ""){
        item_status={
            title:"NORMAL",
            severity:"secondary"
        }
    }

    if (item.status == "pending" ){
        item_status={
            title:"PENDING",
            severity:"secondary"
        }
    }

    if (item.status == "refunded" ){
        item_status={
            title:"REFUNDED",
            severity:"danger"
        }
    }

    statuses[item.id] = item_status
})


return statuses

})


</script>