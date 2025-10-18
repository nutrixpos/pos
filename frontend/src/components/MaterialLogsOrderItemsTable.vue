<template>
    <DataTable :value="props.items">
        <Column header="" class="w-1rem">
            <template #body="{index}">
                <Tag v-if="props.order_item_index == index" icon="pi pi-arrow-right" style="background-color:#FFDC00" />
            </template>
        </Column>
        <Column field="product.name" header="Name"></Column>
        <Column header="Materials">
            <template #body="slotProps">
                <ul>
                    <li v-for="(material,index) in slotProps.data.materials" :key="index">
                        <div>
                            {{ material.material.name }}: {{ material.quantity }} {{ material.material.unit }}
                        </div>
                    </li>
                </ul>
            </template>
        </Column>
    </DataTable>
</template>

<script setup lang="ts">
import {defineProps} from 'vue'
import { DataTable,Column,Tag } from 'primevue';

const props = defineProps({
    items: {
        type: Array,
        required: true
    },
    order_item_index: {
        type: Number,
        required: false
    }
})

</script>