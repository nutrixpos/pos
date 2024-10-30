<template>
    <div>
        Access denied
        <Button label="Logout" @click="proxy.$zitadel.oidcAuth.signOut()" />
    </div>
</template>
<script setup lang="ts">
import Button from 'primevue/button'
import { getCurrentInstance  } from "vue";
const { proxy } = getCurrentInstance();
</script>