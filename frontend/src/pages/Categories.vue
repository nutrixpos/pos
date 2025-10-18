<template>
    <div class="w-full">
        <div class="grid mx-2">
            <div class="col-12 flex">
                <div class="gird w-full">
                    <div class="col-12">
                        <h3>{{$t('category',3)}}</h3>
                    </div>
                    <div class="col-12">
                        <DataTable @page="updatCategoriesTableRowsPerPage" :lazy="true" :totalRecords="categoriesTableTotalRecords" :loading="isCategoriesTableLoading"  paginatorPosition="both"  paginator :rows="categoriesTableRowsPerPage" :rowsPerPageOptions="[50, 100, 500]" :value="categories" stripedRows tableStyle="min-width: 50rem;max-height:50vh;" class="w-full pr-2">
                            <template #header>
                                <div class="flex justify-start">
                                    <Button icon="pi pi-plus" :label="$t('add_category')"  rounded raised @click="categoryAddDialog=true" />
                                </div>
                            </template>
                            <Column sortable field="name" :header="$t('name')"></Column>
                            <Column field="products.length" :header="$t('product',3)">
                                <template #body="slotProps">
                                    {{ slotProps.data.products != null ? slotProps.data.products.length : 0 }}
                                </template>
                            </Column>
                            <Column :header="$t('actions')">
                                <template #body="slotProps">
                                    <ConfirmPopup></ConfirmPopup>
                                    <ButtonGroup>
                                        <Button icon="pi pi-pencil" severity="secondary" aria-label="Edit" @click="prepareCategoryToEdit(slotProps.data)"/>
                                        <Button icon="pi pi-trash" severity="danger" aria-label="Remove" @click="confirmDeleteCategory($event,slotProps.data.id)"/>
                                    </ButtonGroup>
                                </template>
                            </Column>
                        </DataTable>
                        <Dialog v-model:visible="categoryAddDialog" modal :header="`Add new category`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '90vw', '575px': '90vw' }">
                            <div class="flex flex-column gap-2 w-5">
                                <label for="name">Name</label>
                                <InputText id="name" v-model="new_category.name" aria-describedby="name" />
                            </div>
                            <div class="flex flex-column gap-2 w-10 mt-3">
                                <label for="name">Products</label>
                                <DataTable :value="new_category.products" stripedRows class="w-full pr-2">
                                    <Column field="name" header="Name"></Column>
                                    <Column header="Actions">
                                        <template #body="slotProps">
                                            <ButtonGroup>
                                                <Button icon="pi pi-times" severity="secondary" aria-label="Remove" @click="new_category.products.splice(slotProps.data.index,1)" />
                                            </ButtonGroup>
                                        </template>
                                    </Column>
                                    <template #header>
                                        <div class="flex justify-start">
                                            <Button icon="pi pi-plus" :label="$t('add_product')"  rounded raised @click="new_cateogry_product_dialog=true" />
                                        </div>
                                    </template>
                                </DataTable>
                            </div>
                            <template #footer>
                                <ButtonGroup>
                                    <Button label="Cancel"  severity="secondary" aria-label="Cancel"  />
                                    <Button class="ml-2" severity="primary" @click="submitCategory" label="Save" aria-label="Save" />
                                </ButtonGroup>
                            </template>
                        </Dialog>
                        <Dialog v-model:visible="categoryEditDialog" modal :header="`Edit category ${categoryToEdit.name}`" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '90vw', '575px': '90vw' }">
                            <div class="flex flex-column gap-2 w-5">
                                <label for="name">Name</label>
                                <InputText id="name" v-model="categoryToEdit.name" aria-describedby="name" />
                            </div>
                            <div class="flex flex-column gap-2 w-10 mt-3">
                                <label for="name">Products</label>
                                <DataTable :value="categoryToEdit.products" stripedRows class="w-full pr-2">
                                    <Column field="name" header="Name"></Column>
                                    <Column header="Actions">
                                        <template #body="slotProps">
                                            <ButtonGroup>
                                                <Button icon="pi pi-times" severity="secondary" aria-label="Remove" @click="categoryToEdit.products.splice(slotProps.data.index,1)" />
                                            </ButtonGroup>
                                        </template>
                                    </Column>
                                    <template #header>
                                        <div class="flex justify-start">
                                            <Button icon="pi pi-plus" label="Add Product"  rounded raised @click="add_editted_product_dialog=true" />
                                        </div>
                                    </template>
                                </DataTable>
                            </div>
                            <template #footer>
                                <ButtonGroup>
                                    <Button label="Cancel"  severity="secondary" aria-label="Cancel"  />
                                    <Button class="ml-2" severity="primary" @click="updateCategory" label="Save" aria-label="Save" />
                                </ButtonGroup>
                            </template>
                        </Dialog>
                        <Dialog v-model:visible="new_cateogry_product_dialog">
                            <PickProduct @returnProduct="(product) => newCategoryAddProduct(product)" />
                        </Dialog>
                        <Dialog v-model:visible="add_editted_product_dialog">
                            <PickProduct @returnProduct="(product) => addEdittedProduct(product)" />
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
import PickProduct from '@/components/PickProduct.vue';
import axios from 'axios'
import { ref,getCurrentInstance } from "vue";
import { useToast } from "primevue/usetoast";
import { useConfirm } from "primevue/useconfirm";

