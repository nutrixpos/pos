<template>
    <div class="flex flex-column">
        <h3><span class="pi pi-language"></span> {{ t('language',3) }}</h3>
        <Dropdown @change="changedLanguage" v-model="selectedLang" :options="languages" optionLabel="language" placeholder="Select a Language" class="w-full md:w-12" />
        <Button @click="applyLang" v-if="changedLang" class="mt-2" type="button" label="Apply" severity="secondary"></Button>
    </div>
</template>

<script setup lang="ts">
import {ref,getCurrentInstance} from 'vue'
import Dropdown from 'primevue/dropdown';
import axios from 'axios';
import { useI18n } from 'vue-i18n'
import Button from 'primevue/button';
import { globalStore } from '../store';

const store = globalStore()
const { proxy } = getCurrentInstance();


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

const applyLang = () => {
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/languages/${selectedLang.value.code}`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then(response => {

        setLocaleMessage(selectedLang.value.code, response.data.data.pack);
        locale.value = response.data.data.code
        store.setOrientation('rtl')
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

</script>