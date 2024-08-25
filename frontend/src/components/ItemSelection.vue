<template>
    <h4>{{ props.item.recipe_name }}</h4>
    <div class="flex my-3 py-2 justify-content-between" style="border-bottom:1px solid gray" v-for="(component,index) in props.item.components" :key="index">
        {{ component.name }}
        <div class="flex">
            <InputText type="text" v-model="itemsComponentQuantity[index]"  size="small"/>
            <span class="ml-2 mt-2">{{ component.unit }}</span>
        </div>
        <Dropdown v-if="props.item.components[index].entries != null && props.item.components[index].entries.length > 0" v-model="itemsEntrySelection[index]"  :options="props.item.components[index].entries" optionLabel="label" placeholder="Select option" class="w-6" />
    </div>
    <div v-for="(subrecipe,index) in props.item.subrecipes" :key="index" class="m-0">
        <ItemSelection :item="subrecipe" />
    </div>
</template>

<script setup lang="ts">
import {ref, defineProps, defineEmits, watch,PropType} from 'vue'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'

interface Component {
    name: string,
    unit: string,
    defaultquantity: number,
    entries: Array<any>,
}

interface ItemSelection {
    name: string,
    Id: string,
    components: Array<Component>,
    subrecipes: Array<ItemSelection>,
}

const props = defineProps({
    item :{
        type: Object as PropType<ItemSelection>,
        required: true
    }
})

const itemsEntrySelection = ref([])
const itemsComponentQuantity = ref([])

const subRecipesEntriesSelection = ref([[]])
const subRecipesComponentQuantity = ref([[]])

const emit = defineEmits(['update:itemsEntrySelection','update:itemsComponentQuantity'])


const init = () => {
    props.item.components.forEach((component,) => {
        itemsComponentQuantity.value.push(component.defaultquantity)
        var entries = component.entries.map(entry => {
                    return {
                        ...entry,
                        label:entry.company + " - " + entry.quantity + " " + entry.unit
                    }
                })
        component.entries = entries
        itemsEntrySelection.value.push(component.entries.length > 0 ? component.entries[0] : null)
    })


    props.item.subrecipes.forEach(() => {
        subRecipesEntriesSelection.value.push([])
        subRecipesComponentQuantity.value.push([])
    })
}

init();

watch(itemsEntrySelection, () => {
    emit('update:itemsEntrySelection',itemsEntrySelection.value)
    emit('update:itemsComponentQuantity',itemsComponentQuantity.value)
})


</script>