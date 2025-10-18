<template>
    <DataTable v-model:expandedRows="expandedSalesLogOrderItemComponents" :value="items">
        <Column expander style="width: 5rem" />
        <Column sortable field="ItemName" :header="$t('name')">
            <template #body="slotProps">
                <Badge value="REFUNDED" v-if="items_refunds[slotProps.data.item_id]" severity="danger"  />
                {{ slotProps.data.item_name }}
            </template>
        </Column>
        <Column sortable field="Cost" :header="$t('cost')">
            <template #body="slotProps">
                <div class="flex gap-2 align-items-center">
                    <div>{{ slotProps.data.cost }} </div>
                    <Badge v-if="items_refunds[slotProps.data.item_id]?.inventory_refunds > 0" :value="`-${items_refunds[slotProps.data.item_id]?.inventory_refunds}`" severity="success" class="mr-2" />
                </div>
            </template>
        </Column>
        <Column sortable field="Quantity" :header="$t('quantity')">
            <template #body="slotProps">
                <div>{{ slotProps.data.quantity != "0" ? slotProps.data.quantity : "1" }}</div>
            </template></Column>
            <Column sortable field="sale_price" :header="$t('sale_price')"></Column>
            <Column sortable field="refunds" :header="$t('refunds')">
                <template #body="slotProps">
                    <Badge :value="`${items_refunds[slotProps.data.item_id]?.refund_amount}`" severity="danger" v-if="items_refunds[slotProps.data.item_id]?.is_refunded" />
                </template>
            </Column>
        <Column sortable field="profit" :header="$t('profit')">
            <template #body="slotProps">
                <div :style="`${ (slotProps.data.sale_price - slotProps.data.cost - (items_refunds[slotProps.data.item_id]?.refund_amount || 0 ) + (items_refunds[slotProps.data.item_id]?.inventory_refunds || 0 )) > 0 ? 'color:green' : 'color:red' }`">{{ slotProps.data.sale_price - slotProps.data.cost - (items_refunds[slotProps.data.item_id]?.refund_amount || 0 ) + (items_refunds[slotProps.data.item_id]?.inventory_refunds || 0 ) }}</div>
            </template>
        </Column>
        <template #expansion="slotProps">
            <DataTable :value="slotProps.data.components">
                <Column sortable field="component_name" :header="$t('name')"></Column>
                <Column sortable field="cost" :header="$t('cost')">
                    <template #body="slotProps">
                        <div class="flex gap-2 align-items-center">
                            {{ slotProps.data.cost }}
                            <Badge :value="`-${items_refunds[slotProps.data.item_id]?.material_refunds[slotProps.data.component_id]?.inventory_return_qty * items_refunds[slotProps.data.item_id]?.material_refunds[slotProps.data.component_id]?.cost_per_unit}`" severity="success" v-if="items_refunds[slotProps.data.item_id]?.material_refunds[slotProps.data.component_id]?.inventory_return_qty > 0" />
                        </div>
                    </template>
                </Column>
                <Column sortable field="quantity" :header="$t('quantity')">
                    <template #body="slotProps">
                        <div class="flex gap-2 align-items-center">
                            {{ slotProps.data.quantity }}
                            <Badge :value="`-${items_refunds[slotProps.data.item_id]?.material_refunds[slotProps.data.component_id]?.inventory_return_qty}`" severity="success" v-if="items_refunds[slotProps.data.item_id]?.material_refunds[slotProps.data.component_id]?.inventory_return_qty > 0" />
                        </div>
                    </template>
                </Column>
                <Column sortable field="entry_id" :header="$t('entry')"></Column>
            </DataTable>
            <SalesLogTableItems v-if="slotProps.data.DownstreamCost != null" :items="slotProps.data.DownstreamCost" />    
        </template>
    </DataTable>
</template>

<script setup lang="ts">

import {defineProps,ref} from 'vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import {Badge} from 'primevue';

const props = defineProps(['items','order_refunds'])
const expandedSalesLogOrderItemComponents = ref([])

const items_refunds = ref({})
const items = ref([])

const init = () => {

    items.value = props.items

    items.value.forEach((item,itemIndex) => {
        item.components.forEach((component,componentIndex) => {
            items.value[itemIndex].components[componentIndex].item_id = item.item_id
        })
    })

    if (props.order_refunds.refunds?.length > 0){
        props.items.forEach((item) => {
    
            props.order_refunds.refunds.forEach((order_refund) => {
    
                if (order_refund.order_item_id == item.item_id){
    
    
                    items_refunds.value[item.item_id] = {
                        is_refunded: true,
                        refund_amount: order_refund.amount,
                        refund_reason: order_refund.reason,
                        inventory_refunds: 0,
                        item_cost: order_refund.item_cost,
                        destination: order_refund.destination,
                        material_refunds: {},
                    }
                    
    
                    if (order_refund.destination == "inventory")
                        items_refunds.value[item.item_id].inventory_refunds = order_refund.item_cost

                    if (order_refund.destination == "custom"){
                        order_refund.material_refunds.forEach((material_refund) => {
                            items_refunds.value[item.item_id].inventory_refunds += material_refund.inventory_return_qty * material_refund.cost_per_unit
                            items_refunds.value[item.item_id].material_refunds[material_refund.material_id] = material_refund
                        })
                    }
    
                    
    
                }
            })
        })
    }

}

init()

</script>