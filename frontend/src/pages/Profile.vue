<template>
    <div class="flex flex-column align-items-center justify-content-center h-full py-5">
        <Card class="w-full md:w-6">
            <template #content>
                <div class="flex flex-column gap-3">
                    <div class="flex align-items-center gap-3">
                        <i class="pi pi-user text-4xl"></i>
                    </div>
                    <div class="flex flex-column gap-1">
                        <label>{{$t('username')}}</label>
                        <InputText v-model="username" disabled />
                    </div>
                    <div class="flex flex-column gap-1">
                        <label>{{$t('email')}}</label>
                        <InputText v-model="email" disabled />
                    </div>
                    <div class="flex flex-column gap-1">
                        <label>{{$t('role',3)}}</label>
                        <div class="flex flex-wrap gap-2">
                            <Chip v-for="(role,index) in roles" :key="index" :label="role" style="height: 1.5rem;" class="m-1" />
                        </div>
                    </div>
                    <Divider />
                    <div class="flex flex-column gap-2">
                        <Button :label="$t('change_password')" icon="pi pi-key" @click="passwordDialog = true" />
                    </div>
                </div>
            </template>
        </Card>

        <Dialog v-model:visible="passwordDialog" modal :header="$t('change_password')" :style="{ width: '25rem' }">
            <div class="flex flex-column gap-2">
                <label for="currentPassword">{{$t('current_password')}}</label>
                <InputText id="currentPassword" v-model="currentPassword" type="password" />
            </div>
            <div class="flex flex-column gap-2 mt-2">
                <label for="newPassword">{{$t('new_password')}}</label>
                <InputText id="newPassword" v-model="newPassword" type="password" />
            </div>
            <div class="flex flex-column gap-2 mt-2">
                <label for="confirmPassword">{{$t('confirm_password')}}</label>
                <InputText id="confirmPassword" v-model="confirmPassword" type="password" />
            </div>
            <template #footer>
                <ButtonGroup>
                    <Button :label="$t('cancel')" severity="secondary" @click="passwordDialog = false" />
                    <Button class="ml-2" severity="primary" @click="changePassword" :label="$t('save')" :loading="saving" />
                </ButtonGroup>
            </template>
        </Dialog>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useI18n } from 'vue-i18n';
import { useRouter } from 'vue-router';
import axios from "axios";
import InputText from "primevue/inputtext";
import Card from "primevue/card";
import Chip from "primevue/chip";
import Divider from "primevue/divider";
import Button from "primevue/button";
import ButtonGroup from "primevue/buttongroup";
import Dialog from "primevue/dialog";
import { useToast } from "primevue/usetoast";
import auth from "@/services/auth";

const { t } = useI18n({ useScope: 'global' });
const router = useRouter();
const toast = useToast();

const username = ref('');
const email = ref('');
const roles = ref<string[]>([]);
const passwordDialog = ref(false);
const currentPassword = ref('');
const newPassword = ref('');
const confirmPassword = ref('');
const saving = ref(false);

const backendUrl = `http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}`;

onMounted(() => {
    const user = auth.currentUser.value;
    if (user) {
        username.value = user.username;
        email.value = user.email;
        roles.value = user.roles;
    } else {
        router.push('/login');
    }
});

const changePassword = async () => {
    if (!currentPassword.value || !newPassword.value || !confirmPassword.value) {
        toast.add({ severity: 'warn', summary: 'Warning', detail: 'Please fill all fields' });
        return;
    }

    if (newPassword.value !== confirmPassword.value) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Passwords do not match' });
        return;
    }

    saving.value = true;

    try {
        await axios.patch(`${backendUrl}/api/auth/password`, {
            current_password: currentPassword.value,
            password: newPassword.value
        }, {
            headers: {
                Authorization: `Bearer ${auth.accessToken.value}`
            }
        });

        toast.add({ severity: 'success', summary: 'Success', detail: 'Password changed successfully' });
        passwordDialog.value = false;
        currentPassword.value = '';
        newPassword.value = '';
        confirmPassword.value = '';
    } catch (error: any) {
        toast.add({ severity: 'error', summary: 'Error', detail: error.response?.data || 'Failed to change password' });
    } finally {
        saving.value = false;
    }
};
</script>