import { createRouter, createWebHistory } from 'vue-router'
import auth from "@/services/auth";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
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
          {path: 'users', component: () => { 
            if (auth.hasRole("superuser")) {
              return import('@/pages/Users.vue')
            }
            return import('@/pages/NoAccessView.vue')
          }},
          {path: 'hubsync', component: () => {
            if (auth.hasRole("superuser")) {
              return import('@/pages/Hubsync.vue')
            }
            return import('@/pages/NoAccessView.vue')
          }},
        ],
      },
  ],
})

export default router
