<template>
     <div class="w-full">
        <div class="grid mx-2">
            <div class="col-12 flex">
                <div class="gird w-full">
                    <div class="col-12">
                        <h3>{{$t('settings')}}</h3>
                    </div>
                    <div class="col-12 flex-column flex">
                        <h4><i class="fa fa-boxes-stacked"></i> {{$t('inventory_item',3)}}</h4>
                        <div class="flex align-items-center gap-2">
                            <span>{{ $t('stock_alert_threshold') }}</span>
                            <InputText class="ml-3" v-model.number="stock_alert_treshold" type="number" />
                        </div>
                        
                        <Divider />
                        <h4 class="mb-2"><i class="pi pi-box"></i> {{$t('order',3)}}</h4>
                        <span class="mt-1">{{$t('queues')}} :</span>
                        <div class="flex align-items-center">
                            <div class="flex flex-column">
                                <div v-for="(queue,index) in order_queues" :key="index" class="flex align-items-center mt-2">
                                    <span>{{$t('prefix')}} :</span>
                                    <InputText v-model="order_queues[index].prefix" class="mx-2" />
                                    <span>{{ $t('next') }} :</span>
                                    <InputText v-model.number="order_queues[index].next"  class="mx-2 "/>
                                    <Button severity="secondary" aria-label="Remove" icon="pi pi-times" @click="order_queues.splice(index,1)" />
                                </div>
                                <div class="flex align-items-center mt-3">
                                    <span>{{$t('prefix')}} :</span>
                                    <InputText v-model="new_queue_prefix" class="mx-2" />
                                    <span>{{$t('next')}} :</span>
                                    <InputText v-model.number="new_queue_next"  class="mx-2 "/>

                                    <Button :label="$t('add')" @click="order_queues.push({prefix:new_queue_prefix,next:new_queue_next}); new_queue_prefix = ''; new_queue_next = 1" />
                                </div>
                            </div>
                        </div>
                        <span class="mt-5">{{$t('default_cost_calculation_method')}} :</span>
                        <div class="flex flex-wrap flex-column gap-4 my-3">
                            <div class="flex items-center gap-2">
                                <RadioButton v-model="default_cost_calculation_method" inputId="exact" name="exact" value="exact" />
                                <label for="exact">{{$t('exact')}}</label>
                                <Badge value="ℹ" size="small" severity="info" v-tooltip.top="'Exact cost of materials will be used for each product during order creation chosen from specific material entry'"></Badge>
                            </div>
                            <div class="flex items-center gap-2">
                                <RadioButton v-model="default_cost_calculation_method" inputId="average" name="average" value="average" />
                                <label for="average">{{$t('average')}}</label>
                                <Badge value="ℹ" size="small" severity="info" v-tooltip.top="'Average cost of materials available in the inventory will be used for each product during order creation, material/s with earliest expiration will be consumed first'"></Badge>
                            </div>
                        </div>

                        <Divider />
                        <div class="flex flex-column">
                            <h3><span class="pi pi-print"></span> {{ t('printer',2) }}</h3>
                            <span class="mt-2 font-bold">{{$t('client_receipt_printer')}}</span>
                            <div class="flex align-items-center mt-3">
                                    <span>{{t("host",1)}}:</span>
                                    <InputText v-model="client_receipt_printer_host"  class="mx-2" />
                            </div>
                            <span class="mt-5 font-bold">{{ $t('kitchen_receipt_printer') }}</span>
                            <div class="flex align-items-center mt-3">
                                    <span>{{t("host",1)}}:</span>
                                    <InputText v-model="kitchen_receipt_printer_host"  class="mx-2" />
                            </div>
                        </div>

                        <Divider />
                        <h4 class="mb-2"><i class="pi pi-wallet"></i> {{$t('payment',1)}}</h4>
                        <div class="mt-3">
                            {{$t('payment_source',3)}} :
                        </div>
                        <div class="flex align-items-center">
                            <div class="flex flex-column">
                                <div v-for="(source,index) in payment_sources" :key="index" class="flex align-items-center mt-2">
                                    <InputText v-model="payment_sources[index].name" class="mx-2" />
                                    <Button severity="secondary" aria-label="Remove" icon="pi pi-times" @click="payment_sources.splice(index,1)" />
                                </div>
                                <div class="flex align-items-center mt-3">
                                    <InputText v-model="new_payment_source" class="mx-2" />
                                    <Button :label="$t('add')" @click="payment_sources.push({name:new_payment_source}); new_payment_source = ''" />
                                </div>
                            </div>
                        </div>

                        <Divider />
                        <div class="flex flex-column">
                            <h3><span class="pi pi-language"></span> {{ t('language',3) }}</h3>
                            <Select @change="changedLanguage" v-model="selectedLang" :options="languages" optionLabel="language" placeholder="Select a Language" class="w-full md:w-3" />
                            <Button @click="applyLang" v-if="changedLang" class="mt-2 md:w-3" type="button" label="Apply" severity="secondary"></Button>
                        </div>

                        <Divider />
                        <div class="flex flex-column">
                            <h3><i class="pi pi-shop"></i> Shop Mode</h3>
                            <p class="mt-0 mb-3" style="color:#94a3b8;font-size:0.9rem;">Controls which features are available. Restart the app after changing.</p>
                            <div class="flex flex-wrap gap-4">
                                <div class="flex align-items-center gap-2">
                                    <RadioButton v-model="shop_mode" inputId="mode_kitchen" name="shop_mode" value="kitchen" />
                                    <label for="mode_kitchen" class="flex align-items-center gap-2" style="cursor:pointer">
                                        <i class="fa fa-kitchen-set" style="color:#f97316"></i>
                                        <span>Kitchen</span>
                                    </label>
                                </div>
                                <div class="flex align-items-center gap-2">
                                    <RadioButton v-model="shop_mode" inputId="mode_retail" name="shop_mode" value="retail" />
                                    <label for="mode_retail" class="flex align-items-center gap-2" style="cursor:pointer">
                                        <i class="pi pi-shopping-cart" style="color:#06b6d4"></i>
                                        <span>Retail</span>
                                    </label>
                                </div>
                            </div>
                        </div>

                        <div class="mt-6">
                            <Button :label="$t('save')" @click="saveSettings()" />
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import axios from 'axios'
import InputText from 'primevue/inputtext';
import Divider from 'primevue/divider';
import Button from 'primevue/button';
import { useToast } from "primevue/usetoast";
import {getCurrentInstance,ref} from 'vue'
import Dropdown from 'primevue/dropdown';
import { useI18n } from 'vue-i18n'
import { globalStore } from '../stores';
import {RadioButton,Avatar,Badge, Select} from 'primevue';
import auth from '../services/auth';

