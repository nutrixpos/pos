<template>
    <div :style="`display:${state == 'finished' ? 'none' : 'block'}`">
        <Card style="width: 20rem; overflow: hidden;">
            <template #header>
                <!-- <h1 class="m-2">#{{props.number}}</h1> -->
                <!-- 2024-06-20T14:31:39.946Z -->
                <div class="grid mt-1 p-2">
                <!-- <div class="flex gap-3 mt-1 p-2 justify-content-center align-items-center"> -->
                    <div class="col-3 flex justify-content-center align-items-center">
                        <p class="px-2"><strong>{{timePassed}}</strong></p>
                    </div>
                    <div class="col-9 flex justify-content-center align-items-center">
                        <ButtonGroup v-if="state != 'in_progress'"  class="w-full">
                            <Button label="Start" iconPos="right" icon="pi pi-play" class="w-full" @click="prepareOrder" severity="info" />
                        </ButtonGroup>
                        <ButtonGroup v-if="state == 'in_progress'" class="w-full">
                            <Button  icon="pi pi-trash" class="w-3" severity="secondary" />
                            <ConfirmPopup></ConfirmPopup>
                            <Button icon="pi pi-check" class="w-9" aria-label="Finish" label="Finish" @click="confirmFinish($event)" severity="success" iconPos="right" />
                        </ButtonGroup>
                    </div>
                </div>
            </template>
            <template #content>
                <div class="flex" v-for="(item,index) in props.order.items" :key="index">
                    <div class="w-full flex my-1 flex-column">
                        <Divider v-if="index > 0" />
                        <div class="w-full flex ">
                            <div :style="`width:.2rem;background-color:${item.comment != '' ? 'orange' : 'silver'}`" class="mr-2"></div>
                            <div class="flex flex-column w-full justify-content-center my-2">
                                <div class="flex justify-content-between align-items-center">
                                    <h3 class="m-0">{{item.name}}</h3>
                                    <Button icon="pi pi-book" severity="contrast" @click="showRecipe(item.recipe)" text rounded aria-label="Star" />
                                </div>
                                <!-- <h1 class="m-0" style="color:blue">x{{item.quantity}}</h1> -->
                                <p v-if="item.comment != ''" class="mt-1 mb-0">
                                    {{ item.comment }}
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </template>
        </Card>
        <Dialog v-model:visible="recipe_visible" modal :header="`Recipe:  ${recipe.Name}`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
           <ul>
              <li class="flex justify-content-between w-6 m-2" v-for="(component,index) in recipe.Components" :key="index"><strong>{{component.Name}}:</strong> &nbsp;&nbsp;&nbsp;&nbsp;{{component.Quantity}} {{component.Unit}}</li>
           </ul>
        </Dialog>
        <Dialog v-model:visible="visible" modal :header="`Order #${props.number}`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
            <!-- <Dialog v-model:visible="visible" modal :header="props.order.items[currentItemIndex].name+` #${currentItemIndex+1}`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"> -->
            <Stepper @update:activeStep="(number) => {currentItemIndex = number}">
                <StepperPanel v-for="item,index in props.order.items" :key="index" :header="item.name">
                    <template #content="{ prevCallback, nextCallback }">
                        <Message v-if="props.order.items[currentItemIndex].comment != ''" severity="warn">{{ props.order.items[currentItemIndex].comment }}</Message>
                        <div class="flex my-3 py-2 justify-content-between" style="border-bottom:1px solid gray" v-for="(component,index) in itemComponentEntries[currentItemIndex].components" :key="index">
                            {{ component.name }}
                            <div class="flex">
                                <InputText type="text" v-model="itemsComponentQuantity[currentItemIndex][index]"  size="small"/>
                                <span class="ml-2 mt-2">{{ component.unit }}</span>
                            </div>
                            <Dropdown v-if="itemComponentEntries[currentItemIndex].components[index].entries != null && itemComponentEntries[currentItemIndex].components[index].entries.length > 0" v-model="itemsEntrySelection[currentItemIndex][index]"  :options="itemComponentEntries[currentItemIndex].components[index].entries" optionLabel="label" placeholder="Select option" class="w-6" />
                        </div>
                        <div class="flex pt-4 justify-content-between">
                            <Button label="Back" severity="secondary" :disabled="currentItemIndex==0" icon="pi pi-arrow-left" @click="prevCallback" />
                            <Button :label="currentItemIndex == props.order.items.length-1 ? 'Go' : 'Next'" :icon="currentItemIndex != props.order.items.length-1 ? 'pi pi-arrow-right' : ''" iconPos="right" @click="if (currentItemIndex == props.order.items.length-1) {startOrder(); visible=false;} else nextCallback()" />
                        </div>
                    </template>
                </StepperPanel>
            </Stepper>
        </Dialog>
    </div>
</template>

<script setup>
import {ref, defineProps} from 'vue'

