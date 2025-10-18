<template>
     <div class="w-full">
        <div class="grid mx-2">
            <div class="col-12 flex">
                <div class="gird w-full">
                    <div class="col-12">
                        <h3>{{$t('product',3)}}</h3>
                    </div>
                    <div class="col-12">
                        <DataTable @page="updatProductsTableRowsPerPage" :lazy="true" :totalRecords="productsTableTotalRecords" :loading="isProductsTableLoading"  paginatorPosition="both"  paginator :rows="productsTableRowsPerPage" :rowsPerPageOptions="[50, 100, 500]" :value="products" stripedRows tableStyle="min-width: 50rem;max-height:50vh;" class="w-full pr-2">
                            <template #header>
                                <div class="flex justify-start">
                                    <Button icon="pi pi-plus" :label="$t('add_product')"  rounded raised @click="productAddDialog=true" />
                                </div>
                            </template>
                            <Column sortable field="name" :header="$t('name')"></Column>
                            <Column field="ready" :header="$t('ready')"></Column>
                            <Column field="materials.length" :header="$t('material',3)"></Column>
                            <Column field="sub_products.length" :header="$t('subproduct',3)">
                                <template #body="slotProps">
                                    {{ slotProps.data.sub_products != null ? slotProps.data.sub_products.length : 0 }}
                                </template>
                            </Column>
                            <Column :header="$t('actions')">
                                <template #body="slotProps">
                                    <ConfirmPopup></ConfirmPopup>
                                    <ButtonGroup>
                                        <Button icon="pi pi-pencil" severity="secondary" aria-label="Edit" @click="prepareProductToEdit(slotProps.data)"/>
                                        <Button icon="pi pi-trash" severity="danger" aria-label="Remove" @click="confirmDeleteProduct($event,slotProps.data.id)"/>
                                    </ButtonGroup>
                                </template>
                            </Column>
                        </DataTable>
                        <Dialog v-model:visible="productAddDialog" modal :header="$t('add_product')" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '90vw', '575px': '90vw' }">
                            <Form v-slot="$form" :initialValues="addproduct_initials" :resolver="add_product_resolver" class="flex flex-column gap-4 w-full sm:w-56">
                                <div class="flex flex-column gap-2 w-5">
                                    <label for="name">{{$t('name')}}</label>
                                    <InputText id="name" name="name" aria-describedby="name" />
                                </div>
                                <div class="flex flex-column gap-2 w-5 mt-2">
                                    <label for="price">{{$t('price')}}</label>
                                    <InputText id="name" name="price" type="number" aria-describedby="price" />
                                    <Message v-if="$form.price?.invalid" severity="error" size="small" variant="simple">{{ $form.price.error?.message }}</Message>
                                </div>
                                <div class="flex flex-column gap-2 mt-2">
                                    <label for="name">{{$t('image')}}</label>
                                    <FileUpload @before-send="beforeNewProductImageUpload" ref="newProductImageUpload" name="image" :fileLimit="1" :showCancelButton="false" :chooseLabel="$t('choose')" :showUploadButton="false"  :url="`${backendUrl}/api/products/${new_product_id}/image`" @upload="onAdvancedUpload($event)" :multiple="false" accept="image/*" :maxFileSize="1000000">
                                        <template #empty>
                                            <span>Drag and drop files to here to upload.</span>
                                        </template>
                                    </FileUpload>
                                </div>
                                <div class="flex flex-column gap-2 w-5 mt-2">
                                    <label for="name">{{$t('ready')}}</label>
                                    <InputText name="ready" type="number" aria-describedby="ready" />
                                    <Message v-if="$form.ready?.invalid" severity="error" size="small" variant="simple">{{ $form.ready.error?.message }}</Message>
                                </div>
                                <div class="flex flex-column gap-2 w-10 mt-3">
                                    <label for="name">{{$t('material',3)}}</label>
                                    <DataTable :value="materials" stripedRows class="w-full pr-2">
                                        <Column field="name" :header="$t('name')"></Column>
                                        <Column field="quantity" :header="$t('quantity')">
                                            <template #body="slotProps">
                                                <InputText type="number" v-model.number="materials[slotProps.data.index].quantity" aria-describedby="quantity" />
                                            </template>
                                        </Column>
                                        <Column field="unit" :header="$t('unit')"></Column>
                                        <Column :header="$t('actions')">
                                            <template #body="slotProps">
                                                <ButtonGroup>
                                                    <Button icon="pi pi-times" severity="secondary" aria-label="Remove" @click="removeMaterial(slotProps.data.index)" />
                                                </ButtonGroup>
                                            </template>
                                        </Column>
                                        <template #header>
                                            <div class="flex justify-start">
                                                <Button icon="pi pi-plus" :label="$t('add_material')"  rounded raised @click="add_material_dialog=true" />
                                            </div>
                                        </template>
                                    </DataTable>
                                </div>
                                <div class="flex flex-column gap-2 w-10 mt-3">
                                    <label for="name">{{$t('subproduct',3)}}</label>
                                    <DataTable :value="sub_products" stripedRows class="w-full pr-2">
                                        <Column field="name" :header="$t('name')"></Column>
                                        <Column field="quantity" :header="$t('quantity')">
                                            <template #body="slotProps">
                                                <InputText type="number" v-model.number="sub_products[slotProps.data.index].quantity" aria-describedby="quantity" />
                                            </template>
                                        </Column>
                                        <Column :header="$t('actions')">
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
                                <div class="flex flex-column gap-2 w-10 mt-3">
                                    <label for="name">{{$t('enable_inventory_consumption',1)}}</label>
                                    <ToggleSwitch v-model="newproduct_enable_inventory_consumption" />
                                </div>

                                <div style="opacity: 0;">
                                    {{ add_product_form=$form }}
                                </div>
                            </Form>
                            
                            <template #footer>
                                <ButtonGroup>
                                    <Button :label="$t('cancel')"  severity="secondary" aria-label="Cancel"  />
                                    <Button class="ml-2" severity="primary" @click="submitProduct" :label="$t('save')" aria-label="Save" />
                                </ButtonGroup>
                            </template>
                        </Dialog>
                        <Dialog v-model:visible="productEditDialog" modal :header="`${$t('edit_product')} ${productToEdit.name}`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '90vw', '575px': '90vw' }">
                            <div class="flex flex-column gap-2 w-5">
                                <label for="name">{{$t('name')}}</label>
                                <InputText id="name" v-model="productToEdit.name" aria-describedby="name" />
                            </div>
                            <div class="flex flex-column gap-2 w-5 mt-2">
                                <label for="price">{{$t('price')}}</label>
                                <InputText id="price" v-model.number="productToEdit.price" aria-describedby="price" />
                            </div>
                            <div class="flex flex-column gap-2 w-5 mt-2">
                                <label for="name">{{$t('ready')}}</label>
                                <InputText v-model.number="productToEdit.ready" type="number" aria-describedby="ready" />
                            </div>
                            <div class="flex flex-column gap-2 w-10 mt-3">
                                <Image :src="backend_host+'/public/'+productToEdit.image_url" alt="Image" width="250" />
                            </div>
                            <div class="flex flex-column gap-2 mt-2">
                                <FileUpload @before-send="beforeEditProductImageUpload" :chooseLabel="$t('edit')" ref="editProductImageUpload" name="image" :fileLimit="1" :showCancelButton="false" :showUploadButton="false"  :url="`${backendUrl}/api/products/${productToEdit.id}/image`" :multiple="false" accept="image/*" :maxFileSize="1000000">
                                    <template #empty>
                                        <span>Drag and drop files to here to upload.</span>
                                    </template>
                                </FileUpload>
                            </div>
                            <div class="flex flex-column gap-2 w-10 mt-3">
                                <label for="name">{{$t('material',3)}}</label>
                                <DataTable :value="productToEdit.materials" stripedRows class="w-full pr-2">
                                    <Column field="name" :header="$t('name')"></Column>
                                    <Column field="quantity" :header="$t('quantity')">
                                        <template #body="slotProps">
                                            <InputText type="number" v-model.number="productToEdit.materials[slotProps.data.index].quantity" aria-describedby="quantity" />
                                        </template>
                                    </Column>
                                    <Column field="unit" :header="$t('unit')"></Column>
                                    <Column :header="$t('actions')">
                                        <template #body="slotProps">
                                            <ButtonGroup>
                                                <Button icon="pi pi-times" severity="secondary" aria-label="Remove" @click="removeEditMaterial(slotProps.data.index)" />
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
                                <label for="name">{{$t('subproduct',3)}}</label>
                                <DataTable :value="productToEdit.sub_products" stripedRows class="w-full pr-2">
                                    <Column field="name" :header="$t('name')"></Column>
                                    <Column field="quantity" :header="$t('quantity')">
                                        <template #body="slotProps">
                                            <InputText type="number" v-model.number="productToEdit.sub_products[slotProps.data.index].quantity" aria-describedby="quantity" />
                                        </template>
                                    </Column>
                                    <Column :header="$t('actions')">
                                        <template #body="slotProps">
                                            <ButtonGroup>
                                                <Button icon="pi pi-times" severity="secondary" aria-label="Remove" @click="productToEdit.sub_products.splice(slotProps.data.index,1)" />
                                            </ButtonGroup>
                                        </template>
                                    </Column>
                                    <template #header>
                                        <div class="flex justify-start">
                                            <Button icon="pi pi-plus" :label="$t('add_subproduct')"  rounded raised @click="edit_subproduct_dialog=true" />
                                        </div>
                                    </template>
                                </DataTable>
                            </div>
                            <div class="flex flex-column gap-2 w-10 mt-3">
                                <label for="name">{{$t('enable_inventory_consumption',1)}}</label>
                                <ToggleSwitch v-model="productToEdit.enable_inventory_consumption" />
                            </div>
                            <template #footer>
                                <ButtonGroup>
                                    <Button :label="$t('cancel')"  severity="secondary" aria-label="Cancel"  />
                                    <Button class="ml-2" severity="primary" @click="updateProduct" :label="$t('save')" aria-label="Save" />
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
        <ConfirmPopup></ConfirmPopup>
    </div>
