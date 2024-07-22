<template>
    <Menubar :model="items" />
    <div class="grid">
      <div class="col-2">
          <Listbox  v-model="selectedCategory" :options="categories" optionLabel="category" class="w-full mt-2">
              <template #option="slotProps">
                  <div class="flex align-items-center">
                      <fa :icon="slotProps.option.icon" class="mr-2" />
                      <div>{{ slotProps.option.name }}</div>
                  </div>
              </template>
          </Listbox>
      </div>
      <div class="col-8 flex pt-3">
          <Panel header="Meals" style="width:100%;">
                <InputText v-model="searchtext" placeholder="Search" class="mb-4" />
                <div class="flex flex-wrap justify-content-center">
                    <MealCard  v-for="(item,index) in products" :key="index" :name="item.name" class="m-2 w-5" @addwithcomment="visible=true;namewithcomment=item.name" @add="orderItems.push({name:item.name,comment})"/>
                </div>
          </Panel>
      </div>
      <div class="col-2 flex pt-3">
        <Panel header="Order Items" class="w-12">
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


const searchtext = ref("")

const addWithComment = () => {
    orderItems.value.push({name:namewithcomment,comment})
    visible.value=false
}

const goOrder = () => {

    axios.post("http://localhost:8000/api/submitorder",
        {
            items:orderItems.value   
        }
    ).then((response) => {
        console.log(response.data)
        toast.add({ severity: 'success', summary: 'Success', detail: 'Order in progress !', life: 3000,group:'br' });
    })

    orderItems.value = []
};


watch(searchtext, (newSearchText) => {
  console.log(`x is ${newSearchText}`)
})


const orderItems = ref([
])
  
const products = ref([
    {name: "Fried Chicken Single"},
    {name: "Fried chicken double"},
    {name: "Fried chicken triple"},
    {name: "Burger single"},
    {name: "Burger double"},
    {name: "Burger triple"},
    {name: "Crepe Mix Chicken"},
    {name: "Crepe Kofta"},
    {name: "Crepe Pane"},
])


  const items = ref([
      {
          label: 'Home',
          icon: 'pi pi-home'
      },
      {
          label: 'Features',
          icon: 'pi pi-star'
      },
      {
          label: 'Projects',
          icon: 'pi pi-search',
          items: [
              {
                  label: 'Components',
                  icon: 'pi pi-bolt'
              },
              {
                  label: 'Blocks',
                  icon: 'pi pi-server'
              },
              {
                  label: 'UI Kit',
                  icon: 'pi pi-pencil'
              },
              {
                  label: 'Templates',
                  icon: 'pi pi-palette',
                  items: [
                      {
                          label: 'Apollo',
                          icon: 'pi pi-palette'
                      },
                      {
                          label: 'Ultima',
                          icon: 'pi pi-palette'
                      }
                  ]
              }
          ]
      },
      {
          label: 'Contact',
          icon: 'pi pi-envelope'
      }
  ]);
  
  const selectedCategory = ref({ name: 'All', icon:'star-of-life' });
  const categories = ref([
      { name: 'All', icon:'star-of-life' },
    //   { name: 'Towers', icon:'burger' },
    //   { name: 'Sandwitches', icon: 'hotdog' },
    //   { name: 'Crepe', icon: 'play' },
    //   { name: 'Dessert', icon: 'cake-candles' },
    //   { name: 'Pizza', icon: 'pizza-slice' },
  ]);
  
  </script>
  
  <style>
  </style>
  