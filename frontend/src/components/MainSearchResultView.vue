<template>
    <div>
        <Message severity="secondary">
            <template #container>
                <div class="px-3 py-2 flex justify-content-between align-items-center">
                    <span>{{ props.order.display_id }}</span>
                    <Badge :value="order_status.title" :severity="order_status.severity" />
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

const order_status : any = computed(() => {

if (props.order.state == "" || props.order.state == "pending" ){
    return {
        title:"PENDING",
        severity:"info"
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