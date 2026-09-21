<template>
    <div class="w-full">
        <div class="grid mx-2">
            <div class="col-12 flex justify-content-between align-items-center">
                <h3>{{ $t('purchase_order', 3) }}</h3>
                <Button icon="pi pi-plus" :label="$t('new_purchase_order')" @click="showNewPODialog" rounded raised />
            </div>
            <div class="col-12">
                <TabView v-model:activeIndex="active_tab">
                    <TabPanel value="0" :header="$t('purchase_order', 3)">
                        <DataTable :value="purchase_orders" stripedRows class="w-full" :loading="is_po_loading">
                            <Column field="display_id" :header="$t('display_id')"></Column>
                            <Column field="supplier" :header="$t('supplier')"></Column>
                            <Column field="total" :header="$t('total')">
                                <template #body="slotProps">
                                    {{ slotProps.data.total?.toFixed(2) }}
                                </template>
                            </Column>
                            <Column :header="$t('status')">
                                <template #body="slotProps">
                                    <Tag :value="$t(statusLabel(slotProps.data.status))" :severity="statusSeverity(slotProps.data.status)" />
                                </template>
                            </Column>
                            <Column field="created_at" :header="$t('date')">
                                <template #body="slotProps">
                                    {{ formatDate(slotProps.data.created_at) }}
                                </template>
                            </Column>
                            <Column :header="$t('actions')" style="width: 22rem">
                                <template #body="slotProps">
                                    <ButtonGroup>
                                        <Button v-if="slotProps.data.status === 'open' || slotProps.data.status === 'partially_received'" icon="pi pi-arrow-down" :label="$t('receive')" severity="success" @click="openReceiveDialog(slotProps.data)" />
                                        <Button icon="pi pi-eye" :label="$t('grn')" severity="secondary" @click="showGRNsForPO(slotProps.data)" />
                                        <Button v-if="slotProps.data.status === 'open'" icon="pi pi-times" :label="$t('delete')" severity="danger" @click="cancelPurchaseOrder(slotProps.data)" />
                                    </ButtonGroup>
                                </template>
                            </Column>
                        </DataTable>
                    </TabPanel>
                    <TabPanel value="1" :header="$t('grn', 3)">
                        <DataTable :value="grns" stripedRows class="w-full" :loading="is_grn_loading">
                            <Column field="display_id" :header="$t('grn_display_id')"></Column>
                            <Column field="purchase_order_display_id" :header="$t('purchase_order')"></Column>
                            <Column field="supplier" :header="$t('supplier')"></Column>
                            <Column field="received_at" :header="$t('date')">
                                <template #body="slotProps">
                                    {{ formatDate(slotProps.data.received_at) }}
                                </template>
                            </Column>
                            <Column :header="$t('status')">
                                <template #body>
                                    <Tag :value="$t('received')" severity="success" />
                                </template>
                            </Column>
                            <Column :header="$t('actions')" style="width: 10rem">
                                <template #body="slotProps">
                                    <Button icon="pi pi-eye" :label="$t('view')" severity="secondary" @click="showGRNDetail(slotProps.data)" />
                                </template>
                            </Column>
                        </DataTable>
                    </TabPanel>
                </TabView>
            </div>
        </div>

        <Dialog v-model:visible="new_po_dialog" modal :header="$t('new_purchase_order')" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '90vw', '575px': '90vw' }">
            <div class="grid p-2">
                <div class="col-12 md:col-4 flex flex-column gap-2">
                    <label for="supplier">{{ $t('supplier') }}</label>
                    <InputText id="supplier" v-model="new_po_supplier" />
                </div>
                <div class="col-12 md:col-4 flex flex-column gap-2">
                    <label for="notes">{{ $t('comment') }}</label>
                    <InputText id="notes" v-model="new_po_notes" />
                </div>
                <div class="col-12 md:col-4 flex align-items-end">
                    <div class="flex align-items-center">
                        <Checkbox v-model="new_po_auto_receive" inputId="auto_receive" :binary="true" />
                        <label for="auto_receive" class="ml-2">{{ $t('auto_receive') }}</label>
                    </div>
                </div>
            </div>

            <div class="grid p-2">
                <div class="col-12">
                    <h4>{{ $t('items') }}</h4>
                </div>
                <div class="col-12">
                    <DataTable :value="new_po_items" stripedRows class="w-full">
                        <Column :header="$t('material')" style="width: 14rem">
                            <template #body="slotProps">
                                <Dropdown v-model="slotProps.data.material_id" :options="material_options" optionLabel="name" optionValue="id" :placeholder="$t('choose')" filter class="w-full" @change="onMaterialSelect(slotProps.data, $event)" />
                            </template>
                        </Column>
                        <Column :header="$t('unit')" style="width: 5rem">
                            <template #body="slotProps">
                                {{ slotProps.data.unit }}
                            </template>
                        </Column>
                        <Column :header="$t('company')" style="width: 8rem">
                            <template #body="slotProps">
                                <InputText v-model="slotProps.data.company" class="w-full" />
                            </template>
                        </Column>
                        <Column :header="$t('sku')" style="width: 6rem">
                            <template #body="slotProps">
                                <InputText v-model="slotProps.data.sku" class="w-full" />
                            </template>
                        </Column>
                        <Column :header="$t('quantity')" style="width: 5rem">
                            <template #body="slotProps">
                                <InputNumber v-model="slotProps.data.quantity" :min="0" mode="decimal" :minFractionDigits="0" :maxFractionDigits="2" class="w-full" />
                            </template>
                        </Column>
                        <Column :header="$t('total_price')" style="width: 8rem">
                            <template #body="slotProps">
                                <InputNumber v-model="slotProps.data.total_price" :min="0" mode="decimal" :minFractionDigits="2" :maxFractionDigits="2" class="w-full" />
                                <small v-if="slotProps.data.quantity > 0" class="block text-color-secondary">
                                    {{ slotProps.data.total_price > 0 ? $t('unit_price') + ': ' + (slotProps.data.total_price / slotProps.data.quantity).toFixed(2) : '' }}
                                </small>
                            </template>
                        </Column>
                        <Column :header="$t('expiration_date')" style="width: 12rem">
                            <template #body="slotProps">
                                <Calendar v-model="slotProps.data.expiration_date" showIcon class="w-full" />
                            </template>
                        </Column>
                        <Column :header="$t('actions')" style="width: 4rem" class="text-center">
                            <template #body="slotProps">
                                <Button icon="pi pi-times" severity="danger" text rounded @click="new_po_items.splice(new_po_items.findIndex(el => el === slotProps.data), 1)" />
                            </template>
                        </Column>
                    </DataTable>
                </div>
                <div class="col-12">
                    <Button icon="pi pi-plus" :label="$t('add_item')" class="my-1" severity="info" @click="addNewPOItem" />
                </div>
            </div>

            <template #footer>
                <ButtonGroup>
                    <Button :label="$t('cancel')" @click="new_po_dialog = false" severity="secondary" />
                    <Button class="ml-2" severity="primary" :label="$t('submit')" :disabled="is_submitting_po" @click="submitNewPurchaseOrder" />
                </ButtonGroup>
            </template>
        </Dialog>

        <Dialog v-model:visible="receive_dialog" modal :header="`${$t('receive')} - ${receive_po?.display_id}`" :style="{ width: '75rem' }">
            <DataTable :value="receive_items" stripedRows tableStyle="min-width: 50rem" class="w-full">
                <Column :header="$t('receive')" style="width: 6rem">
                    <template #body="slotProps">
                        <Checkbox v-model="slotProps.data.selected" :binary="true" />
                    </template>
                </Column>
                        <Column field="material_name" :header="$t('material')"></Column>
                        <Column field="unit" :header="$t('unit')"></Column>
                        <Column field="company" :header="$t('company')"></Column>
                        <Column field="quantity" :header="$t('quantity')">
                            <template #body="slotProps">
                                {{ slotProps.data.quantity?.toFixed(2) }}
                            </template>
                        </Column>
                        <Column :header="$t('remaining')">
                    <template #body="slotProps">
                        {{ slotProps.data.remaining?.toFixed(2) }}
                    </template>
                </Column>
                <Column :header="$t('received_quantity')" style="min-width: 10rem">
                    <template #body="slotProps">
                        <InputNumber v-model="slotProps.data.received_quantity" :min="0" :max="slotProps.data.remaining" mode="decimal" :minFractionDigits="0" :maxFractionDigits="2" class="w-full" :disabled="!slotProps.data.selected" />
                    </template>
                </Column>
            </DataTable>
            <template #footer>
                <ButtonGroup>
                    <Button :label="$t('cancel')" @click="receive_dialog = false" severity="secondary" />
                    <Button class="ml-2" severity="success" :label="$t('receive')" :disabled="is_submitting_receive" @click="submitReceive" />
                </ButtonGroup>
            </template>
        </Dialog>

        <Dialog v-model:visible="grn_detail_dialog" modal :header="grn_detail?.display_id" :style="{ width: '75rem' }">
            <div class="grid p-2">
                <div class="col-6">
                    <span>{{ $t('purchase_order') }}: <b>{{ grn_detail?.purchase_order_display_id }}</b></span>
                </div>
                <div class="col-6">
                    <span>{{ $t('supplier') }}: <b>{{ grn_detail?.supplier }}</b></span>
                </div>
                <div class="col-12">
                    <DataTable :value="grn_detail?.items" stripedRows tableStyle="min-width: 50rem" class="w-full">
                        <Column field="material_name" :header="$t('material')"></Column>
                        <Column field="unit" :header="$t('unit')"></Column>
                        <Column field="company" :header="$t('company')"></Column>
                        <Column field="quantity" :header="$t('quantity')"></Column>
                        <Column :header="$t('total_price')">
                            <template #body="slotProps">
                                {{ (slotProps.data.quantity * slotProps.data.purchase_price)?.toFixed(2) }}
                            </template>
                        </Column>
                        <Column field="entry_id" :header="$t('entry_id')"></Column>
                    </DataTable>
                </div>
            </div>
        </Dialog>
    </div>
