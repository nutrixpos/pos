<template>
    <div :style="`position:relative;${ props.item.availability <= 0 ? 'filter: grayscale(100%);' : ''}'`">
        <div v-if="props.item.availability <= 0" class="w-full h-full" style="position:absolute;z-index:99"></div>
        <Card style="overflow: hidden;height:13rem;cursor: pointer;" @click="$emit('add')">
            <template #header>
                <div class="flex flex-column" style="position:relative;">
                    <Button icon="pi pi-ellipsis-h" @click.stop="toggle" severity="secondary" aria-label="Save" style="width: 2rem; height: 2rem; position:absolute;top:0;right:0;" size="small" class="m-1" />
                    <div id='logo' style='background:url(https://girlheartfood.com/wp-content/uploads/2020/06/Crispy-Chicken-Burger-10.jpg) ;height:7rem;background-size:cover;background-position:center;' class="w-full"></div>
                    <div class="flex align-items-center" style="height: 3rem;">
                        <h4 class="m-0 p-1">{{props.item.name}}</h4>
                    </div>
                    <p class="mx-1 my-1" style="color:green"><strong>{{props.item.price}} EGP</strong></p>   

                    <div class="text-center flex align-items-center justify-content-center" style="background-color:#ffd589;">
                        <p class="m-0" style="font-size:0.9rem;">{{ props.item.availability ? Math.max(0, props.item.availability) : "..." || "..." }} can</p>
                    </div>

                    <OverlayPanel ref="op">
                        <div class="flex flex-column gap-3 w-25rem">
                            <Button label="Add with comment" @click="$emit('addwithcomment')"   icon="pi pi-comment" />
                        </div>
                    </OverlayPanel>
                </div>
            </template>
        </Card>
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
import {ref, defineProps} from 'vue'

import Card from 'primevue/card';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog'
import InputText from 'primevue/inputtext'

import OverlayPanel from 'primevue/overlaypanel';


const op = ref();


const props = defineProps(['item'])


const toggle = (event) => {
    op.value.toggle(event);
}

const visible = ref(false)
</script>