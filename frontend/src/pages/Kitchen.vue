<template>
    <div>
        <h1 class="m-3">Kitchen</h1>
        <div class="flex flex-wrap">
            <QueueOrder @openedDialog="openedDialogs++" @closedDialog="openedDialogs--" v-for="(order,index) in orders" :key="index" :order="order" :number="index+1" class="queue-order"/>
        </div>
    </div>
</template>

<script setup>
import QueueOrder from '@/components/QueueOrder.vue'
import axios from 'axios';
import {ref} from 'vue';


const orders = ref([])
const openedDialogs = ref(0)


const loadOrders =  () => {
    axios.get(`http://${process.env.VUE_APP_BACKEND_HOST}/api/orders`)
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