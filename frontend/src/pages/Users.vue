<template>
     <div class="w-full">
        <div class="grid mx-2">
            <div class="col-12 flex">
                <div class="gird w-full">
                    <div class="col-12">
                        <h3>{{$t('user',3)}}</h3>
                    </div>
                    <div class="col-12">
                        <DataTable :value="users" stripedRows tableStyle="min-width: 50rem;max-height:50vh;" class="w-full pr-2">
                            <template #header>
                                <div class="flex justify-start">
                                    <Button icon="pi pi-plus" :label="$t('add_user')"  rounded raised @click="userAddDialog=true" />
                                </div>
                            </template>
                            <Column sortable field="username" :header="$t('username')"></Column>
                            <Column sortable field="email" :header="$t('email')"></Column>
                            <Column field="roles" :header="$t('role',3)">
                                <template #body="slotProps">
                                    <Chip v-for="(role,index) in slotProps.data.roles" :key="index" :label="role" style="height: 1.5rem;" class="m-1" />
                                </template>
                            </Column>
                            <Column :header="$t('actions')">
                                <template #body="slotProps">
                                    <ConfirmPopup></ConfirmPopup>
                                    <ButtonGroup>
                                        <Button icon="pi pi-trash" severity="danger" aria-label="Remove" @click="confirmDeleteUser($event,slotProps.data.id)"/>
                                    </ButtonGroup>
                                </template>
                            </Column>
                        </DataTable>
                        <Dialog v-model:visible="userAddDialog" modal :header="$t('add_user')" :style="{ width: '75rem' }" :breakpoints="{ '1199px': '90vw', '575px': '90vw' }">
                            <div class="flex flex-column gap-2 w-5">
                                <label for="username">{{$t('username')}}</label>
                                <InputText id="username" v-model="newUser.username" aria-describedby="username" />
                            </div>
                            <div class="flex flex-column gap-2 w-5 mt-2">
                                <label for="email">{{$t('email')}}</label>
                                <InputText id="email" v-model="newUser.email" aria-describedby="email" />
                            </div>
                            <div class="flex flex-column gap-2 w-5 mt-2">
                                <label for="password">{{$t('password')}}</label>
                                <InputText id="password" v-model="newUser.password" type="password" aria-describedby="password" />
                            </div>
                            <div class="flex flex-column gap-2 w-10 mt-3">
                                <label for="roles">{{$t('role',3)}}</label>
                                <div class="flex flex-wrap gap-2">
                                    <div v-for="role in availableRoles" :key="role" class="flex align-items-center">
                                        <Checkbox v-model="newUser.roles" :inputId="role" :value="role" />
                                        <label :for="role" class="ml-2">{{ role }}</label>
                                    </div>
                                </div>
                            </div>
                            <template #footer>
                                <ButtonGroup>
                                    <Button :label="$t('cancel')"  severity="secondary" aria-label="Cancel" @click="userAddDialog=false" />
                                    <Button class="ml-2" severity="primary" @click="submitUser" :label="$t('save')" aria-label="Save" />
                                </ButtonGroup>
                            </template>
                        </Dialog>
                    </div>
                </div>
            </div>
        </div>
     </div>
</template>

<script setup lang="ts">
import {ref,onMounted} from "vue";
import { useI18n } from 'vue-i18n'
import { globalStore } from '@/stores';
import axios from "axios";
import { getCurrentInstance } from "vue";
import DataTable from "primevue/datatable";
import Column from "primevue/column";
import Button from "primevue/button";
import ButtonGroup from "primevue/buttongroup";
import Dialog from "primevue/dialog";
import InputText from "primevue/inputtext";
import Chip from "primevue/chip";
import Checkbox from "primevue/checkbox";
import ConfirmPopup from "primevue/confirmpopup";
import { useConfirm } from "primevue/useconfirm";
import { useToast } from "primevue/usetoast";

const { proxy } = getCurrentInstance();
const confirm = useConfirm();
const toast = useToast();
const { t } = useI18n({ useScope: 'global' })

const store = globalStore()
const backendUrl = `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}`;

const users = ref([])
const userAddDialog = ref(false)
const availableRoles = ['admin', 'cashier', 'chef', 'superuser']
const newUser = ref({
    username: '',
    email: '',
    password: '',
    roles: []
})

const auth = proxy.$auth;

const loadUsers = () => {
    axios.get(`${backendUrl}/api/auth/users`, {
        headers: {
            Authorization: `Bearer ${auth.accessToken.value}`
        }
    })
    .then(response => {
        users.value = response.data;
    })
    .catch(() => {
        toast.add({severity:'error', summary: 'Error', detail: 'Failed to load users'});
    });
}

const submitUser = () => {
    if (!newUser.value.username || !newUser.value.password || !newUser.value.email) {
        toast.add({severity: 'warn', summary: 'Warning', detail: 'Please fill all required fields'});
        return;
    }

    axios.post(`${backendUrl}/api/auth/users`, newUser.value, {
        headers: {
            Authorization: `Bearer ${auth.accessToken.value}`
        }
    })
    .then(() => {
        toast.add({ severity: 'success', summary: 'User Added', detail: t('done'),group:'br' });
        userAddDialog.value = false;
        newUser.value = {
            username: '',
            email: '',
            password: '',
            roles: []
        };
        loadUsers();
    })
    .catch(error => {
        toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data?.message || 'An error occurred',group:'br' });
    });
}

const deleteUser = (user_id: string) => {
    axios.delete(`${backendUrl}/api/auth/users?id=${user_id}`, {
        headers: {
            Authorization: `Bearer ${auth.accessToken.value}`
        }
    }).then(() => {
        toast.add({ severity: 'success', summary: 'Success', detail: 'User deleted successfully',group:'br',life:3000 });
        loadUsers();
    }).catch(() => {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to delete user',group:'br',life:3000 });
    });
}


const confirmDeleteUser = (event, user_id) => {
    confirm.require({
        target: event.currentTarget,
        message: 'Are you sure you want to delete this user ?',
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
            deleteUser(user_id)
        },
        reject: () => {
        }
    });
}

onMounted(() => {
    loadUsers();
});
</script>