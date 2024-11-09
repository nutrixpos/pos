<template>
    <div>
        <h1 class="m-3">Kitchen</h1>
        <div class="flex flex-wrap">
            <QueueOrder @finished="orderFinished(index)" @openedDialog="openedDialogs++" @closedDialog="openedDialogs--" v-for="(order,index) in orders" :key="index" :order="order" :number="index+1" class="queue-order"/>
        </div>
    </div>
</template>

<script setup>
import QueueOrder from '@/components/QueueOrder.vue'
import axios from 'axios';
import {ref,getCurrentInstance} from 'vue';

const { proxy } = getCurrentInstance();


const orders = ref([])
const openedDialogs = ref(0)

const orderFinished = (index) => {
    orders.value.splice(index,1)
}


const loadOrders =  () => {
    axios.get(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/orders`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel.oidcAuth.accessToken}`
        }
    })
    .then((result)=>{
        orders.value = result.data
    })
};


setInterval(() => {
    if (openedDialogs.value == 0)
        loadOrders()
}, 3000);


loadOrders()

</script>

<style>
.queue-order {
    margin:5px;
}
</style>