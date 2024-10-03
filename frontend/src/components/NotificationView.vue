<template>
    <div>
        <Message :severity="notification.severity">
            <div>
                <h4 class="m-0 mb-1">{{ notification.topic_name }}</h4>
                {{ notification.description }}
            </div>
        </Message>
    </div>
</template>

<script setup lang="ts">
import {defineProps} from 'vue'
import Message from 'primevue/message'

import { Notification } from '@/classes/Notification'

defineProps({
    notification: {
        type: Notification,
        required: true
    }
})


</script>