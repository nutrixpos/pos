<template>
    <div class="flex justify-content-between align-items-center">
        <h4>{{ model.recipe_name }}</h4>
        <InputText :disabled="!model.isConsumeFromReady" type="number" v-model="model.Quantity"  size="small"/>
        <div class="flex align-items-center justify-content-center">
            <span class="mx-2">From Ready</span>
            <InputSwitch v-model="model.isConsumeFromReady" :disabled="!model.canChangeReadyToggle" />
            <span class="mx-2">
                <p style="font-size: 0.9rem;">{{model.Ready}} Ready</p>
            </span>
        </div>
    </div>
    <div v-if="!model.isConsumeFromReady">
        <div class="flex my-3 py-2 justify-content-between" style="border-bottom:1px solid gray" v-for="(component,index) in model.Components" :key="index">
            {{ component.name }}
            <div class="flex">
                <InputText type="number" v-model="model.Selections[index].Quantity"  size="small"/>
                <span class="ml-2 mt-2">{{ component.unit }}</span>
            </div>
            <Dropdown v-if="model.Components[index].entries != null && model.Components[index].entries.length > 0" v-model="model.Selections[index].Entry"  :options="model.Components[index].entries" optionLabel="label" placeholder="Select option" class="w-6" />
        </div>
    </div>
    <div v-for="(subrecipe,index) in model.SubRecipes" :key="index" class="m-0">
        <ItemSelection v-model="model.SubRecipes[index]" />
    </div>
</template>

<script setup lang="ts">
import {defineModel} from 'vue'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import InputSwitch from 'primevue/inputswitch';
import { RecipeSelections, ComponentSelection} from '@/classes/ItemSelection'

const model = defineModel<RecipeSelections>({
    required: true})


const init = () => {

    model.value.Selections = []

    model.value.Components.forEach((component) => {
        const selection = new ComponentSelection()
        selection.ComponentId = component.component_id
        selection.Quantity = component.defaultquantity
        selection.Name = component.name
        selection.Unit = component.unit
        selection.Entry = component.entries[0]

        model.value.Selections.push(selection)
    })
}


init()

</script>