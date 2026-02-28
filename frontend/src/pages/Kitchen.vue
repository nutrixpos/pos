<template>
    <div v-if="!loading" style="overflow-x: hidden;">
        <div class="grid">
            <div class="col-12 pt-5 pl-3">
                <div v-if="orders.length <= 0" style="height:100vh;width:100%;" class="flex justify-content-center align-items-center">
                    <h1 style="color:cadetblue">0 {{$t('order',3)}}</h1>
                </div>
                <div class="flex flex-wrap">
                    <div v-for="column_orders,index in dynamic_columns" :key="index">
                        <div class="flex flex-column gap-1">
                            <QueueOrder @finished="orderFinished(order)" @openedDialog="openedDialogs++" @closedDialog="openedDialogs--" v-for="(order,i) in column_orders" :key="i" :order="order" :number="i+1" class="queue-order"/>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    <div style="width:100vw;height:100vh;display:flex;justify-content:center;align-items:center" v-if="loading">
      <ProgressSpinner style="width: 35px; height: 35px;stroke:blue !important;" strokeWidth="6" fill="transparent"
      animationDuration=".5s" aria-label="Custom ProgressSpinner" />
    </div>
</template>

<script setup lang="ts">
import QueueOrder from '@/components/QueueOrder.vue'
import axios from 'axios';
import {ref,getCurrentInstance,computed,useTemplateRef, nextTick,watch} from 'vue';
import { useToast } from "primevue/usetoast";
import { Notification} from '@/classes/Notification';
import { globalStore } from '@/stores';
import { useI18n } from 'vue-i18n'
import ProgressSpinner from "primevue/progressspinner";

const store = globalStore()

const { proxy } = getCurrentInstance();
const toast = useToast();


const notifications = ref<Notification[]>([])

const orders = ref([])
const openedDialogs = ref(0)

const chat_text = ref("")
const chats = ref<any[]>([])
const chat_container = useTemplateRef("chat_container")


const user : any = computed(() => {

return proxy.$zitadel?.oidcAuth.userProfile

})


let socket : WebSocket
const orientation = computed(() => store.currentOrientation)


const loading = ref(true)
const { locale,setLocaleMessage } = useI18n({ useScope: 'global' })

const dynamic_columns = ref<Array<Array<any>>>([])
const prepareLayout = () => {
    const screenWidth = window.innerWidth;
    dynamic_columns.value = []
    for(var i=0;i<parseInt(`${screenWidth / 16 / 20}rem`);i++){
        dynamic_columns.value.push([])
    }
}

watch(
    () => orders.value,
    (newVal) => {
        displayOrders()
    },
    {
        deep: true,
    }
);

const loadLanguage = async () => {

    await axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/settings`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        },
    })
    .then(async (response)=>{
        await axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/languages/${response.data.data.language.code}`, {
            headers: {
                Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
            }
        })
        .then(response2 => {

            setLocaleMessage(response2.data.data.code, response2.data.data.pack);
            locale.value = response2.data.data.code
            store.setOrientation(response2.data.data.orientation)
            loading.value = false
        })
        .catch((err) => {
            console.log(err)
        });
        loading.value = false
    })
    .catch((err) => {
        console.log(err)
    });

}


const startWebsocket = () => {
    socket = new WebSocket(`ws://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/ws`);
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
                toast.add({ severity: 'success', summary: 'Order Finished', detail: `order with id ( ${data.order_id} ) finished and is ready to be served !`, life: 3000,group:'br' });

                const notification = new Notification();
                notification.description = `order with id #${data.order_id} finished and is ready to be served !`
                notification.severity = "success"
                notification.topic_name = "Order Finished"
                notification.type = "topic_message"
                notifications.value.push(notification);


                orders.value = orders.value.filter(o => o.id !== data.order_id)
            } else if (data.topic_name == "order_submitted") {
    
                orders.value.push(data.order)
                displayOrders()
            }else {
                const notification = new Notification();
                notification.description = data.message
                notification.severity = data.severity
                notification.topic_name = data.topic_name
                notification.type = data.type
                notifications.value.push(notification);

                toast.removeGroup('br')
                toast.add({ severity: data.severity, summary: data.topic_name, detail: data.message, life: 30000,group:'br' });

                chats.value.push({
                    type:"notification",
                    message: data.message,
                    topic_name: data.topic_name
                })
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

const orderFinished = (order) => {
    orders.value = orders.value.filter(o => o.id !== order.id)
}


const loadOrders =  () => {
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders?filter[state]=!finished&filter[state]=!stashed`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then((result)=>{
        if (result.data.data == null){
            orders.value = []
        }else {
            orders.value = result.data.data
            displayOrders()
        }
    })
};

const displayOrders = () => {

    for (var i=0;i<dynamic_columns.value.length;i++){
        dynamic_columns.value[i] = []
    }

    for(var i=0;i<orders.value.length;i++){
        dynamic_columns.value[i % dynamic_columns.value.length].push(orders.value[i])
    }
}


prepareLayout()
loadOrders()
loadLanguage()
startWebsocket()

</script>

<style>
.queue-order {
    margin:5px;
}
</style>