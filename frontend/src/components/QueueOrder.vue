<template>
    <div :style="`display:${display};${state == 'cancelled' ? 'position:relative;' : ''}`">
        <div v-if="state =='cancelled'" style="position:absolute;height:100%;width:100%;z-index:9" class="flex justify-content-center align-items-center">
            <div class="flex flex-column justify-content-center align-items-center">
                <span style="color:white;font-size:1.3rem;font-weight:800;text-decoration: line-through;background-color: rgb(0,119,255);">{{ props.order.display_id }}</span>
                <h2 style="color:tomato;font-weight:800">{{$t('cancelled')}}</h2>
            </div>
        </div>
        <Card :style="`width: 20rem; overflow: hidden;${state == 'cancelled' ? 'filter:grayscale(1) blur(0.2rem);' : ''}`">
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
                            <Button :loading="loading" :label="$t('start')" iconPos="right" icon="pi pi-play" class="w-full" @click="prepareOrder" severity="info" />
                        </ButtonGroup>
                        <ButtonGroup v-if="state == 'in_progress'" class="w-full">
                            <Button  icon="pi pi-trash" class="w-3" severity="secondary" />
                            <ConfirmPopup></ConfirmPopup>
                            <Button icon="pi pi-check" class="w-9" aria-label="Finish" :label="$t('finish')" @click="confirmFinish($event)" severity="success" iconPos="right" />
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
            <Stepper :value="1">
                <StepList>
                    <Step v-for="item,index in items" :key="index" :value="index+1">
                        {{ item.product.name }}
                    </Step>
                </StepList>
                <StepPanels>
                    <StepPanel v-for="item,index in items" :key="index" :value="index+1" v-slot="{ activateCallback }">
                        <Message v-if="item.comment != ''" severity="warn">{{ item.comment }}</Message>
                        <OrderItemView v-model="items[index]" />
                        <div class="flex pt-6 justify-end" >
                            <Button v-if="index == 0 && index != items.length-1" label="Next" icon="pi pi-arrow-right" iconPos="right" @click="activateCallback(index+2)" />
                            <Button v-if="index != items.length-1 && index != 0" label="Back" severity="secondary" icon="pi pi-arrow-left" @click="activateCallback(index)" />
                            <Button v-if="index == items.length-1" @click="startOrder(); visible=false;" :label="$t('start')" :disabled="!isValid" severity="success" />
                        </div>
                    </StepPanel>
                </StepPanels>
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
import StepList from 'primevue/steplist';
import StepPanels from 'primevue/steppanels';
import Step from 'primevue/step';
import StepPanel from 'primevue/steppanel';
import Message from 'primevue/message';
import Divider from 'primevue/divider';
import { useConfirm } from "primevue/useconfirm";
import ConfirmPopup from 'primevue/confirmpopup';
import { useToast } from "primevue/usetoast";
import OrderItemView from "./OrderItemView.vue";
import {OrderItem, Product} from '@/classes/OrderItem'
import {getCurrentInstance,nextTick} from 'vue'


const { proxy } = getCurrentInstance();


const loading = ref(false)
const toast = useToast();

const items = ref<OrderItem[]>([])


const confirm = useConfirm();
const product_details_visible= ref(false)
const product_for_details = ref<Product>()

const started_at = ref("")


const state = ref("pending")


// const orderItemSelectedOptions = ref({})



const visible = ref(false)
const props = defineProps(['order','number'])

const timePassed = ref("")


const emit = defineEmits(['openedDialog', 'closedDialog','finished']);


watch(props.order.state, (newVal) => {
    state.value = newVal
})


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
        accept: async () => {
            await finishOrder();
            toast.add({ severity: 'success', summary: 'Order finished', detail: 'Good job 🎉', life: 3000,group:'br' });
        },
        reject: () => {
            toast.add({ severity: 'error', summary: 'Failed', detail: 'Failed finishing order !', life: 3000,group:'br' });  
        }
    });
};


const finishOrder = () => {

    axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders/${props.order.id}/finish`,
        {},
        {
            headers: {
                Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
            }
        }
        ).then(() => {
            emit('closedDialog');
            state.value = 'finished'
        })
}



const startOrder =  () => {

    axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders/${props.order.id}/start`,
        {
            "data" : items.value
        },
        {
            headers: {
                Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
            }
        }
        ).then((response) => {
            started_at.value = response.data.started_at
            state.value = "in_progress"
        }).catch((error) => {
            toast.removeGroup('br')
            toast.add({ severity: 'error', summary: 'Error', detail: error.response.data.body, life: 5000,group:'br' });
        })
}



const prepareOrder = async () => {

    items.value = []
    loading.value = true


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
    }
    
    await nextTick()
    visible.value = true  

    loading.value = false
}


const display : boolean = computed(() => {

let valid = 'block'

if (state.value == 'finished')
    valid = 'none'

if (state.value == 'cancelled'){
    if (props.order.submitted_at) {
        const submittedAt = moment(props.order.submitted_at);
        const now = moment();
        const duration = moment.duration(now.diff(submittedAt));
        
        if (duration.asHours() >= 1) {
            valid = 'none'
        }
    }
}


return valid

})

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

        updateElapsedTime();
    }

    state.value = props.order.state
}


init();

</script>

<style>
</style>