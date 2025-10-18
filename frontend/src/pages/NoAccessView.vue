<template>
    <div class="flex justify-content-start align-items-center flex-column" style="height:100vh;width:100vw;">
        <i class="pi pi-ban mt-8" style="font-size:4rem;color:#8c737c"></i>
        <h2>Access denied</h2>
        <Button label="Relogin" severity="info" class="mt-2" @click="proxy.$zitadel?.oidcAuth.signOut()" />
    </div>
</template>
<script setup lang="ts">
import Button from 'primevue/button'
import { getCurrentInstance  } from "vue";
const { proxy } = getCurrentInstance();
</script>