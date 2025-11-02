<template>
    <div :style="`position:relative;${ props.item.availability < 1 ? 'filter: grayscale(100%);' : ''}'`">
        <div v-if="props.item.availability < 1" class="w-full h-full" style="position:absolute;z-index:99;cursor: not-allowed;"></div>
        <div class="mealcard" :style="`overflow: hidden;cursor: pointer;background-color:${store.getColorMode == 'light' ? 'white' : '#27272A'}`" @click="$emit('add')">
            <div class="flex flex-column" style="position:relative;">
                <Button icon="pi pi-ellipsis-h" @click.stop="toggle" severity="secondary" aria-label="Save" style="width: 2rem; height: 2rem; position:absolute;top:0;right:0;" size="small" class="m-1" />
                <div id='logo' :style="`background:url(${backend_host}/public/${props.item.image_url}) ;height:7rem;background-size:cover;background-position:center;`" class="w-full"></div>
                <div class="flex align-items-center" style="height: 3rem;">
                    <h4 class="m-0 p-1">{{props.item.name}}</h4>
                </div>
                <p class="mx-1 my-1" :style="`color:${store.getColorMode == 'light' ? 'green' : 'lightgreen'}`"><strong>{{props.item.price}} {{$t('egp')}}</strong></p>   

                <div class="text-center gap-1 flex align-items-center justify-content-center" style="background-color:#ffd589;">
                    <i v-tooltip.top="$t('inventory_consumption_disabled')" v-if="!props.item.enable_inventory_consumption" class="fa fa-unlink" style="font-size:1rem;color:indigo;z-index:99"></i>
                    <p class="m-0" :style="`font-size:0.9rem;color:${store.getColorMode == 'light' ? '' : 'black'}`">{{ props.item.availability != undefined ? Math.max(0, props.item.availability) : "..." || "..." }} {{$t('possible')}}</p>
                </div>

                <OverlayPanel ref="op">
                    <div class="flex flex-column gap-3 w-25rem">
                        <Button label="Add with comment" @click="$emit('addwithcomment')"   icon="pi pi-comment" />
                    </div>
                </OverlayPanel>
            </div>
        </div>
        <Dialog v-model:visible="visible" modal header="Edit Profile" :style="{ width: '25rem' }">
            <span class="p-text-secondary block mb-5">Update your information.</span>
            <div class="flex align-items-center gap-3 mb-3">
                <label for="username" class="font-semibold w-6rem">Username</label>
                <InputText id="username" class="flex-auto" autocomplete="off" />
            </div>
            <div class="flex align-items-center gap-3 mb-5">
                <label for="email" class="font-semibold w-6rem">Email</label>
                <InputText id="email" class="flex-auto" autocomplete="off" />
            </div>
            <div class="flex justify-content-end gap-2">
                <Button type="button" label="Cancel" severity="secondary" @click="visible = false"></Button>
                <Button type="button" label="Save" @click="visible = false"></Button>
            </div>
        </Dialog>
    </div>
</template>

<script setup>
import {ref, defineProps,computed} from 'vue'

import Card from 'primevue/card';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'
import OverlayPanel from 'primevue/overlaypanel';
import { globalStore } from '@/stores';

const backend_host = computed(() => {
    return `http://${import.meta.env.VITE_APP_BACKEND_HOST}`;
});

const store = globalStore()

const op = ref();


const props = defineProps(['item'])


const toggle = (event) => {
    op.value.toggle(event);
}

const visible = ref(false)
</script>

<style>
.mealcard {
    border-radius: 10px;
    box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.1), 0 1px 5px 0 rgba(0, 0, 0, 0.08);
}
</style>