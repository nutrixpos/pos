import { createApp } from 'vue'
import App from './App.vue'


import PrimeVue from 'primevue/config';
import 'primevue/resources/themes/aura-light-green/theme.css'
import 'primevue/resources/primevue.min.css'
import 'primeflex/primeflex.css'
import 'primeicons/primeicons.css'


import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import { library } from "@fortawesome/fontawesome-svg-core";
import { fas } from '@fortawesome/free-solid-svg-icons';
import ToastService from 'primevue/toastservice';
import ConfirmationService from 'primevue/confirmationservice';



library.add(fas)



import {  createWebHistory, createRouter } from 'vue-router'

import Home from '@/pages/Home.vue'
import Kitchen from '@/pages/Kitchen.vue'
import Analytics from '@/pages/Analytics.vue'

const routes = [
  { path: '/', alias:['/home'], component: Home },
  { path: '/kitchen', component: Kitchen },
  { path: '/analytics', component: Analytics },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})



const app = createApp(App)
app
.use(router)
.use(PrimeVue)
.use(ToastService)
.use(ConfirmationService)
.component('fa', FontAwesomeIcon)
.mount('#app')

 