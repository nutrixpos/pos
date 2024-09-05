<template>
    <DataTable v-model:expandedRows="expandedSalesLogOrderItemComponents" :value="props.items">
        <Column expander style="width: 5rem" />
        <Column sortable field="ItemName" header="Name"></Column>
        <Column sortable field="Cost" header="Cost"></Column>
        <Column sortable field="SalePrice" header="Sale"></Column>
        <Column sortable field="profit" header="Profit">
            <template #body="slotProps">
                <div :style="`${ (slotProps.data.SalePrice - slotProps.data.Cost) > 0 ? 'color:green' : 'color:red' }`">{{ slotProps.data.SalePrice - slotProps.data.Cost }}</div>
            </template>
        </Column>
        <template #expansion="slotProps">
            <DataTable :value="slotProps.data.Components">
                <Column sortable field="ComponentName" header="Component Name"></Column>
                <Column sortable field="Cost" header="Cost"></Column>
                <Column sortable field="Quantity" header="Quantity"></Column>
                <Column sortable field="EntryId" header="Entry"></Column>
            </DataTable>
            <SalesLogTableItems :items="slotProps.data.DownstreamCost" />    
        </template>
    </DataTable>
</template>

<script setup lang="ts">

import {defineProps,ref} from 'vue'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'

const props = defineProps(['items'])
const expandedSalesLogOrderItemComponents = ref([])
</script>