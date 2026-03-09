<template>
  <main :style="`height:100%;direction:${orientation}`">
      <RouterView v-if="!setting_up" />
      <Toast  group="br" position="top-center">
        <template #content="{ message }">
            <section class="flex p-3 gap-3 w-full bg-black-alpha-90">
                <div style="border-radius:100%;background-color:white;height:1.9rem;width:2rem;" class="p-1 flex justify-content-center align-items-center opacity-95">
                  <i v-if="message.severity == `success`"  :class="`pi pi-check text-${message.severity}-500`"></i>
                  <i v-if="message.severity == `info`"  :class="`pi pi-info text-${message.severity}-500`"></i>
                  <i v-if="message.severity == `warning`"  :class="`pi pi-exclamation-triangle text-${message.severity}-500`"></i>
                  <i v-if="message.severity == `error`"  :class="`pi pi-times text-${message.severity}-500`"></i>
                </div>
                <div class="flex flex-column gap-3 w-full">
                    <p class="m-0 font-semibold text-base text-white opacity-80">{{ message.summary }}</p>
                    <p class="m-0 text-white text-base text-700">{{ message.detail }}</p>
                </div>
            </section>
        </template>
      </Toast>
  </main>
</template>

<script setup>
import { computed, getCurrentInstance,ref } from 'vue';
import Toast from 'primevue/toast';
import { globalStore } from '@/stores';
import axios from 'axios';
import { useRouter } from 'vue-router'
import { useToast } from "primevue/usetoast";

const toast = useToast();

const router = useRouter()

const setting_up = ref(true)

const { proxy } = getCurrentInstance();
const store = globalStore()
const orientation = computed(() => store.currentOrientation)

const getSettings = () => {
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/settings`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        },
    })
    .then((response)=>{
        console.log(response.data.data)
        store.setSettings(response.data.data)
    })
    .catch((err) => {
        console.log(err)
    });
}

const getConfigStatus = () => {
    return axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/api/setup/status`, {
        headers: {
            Authorization: `bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        },
    })
}

const init = () => {
    getConfigStatus()
    .then((response) => {
        if(response.data.setup) {
            getSettings()
            setting_up.value = false
        }else {
            router.push({ path: '/setup' }).finally(() => {
                setting_up.value = false
            })
        }
    })
    .catch((err) => {
        toast.add({severity:'error', summary: 'Error', detail: err, life: 3000, group: 'br'});
    });

}

init()

</script>

<style lang="scss">
@use '@/assets/styles.scss';

body {
    font-family: sans-serif; /* Replace with your desired font */
    height: 100vh;
    margin:0px;
    background-color: rgb(247, 247, 247);
}
.my-app-dark {
    body {
        font-family: sans-serif; /* Replace with your desired font */
        height: 100vh;
        background-color: #232327;
    }
}

</style>