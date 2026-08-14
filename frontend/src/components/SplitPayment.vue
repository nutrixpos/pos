<template>
    <div class="flex flex-column gap-2 w-full">
        <div v-for="(payment, index) in model" :key="index" class="flex align-items-center gap-2 w-full">
            <Select v-model="model[index].source" :options="availableOptions(index)" optionLabel="name" optionValue="name" :placeholder="$t('payment_source')" class="w-6" />
            <InputText type="number" :placeholder="`0.00 ${$t('egp')}`" v-model.number="model[index].amount" class="w-3" :invalid="!(payment.amount > 0)" />
            <Button icon="pi pi-times" severity="secondary" aria-label="Remove" @click="removePayment(index)" />
        </div>
        <div class="flex align-items-center justify-content-between w-full">
            <Button :label="$t('add_payment_method')" severity="secondary" icon="pi pi-plus" size="small" :disabled="unusedSources.length === 0" @click="addPayment()" />
            <Button :label="$t('fill_remaining')" severity="secondary" size="small" icon="pi pi-fill" :disabled="remaining <= 0 || model.length === 0" @click="fillRemaining()" />
        </div>
        <div class="flex justify-content-between w-full">
            <span>{{ $t('remaining') }}:</span>
            <strong :style="`color: ${remainingValid ? 'inherit' : 'red'}`">{{ remaining.toFixed(2) }} {{ $t('egp') }}</strong>
        </div>
    </div>
</template>

<script setup lang="ts">
import { defineModel, computed, watch } from 'vue'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import { Select } from 'primevue'
import type { OrderPayment } from '@/classes/Order'

const props = defineProps({
    total: {
        type: Number,
        required: true
    },
    payment_sources: {
        type: Array as () => any[],
        required: true
    }
})

const emit = defineEmits(['validity-change'])

const model = defineModel<OrderPayment[]>({ required: true })

const usedSources = (excludeIndex: number): Set<string> => {
    const used = new Set<string>()
    model.value.forEach((payment, index) => {
        if (index !== excludeIndex && payment.source) {
            used.add(payment.source)
        }
    })
    return used
}

const unusedSources = computed(() => {
    const used = usedSources(-1)
    return props.payment_sources.filter((source) => !used.has(source.name))
})

const availableOptions = (index: number) => {
    const used = usedSources(index)
    return props.payment_sources.filter((source) => !used.has(source.name))
}

const totalAmount = computed(() => {
    return model.value.reduce((sum, payment) => sum + (payment.amount || 0), 0)
})

const remaining = computed(() => {
    return props.total - totalAmount.value
})

const remainingValid = computed(() => {
    return Math.abs(remaining.value) <= 0.01
})

const isValid = computed(() => {
    if (model.value.length === 0) return false
    for (const payment of model.value) {
        if (payment.source === "" || payment.amount <= 0) return false
    }
    return remainingValid.value
})

const addPayment = () => {
    if (unusedSources.value.length === 0) return
    model.value.push({
        source: unusedSources.value[0].name,
        amount: 0
    })
}

const removePayment = (index: number) => {
    model.value.splice(index, 1)
}

const fillRemaining = () => {
    if (model.value.length === 0) return
    const last = model.value[model.value.length - 1]
    last.amount = Number(((last.amount || 0) + remaining.value).toFixed(2))
}

watch([model, remainingValid], () => {
    emit('validity-change', isValid.value)
}, { deep: true })

emit('validity-change', isValid.value)
</script>
