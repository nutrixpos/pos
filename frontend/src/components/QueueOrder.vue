<template>
    <div>
        <Card style="width: 25rem; overflow: hidden">
            <template #header>
                <!-- <h1 class="m-2">#{{props.number}}</h1> -->
                <!-- 2024-06-20T14:31:39.946Z -->
                <div class="grid mt-1 p-2">
                <!-- <div class="flex gap-3 mt-1 p-2 justify-content-center align-items-center"> -->
                    <div class="col-3 flex justify-content-center align-items-center">
                        <p class="px-2"><strong>{{timePassed}}</strong></p>
                    </div>
                    <div class="col-9 flex justify-content-center align-items-center">
                        <Button v-if="state != 'in_progress'" label="Start" class="w-full" @click="prepareOrder" severity="info" style="height:2.5rem !important;" />
                        <ButtonGroup v-if="state == 'in_progress'">
                            <Button  icon="pi pi-trash" severity="warning" />
                            <Button  icon="pi pi-info-circle" aria-label="Info" label="Info" severity="secondary" iconPos="right" />
                            <Button icon="pi pi-check" aria-label="Finish" label="Finish" severity="success" iconPos="right" />
                        </ButtonGroup>
                    </div>
                </div>
            </template>
            <template #content>
                <div class="flex" v-for="(item,index) in props.order.items" :key="index">
                    <Avatar shape="circle" image="https://girlheartfood.com/wp-content/uploads/2020/06/Crispy-Chicken-Burger-10.jpg" class="mr-2" size="xlarge" />
                    <div class="flex flex-column w-9">
                        <h1 class="m-0">{{item.name}}</h1>
                        <!-- <h1 class="m-0" style="color:blue">x{{item.quantity}}</h1> -->
                        <p>
                            {{ item.comment }}
                        </p>
                    </div>
                </div>
            </template>
        </Card>
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
import Avatar from 'primevue/avatar';
import Dialog from 'primevue/dialog'
import moment from 'moment';
import axios from 'axios';
import Dropdown from 'primevue/dropdown';
import InputText from 'primevue/inputtext';
import Stepper from 'primevue/stepper';
import StepperPanel from 'primevue/stepperpanel';
import Message from 'primevue/message';

const itemComponentEntries = ref([])
const itemsEntrySelection = ref([[]])
const itemsComponentQuantity = ref([[]])

const state = ref("pending")
const started_at = ref("")


// const orderItemSelectedOptions = ref({})

const currentItemIndex = ref(0)
const counter = ref(0)



const visible = ref(false)
const props = defineProps(['order','number'])

const timePassed = ref("")



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




const startOrder =  () => {

    var ingredients = []

    itemsComponentQuantity.value.forEach((item,itemIndex) => {

        ingredients.push([])

        item.forEach((quantity,componentIndex) => {
            ingredients[itemIndex][componentIndex] = {
                name : itemComponentEntries.value[itemIndex].components[componentIndex].name,
                quantity
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

        axios.post("http://localhost:8000/api/prepareitem",
        {
            name:item.name
        }
        ).then((response) => {
            var components = []

            response.data.forEach((component) => {


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