</template>

<script setup lang="ts">
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import axios from 'axios'
import Button from 'primevue/button'
import ButtonGroup from 'primevue/buttongroup'
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Dropdown from 'primevue/dropdown'
import Checkbox from 'primevue/checkbox'
import Calendar from 'primevue/calendar'
import Tag from 'primevue/tag'
import TabView from 'primevue/tabview';
import TabPanel from 'primevue/tabpanel';
import { ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useToast } from "primevue/usetoast";
import { PurchaseOrder, PurchaseOrderItem, GRN } from '@/classes/PurchaseOrder';
import { Material } from '@/classes/OrderItem';
import auth from '@/services/auth';

const toast = useToast();
const route = useRoute();
const router = useRouter();

const active_tab = ref(0)

const purchase_orders = ref<PurchaseOrder[]>([])
const grns = ref<GRN[]>([])
const material_options = ref<Material[]>([])
const is_po_loading = ref(false)
const is_grn_loading = ref(false)
const is_submitting_po = ref(false)
const is_submitting_receive = ref(false)

const new_po_dialog = ref(false)
const new_po_supplier = ref("")
const new_po_notes = ref("")
const new_po_auto_receive = ref(true)
const new_po_items = ref<any[]>([])

const receive_dialog = ref(false)
const receive_po = ref<PurchaseOrder>()
const receive_items = ref<any[]>([])

