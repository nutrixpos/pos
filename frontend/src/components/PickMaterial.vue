<template>
    <div>
        <DataTable  filterDisplay="row" :loading="loading" v-model:filters="filters" :globalFilterFields="['name']" :value="materials" stripedRows tableStyle="min-width: 50rem" class="w-full pr-5">
            <template #header>
                <div class="flex justify-content-start">
                    <IconField iconPosition="left">
                        <InputIcon>
                            <i class="pi pi-search" />
                        </InputIcon>
                        <InputText v-model="filters['name'].value" placeholder="Search by name" />
                    </IconField>
                </div>
            </template>
            <Column field="name" :header="$t('name')"></Column>
            <Column field="quantity" :header="$t('quantity')"></Column>
            <Column field="unit" :header="$t('unit')"></Column>
            <Column :header="$t('actions')">
                <template #body="slotProps">
                    <ButtonGroup>
                        <Button icon="pi pi-plus" :label="$t('add')" severity="secondary" aria-label="Ddd" @click="returnMaterial(slotProps.data)" />
                    </ButtonGroup>
                </template>
            </Column>
        </DataTable>
    </div>
</template>

<script setup lang="ts">
import {ref,defineEmits,getCurrentInstance} from 'vue'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext';
import InputIcon from 'primevue/inputicon'
import IconField from 'primevue/iconfield'
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import axios from 'axios'
import { Material, MaterialEntry } from '@/classes/OrderItem';
const { proxy } = getCurrentInstance();


const materials = ref([])
const loading = ref(false)

const filters = ref({
    name: { value: null},
});


const emit = defineEmits(['returnMaterial'])


const returnMaterial = (material: Material) => {
    emit('returnMaterial', material)
}


const GetMaterials = () => {
    loading.value = true
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/materials`,{
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then((response) => {

        
        response.data.data.forEach((component:Material,materialIndex: number) => {

            component.entries?.forEach((entry: MaterialEntry,entryIndex: number) => {
                if (entry.quantity < 0){
                    response.data.data[materialIndex].entries.splice(entryIndex, 1)
                    return
                }
                    
                component.quantity += entry.quantity
            });
        });


        materials.value = response.data.data
        loading.value = false
    })
}

GetMaterials()


</script>