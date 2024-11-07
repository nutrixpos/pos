<template>
    <div>
        <DataTable  filterDisplay="row" :loading="loading" v-model:filters="filters" :globalFilterFields="['name']" :value="products" stripedRows tableStyle="min-width: 50rem" class="w-full pr-5">
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
            <Column field="name" header="Name"></Column>
            <Column header="Actions">
                <template #body="slotProps">
                    <ButtonGroup>
                        <Button icon="pi pi-plus" label="Add" severity="secondary" aria-label="Ddd" @click="returnProduct(slotProps.data)" />
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
import { FilterMatchMode } from 'primevue/api';
const { proxy } = getCurrentInstance();


const products = ref([])
const loading = ref(false)

const filters = ref({
    name: { value: null, matchMode: FilterMatchMode.CONTAINS },
});


const emit = defineEmits(['returnProduct'])


const returnProduct = (product) => {
    emit('returnProduct', product)
}


const GetProducts = (first=0,rows=9999999999) => {
    loading.value = true
    axios.get(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/products?first_index=${first}&rows=${rows}`,{
        headers: {
            Authorization: `Bearer ${proxy.$zitadel.oidcAuth.accessToken}`
        }
    })
    .then((response) => {
        products.value = response.data.products
        loading.value = false
    })
}

GetProducts()


</script>