const grn_detail_dialog = ref(false)
const grn_detail = ref<GRN>()

const statusSeverity = (status: string) => {
    switch (status) {
        case 'open': return 'info'
        case 'partially_received': return 'warn'
        case 'received': return 'success'
        case 'cancelled': return 'danger'
        default: return 'info'
    }
}

const statusLabel = (status: string) => {
    switch (status) {
        case 'open': return 'open'
        case 'partially_received': return 'partially_received'
        case 'received': return 'received'
        case 'cancelled': return 'cancelled'
        default: return 'open'
    }
}

const formatDate = (date: string) => {
    if (!date) return ""
    return new Date(date).toLocaleString()
}

const addNewPOItem = () => {
    const item = new PurchaseOrderItem()
    new_po_items.value.push(item)
}

const onMaterialSelect = (item: any, event: any) => {
    const material = material_options.value.find(m => m.id === event.value)
    if (material) {
        item.material_name = material.name
        item.unit = material.unit
    }
}

const showNewPODialog = () => {
    new_po_supplier.value = ""
    new_po_notes.value = ""
    new_po_auto_receive.value = true
    new_po_items.value = []
    if (new_po_items.value.length === 0) {
        addNewPOItem()
    }
    new_po_dialog.value = true
}

const submitNewPurchaseOrder = () => {
    const items = new_po_items.value
        .filter(i => i.material_id && i.quantity > 0)
        .map(i => ({
            material_id: i.material_id,
            company: i.company,
            sku: i.sku,
            quantity: i.quantity,
            purchase_price: i.quantity > 0 ? i.total_price / i.quantity : 0,
            expiration_date: i.expiration_date
        }))

    if (items.length === 0) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Purchase order must have at least one item', life: 3000, group: 'br' })
        return
    }

    is_submitting_po.value = true

    axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/purchase-orders`, {
        data: {
            supplier: new_po_supplier.value,
            notes: new_po_notes.value,
            auto_receive: new_po_auto_receive.value,
            items
        }
    }, {
        headers: {
            Authorization: `Bearer ${auth.accessToken.value}`
        }
    })
    .then(() => {
        toast.add({ severity: 'success', summary: 'Success', detail: 'Purchase order saved successfully', life: 3000, group: 'br' })
        new_po_dialog.value = false
        loadPurchaseOrders()
        loadGRNs()
    })
    .catch((error) => {
        toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data || error.message, life: 3000, group: 'br' })
    })
    .finally(() => {
        is_submitting_po.value = false
    })
}

const openReceiveDialog = (po: PurchaseOrder) => {
    receive_po.value = po
    receive_items.value = po.items
        .filter(i => (i.quantity - i.received_quantity) > 0)
        .map(i => ({
            item_id: i.item_id,
            material_name: i.material_name,
            unit: i.unit,
            company: i.company,
            quantity: i.quantity,
            remaining: i.quantity - i.received_quantity,
            received_quantity: i.quantity - i.received_quantity,
            selected: true
        }))
    receive_dialog.value = true
}

const submitReceive = () => {
    if (!receive_po.value) return

    const items = receive_items.value
        .filter(i => i.selected && i.received_quantity > 0)
        .map(i => ({
            item_id: i.item_id,
            quantity: i.received_quantity
        }))

    if (items.length === 0) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Nothing to receive', life: 3000, group: 'br' })
        return
    }

    const purchase_order_id = receive_po.value.id

    is_submitting_receive.value = true

    axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/purchase-orders/${purchase_order_id}/receive`, {
        data: items
    }, {
        headers: {
            Authorization: `Bearer ${auth.accessToken.value}`
        }
    })
    .then((response) => {
        toast.add({ severity: 'success', summary: 'Success', detail: `${response.data.data.display_id} created`, life: 3000, group: 'br' })
        receive_dialog.value = false
        loadPurchaseOrders()
        loadGRNs()
    })
    .catch((error) => {
        toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data || error.message, life: 3000, group: 'br' })
    })
    .finally(() => {
        is_submitting_receive.value = false
    })
}

