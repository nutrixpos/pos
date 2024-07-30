<template>
    <main style="height:100%;">
        <RouterView />
        <Toast  group="br" position="bottom-right" />
    </main>
</template>

<script setup>
import Toast from 'primevue/toast';
</script>

<style>
#app{
    height: 100vh;
    overflow: hidden;
}
</style>