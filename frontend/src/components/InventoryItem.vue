<template>
    <div class="grid container mx-1 p-0" style="cursor: pointer;" @click="visible=true">
        <div class="col-9 flex justify-content-center align-items-center">
            Motzarilla
        </div>
        <div class="col-3 flex justify-content-center align-items-center quantity-container">
            2kg
        </div>
        <Dialog v-model:visible="visible" modal header="Edit Profile" :style="{ width: '25rem' }">
            <div>
                <InventoryItem class="my-2" />
            </div>
            <div class="flex justify-content-end gap-2">
                <Button type="button" label="Close" severity="secondary" @click="visible = false"></Button>
                <Button type="button" label="History" @click="visible = false"></Button>
            </div>
        </Dialog>
    </div>
</template>

<script setup>
import {ref} from 'vue'
import Button from 'primevue/button'
import Dialog from 'primevue/dialog'
import InventoryItem from '@/components/InventoryItem.vue'

const visible = ref(false)
</script>

<style scoped>
.container {
    border: 1px solid #ccc;
    border-radius: 5px;
    padding: 5px;
    margin:0px;
}
.quantity-container{
    border-top-right-radius:5px;
    border-bottom-right-radius:5px;
    background-color:bisque;
    height:3rem;
}
</style>