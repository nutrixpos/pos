<template>
    <div>
        <div class="grid w-12">

            <div class="col-12">
                
            </div>

            <div class="flex col-12">
                <Badge :value="order_status.title" :severity="order_status.severity" />
                <Badge :value="payment_status.title" :severity="payment_status.severity" class="mx-1" />
            </div>

            <div class="col-1">{{$t('item',3)}}</div>
            <div class="col-11">
                <OrderItemsInfo :order="order" @updated="emit('updated')" />
            </div>

            <div class="col-6">{{ $t('display_id') }}</div>
            <div class="col-6">{{ props.order.display_id }}</div>

            <div class="col-6">{{ $t('status') }}</div>
            <Badge :value="order_status.title" :severity="order_status.severity" />

            <div class="col-6">{{ $t('submitted_at') }}</div>
            <div class="col-6">{{ props.order.submitted_at }}</div>

            <div class="col-6">{{$t('started_at')}}</div>
            <div class="col-6">{{ props.order.started_at }}</div>

            <div class="col-6">{{ $t('comment') }}</div>
            <div class="col-6">{{ props.order.comment }}</div>

            <div class="col-6">{{ $t('cost') }}</div>
            <div class="col-6">{{ props.order.cost }}</div>

            <div class="col-6">{{ $t('sale_price') }}</div>
            <div class="col-6">{{ props.order.sale_price }}</div>

            <div class="col-6">{{$t('pay_later')}}?</div>
            <div class="col-6">{{ props.order.is_pay_later }}</div>

            <div class="col-6">{{$t('paid')}}</div>
            <div class="col-6">{{ props.order.is_paid }}</div>

            <div class="col-6">{{$t('discount')}}</div>
            <div class="col-6">{{ props.order.discount }}</div>

            <div class="col-12 flex flex-column">
                <h4>Actions</h4>
                <div class="flex gap-3">
                    <ButtonGroup class="flex">
                        <Button icon="fa fa-print" severity="secondary" :label="$t('client_receipt')" @click="PrintClientReceipt()" />
                        <Button icon="fa fa-print" severity="secondary" :label="$t('kitchen_receipt')" @click="PrintKitchenReceipt()" />
                        <Button v-if="!props.order.is_paid" icon="fa fa-hand-holding-dollar" severity="secondary" :label="$t('collect_money')" @click="collectedMoney()"/>
                        <Button v-if="props.order.state.toUpperCase() != 'CANCELLED' && props.order.state.toUpperCase() != 'FINISHED'" icon="fa fa-check" severity="secondary" :label="$t('finish')" @click="finishOrder()"/>
                        <Button v-if="props.order.state.toUpperCase() != 'CANCELLED' && props.order.state.toUpperCase() != 'FINISHED'" severity="secondary" size="small" aria-label="Cancel order" @click.stop="confirmCancelOrder($event)">
                            {{$t('cancel')}} {{ $t('order') }}
                        </Button>
                        <ConfirmPopup></ConfirmPopup>
                    </ButtonGroup>
                    <ButtonGroup class="flex">
                        <Button icon="fa fa-book" severity="secondary" :label="$t('logs')" @click="getOrderLogs()" />
                    </ButtonGroup>
                </div>

            </div>
            <Dialog v-model:visible="order_logs_dialog" modal :header="`#${props.order.display_id} logs`" class="xs:w-12 md:w-10 lg:w-8">
                <pre>
                    {{order_logs}}
                </pre>
            </Dialog>
        </div>
    </div>
</template>

<script setup lang="ts">
import {defineProps,getCurrentInstance,computed,defineEmits,ref} from 'vue'
import Badge from 'primevue/badge';
import { ButtonGroup } from 'primevue';
import Button from 'primevue/button';
import { useConfirm } from "primevue/useconfirm";
import ConfirmPopup from 'primevue/confirmpopup';
import axios from 'axios'
import { useToast } from "primevue/usetoast";
import OrderItemsInfo from './OrderItemsInfo.vue';
import Order from '@/classes/Order';
import Dialog from 'primevue/dialog'

const toast = useToast()

const order_logs = ref([])
const order_logs_dialog = ref(false)


const confirm = useConfirm();
const { proxy } = getCurrentInstance();
const emit = defineEmits(['amount_collected','finished','updated','cancelled'])


const getOrderLogs = () => {
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders/${props.order.id}/logs`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then((response)=>{
        order_logs.value = response.data.data
        order_logs_dialog.value = true
    })
    .catch((err) => {
        toast.add({ severity: 'error', summary: 'Error '+err.response.status, detail: err.response.data,group:'br' });
    });
}

const finishOrder = () => {
    axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders/${props.order.id}/finish`,{}, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`,
        }
    })
    .then(()=>{
        toast.add({ severity: 'success', summary: 'Success', detail: 'Order finished',group:'br' });
        emit('finished')
    })
    .catch(() => {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to perform the request',group:'br' });
    });
}

const collectedMoney = () => {
    axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders/${props.order.id}/pay`,{}, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`,
        }
    })
    .then(()=>{
        toast.add({ severity: 'success', summary: 'Success', detail: 'Money collected',group:'br' });
        emit('amount_collected')
    })
    .catch(() => {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to perform the request',group:'br' });
    });
}

const PrintKitchenReceipt = () => {
    axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders/${props.order.id}/printkitchenreceipt`,{}, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`,
            "Accept-Language": proxy.$i18n.locale
        }
    })
    .then(()=>{
        toast.add({ severity: 'success', summary: 'Success', detail: 'Print signal sent',group:'br' });
    })
    .catch(() => {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to perform the request',group:'br' });
    });
}

const PrintClientReceipt = () => {
    axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders/${props.order.id}/printclientreceipt`,{}, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`,
            "Accept-Language": proxy.$i18n.locale
        }
    })
    .then(()=>{
        toast.add({ severity: 'success', summary: 'Success', detail: 'Print signal sent',group:'br' });
    })
    .catch(() => {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to perform the request',group:'br' });
    });
}

const confirmCancelOrder = (event) => {
    confirm.require({
        target: event.currentTarget,
        message: 'Are you sure you want to cancel this order ?',
        icon: 'pi pi-exclamation-triangle',
        rejectProps: {
            label: 'Cancel',
            severity: 'secondary',
            outlined: true
        },
        acceptProps: {
            label: 'Yes'
        },
        accept: () => {
            
            axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders/${props.order.id}/cancel`, {},{
                headers: {
                    Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
                }
            })
            .then(()=>{
                toast.add({ severity: 'success', summary: 'Success', detail: 'Order cancelled successfully',group:'br' });
                emit('cancelled')
            })
            .catch(() => {
                toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to cancel order',group:'br' });
            });

            
        },
        reject: () => {
            
        }
    });
  }

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

return {}

})


</script>