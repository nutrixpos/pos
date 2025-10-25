<template>
     <div class="w-full">
        <div class="grid mx-2">
            <div class="col-12 flex">
                <div class="gird w-full">
                    <div class="col-12">
                        <h3>{{$t('settings')}}</h3>
                    </div>
                    <div class="col-12 flex-column flex">
                        <h4><i class="fa fa-carrot"></i> {{$t('material',3)}}</h4>
                        <div class="flex align-items-center">
                            <span>stock_alert_treshold</span>
                            <InputText class="ml-3" v-model.number="stock_alert_treshold" type="number" />
                        </div>
                        
                        <Divider />
                        <h4 class="mb-2"><i class="pi pi-box"></i> {{$t('order',3)}}</h4>
                        <span class="mt-1">Queues:</span>
                        <div class="flex align-items-center">
                            <div class="flex flex-column">
                                <div v-for="(queue,index) in order_queues" :key="index" class="flex align-items-center mt-2">
                                    <span>Prefix:</span>
                                    <InputText v-model="order_queues[index].prefix" class="mx-2" />
                                    <span>Next:</span>
                                    <InputText v-model.number="order_queues[index].next"  class="mx-2 "/>
                                    <Button severity="secondary" aria-label="Remove" icon="pi pi-times" @click="order_queues.splice(index,1)" />
                                </div>
                                <div class="flex align-items-center mt-3">
                                    <span>Prefix:</span>
                                    <InputText v-model="new_queue_prefix" class="mx-2" />
                                    <span>Next:</span>
                                    <InputText v-model.number="new_queue_next"  class="mx-2 "/>

                                    <Button label="Add" @click="order_queues.push({prefix:new_queue_prefix,next:new_queue_next}); new_queue_prefix = ''; new_queue_next = 1" />
                                </div>
                            </div>
                        </div>
                        <span class="mt-5">Default cost calculation method:</span>
                        <div class="flex flex-wrap flex-column gap-4 my-3">
                            <div class="flex items-center gap-2">
                                <RadioButton v-model="default_cost_calculation_method" inputId="exact" name="exact" value="exact" />
                                <label for="exact">Exact</label>
                                <Badge value="ℹ" size="small" severity="info" v-tooltip.top="'Exact cost of materials will be used for each product during order creation chosen from specific material entry'"></Badge>
                            </div>
                            <div class="flex items-center gap-2">
                                <RadioButton v-model="default_cost_calculation_method" inputId="average" name="average" value="average" />
                                <label for="average">Average</label>
                                <Badge value="ℹ" size="small" severity="info" v-tooltip.top="'Average cost of materials available in the inventory will be used for each product during order creation, FIFO will be used for material consumption'"></Badge>
                            </div>
                        </div>

                        <Divider />
                        <div class="flex flex-column">
                            <h3><span class="pi pi-print"></span> {{ t('printer',1) }}</h3>
                            <div class="flex align-items-center mt-3">
                                    <span>{{t("host",1)}}:</span>
                                    <InputText v-model="receipt_printer_host"  class="mx-2" />
                            </div>
                        </div>

                        <Divider />
                        <div class="flex flex-column">
                            <h3><span class="pi pi-language"></span> {{ t('language',3) }}</h3>
                            <Dropdown @change="changedLanguage" v-model="selectedLang" :options="languages" optionLabel="language" placeholder="Select a Language" class="w-full md:w-3" />
                            <Button @click="applyLang" v-if="changedLang" class="mt-2 md:w-3" type="button" label="Apply" severity="secondary"></Button>
                        </div>

                        <div class="mt-6">
                            <Button label="Save" @click="saveSettings()" />
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
import {RadioButton,Avatar,Badge} from 'primevue';

const { proxy } = getCurrentInstance();

const stock_alert_treshold = ref(0)
const order_queues = ref<any>({})

const new_queue_prefix = ref("")
const new_queue_next = ref(1)
const receipt_printer_host = ref()

const default_cost_calculation_method = ref("average")

const toast = useToast();

const store = globalStore()


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
                receipt_printer: {
                    host: receipt_printer_host.value
                }
            }
        },
        {
            headers: {
                Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
            }
        }
    )
    .then(()=>{
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
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        },
    })
    .then((response)=>{
        console.log(response.data.data)
        stock_alert_treshold.value = response.data.data.inventory.stock_alert_treshold
        order_queues.value = response.data.data.orders.queues
        default_cost_calculation_method.value = response.data.data.orders.default_cost_calculation_method
        selectedLang.value = response.data.data.language
        receipt_printer_host.value = response.data.data.receipt_printer.host
    })
    .catch((err) => {
        console.log(err)
    });
}

const applyLang = () => {
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/languages/${selectedLang.value.code}`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
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
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
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