</template>

<script setup lang="ts">
import ConfirmPopup from 'primevue/confirmpopup';
import DataTable from 'primevue/datatable';
import Dialog from 'primevue/dialog';
import FileUpload from 'primevue/fileupload';
import InputText from 'primevue/inputtext';
import Column from 'primevue/column';
import Button from 'primevue/button';
import ButtonGroup from 'primevue/buttongroup';
import axios from 'axios'
import { ref,getCurrentInstance,computed, nextTick,reactive } from "vue";
import { useToast } from "primevue/usetoast";
import PickMaterial from '@/components/PickMaterial.vue';
import PickProduct from '@/components/PickProduct.vue';
import { useConfirm } from "primevue/useconfirm";
import { Image, Message, ToggleSwitch } from 'primevue';
import { Form } from '@primevue/forms';
// import { Material } from '@/classes/OrderItem';

const { proxy } = getCurrentInstance();
const confirm = useConfirm();

const productToEdit = ref<any>({})
const productEditDialog =  ref(false)
const edit_material_dialog = ref(false)
const edit_subproduct_dialog = ref(false)


const add_product_form = ref({})

const productAddDialog = ref(false)


const newproduct_enable_inventory_consumption = ref(true)

const addproduct_initials = reactive({
    name: '',
    ready: 0,
    id: '',
    price: 0
});

