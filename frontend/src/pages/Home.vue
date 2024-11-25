<template>
    <div class="flex flex-column m-2" style="height: 100%;">
        <Toolbar>
            <template #start>
                <router-link v-for="(item,index) in items" :key="index" :to="item.link">
                    <Button :icon="item.icon" :label="item.label"  text severity="secondary" />
                </router-link>
            </template>

            <template #center>
                <IconField iconPosition="left">
                    <InputIcon>
                        <i class="pi pi-search" />
                    </InputIcon>
                    <InputText placeholder="Search" v-model="mainSearchText" @keyup.stop="mainSearchTextChanged" />
                </IconField>
                <OverlayPanel ref="mainsearch_op" class="w-5 lg:w-3" style="max-height:60vh;overflow-y: auto;">
                    <MainSearchResultView @view-order-pressed="order_to_show = result; order_details_dialog=true" v-for="(result,index) in mainSearchResult" :key="index" :order="result" />
                </OverlayPanel>
            </template>

            <template #end>

                <Button  severity="secondary" size="large"  text rounded aria-label="PayLater" @click.stop="paylater_toggle">
                    <span class="p-button-icon pi pi-hourglass"></span>
                    <Badge :value="payLaterOrders.length" class="p-badge-warning"  />
                </Button>
                <OverlayPanel ref="paylater_orders_op" class="w-5 lg:w-3" style="max-height:60vh;overflow-y: auto;">
                    <h4 class="my-0 mx-2" style="color:#c2c2c2">Paylater Orders</h4>
                    <PayLaterOrder @order_paid="PayLaterOrderPaid(index)" :order="order" v-for="(order,index) in payLaterOrders" :key="index" />
                </OverlayPanel>
                <Button  severity="secondary" size="large"  text rounded aria-label="Stashed"  @click.stop="chats_toggle">
                    <span class="p-button-icon pi pi-comments"></span>
                    <Badge class="p-badge-danger" v-if="has_new_message"  />
                </Button>
                <OverlayPanel ref="chats_op" class="w-5 lg:w-4" style="max-height:90vh;overflow-y: auto;">
                    <Panel header="Chats">

                        <div style="height:40vh;overflow-y: auto;" ref="chat_container">                            
                            <div v-for="(chat,index) in chats" :key="index" :class="`flex ${chat.user_sub == user.sub ? 'justify-content-end' : ''}`">
                                <InlineMessage severity="success" v-if="chat.user_sub != user.sub" class="mt-2">
                                    <template #default>
                                        <div class="px-3 flex flex-column">
                                            <strong>{{ chat.sender_name }}</strong>
                                            <span class="pt-2 px-2">{{ chat.message }}</span>
                                        </div>
                                    </template>
                                </InlineMessage>
    
                                <InlineMessage severity="info" v-if="chat.user_sub == user.sub" class="mt-2">
                                    <template #default>
                                        <div class="px-3 flex flex-column">
                                            <strong>{{ chat.sender_name }}</strong>
                                            <span class="pt-2 px-2">{{ chat.message }}</span>
                                        </div>
                                    </template>
                                </InlineMessage>
                            </div>
                        </div>

                    </Panel>

                    <InputGroup class="mt-2">
                        <InputText v-model="chat_text" placeholder="Write message.." @keyup.enter="SendChatMessage(chat_text)" />
                        <Button icon="pi pi-send" severity="info" @click="SendChatMessage(chat_text)" />
                    </InputGroup>
                </OverlayPanel>
                <Button  severity="secondary" size="large"  text rounded aria-label="Stashed" label="Stashed"  @click.stop="stashed_toggle">
                    <span class="p-button-icon pi pi-bookmark"></span>
                    <Badge :value="stashedOrders.length" class="p-badge-success"  />
                </Button>
                <OverlayPanel ref="stashed_orders_op" class="w-5 lg:w-3" style="max-height:60vh;overflow-y: auto;">
                    <h4 class="my-0 mx-2" style="color:#c2c2c2">Stashed Orders</h4>
                    <StashedOrder :order="order" v-for="(order,index) in stashedOrders" :key="index" @back_to_checkout="BackStashedOrderToCheckout(index)" />
                </OverlayPanel>
                <Button  severity="secondary" size="large"  text rounded aria-label="Notifications" @click.stop="notifications_toggle">
                    <span class="p-button-icon pi pi-bell"></span>
                    <Badge :value="notifications_severity_counter[0]" class="p-badge-success"  />
                    <Badge :value="notifications_severity_counter[1]" class="p-badge-info"  />
                    <Badge :value="notifications_severity_counter[2]" class="p-badge-warning" />
                    <Badge :value="notifications_severity_counter[3]" class="p-badge-danger" />
                </Button>
                <OverlayPanel ref="notifications_op" class="w-3" style="max-height:60vh;overflow-y: auto;">
                    <h4 class="my-0 mx-2" style="color:#c2c2c2">Notifications</h4>
                    <Button text label="Clear all" severity="secondary" @click="clearNotifications()"/>
                    <NotificationView :notification="notification" v-for="(notification,index) in notifications" :key="index" />
                </OverlayPanel>
                <Button  severity="secondary" size="large"  text rounded aria-label="Profile" label="Profile" @click.stop="user_profile_toggle">
                    <span style="font-size:0.9rem;" class="mr-2">{{ user.name }}</span>
                    <span class="p-button-icon pi pi-user"></span>
                </Button>
                <OverlayPanel ref="user_profile_op" class="lg:w-2 md:w-3">
                    <div class="flex flex-column">
                        <span>Welcome <strong>{{ user.name }}</strong></span>
                        <div class="mt-2">
                            <Chip v-for="(role,index) in roles" :key="index" :label="role" style="height: 1.5rem;" class="m-1" />
                        </div>
                        <Button class="mt-5" icon="pi pi-sign-out" severity="secondary" text aria-label="Signout" label="Signout" @click="proxy.$zitadel.oidcAuth.signOut()" />
                    </div>
                </OverlayPanel>
            </template>
        </Toolbar>
        <div class="grid" style="flex-grow:1;">
            <div class="col-2">
                <Listbox  v-model="selectedCategory" :options="categories" optionLabel="name" class="w-full mt-2" filter>
                    <template #option="slotProps">
                        <div class="flex align-items-center">
                            <div>{{ slotProps.option.name }}</div>
                        </div>
                    </template>
                </Listbox>
            </div>
            <div class="lg:col-7 col-6 flex pt-3 pb-3">
                <Panel header="Recipes" style="width:100%;">
                    <InputText v-model="searchtext" placeholder="Search" class="mb-4" />
                    <div class="flex flex-wrap">
                        <MealCard  v-for="(item,index) in products" :key="index" :item="item" class="m-2" style="width:9rem;" @addwithcomment="visible=true;idwithcomment=item.id; namewithcomment=item.name" @add="addItem(item)"/>
                    </div>
                </Panel>
            </div>
            <div class="col-4 lg:col-3 flex pt-3 pb-3">
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
                        <div class="flex flex-column flex-wrap justify-content-between">
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
                            <ButtonGroup class="mb-4">
                                <Button icon="pi pi-bookmark" @click="stashOrder" severity="secondary" />
                                <Button class="ml-2" icon="pi pi-hourglass" @click="payLaterStart" severity="secondary" />
                            </ButtonGroup>
                            <Button label="Checkout" :disabled="!is_order_valid" @click="goOrder()" />
                        </div>
                    </div>
                </Panel>
            </div>
            <Dialog v-model:visible="edit_item_dialog" modal header="Edit item" class="xs:w-12 md:w-10 lg:w-8">
                <OrderItemView v-model="orderItems[itemToEditIndex]"  />
            </Dialog>
            <Dialog v-model:visible="order_details_dialog" modal header="Order details" class="xs:w-12 md:w-10 lg:w-8">
                <OrderView :order="order_to_show" @order-cancelled="order_details_dialog=false" />
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
  import InputIcon from 'primevue/inputicon';
  import IconField from 'primevue/iconfield';
  import Toolbar from 'primevue/toolbar';
  import Dialog from 'primevue/dialog';
  import Listbox from 'primevue/listbox';
  import Panel from 'primevue/panel';
  import InputText from 'primevue/inputtext';
  import Chip from 'primevue/chip';
  import Button from 'primevue/button';
  import { useToast } from "primevue/usetoast";
  import axios from 'axios'
  import OrderItemView from '@/components/OrderItemView.vue'
  import {OrderItem} from '@/classes/OrderItem'
  import Order from '@/classes/Order'
  import Divider from 'primevue/divider';
  import Slider from 'primevue/slider';
  import Badge from 'primevue/badge'
  import NotificationView from '@/components/NotificationView.vue';
  import OverlayPanel from 'primevue/overlaypanel';
  import { Notification} from '@/classes/Notification';
  import { ref,watch,computed,getCurrentInstance, nextTick, useTemplateRef  } from "vue";
  import StashedOrder from '@/components/StashedOrder.vue'
  import PayLaterOrder from '@/components/PayLaterOrder.vue'
  import InputGroup from 'primevue/inputgroup';
  import InlineMessage from 'primevue/inlinemessage'
  import MainSearchResultView from '@/components/MainSearchResultView.vue';
  import OrderView from '@/components/OrderView.vue';


  const { proxy } = getCurrentInstance();





