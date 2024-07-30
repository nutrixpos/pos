<template>
    <div class="flex flex-column m-2" style="height: 100%;">
        <Menubar :model="items">
            <template #item="{ item, props }">
                <a v-ripple class="flex align-items-center" v-bind="props.action" :href="item.link">
                    <span :class="item.icon" />
                    <span class="ml-2">{{ item.label }}</span>
                </a>
            </template>
        </Menubar>
        <div class="grid" style="flex-grow:1;">
            <div class="col-2">
                <Listbox  v-model="selectedCategory" :options="categories" optionLabel="name" class="w-full mt-2" filter>
                    <template #option="slotProps">
                        <div class="flex align-items-center">
                            <fa :icon="slotProps.option.icon" class="mr-2" />
                            <div>{{ slotProps.option.name }}</div>
                        </div>
                    </template>
                </Listbox>
            </div>
            <div class="col-8 flex pt-3 pb-3">
                <Panel header="Recipes" style="width:100%;">
                    <InputText v-model="searchtext" placeholder="Search" class="mb-4" />
                    <div class="flex flex-wrap">
                        <MealCard  v-for="(item,index) in products" :key="index" :name="item.name" :price="item.price" class="m-2" style="width:9rem;" @addwithcomment="visible=true;namewithcomment=item.name" @add="orderItems.push({name:item.name,comment})"/>
                    </div>
                </Panel>
            </div>
            <div class="col-2 flex pt-3 pb-3">
                <Panel header="Order Items" class="w-12" >
                    <div class="flex flex-column">
                        <Button label="Go" @click="goOrder" />
                        <div v-for="(item,index) in orderItems" :key="index">
                            <div class="flex justify-content-between align-items-center">
                                <p><strong>{{ item.name }}</strong></p>
                                <Button icon="pi pi-times" size="small" style="width:2rem;height: 2rem;" aria-label="Remove" severity="secondary" @click="orderItems.splice(index,1)" />
                            </div>
                            <p class="m-0">{{ item.comment }}</p>
                        </div>
                    </div>
                </Panel>
            </div>
            <Dialog v-model:visible="visible" modal header="Add Comment" :style="{ width: '25rem' }">
                <InputText v-model="comment" placeholder="Comment" class="mb-4" />
                <div class="flex justify-content-end gap-2">
                    <Button type="button" label="Close" severity="secondary"></Button>
                    <Button type="button" label="Add" @click="addWithComment()"></Button>
                </div>
            </Dialog>
        </div>
    </div>
</template>

<script setup>
  import Menubar from 'primevue/menubar';
  import Dialog from 'primevue/dialog';
  import Listbox from 'primevue/listbox';
  import Panel from 'primevue/panel';
  import InputText from 'primevue/inputtext';
  import Button from 'primevue/button';
  import { useToast } from "primevue/usetoast";
  import axios from 'axios'


  const toast = useToast();

  
  import MealCard from '@/components/MealCard.vue';
  
  import { ref,watch } from "vue";
  
  
  const comment = ref("")
  const namewithcomment = ref("")
  const visible = ref(false)
  const selectedCategory = ref();
  
  
const searchtext = ref("")
const categories = ref([])

const orderItems = ref([
])

const addWithComment = () => {
    orderItems.value.push({name:namewithcomment,comment:comment.value})
    visible.value=false
    comment.value = ""
}


const getCategories = async () => {
    const response = await axios.get("http://localhost:8000/api/categories")
    categories.value = categories.value.concat(response.data)
    selectedCategory.value = categories.value[0]
}

getCategories();


const goOrder = () => {

    if (orderItems.value.length > 0){
        axios.post("http://localhost:8000/api/submitorder",
            {
                items:orderItems.value   
            }
        ).then((response) => {
            console.log(response.data)
            toast.add({ severity: 'success', summary: 'Success', detail: 'Order in progress !', life: 3000,group:'br' });
        })
    
        orderItems.value = []
    }
};


watch(searchtext, (newSearchText) => {
  console.log(`x is ${newSearchText}`)
})


  
const products = ref([
])


// const showAllItems = () => {
//     categories.value.forEach((category) => {
//         if (category.recipes){
//             category.recipes.forEach((recipe) => {
//                 products.value.push({
//                     name:recipe.Name
//                 })
//             })
//         }
//     })
// }

watch(selectedCategory, (category) => {
    if (category != null){
        products.value = []
        category.recipes.forEach((recipe) => {
            products.value.push({
                name:recipe.Name,
                price:recipe.Price
            })
        })
    }
})


  const items = ref([
      {
          label: 'Cashier',
          icon: 'pi pi-home',
          link: 'home',
      },
      {
          label: 'Kitchen',
          icon: 'pi pi-star',
          link:'kitchen'
      },
      {
          label: 'Analytics',
          icon: 'pi pi-search',
          link: 'analytics',
      }
  ]);
  
  
  </script>
  
  <style>
  html,
  body {
    height: 100%;
    margin:0;
  }
  </style>
  