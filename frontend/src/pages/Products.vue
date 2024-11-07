<template>
     <div class="w-full">
        <div class="grid mx-2">
            <div class="col-12 flex">
                <div class="gird w-full">
                    <div class="col-12">
                        <h3>Products</h3>
                    </div>
                    <div class="col-12">
                        <DataTable @page="updatProductsTableRowsPerPage" :lazy="true" :totalRecords="productsTableTotalRecords" :loading="isProductsTableLoading"  paginatorPosition="both"  paginator :rows="productsTableRowsPerPage" :rowsPerPageOptions="[50, 100, 500]" :value="products" stripedRows tableStyle="min-width: 50rem;max-height:50vh;" class="w-full pr-2">
                            <template #header>
                                <div class="flex justify-start">
                                    <Button icon="pi pi-plus" label="Add Product"  rounded raised @click="productAddDialog=true" />
                                </div>
                            </template>
                            <Column sortable field="name" header="Name"></Column>
                            <Column field="ready" header="Ready"></Column>
                            <Column field="materials.length" header="Materials"></Column>
                            <Column field="sub_products.length" header="Sub Products">
                                <template #body="slotProps">
                                    {{ slotProps.data.sub_products != null ? slotProps.data.sub_products.length : 0 }}
                                </template>
                            </Column>
                            <Column header="Actions">
                                <template #body="slotProps">
                                    <ConfirmPopup></ConfirmPopup>
                                    <ButtonGroup>
                                        <Button icon="pi pi-pencil" severity="secondary" aria-label="Edit" @click="prepareProductToEdit(slotProps.data)"/>
                                        <Button icon="pi pi-trash" severity="danger" aria-label="Remove" @click="confirmDeleteProduct($event,slotProps.data.id)"/>
                                    </ButtonGroup>
                                </template>
                            </Column>
                        </DataTable>
                        <Dialog v-model:visible="productAddDialog" modal :header="`Add new product`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
                            <div class="flex flex-column gap-2 w-5">
                                <label for="name">Name</label>
                                <InputText id="name" v-model="new_product_name" aria-describedby="name" />
                            </div>
                            <div class="flex flex-column gap-2 w-5 mt-2">
                                <label for="name">Ready Quantity</label>
                                <InputText v-model="new_product_ready" type="number" aria-describedby="ready" />
                            </div>
                            <div class="flex flex-column gap-2 w-10 mt-3">
                                <label for="name">Materials</label>
                                <DataTable :value="materials" stripedRows class="w-full pr-2">
                                    <Column field="name" header="Name"></Column>
                                    <Column field="quantity" header="Quantity">
                                        <template #body="slotProps">
                                            <InputText type="number" v-model.number="materials[slotProps.data.index].quantity" aria-describedby="quantity" />
                                        </template>
                                    </Column>
                                    <Column field="unit" header="Unit"></Column>
                                    <Column header="Actions">
                                        <template #body="slotProps">
                                            <ButtonGroup>
                                                <Button icon="pi pi-times" severity="secondary" aria-label="Remove" @click="materials.splice(slotProps.data.index,1)" />
                                            </ButtonGroup>
                                        </template>
                                    </Column>
                                    <template #header>
                                        <div class="flex justify-start">
                                            <Button icon="pi pi-plus" label="Add Material"  rounded raised @click="add_material_dialog=true" />
                                        </div>
                                    </template>
                                </DataTable>
                            </div>
                            <div class="flex flex-column gap-2 w-10 mt-3">
                                <label for="name">Sub Products</label>
                                <DataTable :value="sub_products" stripedRows class="w-full pr-2">
                                    <Column field="name" header="Name"></Column>
                                    <Column field="quantity" header="Quantity">
                                        <template #body="slotProps">
                                            <InputText type="number" v-model.number="sub_products[slotProps.data.index].quantity" aria-describedby="quantity" />
                                        </template>
                                    </Column>
                                    <Column header="Actions">
                                        <template #body="slotProps">
                                            <ButtonGroup>
                                                <Button icon="pi pi-times" severity="secondary" aria-label="Remove" @click="sub_products.splice(slotProps.data.index,1)" />
                                            </ButtonGroup>
                                        </template>
                                    </Column>
                                    <template #header>
                                        <div class="flex justify-start">
                                            <Button icon="pi pi-plus" label="Add Subproduct"  rounded raised @click="add_subproduct_dialog=true" />
                                        </div>
                                    </template>
                                </DataTable>
                            </div>
                            <template #footer>
                                <ButtonGroup>
                                    <Button label="Cancel"  severity="secondary" aria-label="Cancel"  />
                                    <Button class="ml-2" severity="primary" @click="submitProduct" label="Save" aria-label="Save" />
                                </ButtonGroup>
                            </template>
                        </Dialog>
                        <Dialog v-model:visible="productEditDialog" modal :header="`Edit product ${productToEdit.name}`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '50vw', '575px': '90vw' }">
                            <div class="flex flex-column gap-2 w-5">
                                <label for="name">Name</label>
                                <InputText id="name" v-model="productToEdit.name" aria-describedby="name" />
                            </div>
                            <div class="flex flex-column gap-2 w-5 mt-2">
                                <label for="name">Ready Quantity</label>
                                <InputText v-model="productToEdit.ready" type="number" aria-describedby="ready" />
                            </div>
                            <div class="flex flex-column gap-2 w-10 mt-3">
                                <label for="name">Materials</label>
                                <DataTable :value="productToEdit.materials" stripedRows class="w-full pr-2">
                                    <Column field="name" header="Name"></Column>
                                    <Column field="quantity" header="Quantity">
                                        <template #body="slotProps">
                                            <InputText type="number" v-model.number="productToEdit.materials[slotProps.data.index].quantity" aria-describedby="quantity" />
                                        </template>
                                    </Column>
                                    <Column field="unit" header="Unit"></Column>
                                    <Column header="Actions">
                                        <template #body="slotProps">
                                            <ButtonGroup>
                                                <Button icon="pi pi-times" severity="secondary" aria-label="Remove" @click="productToEdit.materials.splice(slotProps.data.index,1)" />
                                            </ButtonGroup>
                                        </template>
                                    </Column>
                                    <template #header>
                                        <div class="flex justify-start">
                                            <Button icon="pi pi-plus" label="Add Material"  rounded raised @click="edit_material_dialog=true" />
                                        </div>
                                    </template>
                                </DataTable>
                            </div>
                            <div class="flex flex-column gap-2 w-10 mt-3">
                                <label for="name">Sub Products</label>
                                <DataTable :value="productToEdit.sub_products" stripedRows class="w-full pr-2">
                                    <Column field="name" header="Name"></Column>
                                    <Column field="quantity" header="Quantity">
                                        <template #body="slotProps">
                                            <InputText type="number" v-model.number="productToEdit.sub_products[slotProps.data.index].quantity" aria-describedby="quantity" />
                                        </template>
                                    </Column>
                                    <Column header="Actions">
                                        <template #body="slotProps">
                                            <ButtonGroup>
                                                <Button icon="pi pi-times" severity="secondary" aria-label="Remove" @click="productToEdit.sub_products.splice(slotProps.data.index,1)" />
                                            </ButtonGroup>
                                        </template>
                                    </Column>
                                    <template #header>
                                        <div class="flex justify-start">
                                            <Button icon="pi pi-plus" label="Add Subproduct"  rounded raised @click="edit_subproduct_dialog=true" />
                                        </div>
                                    </template>
                                </DataTable>
                            </div>
                            <template #footer>
                                <ButtonGroup>
                                    <Button label="Cancel"  severity="secondary" aria-label="Cancel"  />
                                    <Button class="ml-2" severity="primary" @click="updateProduct" label="Save" aria-label="Save" />
                                </ButtonGroup>
                            </template>
                        </Dialog>
                        <Dialog v-model:visible="add_material_dialog">
                            <PickMaterial @returnMaterial="(material) => addMaterial(material)" />
                        </Dialog>
                        <Dialog v-model:visible="add_subproduct_dialog">
                            <PickProduct @returnProduct="(product) => addSubProduct(product)" />
                        </Dialog>

                        <Dialog v-model:visible="edit_material_dialog">
                            <PickMaterial @returnMaterial="(material) => addEdittedMaterial(material)" />
                        </Dialog>
                        <Dialog v-model:visible="edit_subproduct_dialog">
                            <PickProduct @returnProduct="(product) => addEdittedSubProduct(product)" />
                        </Dialog>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import ConfirmPopup from 'primevue/confirmpopup';
