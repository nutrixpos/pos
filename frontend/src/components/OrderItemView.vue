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
        <Button label="Add Material" @click="new_component_dialog = true" />
        <div class="flex my-3 py-2 justify-content-between align-items-center" style="border-bottom:1px solid gray" v-for="(material,index) in model.materials" :key="index">

            <Button icon="pi pi-times" size="small" style="width:2rem;height: 2rem;" aria-label="Remove" severity="secondary" @click="model.RemoveMaterialByIndex(index)" />

            {{ material.material.name }}
            <div class="flex">
                <InputText type="number"  @change="model.UpdateMaterialEntryCost(index)" v-model.number="model.materials[index].quantity" size="small"/>
                <span class="ml-2 mt-2">{{ material.material.unit }}</span>
            </div>
            <Dropdown @change="model.UpdateMaterialEntryCost(index)"  v-model="model.materials[index].entry"  :options="model.materials[index].entries" optionLabel="label" placeholder="Select option" class="w-6" />
            Cost: {{ material.entry.cost }}
        </div>
    </div>
    <div v-if="model.sub_items != null">
        <div v-for="(subitem,index) in model.sub_items" :key="index" class="m-0">
            <OrderItemView v-model="model.sub_items[index]" />
        </div>
    </div>
    <Dialog v-model:visible="new_component_dialog" modal header="Add material">
        <PickMaterial @returnMaterial="(material) => addMaterial(material)" />
    </Dialog>
</template>

<script setup lang="ts">
import {defineModel,ref} from 'vue'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import Dropdown from 'primevue/dropdown'
import InputSwitch from 'primevue/inputswitch';
import { Material, OrderItem, OrderItemMaterial } from '@/classes/OrderItem'
import PickMaterial from '@/components/PickMaterial.vue'
import Dialog from 'primevue/dialog'

const model = defineModel<OrderItem>({
    required: true})


const new_component_dialog = ref(false)


const updateEntriesCost = () => {
    model.value.materials.forEach(async (material: OrderItemMaterial,index: number) => {
        await model.value.UpdateMaterialEntryCost(index)
    })
}

updateEntriesCost();

const addMaterial = (material: Material) => {
    model.value.PushMaterial(material)
}


// watch(model.value.materials, () => {
//     updateEntriesCost()
// },
// {deep: true})


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