const add_product_resolver = ({ values }) => {
    const errors = {};

    if (!values.name) {
        errors.name = [{ message: 'Name is required.' }];
    }

    if (!/^\d+(\.\d+)?$/.test(values.ready)) {
        errors.ready = [{ message: 'Ready is not a valid number.' }];
    }

    if (!/^\d+(\.\d+)?$/.test(values.price)) {
        errors.price = [{ message: 'Price is not a valid number.' }];
    }

    return {
        values, // (Optional) Used to pass current form values to submit event.
        errors
    };
};

const new_product_id = ref("")
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

const newProductImageUpload = ref(null)
const editProductImageUpload = ref(null)

const beforeNewProductImageUpload = (event:any) => {
    event.xhr.setRequestHeader("Authorization", "Bearer " + proxy.$zitadel?.oidcAuth.accessToken);
}

const beforeEditProductImageUpload = (event:any) => {
    event.xhr.setRequestHeader("Authorization", "Bearer " + proxy.$zitadel?.oidcAuth.accessToken);
}

const onAdvancedUpload = () => {
    toast.add({ severity: 'info', summary: 'Success', detail: `${proxy.$t('image_uploaded')}`, life: 3000 });
};

const backend_host = computed(() => {
    return `http://${import.meta.env.VITE_APP_BACKEND_HOST}`;
});