const cancelPurchaseOrder = (po: PurchaseOrder) => {
    axios.delete(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/purchase-orders/${po.id}`, {
        headers: {
            Authorization: `Bearer ${auth.accessToken.value}`
        }
    })
    .then(() => {
        toast.add({ severity: 'success', summary: 'Success', detail: 'Purchase order cancelled', life: 3000, group: 'br' })
        loadPurchaseOrders()
    })
    .catch((error) => {
        toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data || error.message, life: 3000, group: 'br' })
    })
}

const showGRNsForPO = (po: PurchaseOrder) => {
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/grns`, {
        params: { "filter[purchase_order_id]": po.id },
        headers: {
            Authorization: `Bearer ${auth.accessToken.value}`
        }
    })
    .then((response) => {
        const po_grns: GRN[] = response.data.data
        if (po_grns.length === 1) {
            active_tab.value = 1
            showGRNDetail(po_grns[0])
        } else {
            active_tab.value = 1
            loadGRNs(po.id)
        }
    })
    .catch(() => {
        active_tab.value = 1
        loadGRNs(po.id)
    })
}

const showGRNDetail = (grn: GRN) => {
    grn_detail.value = grn
    grn_detail_dialog.value = true
}

const loadMaterials = () => {
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/materials`, {
        headers: {
            Authorization: `Bearer ${auth.accessToken.value}`
        }
    })
    .then((response) => {
        material_options.value = response.data.data
    })
}

const loadPurchaseOrders = () => {
    is_po_loading.value = true
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/purchase-orders`, {
        headers: {
            Authorization: `Bearer ${auth.accessToken.value}`
        }
    })
    .then((response) => {
        purchase_orders.value = response.data.data
    })
    .finally(() => {
        is_po_loading.value = false
    })
}

const loadGRNs = (purchase_order_id?: string) => {
    is_grn_loading.value = true
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/grns`, {
        params: purchase_order_id ? { "filter[purchase_order_id]": purchase_order_id } : undefined,
        headers: {
            Authorization: `Bearer ${auth.accessToken.value}`
        }
    })
    .then((response) => {
        grns.value = response.data.data
    })
    .finally(() => {
        is_grn_loading.value = false
    })
}

loadMaterials()
loadPurchaseOrders()
loadGRNs()

if (route.query.grn_id) {
    active_tab.value = 1
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/grns/${route.query.grn_id}`, {
        headers: {
            Authorization: `Bearer ${auth.accessToken.value}`
        }
    })
    .then((response) => {
        grn_detail.value = response.data.data
        grn_detail_dialog.value = true
        router.replace({ path: '/admin/purchase-orders' })
    })
}

if (route.query.new_po) {
    showNewPODialog()
    router.replace({ path: '/admin/purchase-orders' })
}
</script>