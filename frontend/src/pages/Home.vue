<template>
    <div class="flex flex-column m-2" style="height: 100%;">
        <Menubar :model="items">
            <template #item="{ item, props }">
                <a v-ripple class="flex align-items-center" v-bind="props.action" :href="item.link">
                    <span :class="item.icon" />
                    <span class="ml-2">{{ item.label }}</span>
                </a>
            </template>
            <template #end>
                <Button icon="pi pi-bell" severity="secondary" size="large" badge="0"  text rounded aria-label="Notifications" />
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
            <div class="lg:col-8 col-6 flex pt-3 pb-3">
                <Panel header="Recipes" style="width:100%;">
                    <InputText v-model="searchtext" placeholder="Search" class="mb-4" />
                    <div class="flex flex-wrap">
                        <MealCard  v-for="(item,index) in products" :key="index" :item="item" class="m-2" style="width:9rem;" @addwithcomment="visible=true;idwithcomment=item.id; namewithcomment=item.name" @add="addItem(item)"/>
                    </div>
                </Panel>
            </div>
            <div class="col-4 lg:col-2 flex pt-3 pb-3">
                <Panel header="Order Items" class="w-12" :style="`background-color:${is_order_valid ?  'white' : 'var(--red-100)'};border-color: ${is_order_valid ?  '' : 'red'};`">
                    <div class="flex flex-column" style="height:calc(100vh - 10rem)">
                        <div style="height:60vh;overflow: auto;">
                            <div v-for="(item,index) in orderItems" :key="index">
                                <div class="flex justify-content-between align-items-center">
                                    <div style="background-color:red;width:0.3rem;height:2.5rem;" class="mr-2" v-if="!item.isValid">
                                        &nbsp;
                                    </div>
                                    <p class="w-6" style="text-overflow:ellipsis"><strong>{{ item.product.name }}</strong></p>
                                    <p>{{ item.price }} EGP</p>
                                    <div>
                                        <Button icon="pi pi-pencil" size="small" style="width:2rem;height: 2rem;" aria-label="Edit" severity="secondary" @click="itemToEditIndex = index; edit_item_dialog=true" class="mr-1"/>
                                        <Button icon="pi pi-times" size="small" style="width:2rem;height: 2rem;" aria-label="Remove" severity="secondary" @click="orderItems.splice(index,1)" />
                                    </div>
                                </div>
                                <p class="m-0">{{ item.comment }}</p>
                            </div>
                        </div>
                        <div class="flex flex-column flex-wrap justify-content-between h-20rem">
                            <div>
                                <Divider />
                                <div class="flex justify-content-between flex-wrap align-items-center">
                                    <p>Subtotal : </p>
                                    <p style="font-size:1rem">{{ subtotal.toFixed(2) }} <span style="font-size:0.8rem">EGP</span></p>
                                </div>
                                <div class="flex justify-content-between flex-wrap align-items-center">
                                    <p class="my-0">Discount :</p>
                                    <div class="w-7 flex justify-content-end align-items-center">
                                        <InputText v-model="discount" :disabled="orderItems.length == 0" placeholder="0" type="number" class="w-8 h-2rem"  />
                                        <p style="font-size:0.8rem" class="ml-2">EGP</p>
                                    </div>
                                </div>
                                <div class="flex justify-content-center align-items-center">
                                    <Slider v-model="discount_percent" class="w-9 mt-1" style="height:0.6rem;" />
                                    <p class="ml-2" style="font-size:0.8rem">{{ discount_percent.toFixed(2) }} %</p>
                                </div>
                                <div class="flex justify-content-between flex-wrap align-items-center">
                                    <h2>Total : </h2>
                                    <p style="font-size:1.4rem">{{ total.toFixed(2) }} <span style="font-size:1rem">EGP</span></p>
                                </div>
                            </div>
                            <Button label="Checkout" :disabled="!is_order_valid" @click="goOrder" />
                        </div>
                    </div>
                </Panel>
            </div>
            <Dialog v-model:visible="edit_item_dialog" modal header="Edit item" class="xs:w-12 lg:w-6">
                <OrderItemView v-model="orderItems[itemToEditIndex]"  />
            </Dialog>
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

