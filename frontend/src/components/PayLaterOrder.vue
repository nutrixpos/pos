<template>
    <div>
        <Message severity="secondary">
            <template #container>
                <div class="p-3 flex justify-content-between align-items-center">
                    <span>{{ props.order.display_id }}</span>
                    <Badge :value="order_status.title" :severity="order_status.severity" />
                    <span>{{ props.order.items.length }} Item(s)</span>
                    <Button @click="payUnpaidOrder()">
                        {{$t('checkout')}}
                        <Badge :value="order_unpaid_amount_str" class="p-badge-secondary" />
                    </Button>
                </div>
            </template>
        </Message>
    </div>
</template>

<script setup lang="ts">
import {defineProps,computed,getCurrentInstance,defineEmits} from 'vue'
import Message from 'primevue/message';
import Button from 'primevue/button';
import Badge from 'primevue/badge';
import Order from '@/classes/Order'
import axios from 'axios';
import { useToast } from "primevue/usetoast";

const { proxy } = getCurrentInstance();
const toast = useToast();


const emit = defineEmits(['order_paid'])

const payUnpaidOrder = () => {
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders/${props.order.id}/pay`,{
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(() => {
        toast.add({ severity: 'success', summary: 'Success', detail: `Successfully paid order ${props.order.display_id}`, life: 3000,group:'br' });
        emit('order_paid')
    })
};


const props = defineProps({
    order: {
        type: Order,
        required: true
    }
})

const order_unpaid_amount_str : any = computed(() => {
    return props.order.sale_price + " EGP"
})

const order_status : any = computed(() => {

if (props.order.state == "" || props.order.state == "pending" ){
    return {
        title:"PENDING",
        severity:"info"
    }
}

if (props.order.state == "in_progress" ){
    return {
        title:"INPROGRESS",
        severity:"success"
    }
}


if (props.order.state == "cancelled" ){
    return {
        title:"CANCELLED",
        severity:"danger"
    }
}

if (props.order.state == "finished" ){
    return {
        title:"FINISHED",
        severity:"success"
    }
}

return {}

})

</script>