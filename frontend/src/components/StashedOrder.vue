<template>
    <div>
        <Message severity="info">
            <template #container>
                <div class="flex justify-content-between align-items-center p-2">
                    <div class="flex justify-content-start align-items-center">
                        <i class="pi pi-clock"></i>
                        <div class="flex ml-2">
                            <h4 class="m-0 mb-1">{{ props.order.display_id }}</h4>
                            &nbsp;({{ props.order.items.length }} {{$t('item',3)}})
                        </div>
                    </div>
                    <div class="flex align-items-center">
                        <Tag :value="order_status.title" :severity="order_status.severity" />
                        <Button icon="pi pi-arrow-down-right" severity="secondary"  class="ml-2" size="small" @click="emit('back_to_checkout')" />
                        <Button icon="pi pi-times" severity="secondary"  class="ml-2" size="small" />
                    </div>
                </div>
            </template>
        </Message>
    </div>
</template>

<script setup lang="ts">
import {computed,defineProps,defineEmits} from 'vue'
import Message from 'primevue/message'
import Button from 'primevue/button'
import Tag from 'primevue/tag'

import Order from '@/classes/Order'


const emit = defineEmits(['back_to_checkout'])

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

    if (props.order.state == "stashed" ){
        return {
            title:"STASHED",
            severity:"warning"
        }
    }

    return {}

})


</script>