<script setup lang="ts">
  import Menubar from 'primevue/menubar';
  import Dialog from 'primevue/dialog';
  import Listbox from 'primevue/listbox';
  import Panel from 'primevue/panel';
  import InputText from 'primevue/inputtext';
  import Button from 'primevue/button';
  import { useToast } from "primevue/usetoast";
  import axios from 'axios'
  import OrderItemView from '@/components/OrderItemView.vue'
  import {OrderItem} from '@/classes/OrderItem'
  import Divider from 'primevue/divider';
  import Slider from 'primevue/slider';




  const toast = useToast();
  const itemToEditIndex = ref(0)
  const edit_item_dialog = ref(false)
  const is_order_valid = ref(true)

  
  import MealCard from '@/components/MealCard.vue';
  
  import { ref,watch } from "vue";
  
  
  const comment = ref("")
  const subtotal = ref(0)
  const discount = ref(0)
  const discount_percent = ref(0)
  const total = ref(0)
  const namewithcomment = ref("")
  const idwithcomment = ref("")
  const visible = ref(false)
  const selectedCategory = ref();
  
  
const searchtext = ref("")
const categories = ref([])

const orderItems = ref<OrderItem[]>([])


const addItem = async (item) => {

    const new_item = new OrderItem()
    new_item.product.name = item.name
    new_item.product.id = item.id
    await new_item.ReloadDefaults()


    orderItems.value.push(new_item)
}

const addWithComment = async () => {

    const new_item = new OrderItem()
    new_item.product.name = namewithcomment.value
    new_item.comment = comment.value
    new_item.product.id = idwithcomment.value
    await new_item.ReloadDefaults()

    orderItems.value.push(new_item)
    visible.value=false
    comment.value = ""
    idwithcomment.value = ""
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
                items:orderItems.value,
                discount, 
            }
        ).then((response) => {
            console.log(response.data)
            toast.add({ severity: 'success', summary: 'Success', detail: 'Order in progress !', life: 3000,group:'br' });
        })
    
        orderItems.value = []
    }
};


const isUpdatingDiscount = ref(false)
const isUpdatingDiscountPercent = ref(false)


watch(searchtext, (newSearchText) => {
  console.log(`x is ${newSearchText}`)
})

watch(subtotal, (new_subtotal) => {
  total.value = new_subtotal - discount.value
  if (total.value < 0)
    total.value = 0
})

watch(discount, (new_discount) => {
  if (!isUpdatingDiscountPercent.value){
    isUpdatingDiscount.value = true
    total.value = subtotal.value - new_discount

    if (new_discount != 0)
        discount_percent.value = new_discount*100 / subtotal.value
    if (total.value < 0)
    total.value = 0
    }else{
      isUpdatingDiscountPercent.value = false
  }
})
watch(discount_percent, (new_discount_percent) => {
 if (!isUpdatingDiscount.value){
    isUpdatingDiscountPercent.value = true
    discount.value = new_discount_percent * subtotal.value / 100
    total.value = subtotal.value - discount.value
    isUpdatingDiscount.value = false
  }else {
    isUpdatingDiscount.value = false
  }
})


watch(() => orderItems.value, 
  (currentValue) => {
    let x = 0
    let valid = true
    currentValue.forEach((item) => {

        x += item.price
        if (!item.isValid)
            valid=false

    })

    is_order_valid.value = valid
    subtotal.value = x
    discount.value = discount_percent.value * subtotal.value / 100
  },
  {deep: true}
);

  
const products = ref([
])



const refreshAvailabilities = () => {
    var product_ids = ""
    products.value.forEach((product,index) => {
        product_ids += index > 0 ? "," +product.id : product.id
    })

    axios.get("http://localhost:8000/api/recipeavailability?ids="+product_ids)
    .then((response) => {
        products.value.forEach((product,index) => {
            products.value[index].availability = Math.round(response.data.filter((x) => x.recipe_id == product.id)[0].available * 100) / 100
        })
    })
}


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
                id: recipe.id,
                name:recipe.name,
                price:recipe.price
            })
        })
        refreshAvailabilities();
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
          label: 'Admin',
          icon: 'pi pi-cog',
          link: 'admin',
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
  