<template>
    <div>
        <Message severity="secondary">
            <template #container>
                <div class="px-3 py-2 flex justify-content-between align-items-center">
                    <span>{{ props.order.display_id }}</span>
                    <div class="flex justify-content-center align-items-center">
                        <Badge :value="order_status.title" :severity="order_status.severity" />
                        <Badge :value="payment_status.title" :severity="payment_status.severity" class="mx-1" />
                    </div>
                    <span>{{ props.order.items.length }} Item(s)</span>
                    <Button icon="pi pi-book" severity="secondary" aria-label="Info" @click="emit('view-order-pressed')" />
                </div>
            </template>
        </Message>
    </div>
</template>

<script setup lang="ts">
import {defineProps,computed,defineEmits} from 'vue'
import Message from 'primevue/message';
import Button from 'primevue/button';
import Badge from 'primevue/badge';
import Order from '@/classes/Order'


const emit = defineEmits(['view-order-pressed'])

const props = defineProps({
    order: {
        type: Order,
        required: true
    }
})


const payment_status: any = computed(() => {
    if (props.order.is_paid){
        return {
            title:"PAID",
            severity:"success"
        }
    }

    return {
        title:"UNPAID",
        severity:"warn"
    }
})

const order_status : any = computed(() => {

if (props.order.state == "" || props.order.state == "pending" ){
    return {
        title:"PENDING",
        severity:"secondary"
    }
}


if (props.order.state == "cancelled" ){
    return {
        title:"CANCELLED",
        severity:"danger"
    }
}

if (props.order.state == "in_progress" ){
    return {
        title:"INPROGRESS",
        severity:"info"
    }
}

if (props.order.state == "finished" ){
    return {
        title:"FINISHED",
        severity:"success"
    }
}

if (props.order.state == "stashed" ){
    return {
        title:"STASHED",
        severity:"secondary"
    }
}

return {}

})

</script>