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
import '@fortawesome/fontawesome-free/css/all.css';


// library.add(fas)



import {  createWebHistory, createRouter } from 'vue-router'

import { createPinia } from 'pinia'
import zitadelAuth from "@/services/zitadelAuth";
import Tooltip from 'primevue/tooltip';
import StyleClass from 'primevue/styleclass';
import Ripple from 'primevue/ripple';

const insecureRoutes = [
  {
    path: '/no-access', 
    component: ()=>{
        return import('@/pages/NoAccessView.vue')
    },
  },
  { 
    path: '/', alias:['/home'], 
    component: () => {
      return import('@/pages/Home.vue')
    }
  },
  { 
    path: '/kitchen', component: () => {
      return import('@/pages/Kitchen.vue')
    } 
  },
  { 
    path: '/admin', 
    component: () => {
      return import('@/pages/Admin.vue')
    },
    children: [
      {
        path: '',
        redirect: { path: '/admin/inventory' }
      },
      {path: 'inventory', component: () => import('@/pages/Inventory.vue'),},
      {path: 'sales', component: () => import('@/pages/Sales.vue'),},
      {path: 'products', component: ()=> import('@/pages/Products.vue'),},
      {path: 'categories', component: () => import('@/pages/Categories.vue'),},
      {path: 'orders',
      children:[
        {path: '', component: () => import('@/pages/Orders.vue'),},
      ]},
      {path: 'settings', component: () => import('@/pages/Settings.vue'),},
      {path: 'customers', component: () => import('@/pages/Customers.vue'),},
      {path: 'hubsync', component: () => import('@/pages/Hubsync.vue'),},
    ],
  },
]


const zitadelRoutes = [
  {
    path: '/no-access', 
    component: ()=>{
        return import('@/pages/NoAccessView.vue')
    },
  },
  { 
    path: '/', alias:['/home'], 
    meta: {
      authName: zitadelAuth.oidcAuth.authName
    },
    component: () => {

      if (zitadelAuth.hasRole("admin") || zitadelAuth.hasRole("cashier") ) {
        return import('@/pages/Home.vue')
      }
      return import('@/pages/NoAccessView.vue')
    }
  },
  { 
    path: '/kitchen', component: () => {

      if (zitadelAuth.hasRole("admin") || zitadelAuth.hasRole("chef")) {
        return import('@/pages/Kitchen.vue')
      }
      return import('@/pages/NoAccessView.vue')

    } 
  },
  { 
    path: '/admin', 
    component: () => {
      if (zitadelAuth.hasRole("admin")) {
        return import('@/pages/Admin.vue')
      }
      return import('@/pages/NoAccessView.vue')
    },
    children: [
      {
        path: '',
        redirect: { path: '/admin/inventory' }
      },
      {path: 'inventory', component: () => import('@/pages/Inventory.vue'),},
      {path: 'sales', component: () => import('@/pages/Sales.vue'),},
      {path: 'products', component: ()=> import('@/pages/Products.vue'),},
      {path: 'categories', component: () => import('@/pages/Categories.vue'),},
      {path: 'orders',
      children:[
        {path: '', component: () => import('@/pages/Orders.vue'),},
      ]},
      {path: 'settings', component: () => import('@/pages/Settings.vue'),},
      {path: 'customers', component: () => import('@/pages/Customers.vue'),},
      {path: 'hubsync', component: () => import('@/pages/Hubsync.vue'),},
    ],
  },
]

const zitadelRouter = createRouter({
  history: createWebHistory(),
  routes: zitadelRoutes,
})

const insecureRouter = createRouter({
  history: createWebHistory(),
  routes: insecureRoutes,
})

declare module 'vue' {
  interface ComponentCustomProperties {
      $zitadel: typeof zitadelAuth
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
  // components: {
  //   button: {
  //     colorScheme: {
  //       light: {
  //           root: {
  //               primary: {
  //                 color: '#2e4762',
  //                 hover: {
  //                   color: '#2e4762',
  //                 },
  //                 active: {
  //                   color: '#2e4762',
  //                 }
  //               }
  //           },
  //           subtitle: {
  //               color: '#2e4762'
  //           }
  //       },
  //     }
  //   },
  // },
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
                  color: '#2e4762',
                  inverseColor: '#FFDC00',
                  hoverColor: '#365473',
                  activeColor: '#263a51'
              },
              highlight: {
                  background: '#fff6c7',
                  focusBackground: '#FFDC00',
                  color: '#173350',
                  focusColor: '#173350'
              }
          },
          dark: {
              primary: {
                  color: '#FFDC00',
                  inverseColor: '#2e4762',
                  hoverColor: '#ffec54',
                  activeColor: '#ffce1e'
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

if (import.meta.env.VITE_APP_ZITADEL_ENABLED === 'true'){
  zitadelAuth.oidcAuth.useRouter(zitadelRouter)

  zitadelAuth.oidcAuth.startup().then(ok => {
    if (ok) {
          const app = createApp(App).use(createPinia())
          app.config.globalProperties.$zitadel = zitadelAuth

          app
          .use(zitadelRouter)
          .use(PrimeVue,{
              // Default theme configuration
              theme: {
                  preset: Noir,
                  options: {
                      prefix: 'p',
                      darkModeSelector: 'system',
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
    } else {
        console.error('Zitadel startup was not ok')
    }
  })
} else {
  const app = createApp(App).use(createPinia())

  app
  .use(insecureRouter)
  .use(PrimeVue,{
      // Default theme configuration
      theme: {
          preset: Noir,
          options: {
              prefix: 'p',
              darkModeSelector: 'system',
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
}
 
