<template>
    <div class="flex flex-column">
        <div class="flex w-full flex-column">
            Money To refund
            <div class="flex w-full justify-content-start align-items-center">
                <Slider v-model="refund_value" :max="props.item.price" class="w-9 m-1 mx-2" style="height:0.6rem;" />
                <p class="ml-2" style="font-size:0.8rem">{{ refund_value.toFixed(2) }}</p>
            </div>
        </div>
        <div class="flex gap-4">
            {{ props.item.product.name }} Destination: 
            <div class="flex flex-wrap gap-4">
                <div class="flex items-center gap-2">
                    <RadioButton v-model="item_destination" inputId="option2" name="item_destination" value="inventory" />
                    <label for="option2">Inventory</label>
                </div>
                <div class="flex items-center gap-2">
                    <RadioButton v-model="item_destination" inputId="option1" name="item_destination" value="dispose" />
                    <label for="option1">Dispose</label>
                </div>
                <div class="flex items-center gap-2">
                    <RadioButton v-model="item_destination" inputId="option3" name="item_destination" value="waste" />
                    <label for="option3">Waste</label>
                </div>
                <div class="flex items-center gap-2">
                    <RadioButton v-model="item_destination" inputId="option4" name="item_destination" value="custom" />
                    <label for="option4">Custom (materials / products)</label>
                </div>
            </div>
        </div>
        <Divider />
        <Textarea v-model="reason" class="w-4" :placeholder="$t('reason')"/>
        <h3 v-if="item_destination == 'custom'">Materials</h3>
        <div  v-if="item_destination == 'custom'" class="flex flex-column">
            <div v-for="(material,index) in material_refunds" :key="index" class="flex flex-column gap-2">
                <div class="flex gap-3 align-items-center">
                    <h4>{{ material.material.name }}</h4>
                    <h4>{{ material.quantity }}</h4>
                    <h4>{{ material.material.unit }}</h4>
                </div>
                <div class="flex align-items-center gap-2">
                    <label for="inventory_return_quantity">{{$t('inventory_return')}}</label>
                    <InputText v-model.number="material_refunds[index]['inventory_return_qty']" type="number" aria-describedby="inventory_return_qty" />
                </div>
                <div class="flex align-items-center gap-2">
                    <label for="dispose">{{$t('dispose')}}</label>
                    <InputText v-model.number="material_refunds[index]['dispose_qty']" type="number" aria-describedby="dispose_qty" />
                </div>
                <div class="flex align-items-center gap-2">
                    <label for="waste">{{$t('waste')}}</label>
                    <InputText v-model.number="material_refunds[index]['waste_qty']" type="number" aria-describedby="waste_qty" />
                </div>
            </div>
        </div>
        <Divider v-if="item_destination == 'custom'" />
        <div class="flex align-items-center gap-2" v-if="item_destination == 'custom'">
            <h3>Products</h3>
            <Button icon="pi pi-plus" size="small" severity="secondary" @click="pick_product_dialog=true"/>
        </div>
        <div class="flex flex-column gap-3">
            <div class="flex gap-4 align-items-center" v-for="(product,index) in products_to_add" :key="index">
                {{ product.name }}
                <div class="flex align-items-center gap-2">
                    <InputText v-model.number="added_products_quantities[index]" type="number" aria-describedby="added_product_quantity" />
                </div>
                <Button icon="pi pi-times" severity="secondary" @click="products_to_add.splice(index,1); added_products_quantities.splice(index,1)" />
            </div>
        </div>

        <Divider v-if="item_destination == 'custom'" />
        
        <div class="flex justify-content-end">
            <Button :label="$t('submit')" class="w-3" @click="submitRefund" />
        </div>
    </div>
    <Dialog v-model:visible="pick_product_dialog">
        <PickProduct @returnProduct="(product) => products_to_add.push(product)" />
    </Dialog>
</template>

<script setup lang="ts">
import { defineProps,ref,defineEmits,getCurrentInstance } from 'vue';
import { OrderItem, Product } from '@/classes/OrderItem';
import PickProduct from './PickProduct.vue';
import { Slider,RadioButton,Divider,InputText,Button,Dialog,Textarea } from 'primevue';
import axios from 'axios';
import Order from '@/classes/Order.ts'
import { useToast } from "primevue/usetoast";

const toast = useToast()

const { proxy } = getCurrentInstance();

const refund_value = ref(0)

const reason = ref("")

const material_refunds = ref([])
const pick_product_dialog = ref(false)

const products_to_add = ref<Array<Product>>([])
const added_products_quantities = ref([])

const item_destination = ref("inventory")

const props = defineProps({
    item: {
        type: OrderItem,
        required: true
    },
    order: {
        type: Order,
        required: true
    }
})

const emit = defineEmits(['finished'])

const init = () => {
    refund_value.value = props.item.price
    props.item.materials.forEach((material) => {
        material_refunds.value.push({
            material_id: material.material.id,
            entry_id: material.entry.id,
            inventory_return_qty:0,
            waste_qty:0,
            dispose_qty:0,
            quantity: material.quantity,
            material: material.material,
        })
    })
}

const submitRefund = () => {

    const payload = {
        order_id: props.order.id,
        product_id: props.item.product.id,
        refund_value: refund_value.value,
        destination: item_destination.value,
        material_refunds: material_refunds.value
    }

    axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders/${props.order.id}/items/${props.item.id}/refund?reason=${reason.value}`, 
    {
        data : payload
    },
    {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(()=>{
        toast.add({ severity: 'success', summary: 'Success', detail: 'Item refunded',group:'br' });
        emit('finished')
    })
    .catch((err) => {
        toast.add({ severity: 'error', summary: 'Error', detail: err.response.data,group:'br' });
    });
}



init()
</script>