import DataTable from 'primevue/datatable';
import Dialog from 'primevue/dialog';
import InputText from 'primevue/inputtext';
import Column from 'primevue/column';
import Button from 'primevue/button';
import ButtonGroup from 'primevue/buttongroup';
import axios from 'axios'
import { ref,getCurrentInstance } from "vue";
import { useToast } from "primevue/usetoast";
import PickMaterial from '@/components/PickMaterial.vue';
import PickProduct from '@/components/PickProduct.vue';
import { useConfirm } from "primevue/useconfirm";
// import { Material } from '@/classes/OrderItem';

const { proxy } = getCurrentInstance();
const confirm = useConfirm();

const productToEdit = ref<any>({})
const productEditDialog =  ref(false)
const edit_material_dialog = ref(false)
const edit_subproduct_dialog = ref(false)

const productAddDialog = ref(false)
const new_product_name = ref("")
const new_product_ready = ref(0)
const add_subproduct_dialog = ref(false)
// const new_product_materials = ref<Material[]>([])
// const new_product_subproducts = ref([])
const materials = ref<any[]>([])
const sub_products = ref<any[]>([])

const add_material_dialog = ref(false)


const productsTableTotalRecords  = ref(0)
const productsTableRowsPerPage = ref(50)
const isProductsTableLoading = ref(true)
const products = ref([])
const toast = useToast()



