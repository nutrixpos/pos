<template>
    <div>
        <Card style="overflow: hidden">
            <template #content>
                <div class="flex p-2">
                    <Avatar shape="circle" image="https://girlheartfood.com/wp-content/uploads/2020/06/Crispy-Chicken-Burger-10.jpg" class="mr-2" size="xlarge" />
                    <div class="flex align-items-center justify-content-center">
                        <h3>{{props.name}}</h3>
                    </div>
                </div>
            </template>
            <template #footer>
                <div class="flex gap-4 mt-1">
                    <Button class="w-full" label="+Comment" severity="primary" outlined @click="$emit('addwithcomment')" />
                    <Button class="w-full" label="Add" @click="$emit('add')"/>
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
import Avatar from 'primevue/avatar';
import Dialog from 'primevue/dialog'

const props = defineProps(['name'])

const visible = ref(false)
</script>