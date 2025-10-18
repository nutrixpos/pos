<template>
     <div class="w-full">
        <div class="grid mx-2">
            <div class="col-12 flex">
                <div class="gird w-full">
                    <div class="col-12">
                        <h3>Hubsync</h3>
                    </div>
                    <div class="col-12">
                        <div class="col-12 lg:col-4 flex flex-column gap-2">
                            <label for="binary" class="font-bold">Enabled</label>
                            <ToggleSwitch v-model="hubSync.settings.enabled" class="mb-2" />
                        </div>
                        <div class="col-12 lg:col-4 flex flex-column gap-2">
                            <label for="binary" class="font-bold">Server Host</label>
                            <InputText v-model="hubSync.settings.server_host" type="text" class="w-full mb-2" />
                        </div>
                        <div class="col-12 lg:col-4 flex flex-column gap-2">
                            <label for="binary" class="font-bold">Token</label>
                            <InputText v-model="hubSync.settings.token" type="text" class="w-full mb-2" />
                        </div>
                        <div class="col-12 lg:col-4 flex flex-column gap-2">
                            <label for="binary" class="font-bold">Sync Interval (seconds)</label>
                            <InputText v-model.number="hubSync.settings.sync_interval" type="text" class="w-full mb-2" />
                        </div>
                        <div class="col-12 lg:col-4 flex flex-column gap-2">
                            <label for="binary" class="font-bold">Buffer Size</label>
                            <InputNumber v-model.nunber="hubSync.settings.buffer_size" mode="decimal" class="w-full mb-2" />
                        </div>
                        <div class="col-12 lg:col-4 flex flex-column gap-2">
                            <label for="binary" class="font-bold">Last Synced</label>
                            <span>{{ hubSync.last_synced || "-" }}</span>
                        </div>
                        <div class="col-12 lg:col-4 flex flex-column gap-2">
                            <label for="binary" class="font-bold">Sync Progress</label>
                            <span>{{ hubSync.sync_progress || "-" }}</span>
                        </div>
                        <div class="col-12 lg:col-4 flex flex-column gap-2">
                            <label for="binary" class="font-bold">SyncSales</label>
                            <ToggleSwitch v-model="hubSync.settings.sync_sales" class="mb-2" />
                        </div>
                        <div class="col-12 lg:col-4 flex flex-column gap-2">
                            <label for="binary" class="font-bold">SyncInventory</label>
                            <ToggleSwitch v-model="hubSync.settings.sync_inventory" class="mb-2" />
                        </div>
                        
                        <div class="col-12">
                            <Button label="Save" class="w-5rem" @click="save" />                            
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import {ref,getCurrentInstance} from 'vue'
import {InputText,ToggleSwitch,Button, InputNumber} from 'primevue'
import axios from 'axios'
import { useToast } from "primevue/usetoast";
const { proxy } = getCurrentInstance();

const toast = useToast()

const hubSync = ref({
        settings: {
            enabled: false,
            server_host: "",
            token: "",
            buffer_size: 0,
            sync_interval: 0,
            sync_sales: true,
            sync_inventory: false,
        },
        last_synced: "",
        sync_progress: 0,
    })



const save = () => {
    axios.patch(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/hubsync/api/settings`, {
        data: hubSync.value
    },{
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {
        toast.add({ severity: 'success', summary: 'Success', detail: 'Data saved successfully', life: 3000,group:'br' });
    })
    .catch(err => {
        toast.add({ severity: 'error', summary: 'Error '+err.response.status, detail: err.response.data,group:'br' });
    });
}

const getHubsyncInfo = () => {
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/hubsync/api/settings`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {
        hubSync.value = response.data.data
    })
    .catch(err => {
        toast.add({ severity: 'error', summary: 'Error '+err.response.status, detail: err.response.data,group:'br' });
    });
}


getHubsyncInfo()


</script>

<style scoped>
label {
    font-size:1rem;
    color: var(--p-button-text-primary-color)
}
</style>