const backendUrl = computed(() => {
    return `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}`;
});


// const onFileSelect = (event:any) => {
//     const file = event.files[0];
//     const reader = new FileReader();

//     reader.onload = async (e) => {
//         src.value = e.target?.result;
//     };

//     reader.readAsDataURL(file);
// }




const deleteProduct = (product_id: string) => {
    axios.delete(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/products/${product_id}`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    }).then(() => {
        toast.add({ severity: 'success', summary: 'Success', detail: 'Product deleted successfully',group:'br',life:3000 });
        loadProducts();
    }).catch(() => {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to delete product',group:'br',life:3000 });
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

    axios.patch(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/products/${productToEdit.value.id}`, 
    {
        data : productToEdit.value
    },{
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(() => {
        toast.add({ severity: 'success', summary: 'Product Updated', detail: proxy.$t('done'),group:'br',life:3000 });
        productEditDialog.value = false;


        editProductImageUpload.value?.upload()


        productToEdit.value = {}
        loadProducts()
        // Optionally, refresh the product list or clear inputs
    })
    .catch(error => {
        toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.message || 'An error occurred',group:'br' });
    });
}


const submitProduct = () => {

    if (!add_product_form.value.valid)
        return

    const payload = {
        name: add_product_form.value.name.value,
        ready: Number(add_product_form.value.ready.value),
        price: Number(add_product_form.value.price.value),
        materials: materials.value,
        sub_products: sub_products.value,
        enable_inventory_consumption: newproduct_enable_inventory_consumption.value,
    };

    axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/products`,{
        data : payload
    } , {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
        .then(async response => {
            toast.add({ severity: 'success', summary: 'Product Added', detail: proxy.$t('done'),group:'br' });
            new_product_id.value = response.data.data.id
            await nextTick()
            newProductImageUpload.value?.upload()
            productAddDialog.value = false;
            loadProducts()
            // Optionally, refresh the product list or clear inputs
        })
        .catch(error => {
            toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.data || 'An error occurred',group:'br' });
        });
}

const removeMaterial = (index) => {
    materials.value.splice(index,1)
    materials.value.forEach((m,i) => {
        m.index = i
    })
}

const removeEditMaterial = (index) => {
    productToEdit.value.materials.splice(index,1)
    productToEdit.value.materials.forEach((m,i) => {
        m.index = i
    })
}


const addMaterial = (material: any) => {
    add_material_dialog.value = false
    edit_material_dialog.value = false
    let exists = false

    materials.value.forEach((m) => {
        if (m.name == material.name){
            exists = true
        }
    })

    if (!exists){
        material.quantity = 1
        material.index = materials.value.length
        materials.value.push(material)
    }
}

const addSubProduct = (product: any) => {
    product.quantity = 1
    product.index = sub_products.value.length
    sub_products.value.push(product)
}


const addEdittedMaterial = (material: any) => {
    add_material_dialog.value = false
    edit_material_dialog.value = false
    let exists = false
    productToEdit.value.materials.forEach((m) => {
        if (m.name == material.name){
            exists = true
        }
    })

    if (!exists){
        material.quantity = 1
        material.index = productToEdit.value.materials.length
        productToEdit.value.materials.push(material)
    }
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

    if (first == 0){
        first = 1
    }

    const page_number = Math.ceil((first/rows))
    
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/products`, {
        params: {
            "page[number]": page_number,
            "page[size]": rows
        },
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {
        products.value = response.data.data;
        productsTableTotalRecords.value = response.data.meta.total_records;
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