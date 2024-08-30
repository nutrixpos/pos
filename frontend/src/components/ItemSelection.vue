<template>
    <div class="flex justify-content-between align-items-center">
        <h4>{{ model.name }}</h4>
        <InputText :disabled="!model.isConsumeFromReady" type="text" v-model="model.Quantity"  size="small"/>
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
                <InputText type="text" v-model="itemsComponentQuantity[index]"  size="small"/>
                <span class="ml-2 mt-2">{{ component.unit }}</span>
            </div>
            <Dropdown v-if="props.item.components[index].entries != null && props.item.components[index].entries.length > 0" v-model="itemsEntrySelection[index]"  :options="props.item.components[index].entries" optionLabel="label" placeholder="Select option" class="w-6" />
        </div>
    </div>
    <div v-for="(subrecipe,index) in props.item.sub_recipes" :key="index" class="m-0">
        <ItemSelection :isReturn="props.isReturn" :item="subrecipe" v-model="itemDetails.SubRecipes[index]" />
    </div>
</template>

<script setup lang="ts">
import {ref, watch,defineModel} from 'vue'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import InputSwitch from 'primevue/inputswitch';
import { ComponentEntry,Recipe,RecipeSelections} from '@/classes/ItemSelection'




interface Props {
    item: Recipe,
    isReturn: boolean,
}


const isReturnArray = ref<boolean[]>([])
const props = defineProps<Props>()
const model = defineModel<RecipeSelections>()


const itemsEntrySelection = ref<Array<ComponentEntry | null>>([])
const itemsComponentQuantity = ref<string[]>([])



const itemDetails = ref<RecipeSelections>(new RecipeSelections(props.item));
itemDetails.value.SubRecipes = []


if (props.item.sub_recipes != null)
props.item.sub_recipes?.forEach(subrecipe => {

    let newSubRecipe = new RecipeSelections(subrecipe)
    itemDetails.value.SubRecipes?.push(newSubRecipe)

    isReturnArray.value.push(false)
})


const init = () => {
    props.item.components.forEach((component,index) => {
        itemsComponentQuantity.value.push(component.defaultquantity.toString())
        var entries = component.entries.map(entry => {
                    return {
                        ...entry,
                        label:entry.company + " - " + entry.quantity + " " + entry.unit
                    }
                })
        component.entries = entries
        itemsEntrySelection.value[index] = component.entries.length > 0 ? component.entries[0] : null
    })
}

init();

watch(() => props.isReturn, (newval, oldval) => {
      if (newval == true && oldval == false ){
        returnSelection()
      }
    });


const returnSelection = () => {


    // itemDetails.SubRecipes?.forEach((_,index) => {
    //     isReturnArray.value[index] = true
    // })

    itemDetails.value.Id = props.item.recipe_id
    itemDetails.value.recipe_name = props.item.recipe_name

    itemDetails.value.Selections = []
    
    itemsEntrySelection.value.forEach((entrySelection,index) => {
        itemDetails.value.Selections?.push({
            name: props.item.components[index].name,
            EntryId: entrySelection?._id,
            ComponentId: props.item.components[index].component_id,
            unit: props.item.components[index].unit,
            quantity: Number(itemsComponentQuantity.value[index])
        })
    })
    
    model.value = itemDetails.value
}


</script>