import Card from 'primevue/card';
import Button from 'primevue/button';
import ButtonGroup from 'primevue/buttongroup';
import Dialog from 'primevue/dialog'
import moment from 'moment';
import axios from 'axios';
import Dropdown from 'primevue/dropdown';
import InputText from 'primevue/inputtext';
import Stepper from 'primevue/stepper';
import StepperPanel from 'primevue/stepperpanel';
import Message from 'primevue/message';
import Divider from 'primevue/divider';
import { useConfirm } from "primevue/useconfirm";
import ConfirmPopup from 'primevue/confirmpopup';
import { useToast } from "primevue/usetoast";


const toast = useToast();


const confirm = useConfirm();


const itemComponentEntries = ref([])
const itemsEntrySelection = ref([[]])
const itemsComponentQuantity = ref([[]])
const recipe_visible= ref(false)

const state = ref("pending")
const started_at = ref("")


// const orderItemSelectedOptions = ref({})

const currentItemIndex = ref(0)
const counter = ref(0)



const visible = ref(false)
const props = defineProps(['order','number'])

const timePassed = ref("")

const recipe = ref({
    Name:"null"
})


const showRecipe = (itemRecipe) => {
    recipe.value = itemRecipe;
    recipe_visible.value = true;
}


const updateElapsedTime = () => {
    const now = moment();
    timePassed.value =  formatDuration(moment.duration(now.diff(props.order.submitted_at)))
    // moment(String(props.order.submitted_at)).fromNow()
    setInterval(function(){
        const now = moment();
        timePassed.value = formatDuration(moment.duration(now.diff(props.order.submitted_at)))
    },1000)
}

const formatDuration = (duration) => {
    const hours = Math.floor(duration.asHours());
    const minutes = Math.floor(duration.asMinutes()) - hours * 60;
    const seconds = Math.floor(duration.asSeconds()) - minutes * 60 - hours * 3600;
    return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${seconds.toString().padStart(2, '0')}`;
}


const confirmFinish = (event) => {
    confirm.require({
        target: event.currentTarget,
        message: 'Are you sure you want to finish the order ?',
        icon: 'pi pi-exclamation-triangle',
        rejectProps: {
            label: 'Cancel',
            severity: 'secondary',
            outlined: true
        },
        acceptProps: {
            label: 'Yes'
        },
        accept: () => {
            finishOrder();
            toast.add({ severity: 'success', summary: 'Order finished', detail: 'Good job 🎉', life: 3000,group:'br' });
        },
        reject: () => {
            
        }
    });
};


const finishOrder = () => {

    axios.post("http://localhost:8000/api/finishorder",
        {
            "order_id":props.order.id
        }
        ).then(() => {
            state.value = "finished"
        })
}


const startOrder =  () => {

    var ingredients = []

    itemsComponentQuantity.value.forEach((item,itemIndex) => {

        ingredients.push([])

        item.forEach((quantity,componentIndex) => {
            ingredients[itemIndex][componentIndex] = {
                name : itemComponentEntries.value[itemIndex].components[componentIndex].name,
                quantity,
                entry_id: itemsEntrySelection.value[itemIndex][componentIndex] != null ? itemsEntrySelection.value[itemIndex][componentIndex]._id : "null", 
                company: itemsEntrySelection.value[itemIndex][componentIndex] != null ? itemsEntrySelection.value[itemIndex][componentIndex].company : "null"
            }
        })
    })


    axios.post("http://localhost:8000/api/startorder",
        {
            "order_id":props.order.id,
            "ingredients" : ingredients
        }
        ).then((response) => {
            state.value = "in_progress"
            started_at.value = response.data.started_at
        })
}



const prepareOrder = () => {


    itemsComponentQuantity.value = [[]]
    itemsEntrySelection.value = [[]]
    itemComponentEntries.value = []
    currentItemIndex.value = 0
    counter.value = 0 ;


    props.order.items.forEach((item) => {

        axios.get("http://localhost:8000/api/recipetree?id="+item.id,).then((response) => {
            var components = []

            response.data.components.forEach((component) => {


                var entries = component.entries.map(entry => {
                    return {
                        ...entry,
                        label:entry.company + " - " + entry.quantity + " " + entry.unit
                    }
                })

                components.push({
                    index: counter.value,
                    name: component.name,
                    defaultquantity: component.defaultquantity,
                    unit: component.unit,
                    entries
                })


                if (itemsEntrySelection.value.length < (counter.value+1)){
                    itemsEntrySelection.value.push([])
                    itemsComponentQuantity.value.push([])
                }
                
                itemsEntrySelection.value[counter.value].push(entries.length > 0 ? entries[0] : null)
                itemsComponentQuantity.value[counter.value].push(component.defaultquantity)
                

                
            })
            
            itemComponentEntries.value.push({name:item.name,components})
            counter.value++
            visible.value = true
        })  
    })
}


const init = () => {
    if (props.order.started_at != null){
        started_at.value = props.order.started_at
        state.value = props.order.state

        updateElapsedTime();
    }
}


init();

</script>

<style>
</style>