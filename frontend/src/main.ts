import { createApp } from 'vue'
import App from './App.vue'


import '@/assets/styles.scss'
import PrimeVue from 'primevue/config';
import 'primevue/resources/themes/aura-light-green/theme.css'
import 'primevue/resources/primevue.min.css'
import 'primeicons/primeicons.css'


import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome';
import ToastService from 'primevue/toastservice';
import ConfirmationService from 'primevue/confirmationservice';



// library.add(fas)



import {  createWebHistory, createRouter } from 'vue-router'

import Home from '@/pages/Home.vue'
import Kitchen from '@/pages/Kitchen.vue'
import Admin from '@/pages/Admin.vue'
import Inventory from '@/pages/Inventory.vue'
import Sales from '@/pages/Sales.vue'
import Products from '@/pages/Products.vue'
import Categories from '@/pages/Categories.vue'
import Orders from '@/pages/Orders.vue'
import { createPinia } from 'pinia'
import zitadelAuth from "@/services/zitadelAuth";
import Tooltip from 'primevue/tooltip';


const routes = [
  { 
    path: '/', alias:['/home'], 
    meta: {
      authName: zitadelAuth.oidcAuth.authName
    },
    component: () => {

      if (zitadelAuth.hasRole("admin") || zitadelAuth.hasRole("cashier") ) {
        return Home 
      }
      return import('@/pages/NoAccessView.vue')
    }
  },
  { 
    path: '/kitchen', component: () => {

      if (zitadelAuth.hasRole("admin") || zitadelAuth.hasRole("chef")) {
        return Kitchen
      }
      return import('@/pages/NoAccessView.vue')

    } 
  },
  { 
    path: '/admin', 
    component: () => {
      if (zitadelAuth.hasRole("admin")) {
        return Admin
      }
      return import('@/pages/NoAccessView.vue')
    },
    children: [
      {
        path: '',
        redirect: { path: '/admin/inventory' }
      },
      {path: 'inventory', component: Inventory,},
      {path: 'sales', component: Sales,},
      {path: 'products', component: Products,},
      {path: 'categories', component: Categories,},
      {path: 'orders', component: Orders,},
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

declare module 'vue' {
  interface ComponentCustomProperties {
      $zitadel: typeof zitadelAuth
  }
}


zitadelAuth.oidcAuth.useRouter(router)

zitadelAuth.oidcAuth.startup().then(ok => {
  if (ok) {
        const app = createApp(App).use(createPinia())
        app.config.globalProperties.$zitadel = zitadelAuth

        app
        .use(router)
        .use(PrimeVue)
        .use(ToastService)
        .use(ConfirmationService)
        .directive('tooltip', Tooltip)
        .mount('#app')
  } else {
      console.error('Startup was not ok')
  }
})

 