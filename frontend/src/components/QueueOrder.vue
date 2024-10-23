<template>
    <div :style="`display:${state == 'finished' ? 'none' : 'block'}`">
        <Card style="width: 20rem; overflow: hidden;">
            <template #header>
                <!-- <h1 class="m-2">#{{props.number}}</h1> -->
                <!-- 2024-06-20T14:31:39.946Z -->
                <div class="grid mt-1 p-2">
                <!-- <div class="flex gap-3 mt-1 p-2 justify-content-center align-items-center"> -->
                    <h3 class="col-12 m-0 p-0 px-2">
                        {{props.order.display_id}}
                    </h3>
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
                                    <h3 class="m-0">{{item.product.name}}</h3>
                                    <Button icon="pi pi-book" severity="contrast" @click="showProductDetails(item.product)" text rounded aria-label="Star" />
                                </div>
                                <!-- <h1 class="m-0" style="color:blue">x{{item.quantity}}</h1> -->
                                <p v-if="item.comment && item.comment != ''" class="mt-1 mb-0">
                                    {{ item.comment }}
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
            </template>
        </Card>
        <Dialog v-model:visible="product_details_visible" modal :header="`Product:  (${product_for_details?.name}) details`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
           <ul>
              <li class="flex justify-content-between w-6 m-2" v-for="(material,index) in product_for_details?.materials" :key="index"><strong>{{material.name}}:</strong> &nbsp;&nbsp;&nbsp;&nbsp;{{material.quantity}} {{material.unit}}</li>
           </ul>
        </Dialog>
        <Dialog v-model:visible="visible" modal :header="`Order #${props.order.display_id}`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
            <!-- <Dialog v-model:visible="visible" modal :header="props.order.items[currentItemIndex].name+` #${currentItemIndex+1}`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '75vw', '575px': '90vw' }"> -->
            <Stepper @update:activeStep="(number) => {currentItemIndex = number}">
                <StepperPanel v-for="item,index in items" :key="index" :header="item.product.name">
                    <template #content="{ prevCallback, nextCallback }">
                        <Message v-if="props.order.items[currentItemIndex].comment != ''" severity="warn">{{ props.order.items[currentItemIndex].comment }}</Message>
                        <OrderItemView v-model="items[index]" />
                        <div class="flex pt-4 justify-content-between">
                            <Button label="Back" severity="secondary" :disabled="currentItemIndex==0" icon="pi pi-arrow-left" @click="prevCallback" />
                            <Button :disabled="currentItemIndex == props.order.items.length-1 && !isValid" :label="currentItemIndex == props.order.items.length-1 ? 'Go' : 'Next'" :icon="currentItemIndex != props.order.items.length-1 ? 'pi pi-arrow-right' : ''" iconPos="right" @click="if (currentItemIndex == props.order.items.length-1) {startOrder(); visible=false;} else nextCallback()" />
                        </div>
                    </template>
                </StepperPanel>
            </Stepper>
        </Dialog>
    </div>
</template>

<script setup lang="ts">
import {ref, defineProps, watch, defineEmits, computed} from 'vue'

import Card from 'primevue/card';
import Button from 'primevue/button';
import ButtonGroup from 'primevue/buttongroup';
import Dialog from 'primevue/dialog'
import moment from 'moment';
import axios from 'axios';
import Stepper from 'primevue/stepper';
import StepperPanel from 'primevue/stepperpanel';
import Message from 'primevue/message';
import Divider from 'primevue/divider';
import { useConfirm } from "primevue/useconfirm";
import ConfirmPopup from 'primevue/confirmpopup';
import { useToast } from "primevue/usetoast";
import OrderItemView from "./OrderItemView.vue";
import {OrderItem, Product} from '@/classes/OrderItem'


const toast = useToast();

const items = ref<OrderItem[]>([])


const confirm = useConfirm();
const product_details_visible= ref(false)
const product_for_details = ref<Product>()

const state = ref("pending")
const started_at = ref("")


// const orderItemSelectedOptions = ref({})

const currentItemIndex = ref(0)



const visible = ref(false)
const props = defineProps(['order','number'])

const timePassed = ref("")


const emit = defineEmits(['openedDialog', 'closedDialog']);

watch(visible, (newVal) => {
  if (newVal){
    if (!product_details_visible.value){
        emit('openedDialog');
    }
  }
  if (!newVal){
    if (!product_details_visible.value){
        emit('closedDialog');
    }
  }
})

watch(product_details_visible, (newVal) => {
  if (newVal){
    if (!visible.value){
        emit('openedDialog');
    }
  }
  if (!newVal){
    if (!visible.value){
        emit('closedDialog');
    }
  }
})


const showProductDetails = (itemRecipe) => {
    product_for_details.value = itemRecipe;
    product_details_visible.value = true;
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

    axios.post(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/finishorder`,
        {
            "order_id":props.order.id
        }
        ).then(() => {
            state.value = "finished"
        })
}



const startOrder =  () => {

    axios.post(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/startorder`,
        {
            "order_id":props.order.id,
            "items" : items.value
        }
        ).then((response) => {
            state.value = "in_progress"
            started_at.value = response.data.started_at
        }).catch((error) => {
            toast.removeGroup('br')
            toast.add({ severity: 'error', summary: 'Error', detail: error.response.data.body, life: 5000,group:'br' });
        })
}



const prepareOrder = async () => {

    currentItemIndex.value = 0
    items.value = []


    for(var i=0;i<props.order.items.length;i++){
        const item = new OrderItem()
        item.SetProductId(props.order.items[i].product.id)
        await item.RefreshReadyNumber()
        await item.FromItemData(props.order.items[i])// const item = new OrderItem(orderItem.product)
        await item.RefreshProductData()

        item.materials.forEach((material,materialIndex) => {
            item.ValidateMaterialQuantity(materialIndex)
        })
        

        items.value.push(item)
        visible.value = true  
    }
}

const isValid : boolean = computed(() => {

    let valid = true


    items.value.forEach((item)=>{
        if(!item.isValid)
            valid = false
    })

    return valid

})


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