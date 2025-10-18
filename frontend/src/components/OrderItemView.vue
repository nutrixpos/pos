<template>
    <div class="flex justify-content-between align-items-center">
        <h4>{{ model.product.name }}</h4>
        <div>
            x <InputText :disabled="!model.is_consume_from_ready" type="number" v-model.number="model.quantity"  size="small"/>
        </div>
        <div class="flex align-items-center justify-content-center">
            <span class="mx-2">{{$t('from_ready')}}</span>
            <ToggleSwitch @change="model.ValidateItem()" v-model="model.is_consume_from_ready" :disabled="!model.can_change_ready_toggle" />
            <span class="mx-2">
                <p style="font-size: 0.9rem;">{{model.ready}} {{$t('ready')}}</p>
            </span>
        </div>
    </div>
    <div v-if="!model.is_consume_from_ready">
        <Button :label="$t('add_material')" @click="new_component_dialog = true" />
        <div class="flex my-3 py-2 justify-content-between align-items-center" style="border-bottom:1px solid gray" v-for="(material,index) in model.materials" :key="index">

            <Button icon="pi pi-times" size="small" style="width:2rem;height: 2rem;" aria-label="Remove" severity="secondary" @click="removeMaterialByIndex(index)" />

            {{ material.material.name }}
            <div class="flex">
                <InputText type="number" @change="MaterialInputChanged(index)" :invalid="!material.isQuantityValid" v-model.number="model.materials[index].quantity" size="small"/>
                <span class="ml-2 mt-2">{{ material.material.unit }}</span>
            </div>
            <Dropdown @change="EntryDropDownChanged(index)"  v-model="model.materials[index].entry"  :options="model.materials[index].material.entries" optionLabel="label" placeholder="Select option" class="w-6" />
            {{$t('cost')}}: {{ material.entry?.cost * model.quantity }}
        </div>
    </div>
    <div v-if="model.sub_items != null && !model.is_consume_from_ready">
        <div v-for="(subitem,index) in model.sub_items" :key="index" class="m-0">
            <OrderItemView @changed="validateItem()" v-model="model.sub_items[index]" />
        </div>
    </div>
    <Dialog v-model:visible="new_component_dialog" modal header="Add material">
        <PickMaterial @returnMaterial="(material) => addMaterial(material)" />
    </Dialog>
</template>

<script setup lang="ts">
import {defineModel,ref, watch,defineEmits} from 'vue'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import Dropdown from 'primevue/dropdown'
import ToggleSwitch from 'primevue/toggleswitch';
import { Material, OrderItem } from '@/classes/OrderItem'
import PickMaterial from '@/components/PickMaterial.vue'
import Dialog from 'primevue/dialog'

const model = defineModel<OrderItem>({
    required: true})


const emit = defineEmits(['changed'])


const new_component_dialog = ref(false)
const materialValidity = ref<boolean[]>([])


const removeMaterialByIndex = (index:number) => {
    model.value.RemoveMaterialByIndex(index)
    materialValidity.value.splice(index,1)
    emit('changed')
}

const MaterialInputChanged = async (index:number) => {
    await model.value.UpdateMaterialEntryCost(index); 
    await validateMaterialQuantity(index); 
    emit('changed'); 
}

const EntryDropDownChanged = async (index:number) => {
    await model.value.UpdateMaterialEntryCost(index);
    emit('changed'); 
}


const validateMaterialQuantity = (index: number)  => {
   model.value.ValidateMaterialQuantity(index)
}

const validateItem = () => {
    model.value.ValidateItem()
}

const Init = async () => {
    await updateEntriesCost()
    model.value.materials.forEach((_,index) => {
        model.value.ValidateMaterialQuantity(index)
    })
}

const updateEntriesCost = async () => {

    for (var i=0;i<model.value.materials.length;i++){
        await model.value.UpdateMaterialEntryCost(i)
    }
}

Init()

const addMaterial = async (material: Material) => {
    await model.value.PushMaterial(material)
    validateMaterialQuantity(model.value.materials.length - 1)
}


watch(materialValidity, (newValidities) => {
    let valid = true 

    newValidities.forEach((validity) => {
        if (!validity){
            valid = false
        }  
    })

    model.value.isValid = valid
},
{deep:true})


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