const { proxy } = getCurrentInstance();
const confirm = useConfirm();
const toast = useToast()

const new_category = ref<any>({})
const new_cateogry_product_dialog = ref(false)

const categories = ref<any[]>([])
const categoriesTableTotalRecords = ref(0)
const isCategoriesTableLoading = ref(false)
const categoriesTableRowsPerPage = ref(50)
const categoryAddDialog = ref(false)

const categoryToEdit = ref<any>({})
const categoryEditDialog = ref(false)
const add_editted_product_dialog = ref(false)

const prepareCategoryToEdit = (product: any) => {
    categoryToEdit.value = JSON.parse(JSON.stringify(product))

    for (let i=0;i<categoryToEdit.value.products?.length;i++){
        categoryToEdit.value.products[i].index = i
    }


    categoryEditDialog.value = true
}

const updatCategoriesTableRowsPerPage = (event: any) => {
    getCategories(event.first,event.rows)
}

const submitCategory = () => {

    axios.post(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/categories`, {
        data: new_category.value
    }, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    }).then(() => {
        categoryAddDialog.value = false
        new_category.value = {}
        getCategories();
        setTimeout(() => {
            toast.add({ severity: 'success', summary: 'Success', detail: 'Category added successfully',group:'br' });
        }, 1000);
    }).catch(() => {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to add category',group:'br' });
    })
}

const deleteCategory = (category_id: string) => {
    axios.delete(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/categories/${category_id}`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    }).then(() => {
        
        setTimeout(() => {
            toast.add({ severity: 'success', summary: 'Success', detail: 'Category deleted successfully',group:'br' });
        }, 1000);

        getCategories();
    }).catch(() => {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to delete category',group:'br' });
    });
}

const updateCategory = () => {
    axios.patch(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/categories/${categoryToEdit.value.id}`, {
        data: categoryToEdit.value
    }, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    }).then(() => {
        toast.add({ severity: 'success', summary: 'Success', detail: 'Category updated successfully',group:'br' });
        categoryEditDialog.value = false
        getCategories()
    }).catch(() => {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to update category',group:'br' });
    });
}

const confirmDeleteCategory = (event,category_id) => {
    confirm.require({
        target: event.currentTarget,
        message: 'Are you sure you want to delete this category ?',
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


            deleteCategory(category_id)

         
        },
        reject: () => {
        }
    });
  }


const newCategoryAddProduct = (product: any) => {
    product.index = new_category.value.products ?  new_category.value.products.length : 0
    new_category.value.products? new_category.value.products.push(product) : new_category.value.products = [product]
    new_cateogry_product_dialog.value = false
}

const addEdittedProduct = (product: any) => {
    product.index = categoryToEdit.value.products ?  categoryToEdit.value.products.length : 0
    categoryToEdit.value.products? categoryToEdit.value.products.push(product) : categoryToEdit.value.products = [product]
    add_editted_product_dialog.value = false
}

const getCategories = (first=0,rows=categoriesTableRowsPerPage.value) => {


    if (first == 0){
        first = 1
    }

    const page_number = Math.ceil((first/rows))
    
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/categories?page[number]=${page_number}&page[size]=${rows}`, {
        headers: {
            'Authorization': `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {
        categories.value = response.data.data;
        categoriesTableTotalRecords.value = response.data.meta.total_records;
    })
    .catch(() => {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to fetch categories' });
    });

    }



getCategories()


</script>