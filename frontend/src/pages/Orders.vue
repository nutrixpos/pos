<template>
       <div class="w-full">
            <div class="grid mx-2">
                <div class="col-12 flex">
                    <div class="gird w-full">
                        <div class="col-12">
                            <h3>{{ $t('order',3) }}</h3>
                        </div>
                        <div class="col-12">
                            <DataTable @page="updatOrdersTableRowsPerPage" :lazy="true" :totalRecords="ordersTableTotalRecords" :loading="isOrdersTableLoading"  paginatorPosition="both"  paginator :rows="ordersTableRowsPerPage" :rowsPerPageOptions="[50, 100, 500]" :value="orders" stripedRows tableStyle="min-width: 50rem;max-height:50vh;" class="w-full pr-2">
                                <template #header>
                                    <IconField iconPosition="left">
                                        <InputIcon>
                                            <i class="pi pi-search" />
                                        </InputIcon>
                                        <InputText :placeholder="$t('search')" v-model="ordersSearchText" @keyup.stop="(event) => loadOrders(0,100)"/>
                                    </IconField>
                                </template>
                                <Column field="display_id" :header="$t('id')"></Column>
                                <Column :header="$t('status')">
                                    <template #body="slotProps">
                                        <Tag :value="slotProps.data.state" :severity="orderStateSeverity[slotProps.data.state]" />
                                    </template>
                                </Column>
                                <Column field="submitted_at" :header="$t('submitted_at')"></Column>
                                <Column :header="$t('actions')">
                                    <template #body="slotProps">
                                        <ButtonGroup>
                                            <Button v-tooltip.top="'cancel'" icon="pi pi-times" severity="secondary" aria-label="Remove" @click="confirmCancelOrder($event,slotProps.data.display_id,slotProps.data.id)" />
                                            <Button v-tooltip.top="'details'" icon="pi pi-book" severity="secondary" aria-label="Remove" @click="orderToDisplay = slotProps.data; order_details_dialog = true" />
                                        </ButtonGroup>
                                    </template>
                                </Column>
                            </DataTable>
                            <Dialog v-model:visible="order_details_dialog" modal :header="`Order:  (${orderToDisplay?.display_id}) details`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '90vw', '575px': '90vw' }">
                                <OrderView @updated="loadOrders(); order_details_dialog=false" @cancelled="cancelOrderDisplayed()" @finished="finishOrderDisplayed()" @amount_collected="orderToShowAmountCollected()" :order="orderToDisplay" @order-cancelled="order_details_dialog=false" />
                            </Dialog>
                        </div>
                    </div>
                </div>
            </div>
            <ConfirmPopup></ConfirmPopup>
       </div>
</template>

<script setup lang="ts">
import IconField from 'primevue/iconfield';
import InputIcon from 'primevue/inputicon';
import InputText from 'primevue/inputtext';
import ConfirmPopup from 'primevue/confirmpopup';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import {ref, getCurrentInstance} from 'vue'
import Tag from 'primevue/tag';
import axios from 'axios'
import ButtonGroup from 'primevue/buttongroup';
import Button from 'primevue/button';
import { useConfirm } from "primevue/useconfirm";
import Dialog from 'primevue/dialog';
import { useToast } from "primevue/usetoast";
import OrderView from '@/components/OrderView.vue';


const { proxy } = getCurrentInstance();
const confirm = useConfirm();
const toast = useToast()

const ordersTableTotalRecords = ref(0)
const isOrdersTableLoading = ref(false)
const ordersTableRowsPerPage = ref(50)
const orders = ref<any[]>([])
const ordersSearchText = ref("")

const orderToDisplay = ref<any>({})
const order_details_dialog = ref(false)

const orderStateSeverity = ref<any>({
    "pending": "secondary",
    "in_progress": "info",
    "finished": "success",
    "cancelled":"danger"
})

const updatOrdersTableRowsPerPage = (event: any) => {
    loadOrders(event.first,event.rows)
}

const finishOrderDisplayed = () => {
    if (orderToDisplay.value){
        orderToDisplay.value.state = "finished"
    }
}

const cancelOrderDisplayed = () => {
    if (orderToDisplay.value){
        orderToDisplay.value.state = "cancelled"
    }
}

const orderToShowAmountCollected = () => {

    if (orderToDisplay.value){
        orderToDisplay.value.is_paid = true
    }

    loadOrders()
}

const cancelOrder = (order_id: string) => {

    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders/${order_id}/cancel`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(()=>{
        toast.add({ severity: 'success', summary: 'Success', detail: 'Order cancelled successfully',group:'br' });
        orders.value = orders.value.map((order) => {
            if (order.id != order_id)
                return order
            
            return {
                ...order,
                state:"cancelled"
            }
        })
    })
    .catch(() => {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to cancel order',group:'br' });
    });

}

const confirmCancelOrder = (event,display_id,order_id) => {
    confirm.require({
        target: event.currentTarget,
        message: `Are you sure you want to cancel order ${display_id} ?`,
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


            cancelOrder(order_id)

         
        },
        reject: () => {
        }
    });
  }

const loadOrders =  (first=0,rows=100) => {

        let page_number = Math.floor(first/rows) + 1

    if (page_number == 0) {
        page_number = 1
    }


    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders?page[number]=${page_number}&page[size]=${rows}&filter[display_id]=${ordersSearchText.value}`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then((result)=>{
        orders.value = result.data.data
        ordersTableTotalRecords.value = result.data.meta.total_records
    })
};

loadOrders()

</script>