const { proxy } = getCurrentInstance();

const stock_alert_treshold = ref(0)
const order_queues = ref<any>({})

const new_queue_prefix = ref("")
const new_queue_next = ref(1)
const client_receipt_printer_host = ref()
const kitchen_receipt_printer_host = ref()

const default_cost_calculation_method = ref("average")
const shop_mode = ref('')

const toast = useToast();

const store = globalStore()


const new_payment_source = ref("")
const payment_sources = ref<Array<any>>([
    {"name":"Cash"},
    {"name":"Card"},
])

const changedLang = ref(false)

const { t,locale,setLocaleMessage } = useI18n({ useScope: 'global' }) 

const selectedLang : any = ref({"language":"English","code":"en"})
const languages = ref([
    {"language":"English","code":"en"},
])

const changedLanguage = () => {
    if (proxy.$i18n.locale != selectedLang.value.code) {
        changedLang.value = true
    }else {
        changedLang.value = false
    }
}


const saveSettings = () => {
    axios.patch(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/settings`, 
        {
            data: {
                inventory: {
                    stock_alert_treshold:stock_alert_treshold.value
                },
                orders: {
                    queues: order_queues.value,
                    default_cost_calculation_method: default_cost_calculation_method.value
                },
                language:{
                    code: selectedLang.value.code,
                    language: selectedLang.value.language
                },
                client_receipt_printer: {
                    host: client_receipt_printer_host.value
                },
                kitchen_receipt_printer: {
                    host: kitchen_receipt_printer_host.value
                },
                payment_sources: payment_sources.value == null ? [] : payment_sources.value,
                shop_mode: shop_mode.value
            }
        },
        {
            headers: {
                Authorization: `Bearer ${auth.accessToken.value}`
            }
        }
    )
    .then(()=>{
        store.setShopMode(shop_mode.value)
        toast.add({ severity: 'success', summary: 'Settings updated successfully!', detail: 'Done! ', life: 3000,group:'br' });
    })
    .catch((err) => {
        toast.add({ severity: 'error', summary: 'Failed', detail: 'Error updating settings !', life: 3000,group:'br' });  
        console.log(err)
    });
}

const getSettings = () => {
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/settings`, {
        headers: {
            Authorization: `Bearer ${auth.accessToken.value}`
        },
    })
    .then((response)=>{
        console.log(response.data.data)
        stock_alert_treshold.value = response.data.data.inventory.stock_alert_treshold
        order_queues.value = response.data.data.orders.queues
        default_cost_calculation_method.value = response.data.data.orders.default_cost_calculation_method
        selectedLang.value = response.data.data.language
        client_receipt_printer_host.value = response.data.data.client_receipt_printer?.host || ''
        kitchen_receipt_printer_host.value = response.data.data.kitchen_receipt_printer?.host || ''
        payment_sources.value = response.data.data.payment_sources == null ? [] : response.data.data.payment_sources
        shop_mode.value = response.data.data.shop_mode || ''
    })
    .catch((err) => {
        console.log(err)
    });
}

const applyLang = () => {
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/languages/${selectedLang.value.code}`, {
        headers: {
            Authorization: `Bearer ${auth.accessToken.value}`
        }
    })
    .then(response => {

        setLocaleMessage(selectedLang.value.code, response.data.data.pack);
        locale.value = response.data.data.code
        store.setOrientation(response.data.data.orientation)
    })
    .catch(() => {
        proxy.$toast.add({severity:'error', summary: 'Error', detail: "error loading language", life: 3000,grpup:'br'});
    });
}


const getAvailableLanguages = () => {
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/languages`, {
        headers: {
            Authorization: `Bearer ${auth.accessToken.value}`
        }
    })
    .then(response => {
        languages.value = response.data.data
    })
    .catch(error => {
        proxy.$toast.add({severity:'error', summary: 'Error', detail: error.response.data.error, life: 3000,grpup:'br'});
    });
}

getAvailableLanguages();
getSettings()

</script>