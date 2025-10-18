<template>
  <main :style="`height:100%;direction:${orientation}`">
      <RouterView />
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
import { computed } from 'vue';
import Toast from 'primevue/toast';
import { globalStore } from '@/stores';


const store = globalStore()
const orientation = computed(() => store.currentOrientation)

</script>

<style lang="scss">
@use '@/assets/styles.scss';


body {
    font-family: sans-serif; /* Replace with your desired font */
    height: 100vh;
    background-color: rgb(247, 247, 247);
}
</style>