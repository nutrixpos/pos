<template>
    <div class="w-full">
        <div class="grid mx-2">
            <div class="col-12 flex align-items-center gap-2">
                <Button icon="pi pi-arrow-left" severity="secondary" text rounded aria-label="Back" @click="$router.push('/admin/orders')" />
                <h3 class="my-0">{{ $t('order') }}: {{ order?.display_id || $route.params.id }}</h3>
            </div>
            <div class="col-12" v-if="loading">
                <ProgressSpinner style="width: 35px; height: 35px;stroke:blue !important;" strokeWidth="6" fill="transparent" animationDuration=".5s" aria-label="Loading Order" />
            </div>
            <div class="col-12" v-else-if="error" :style="`direction: ${store.orientation == 'rtl' ? 'rtl' : 'ltr'}`">
                <div class="text-red-500">{{ error }}</div>
                <Button class="mt-2" icon="pi pi-refresh" severity="secondary" @click="loadOrder" />
            </div>
            <div class="col-12" v-else-if="order">
                <OrderView @updated="loadOrder" @finished="loadOrder" @cancelled="loadOrder" @amount_collected="loadOrder" :order="order" />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import axios from 'axios'
import Button from 'primevue/button'
import ProgressSpinner from 'primevue/progressspinner'
import OrderView from '@/components/OrderView.vue'
import Order from '@/classes/Order'
import { globalStore } from '@/stores'
import auth from '../services/auth'

const store = globalStore()
const route = useRoute()

const order = ref<Order | null>(null)
const loading = ref(true)
const error = ref("")

const loadOrder = () => {
    loading.value = true
    error.value = ""

    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders/${route.params.id}`, {
        headers: {
            Authorization: `Bearer ${auth.accessToken.value}`
        }
    })
    .then((result) => {
        order.value = result.data.data
    })
    .catch((err) => {
        error.value = err.response?.data || 'Failed to load order'
    })
    .finally(() => {
        loading.value = false
    })
}

loadOrder()
</script>