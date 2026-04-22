import { createApp } from 'vue'
import App from './App.vue'

import PrimeVue from 'primevue/config';
import Aura from '@primeuix/themes/aura';
import 'primeicons/primeicons.css'
import { definePreset } from '@primeuix/themes';



import ToastService from 'primevue/toastservice';
import ConfirmationService from 'primevue/confirmationservice';
import { createI18n } from 'vue-i18n'
import { dt } from '@primevue/themes';
import '@fortawesome/fontawesome-free/css/fontawesome.css';
import '@fortawesome/fontawesome-free/css/regular.css';
import '@fortawesome/fontawesome-free/css/solid.css';
import '@fortawesome/fontawesome-free/css/brands.css';





import {  createWebHistory, createRouter } from 'vue-router'

import { createPinia } from 'pinia'
import auth from "@/services/auth";
import Tooltip from 'primevue/tooltip';
import StyleClass from 'primevue/styleclass';
import Ripple from 'primevue/ripple';

const routes = [
  {
    path: '/setup',
    component: () => import('@/pages/Setup.vue'),
  },
  {
    path: '/admin-setup',
    component: () => import('@/pages/AdminSetup.vue'),
  },
  {
    path: '/login',
    component: () => import('@/pages/Login.vue'),
  },
  {
    path: '/no-access', 
    component: ()=>{
        return import('@/pages/NoAccessView.vue')
    },
  },
  { 
    path: '/', alias:['/home'], 
    component: () => {
      if (!auth.isAuthenticated.value) {
        window.location.href = '/login'
        return import('@/pages/Login.vue')
      }
      if (auth.hasRole("admin") || auth.hasRole("superuser") || auth.hasRole("cashier") ) {
        return import('@/pages/Home.vue')
      }
      return import('@/pages/NoAccessView.vue')
    }
  },
  { 
    path: '/kitchen', component: () => {
      if (!auth.isAuthenticated.value) {
        window.location.href = '/login'
        return import('@/pages/Login.vue')
      }
      if (auth.hasRole("admin") || auth.hasRole("superuser") || auth.hasRole("chef")) {
        return import('@/pages/Kitchen.vue')
      }
      return import('@/pages/NoAccessView.vue')
    } 
  },
{ 
    path: '/admin', 
    component: () => {
      if (!auth.isAuthenticated.value) {
        window.location.href = '/login'
        return import('@/pages/Login.vue')
      }
      if (auth.hasRole("admin") || auth.hasRole("superuser")) {
        return import('@/pages/Admin.vue')
      }
      return import('@/pages/NoAccessView.vue')
    },
    children: [
      {
        path: '',
        redirect: { path: '/admin/inventory' }
      },
      {path: 'inventory', component: () => import('@/pages/Inventory.vue')},
      {path: 'sales', component: () => import('@/pages/Sales.vue')},
      {path: 'products', component: ()=> import('@/pages/Products.vue')},
      {path: 'categories', component: () => import('@/pages/Categories.vue')},
      {path: 'orders',
      children:[
        {path: '', component: () => import('@/pages/Orders.vue')},
      ]},
      {path: 'settings', component: () => import('@/pages/Settings.vue')},
      {path: 'customers', component: () => import('@/pages/Customers.vue')},
      {path: 'hubsync', component: () => {
        if (auth.hasRole("superuser")) {
          return import('@/pages/Hubsync.vue')
        }
        return import('@/pages/NoAccessView.vue')
      }},
    ],
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes: routes,
})

declare module 'vue' {
  interface ComponentCustomProperties {
      $auth: typeof auth
  }
}

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  fallbackLocale: 'en',
  messages: {
    en: {
      "cashier":"Cashier",
      "kitchen":"Kitchen",
      "admin":"Admin",
      "inventory":"Inventory",
      "product": "Product | Products",
      "order":"Order | Orders",
      "order_items":"Order Items",
      "total":"Total",
      "subtotal":"Subtotal",
      "discount":"Discount",
      "egp":"EGP",
      "search":"Search",
      "signout":"Signout",
      "notifications":"Notifications",
      "clear_all":"Clear All",
      "stashed_orders":"Stashed Orders",
      "chats":"Chats",
      "messages":"Messages",
      "write_message":"Write Message",
      "paylater_orders":"Paylater Orders",
      "checkout":"Checkout",
      "category":"Category | Categories",
      "add_component":"Add Component",
      "name":"Name",
      "quantity":"Quantity",
      "unit":"Unit",
      "status":"Status",
      "actions":"Actions",
      "history":"History",
      "list":"List",
      "report":"Report | Reports",
      "settings":"Settings",
      "language":"Language | Languages",
      "sales":"Sales"
    }
  }
})

const Noir = definePreset(Aura, {
  components: {
    progressspinner: {
        colorScheme: {
            light: {
                root: {
                    colorOne: '{primary.900}',
                    colorTwo: '{primary.900}',
                    colorThree: '{primary.900}',
                    colorFour: '{primary.900}'
                }
            },
            dark: {
                root: {
                    colorOne: '{primary.900}',
                    colorTwo: '{primary.900}',
                    colorThree: '{primary.900}',
                    colorFour: '{primary.900}'
                }
            }
        }
    }
  },
  semantic: {
      primary: {
          50: '{zinc.50}',
          100: '{zinc.100}',
          200: '{zinc.200}',
          300: '{zinc.300}',
          400: '{zinc.400}',
          500: '{zinc.500}',
          600: '{zinc.600}',
          700: '{zinc.700}',
          800: '{zinc.800}',
          900: '{zinc.900}',
          950: '{zinc.950}'
      },
      colorScheme: {
          light: {
              primary: {
                  color: '#14977B',
                  inverseColor: '#a5c22f',
                  hoverColor: '#14977B',
                  activeColor: '#14977B'
              },
              highlight: {
                  background: '#DEDB69',
                  focusBackground: '#DEDB69',
                  color: '#173350',
                  focusColor: '#173350'
              }
          },
          dark: {
              primary: {
                  color: '#a5c22f',
                  inverseColor: '#14977B',
                  hoverColor: '#a5c22f',
                  activeColor: '#a5c22f'
              },
              highlight: {
                  background: '#fff6c7',
                  focusBackground: '#FFDC00',
                  color: '#173350',
                  focusColor: '#173350'
              }
          },
      }
  }
});

const app = createApp(App).use(createPinia())
app.config.globalProperties.$auth = auth

app
.use(router)
.use(PrimeVue,{
    theme: {
        preset: Noir,
        options: {
            prefix: 'p',
            darkModeSelector: '.my-app-dark',
            cssLayer: false
        }
    }
})
.use(ToastService)
.use(ConfirmationService)
.use(i18n)
.directive('tooltip', Tooltip)
.directive('styleclass', StyleClass)
.directive('ripple', Ripple)
.mount('#app')
