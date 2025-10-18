<template>
    <div>
        <Message :severity="notification.severity" @close="emit('closed')" class="my-1">
            <div class="flex flex-column">
                <h4 class="m-0 mb-1">{{ props.notification.topic_name }}</h4>
                {{ props.notification.description }}
                <span style="color:gray;font-size:0.8rem" class="mt-1">{{ timePassed }}</span>
            </div>
        </Message>
    </div>
</template>

<script setup lang="ts">
import {defineProps, defineEmits,ref} from 'vue'
import Message from 'primevue/message'
import moment from 'moment';
import { Notification } from '@/classes/Notification'

const timePassed = ref()

const emit = defineEmits(['closed'])


const props = defineProps({
    notification: {
        type: Notification,
        required: true
    }
})

const updateElapsedTime = () => {
    const moment_date = moment(props.notification.date);
    timePassed.value =  moment_date.fromNow()
    // moment(String(props.order.submitted_at)).fromNow()
    setInterval(function(){
        timePassed.value = moment_date.fromNow()
    },10000)
}

const init = () => {
    updateElapsedTime()
}

init()

</script>