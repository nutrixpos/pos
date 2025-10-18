<template>
    <div v-if="!loading" style="overflow-x: hidden;">
        <div class="grid">
            <div class="col-8 pt-5 pl-3">
                <div v-if="orders.length <= 0" style="height:100vh;width:100%;" class="flex justify-content-center align-items-center">
                    <h1 style="color:cadetblue">0 {{$t('order',3)}}</h1>
                </div>
                <div class="flex flex-wrap">
                    <QueueOrder @finished="orderFinished(index)" @openedDialog="openedDialogs++" @closedDialog="openedDialogs--" v-for="(order,index) in orders" :key="index" :order="order" :number="index+1" class="queue-order"/>
                </div>
            </div>
            <div class="col-4" :style="`background-color:${$dt('slate.50').value};height:100vh;position:fixed;${orientation == 'ltr' ? 'right:0' : 'left:0'};top:0`">
                <div class="p-3">
                    <!-- <div style="height:100%;width:100%;" class="flex align-items-center justify-content-center">
                        <i class="pi pi-comments" style="font-size:4rem;color:darkgray"></i>
                        <span class="ml-5" style="color:darkgray;font-size:2rem">/</span>
                        <i class="pi pi-bell ml-5" style="font-size:4rem;color:darkgray"></i>
                    </div> -->
                    <div style="height:calc(100vh - 5rem);overflow-y:auto;" ref="chat_container">
                        <div v-for="(chat,index) in chats" :key="index">
                            <Message severity="success" v-if="chat.type == 'chat_message' && chat.user_sub != user.sub" class="m-1">
                                <template #container>
                                    <div class="p-3 flex flex-column">
                                        <strong>{{ chat.sender_name }}</strong>
                                        <span class="pt-2 px-2">{{ chat.message }}</span>
                                    </div>
                                </template>
                            </Message>

                            <Message severity="info" v-if="chat.type == 'chat_message' && chat.user_sub == user.sub" class="m-1">
                                <template #container>
                                    <div class="p-3 flex flex-column">
                                        <strong>{{ chat.sender_name }}</strong>
                                        <span class="pt-2 px-2">{{ chat.message }}</span>
                                    </div>
                                </template>
                            </Message>

                            <Message severity="warn" v-if="chat.type == 'notification'" class="m-1">
                                <template #container>
                                    <div class="p-3 flex flex-column">
                                        <strong>{{ chat.topic_name }}</strong>
                                        <span class="pt-2 px-2">{{ chat.message }}</span>
                                    </div>
                                </template>
                            </Message>

                        </div>
                    </div>
                    <InputGroup class="mt-2">
                        <InputText v-model="chat_text" :placeholder="$t('write_message')" @keyup.enter="SendChatMessage(chat_text)" />
                        <Button icon="pi pi-send" severity="info" @click="SendChatMessage(chat_text)" />
                    </InputGroup>
                </div>
            </div>
        </div>
    </div>
    <div style="width:100vw;height:100vh;display:flex;justify-content:center;align-items:center" v-if="loading">
      <ProgressSpinner style="width: 50px; height: 50px" strokeWidth="8" fill="transparent"
      animationDuration=".5s" aria-label="Custom ProgressSpinner" />
    </div>
</template>

<script setup lang="ts">
import QueueOrder from '@/components/QueueOrder.vue'
import axios from 'axios';
import {ref,getCurrentInstance,computed,useTemplateRef, nextTick} from 'vue';
import { useToast } from "primevue/usetoast";
import { Notification} from '@/classes/Notification';
import InputGroup from 'primevue/inputgroup';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import Message from 'primevue/message';
import { globalStore } from '@/stores';
import { useI18n } from 'vue-i18n'
import ProgressSpinner from "primevue/progressspinner";
import { $dt } from '@primevue/themes';




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


const SendChatMessage = (msg: string) => {
    if (socket.readyState === WebSocket.OPEN) {
        socket.send(`{"type":"chat_message","message":"${msg}","sender_name":"${user.value.name}","user_sub":"${user.value.sub}","to":"*","date": "${new Date().toLocaleString()}"}`)
    }else {
        console.log("WS closed")
    }
    chat_text.value = ""
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
            } else if (data.topic_name == "order_submitted") {
    
                orders.value.push(data.order)

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


        if (data.type == "chat_message") {

            chats.value.push({
                message:data.message,
                sender_name: data.sender_name,
                user_sub: data.user_sub,
                date: data.date,
                type:"chat_message"
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



const orderFinished = (index) => {
    if (orders.value.length > 0)
        orders.value.splice(index,1)
}


const loadOrders =  () => {
    axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/orders?filter[is_finished]=false&filter[state]=!stashed`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        }
    })
    .then((result)=>{
        if (result.data.data == null){
            orders.value = []
        }else {
            orders.value = result.data.data
        }
    })
};


loadOrders()


loadLanguage()
loadOrders()
startWebsocket()

</script>

<style>
.queue-order {
    margin:5px;
}
</style>