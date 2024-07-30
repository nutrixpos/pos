<template>
    <div>
        <h1 class="m-3">Kitchen</h1>
        <div class="flex flex-wrap">
            <QueueOrder v-for="(order,index) in orders" :key="index" :order="order" :number="index+1" class="queue-order"/>
        </div>
    </div>
</template>

<script setup>
import QueueOrder from '@/components/QueueOrder.vue'
import axios from 'axios';
import {ref} from 'vue';


const orders = ref([])


const loadOrders =  () => {
    axios.get('http://localhost:8000/api/orders')
    .then((result)=>{
        orders.value = result.data
    })
};


loadOrders()

</script>

<style>
.queue-order {
    margin:5px;
}
</style>