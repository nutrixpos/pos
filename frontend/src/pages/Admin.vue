<template>
    <div v-if="!loading">
        <div class="grid p-0 m-0">
            <div class="col-12 p-0">
                <Toolbar style="border-radius: 0px;" class="py-1 lg:py-2">
                    <template #start>
                        <div @click="version_dialog_visible=true" style="text-decoration: none;color:gray">
                            <img src="@/assets/logo.png" alt="logo" style="height:25px" v-if="store.getColorMode == 'light'">
                            <div v-else class="flex justify-content-center align-items-center" style="font-size:1rem;color:white;font-family:'FontAwesome'">
                                nutrix
                            </div>
                        </div>
                        <router-link v-for="(item,index) in items" :key="index" :to="item.link">
                            <Button :icon="item.icon" :label="$t(`${item.label.title}`,item.label.plural ? 3 : 1)"  :text="!item.focused" severity="secondary" />
                        </router-link>
                    </template>

                    <template #end>
                        <Button outlined :icon="`pi pi-${store.getColorMode == 'light' ? 'sun' : 'moon'}`" @click="toggleDarkMode()" />
                        <Button  severity="secondary" size="large"  text rounded aria-label="Profile" label="Profile" @click.stop="user_profile_toggle">
                            <span style="font-size:0.9rem;" class="mr-2">{{ user?.name }}</span>
                            <span class="p-button-icon pi pi-user"></span>
                        </Button>
                        <OverlayPanel ref="user_profile_op" class="lg:w-2 md:w-3">
                            <div class="flex flex-column">
                                <span>Welcome <strong>{{ user?.name }}</strong></span>
                                <div class="mt-2">
                                    <Chip v-for="(role,index) in roles" :key="index" :label="role" style="height: 1.5rem;" class="m-1" />
                                </div>
                                <Button class="mt-5" icon="pi pi-sign-out" severity="secondary" text aria-label="Signout" :label=" $t('signout')" @click="proxy.$zitadel?.oidcAuth.signOut()" />
                            </div>
                        </OverlayPanel>
                    </template>
                </Toolbar>
            </div>
            <div class="col-12">
                <div class="grid">
                    <div class="col-3 xl:col-2">
                        <Tree v-model:expandedKeys="expandedKeys" :value="menu_tree" selectionMode="single" class="w-full" @node-select="(node) => sidemenuNodeSelect(node)">
                            <template #default="slotProps">
                                <div style="text-decoration: none;color: inherit;" class="flex align-items-center w-full">
                                    <div>{{ $t(`${slotProps.node.label.title}`,slotProps.node.label.plural ? 3 : 1) }}</div>
                                </div>
                            </template>
                        </Tree>
                    </div>
                    <div class="col-9 xl:col-10 flex p-0 pt-3 mt-2">
                        <RouterView />
                    </div>
                </div>
            </div>
        </div>
    </div>
    <div style="width:100vw;height:100vh;display:flex;justify-content:center;align-items:center" v-if="loading">
      <ProgressSpinner style="width: 50px; height: 50px" strokeWidth="8" fill="transparent"
      animationDuration=".5s" aria-label="Custom ProgressSpinner" />
    </div>
    <Dialog v-model:visible="version_dialog_visible" header="Nutrix" :style="{ width: '45rem' }">
        <p class="text-justify">
            Nutrix is an open-source restaurant management system
            designed to make managing your restaurant easy and efficient.
            It's built with modern web technologies and provides a simple
            and intuitive interface to manage your menu, orders, customers,
            and more. Nutrix is completely free and open source under the GPL-2 license, meaning
            you have complete control over the system and can modify it
            to suit your needs. With Nutrix, you can focus on what matters
            most - providing great food and service to your customers.
        </p>
        <p>
            For more support & collaboration visit &nbsp;<a style="font-size:large;" href="https://nutrixpos.com" target="_blank"><i class="pi pi-external-link mr-2"></i>https://nutrixpos.com </a>
        </p>
        <p>
            version / commit hash : {{ app_version }}
        </p>
    </Dialog>
</template>