const deleteProduct = (product_id: string) => {
    axios.get(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/deleteproduct`, {
        params: {
            id:product_id
        },
        headers: {
            Authorization: `Bearer ${proxy.$zitadel.oidcAuth.accessToken}`
        }
    }).then(() => {
        toast.add({ severity: 'success', summary: 'Success', detail: 'Product deleted successfully',group:'br' });
        loadProducts();
    }).catch(() => {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to delete product',group:'br' });
    });
}


const confirmDeleteProduct = (event,product_id) => {
    confirm.require({
        target: event.currentTarget,
        message: 'Are you sure you want to delete this product ?',
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


            deleteProduct(product_id)

         
        },
        reject: () => {
        }
    });
  }


const prepareProductToEdit = (product: any) => {
    productToEdit.value = JSON.parse(JSON.stringify(product))


    for (let i=0;i<productToEdit.value.materials?.length;i++){
        productToEdit.value.materials[i].index = i
    }

    for (let i=0;i<productToEdit.value.sub_products?.length;i++){
        productToEdit.value.sub_products[i].index = i
    }


    productEditDialog.value = true
}


const updateProduct = () => {

    axios.post(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/updateproduct`, productToEdit.value, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel.oidcAuth.accessToken}`
        }
    })
    .then(response => {
        toast.add({ severity: 'success', summary: 'Product Updated', detail: response.data.message,group:'br' });
        productEditDialog.value = false;
        productToEdit.value = {}
        loadProducts()
        // Optionally, refresh the product list or clear inputs
    })
    .catch(error => {
        toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.message || 'An error occurred',group:'br' });
    });
}


const submitProduct = () => {
    const payload = {
        name: new_product_name.value,
        ready: new_product_ready.value,
        materials: materials.value,
        sub_products: sub_products.value,
    };

    axios.post(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/addproduct`, payload, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel.oidcAuth.accessToken}`
        }
    })
        .then(response => {
            toast.add({ severity: 'success', summary: 'Product Added', detail: response.data.message,group:'br' });
            productAddDialog.value = false;

            loadProducts()
            // Optionally, refresh the product list or clear inputs
        })
        .catch(error => {
            toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.message || 'An error occurred',group:'br' });
        });
}


const addMaterial = (material: any) => {
    material.quantity = 1
    material.index = materials.value.length
    materials.value.push(material)
}

const addSubProduct = (product: any) => {
    product.quantity = 1
    product.index = sub_products.value.length
    sub_products.value.push(product)
}


const addEdittedMaterial = (material: any) => {
    material.quantity = 1
    material.index = productToEdit.value.materials.length
    productToEdit.value.materials.push(material)
}

const addEdittedSubProduct = (product: any) => {
    product.quantity = 1
    product.index = productToEdit.value.sub_products.length
    productToEdit.value.sub_products.push(product)
}


const updatProductsTableRowsPerPage = (event: any) => {
    loadProducts(event.first,event.rows)
}


const loadProducts = (first=0,rows=productsTableRowsPerPage.value) => {
    isProductsTableLoading.value = true;
    
    axios.get(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/products`, {
        params: {
            first: first,
            rows: rows
        },
        headers: {
            Authorization: `Bearer ${proxy.$zitadel.oidcAuth.accessToken}`
        }
    })
    .then(response => {
        products.value = response.data.products;
        productsTableTotalRecords.value = response.data.total_records;
    })
    .catch(() => {
        toast.add({severity:'error', summary: 'Error', detail: 'Failed to load products'});
    })
    .finally(() => {
        isProductsTableLoading.value = false;
    });
}

loadProducts()


</script>