const toast = useToast();
const itemToEditIndex = ref(0)
const edit_item_dialog = ref(false)
const is_order_valid = ref(true)
const chat_text = ref("")
const chats = ref<any[]>([])
const has_new_message = ref(false)
const chat_container = useTemplateRef("chat_container")
const mainSearchText = ref("")
const order_details_dialog = ref(false)
const order_to_show = ref<Order>()

import MealCard from '@/components/MealCard.vue';


const comment = ref("")
const subtotal = ref(0)
const discount = ref(0)
const discount_percent = ref(0)
const total = ref(0)
const namewithcomment = ref("")
const idwithcomment = ref("")
const visible = ref(false)
const selectedCategory = ref();

const stashedOrders = ref<Order[]>([])
const payLaterOrders = ref<Order[]>([])

const notifications_op = ref();
const stashed_orders_op = ref();
const chats_op = ref();
const user_profile_op = ref();
const paylater_orders_op = ref()
const mainsearch_op = ref()

const mainSearchResult = ref<any[]>([])


const user : any = computed(() => {

    return proxy.$zitadel.oidcAuth.userProfile

})


const mainSearchTextChanged = (event:any) => {

    if (mainSearchText.value != ""){
        mainsearch_op.value.show(event)

        axios.get(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/orders?display_id=${mainSearchText.value}`, {
            headers: {
                Authorization: `Bearer ${proxy.$zitadel.oidcAuth.accessToken}`
            }
        })
        .then(response => {
            mainSearchResult.value = []
            for (var i=0;i<response.data.orders.length;i++){
                mainSearchResult.value.push(response.data.orders[i])
            }
        })
        .catch(error => {
            console.error(error);
        });

    }else{
        mainSearchResult.value = []
        mainsearch_op.value.hide()
    }


    console.log(mainSearchText.value)
}


const PayLaterOrderPaid = (index:number) => {
    payLaterOrders.value.splice(index,1)
}

const payLaterStart = () => {
    goOrder(true)
}

const getUnpaidOrders = () => {
    axios.get(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/unpaidorders`,{
        headers: {
            Authorization: `Bearer ${proxy.$zitadel.oidcAuth.accessToken}`
        }
    })
    .then(response => {
        payLaterOrders.value = response.data.orders
    })
}

const claims : any = computed(() => {

    if (user.value) { 
        return Object.keys(user.value).map(key => ({
          key,
          value: user.value[key]
        }))
      }

      return []

})

const roles : any = computed(()=>{
    if (claims.value.length > 0){

        for (var i=0;i<claims.value.length;i++){
            if (claims.value[i].key == "urn:zitadel:iam:org:project:roles"){
                return Object.keys(claims.value[i].value).map(key => {
                    return key
                })
            }
        }
    }

    return []
})


const BackStashedOrderToCheckout = async (stashed_order_index:number) => {

    const order = stashedOrders.value[stashed_order_index]
    let tmp_subtotal = 0


    for (var index=0;index<order.items.length;index++){
        const tmp_order_item = new OrderItem()
        await tmp_order_item.FromItemData(order.items[index])
        await tmp_order_item.RefreshProductData()
        tmp_order_item.price = order.items[index].product.price
        tmp_subtotal+=order.items[index].price
        order.items[index] = tmp_order_item
    }


    subtotal.value = tmp_subtotal
    discount_percent.value =  isNaN(order.discount * 100 / tmp_subtotal) ? 0 : order.discount * 100 / tmp_subtotal
    orderItems.value = order.items


    // discount.value = order.discount == null || order.discount == undefined ? 0 : order.discount


    axios.post(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/orderremovestash`,{
        order_display_id:order.display_id
    },
    {
        headers:{
            Authorization: `Bearer ${proxy.$zitadel.oidcAuth.accessToken}`
        },
    },)
    .then(()=>{
        stashedOrders.value.splice(stashed_order_index,1)

        stashed_orders_op.value.toggle()
    })
    .catch(() => {
        toast.add({severity:'error', summary: 'Failed to remove stashed order from db', detail: "", life: 3000,group:'br'});
        stashed_orders_op.value.toggle()
    })


}


const user_profile_toggle = (event: any) => {
    user_profile_op.value.toggle(event);
}

const notifications_toggle = (event: any) => {
    notifications_op.value.toggle(event);
}

const stashed_toggle = (event: any) => {
    stashed_orders_op.value.toggle(event);
}

const paylater_toggle = (event: any) => {
    paylater_orders_op.value.toggle(event);
}


const chats_toggle = async (event: any) => {
    has_new_message.value = false
    chats_op.value.toggle(event);
    await nextTick()

    if (chat_container.value != null)
        chat_container.value.scrollTop = chat_container.value?.scrollHeight

}

const notifications = ref<Notification[]>([])


const getStashedOrders = () => {
    axios.get(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/ordergetstashed`,{
        headers:{
            Authorization: `Bearer ${proxy.$zitadel.oidcAuth.accessToken}`
        }
    }).then(async (response) => {

        const dataCopy = JSON.parse(JSON.stringify(response.data))

        for (var i=0;i<dataCopy.length;i++){

            let order = new Order()
            order = JSON.parse(JSON.stringify(dataCopy[i]))
            order.items = []


            for (var j=0;j<dataCopy[i].items.length;j++){

                const item = new OrderItem()
                await item.FromItemData(dataCopy[i].items[j])

                order.items.push(item)
            }

            stashedOrders.value.push(order)
        }
    }).catch(() => {
        toast.add({severity:'error', summary: 'Failed to get stashed orders', detail: "", life: 3000,group:'br'});
    })
}


const stashOrder = () => {

    const order = new Order()
    order.items = orderItems.value
    order.discount = discount.value == null ? 0 : discount.value


    axios.post(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/orderstash`,
    {
        order:order
    },
    {
        headers:{
            Authorization: `Bearer ${proxy.$zitadel.oidcAuth.accessToken}`
        },
    }).then(async (response) => {
        orderItems.value=[]
        discount.value = 0
        discount_percent.value = 0
        total.value =0
        subtotal.value = 0


        console.log(response)


        for (var index=0;index<response.data.order.items.length;index++){
            const temp_order_item = new OrderItem()
            await temp_order_item.FromItemData(response.data.order.items[index])
            await temp_order_item.RefreshProductData()
            temp_order_item.ValidateItem()
            response.data.order.items[index] = temp_order_item
        }

        
        stashedOrders.value.push(response.data.order)
        toast.add({severity:'success', summary: `Order ${order.display_id} stashed successfully !`, detail: "successfully stashed order !", life: 3000,group:'br'});
    }).catch( () => {
        toast.add({severity:'error', summary: 'Error Stashing Item', detail: "", life: 3000,group:'br'});
    })
}


const clearNotifications = () => {
    notifications.value = []
}

let socket : WebSocket


const SendChatMessage = (msg: string) => {
    if (socket.readyState === WebSocket.OPEN) {
        socket.send(`{"type":"chat_message","message":"${msg}","sender_name":"${user.value.name}","user_sub":"${user.value.sub}","to":"*","date": "${new Date().toLocaleString()}"}`)
    }else {
        console.log("WS closed")
    }
    chat_text.value = ""
}


const startWebsocket = () => {
    socket = new WebSocket(`ws://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/ws`);
    socket.onopen = () => {
        console.log("Opened ws connection");
        socket.send(`{"type":"subscribe","topic_name":"all"}`);
    }

    socket.onmessage = async (event) => {
        console.log("Received message: " + event.data);

        const data = JSON.parse(event.data);

        if (data.type == "topic_message") {
            if (data.topic_name == "order_finished"){

                toast.removeGroup('br')
                toast.add({ severity: 'success', summary: 'Order Finished', detail: `order ( ${data.order_id} ) finished and is ready to be served !`, life: 3000,group:'br' });

                const notification = new Notification();
                notification.description = `order (${data.order_id}) finished and is ready to be served !`
                notification.severity = "success"
                notification.topic_name = "Order Finished"
                notification.type = "topic_message"
                notifications.value.push(notification);


                getUnpaidOrders()
            }else {
                const notification = new Notification();
                notification.description = data.message
                notification.severity = data.severity
                notification.topic_name = data.topic_name
                notification.type = data.type
                notifications.value.push(notification);

                toast.removeGroup('br')
                toast.add({ severity: data.severity, summary: data.topic_name, detail: data.message, life: 30000,group:'br' });
            }
        }

        if (data.type == "chat_message") {

            if (!chats_op.value?.visible) {
                has_new_message.value = true
            }

            chats.value.push({
                message:data.message,
                sender_name: data.sender_name,
                user_sub: data.user_sub,
                date: data.date,
            })


            if (chat_container.value != null){
                await nextTick()
                chat_container.value.scrollTop = chat_container.value?.scrollHeight + 100
            }
            
        }

    }
    socket.onerror = (event) => {
        console.log("Error occurred");
        console.log(event);
    }
    socket.onclose = () => {
        console.log("Connection closed");
        const retryConnection = async () => {
            if (socket.readyState !== WebSocket.OPEN) {
                await new Promise(r => setTimeout(r, 5000));
                console.log("Reconnecting to WebSocket...");
                startWebsocket()
            }
        }
        retryConnection();
    }
}

const init = () => {
    startWebsocket()
    getStashedOrders()
    getUnpaidOrders()

    console.log("user:")
    console.log(user.value)

    console.log("claims:")
    console.log(claims.value)
}

init()

// const notifications_severity_counter = ref<number[]>([])

const notifications_severity_counter = computed(() => {
    const counter = [0,0,0,0]
    notifications.value.forEach(notification => {
        switch (notification.severity) {
            case "success":
                counter[0]++;
                break;
            case "info":
                counter[1]++;
                break;
            case "warn":
                counter[2]++;
                break;
            case "error":
                counter[3]++;
                break;
        }
    })
    return counter
})

  
  
const searchtext = ref("")
const categories = ref([])

const orderItems = ref<OrderItem[]>([])


const addItem = async (item) => {

    const new_item = new OrderItem()
    new_item.product.name = item.name
    new_item.product.id = item.id
    await new_item.ReloadDefaults()
    new_item.ValidateAllMaterials()
    new_item.ValidateItem()


    orderItems.value.push(new_item)
}

const addWithComment = async () => {

    const new_item = new OrderItem()
    new_item.product.name = namewithcomment.value
    new_item.comment = comment.value
    new_item.product.id = idwithcomment.value
    await new_item.ReloadDefaults()
    new_item.ValidateItem()

    orderItems.value.push(new_item)
    visible.value=false
    comment.value = ""
    idwithcomment.value = ""
}


const getCategories = async () => {
    const response = await axios.get(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/categories`,{
        headers:{
            Authorization: `Bearer ${proxy.$zitadel.oidcAuth.accessToken}`
        }
    })
    categories.value = categories.value.concat(response.data.categories)
    selectedCategory.value = categories.value[0]
}

getCategories();


const goOrder = (isPaylater:boolean = false) => {

    let order : any =  {
        items:orderItems.value,
        discount:discount.value,
    }

    if (isPaylater) {
        order.is_pay_later = true
    }

    if (orderItems.value.length > 0){
        axios.post(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/submitorder`,
            order,
            {
                headers:{
                    Authorization: `Bearer ${proxy.$zitadel.oidcAuth.accessToken}`
                },
            },
        ).then((response) => {
            console.log(response.data)
            toast.add({ severity: 'success', summary: 'Success', detail: 'Order in progress !', life: 3000,group:'br' });

            if (isPaylater){
                getUnpaidOrders()
            }
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

    axios.get(`http://${process.env.VUE_APP_BACKEND_HOST}${process.env.VUE_APP_MODULE_CORE_API_PREFIX}/api/recipeavailability?ids=`+product_ids,{
        headers:{
            Authorization: `Bearer ${proxy.$zitadel.oidcAuth.accessToken}`
        }
    })
    .then((response) => {
        products.value.forEach((product,index) => {
            products.value[index].availability = Math.round(response.data.filter((x) => x.recipe_id == product.id)[0].available * 100) / 100
        })
    })
}


watch(selectedCategory, (category) => {
    if (category != null){
        products.value = []
        category.products.forEach((recipe) => {
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
  