<script setup lang="ts">
import {ref,getCurrentInstance,computed} from "vue";
import { Toolbar,Dialog } from "primevue";
import Tree from "primevue/tree";
import Button from "primevue/button";
import { useI18n } from 'vue-i18n'
import { globalStore } from '@/stores';
import axios from "axios";
import OverlayPanel from "primevue/overlaypanel";
import ProgressSpinner from "primevue/progressspinner";


const { proxy } = getCurrentInstance();
const store = globalStore()
const user_profile_op = ref();

const app_version = ref("")
const version_dialog_visible = ref(false)

app_version.value = import.meta.env.VITE_APP_APP_VERSION || ""

const user : any = computed(() => {

    return proxy.$zitadel?.oidcAuth.userProfile

})

const toggleDarkMode = () => {
    store.toggleDarkMode()
    document.documentElement.classList.toggle('my-app-dark');
}

const sidemenuNodeSelect = (node) => {
    if (node.link) {
        proxy.$router.push(node.link);
    }
}


// const selected_list_item = ref ({ name: 'Inventory', icon:'inbox', link:'inventory' })

const user_profile_toggle = (event: any) => {
    user_profile_op.value.toggle(event);
}

const expandAll = () => {
    for (let node of menu_tree.value) {
        expandNode(node);
    }

    expandedKeys.value = { ...expandedKeys.value };
};

const expandNode = (node) => {
    if (node.children && node.children.length) {
        expandedKeys.value[node.key] = true;

        for (let child of node.children) {
            expandNode(child);
        }
    }
};

const expandedKeys = ref({});

const menu_tree =ref([
    {
        key: '0',
        label: {
            title:'material',
            plural: true,
        },
        data: 'Materials',
        icon: 'fa fa-carrot',
        link:'/admin/inventory',
    },
    {
        key: '1',
        label: {
            title: 'product',
            plural: true
        },
        data: 'Products',
        icon: 'fa fa-barcode',
        link:'/admin/products',
    },
    {
        key: '2',
        label: {
            title:'category',
            plural: true
        },
        data: 'Categories',
        icon: 'fa fa-sitemap',
        link:'/admin/categories',
    },
    {
        key: '3',
        label: {
            title :'order',
            plural:true
        },
        data: 'Orders',
        icon: 'pi pi-fw pi-box',
        link:'/admin/orders',
        children: [
            {
                key: '3-0',
                label: {
                    title:'list',
                    plural:false
                },
                data: 'List oders',
                icon: 'pi pi-fw pi-list',
                link:'/admin/orders',
            },
        ]
    },
    {
        key: '4',
        label: {
            title:'report',
            plural:true
        },
        data: 'Reports',
        icon: 'pi pi-fw pi-chart-line',
        link:'/admin/sales',
        children: [
            {
                key: '4-0',
                label: {
                    title:'sales',
                    plural:false
                },
                data: 'Sales',
                icon: 'pi pi-fw pi-percentage',
                link:'/admin/sales',
            },
        ]
    },
    {
        key: '5',
        label: {
            title:'settings',
            plural:false
        },
        data: 'Settings',
        icon: 'pi pi-fw pi-cog',
        link:'/admin/settings',
    },
    {
        key: '6',
        label: {
            title:'customer',
            plural:true
        },
        data: 'Customers',
        icon: 'pi pi-users',
        link:'/admin/customers',
    },
    {
        key: '7',
        label: {
            title:'Hubsync',
            plural:false
        },
        data: 'Hubsync',
        icon: 'pi pi-sync',
        link:'/admin/hubsync',
    },
])

const items = ref([
      {
          label: {
              title:'cashier',
              plural:false
          },
          focused: false,
          icon: 'pi pi-desktop',
          link: '/home',
      },
      {
          label: {
              title:'kitchen',
              plural:false
          },
          focused: false,
          icon: 'fa fa-kitchen-set',
          link:'/kitchen'
      },
      {
          label: {
              title:'admin',
              plural:false
          },
          focused: true,
          icon: 'pi pi-cog',
          link: '/admin',
      }
  ]);

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


loadLanguage()
expandAll()
</script>

<style>
html,
body {
height: 100%;
margin:0;
}
</style>