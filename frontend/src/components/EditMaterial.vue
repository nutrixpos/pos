<template>
<div class="grid">
    <div class="flex col-12 flex-column gap-2">
        <label for="name">{{$t('name')}}</label>
        <InputText id="name" v-model="edited_material.name" aria-describedby="name" />
    </div>
    <div class="flex col-12 flex-column gap-2">
        <label for="unit">{{$t('unit')}}</label>
        <InputText id="unit" v-model="edited_material.unit" aria-describedby="unit" />
    </div>
    <div class="col-12 flex">
        <Button :label="$t('cancel')"  severity="secondary" aria-label="Cancel"  />
        <Button class="ml-2" severity="primary" @click="returnMaterial" :label="$t('done')" aria-label="Done" />
    </div>
</div>
</template>

<script setup lang="ts">
import { defineProps,ref,defineEmits } from 'vue';
import { Material } from '@/classes/OrderItem';
import Button from 'primevue/button'
import InputText from 'primevue/inputtext';



const emit = defineEmits(['returnMaterial']);

const props = defineProps({
    material: {
        type: Material,
        required: true
    }
});

const edited_material = ref<Material>(new Material())

const returnMaterial = () => {
    emit('returnMaterial', edited_material.value)
}


const init = () => {
    edited_material.value = props.material
}

init()



</script>