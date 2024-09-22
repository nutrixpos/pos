<template>
    <div class="flex justify-content-between align-items-center">
        <h4>{{ model.product.name }}</h4>
        <InputText :disabled="!model.is_consume_from_ready" type="number" v-model.number="model.quantity"  size="small"/>
        <div class="flex align-items-center justify-content-center">
            <span class="mx-2">From Ready</span>
            <InputSwitch v-model="model.is_consume_from_ready" :disabled="!model.can_change_ready_toggle" />
            <span class="mx-2">
                <p style="font-size: 0.9rem;">{{model.ready}} Ready</p>
            </span>
        </div>
    </div>
    <div v-if="!model.is_consume_from_ready">
        <div class="flex my-3 py-2 justify-content-between" style="border-bottom:1px solid gray" v-for="(material,index) in model.materials" :key="index">
            {{ material.material.name }}
            <div class="flex">
                <InputText type="number" v-model.number="model.materials[index].quantity" size="small"/>
                <span class="ml-2 mt-2">{{ material.material.unit }}</span>
            </div>
            <Dropdown v-if="model.product.materials[index].entries != null && model.product.materials[index].entries.length > 0" v-model="model.materials[index].entry"  :options="model.product.materials[index].entries" optionLabel="label" placeholder="Select option" class="w-6" />
        </div>
    </div>
    <div v-if="model.sub_items != null">
        <div v-for="(subitem,index) in model.sub_items" :key="index" class="m-0">
            <OrderItemView v-model="model.sub_items[index]" />
        </div>
    </div>
</template>

<script setup lang="ts">
import {defineModel} from 'vue'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import InputSwitch from 'primevue/inputswitch';
import { OrderItem } from '@/classes/OrderItem'

const model = defineModel<OrderItem>({
    required: true})


// const init = () => {

//     // model.value.Selections = []

//     if (model.value.Selections.length == 0){

//         model.value.Components.forEach((component) => {
//             const selection = new ComponentSelection()
//                 selection.ComponentId = component.component_id
//                 selection.Quantity = component.defaultquantity
//                 selection.Name = component.name
//                 selection.Unit = component.unit
//                 selection.Entry = component.entries[0]

//                 model.value.Selections.push(selection)
//             })
//     }
// }


// init()

</script>