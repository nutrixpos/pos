<template>
  <main :style="`height:100%;direction:${orientation};`">
      <RouterView v-if="!setting_up && !show_mode_modal" />
      <!-- First-run shop mode selection modal -->
      <div v-if="show_mode_modal" style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; background: linear-gradient(135deg, #14977B, #9CBF34);"></div>
      <Dialog
        v-model:visible="show_mode_modal"
        :closable="false"
        modal
        header=""
        :style="{ width: '56rem', borderRadius: '1.5rem', overflow: 'hidden',border:'0px' }"
        :pt="{ header: { style: 'display:none' }, content: { style: 'padding:0' } }"
      >
        <div class="mode-modal-wrapper">
          <div class="mode-modal-header">
            <div class="mode-modal-logo">
                <img src="@/assets/logo.png" alt="logo" style="height:25px" v-if="store.getColorMode == 'light'">
            </div>  
            <h2 class="mode-modal-title">{{$t('welcome')}}! {{$t('choose_your_shop_mode')}}</h2>
            <p class="mode-modal-subtitle">{{$t('this_setting_controls_which_features_are_available')}}, {{$t('you_can_change_it_later_in_settings')}}</p>
          </div>

          <div class="mode-modal-cards">
            <!-- Kitchen mode card -->
            <div
              class="mode-card"
              :class="{ 'mode-card--selected': selected_mode === 'kitchen' }"
              @click="selected_mode = 'kitchen'"
            >
              <div class="mode-card-icon mode-card-icon--kitchen">
                <i class="fa fa-kitchen-set"></i>
              </div>
              <h3 class="mode-card-title">{{$t('kitchen')}}</h3>
              <ul class="mode-card-features">
                <li><i class="pi pi-check"></i> {{$t('kitchen_display_screen')}}</li>
                <li><i class="pi pi-check"></i> {{$t('product_ready_to_service_tracking')}}</li>
                <li><i class="pi pi-check"></i> {{$t('order_workflow_management')}}</li>
                <li><i class="pi pi-check"></i> {{$t('inventory_tracking_out_of_the_box')}}</li>
              </ul>
            </div>

            <!-- Retail mode card -->
            <div
              class="mode-card"
              :class="{ 'mode-card--selected': selected_mode === 'retail' }"
              @click="selected_mode = 'retail'"
            >
              <div class="mode-card-icon mode-card-icon--retail">
                <i class="pi pi-shopping-cart"></i>
              </div>
              <h3 class="mode-card-title">{{$t('retail')}}</h3>
              <ul class="mode-card-features">
                <li><i class="pi pi-check"></i> {{$t('streamlined_checkout')}}</li>
                <li><i class="pi pi-check"></i> {{$t('no_kitchen_specific_ui')}}</li>
                <li><i class="pi pi-check"></i> {{$t('focused_pos_experience')}}</li>
                <li><i class="pi pi-check"></i> {{$t('inventory_tracking_out_of_the_box')}}</li>
              </ul>
            </div>
          </div>

          <div class="mode-modal-footer">
            <Button
              :label="selected_mode ? `Continue with ${selected_mode === 'kitchen' ? 'Kitchen' : 'Retail'} mode` : 'Select a mode to continue'"
              :disabled="!selected_mode || saving_mode"
              :loading="saving_mode"
              size="large"
              @click="saveShopMode"
              style="min-width: 20rem; border-radius: 2rem;"
            />
          </div>
        </div>
      </Dialog>

      <Toast  group="br" position="top-center">
        <template #content="{ message }">
            <section class="flex p-3 gap-3 w-full bg-black-alpha-90">
                <div style="border-radius:100%;background-color:white;height:1.9rem;width:2rem;" class="p-1 flex justify-content-center align-items-center opacity-95">
                  <i v-if="message.severity == `success`"  :class="`pi pi-check text-${message.severity}-500`"></i>
                  <i v-if="message.severity == `info`"  :class="`pi pi-info text-${message.severity}-500`"></i>
                  <i v-if="message.severity == `warning`"  :class="`pi pi-exclamation-triangle text-${message.severity}-500`"></i>
                  <i v-if="message.severity == `error`"  :class="`pi pi-times text-${message.severity}-500`"></i>
                </div>
                <div class="flex flex-column gap-3 w-full">
                    <p class="m-0 font-semibold text-base text-white opacity-80">{{ message.summary }}</p>
                    <p class="m-0 text-white text-base text-700">{{ message.detail }}</p>
                </div>
            </section>
        </template>
      </Toast>
  </main>
</template>

<script setup>
import { computed, getCurrentInstance, ref, onMounted } from 'vue';
import Toast from 'primevue/toast';
import Dialog from 'primevue/dialog';
import Button from 'primevue/button';
import { globalStore } from '@/stores';
import axios from 'axios';
import { useRouter } from 'vue-router'
import { useToast } from "primevue/usetoast";
import { useI18n } from 'vue-i18n'


onMounted(() => {
    store.applyDarkModeClass();
})


const { t } = useI18n() 
const toast = useToast();
const router = useRouter()
const setting_up = ref(true)
const show_mode_modal = ref(false)
const selected_mode = ref('')
const saving_mode = ref(false)

const { proxy } = getCurrentInstance();
const store = globalStore()
const orientation = computed(() => store.currentOrientation)
const { locale,setLocaleMessage } = useI18n({ useScope: 'global' })

const getSettings = () => {
    return axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/settings`, {
        headers: {
            Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        },
    })
    .then((response) => {
        store.setSettings(response.data.data)
        const mode = response.data.data.shop_mode || ''
        store.setShopMode(mode)
        if (!mode) {
            show_mode_modal.value = true
        }

        axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/languages/${response.data.data.language.code}`, {
            headers: {
                Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
            }
        })
        .then(response2 => {
            setLocaleMessage(response2.data.data.code, response2.data.data.pack);
            locale.value = response2.data.data.code
            store.setOrientation(response2.data.data.orientation)
        })
        .catch((err) => {
            console.log(err)
        });


    })
    .catch((err) => {
        console.log(err)
    });
}

const saveShopMode = () => {
    if (!selected_mode.value) return
    saving_mode.value = true

    // Merge the new shop_mode into the existing settings payload
    const currentSettings = store.getSettings || {}

    axios.patch(`http://${import.meta.env.VITE_APP_BACKEND_HOST}${import.meta.env.VITE_APP_MODULE_CORE_API_PREFIX}/api/settings`,
        {
            data: {
                ...currentSettings,
                shop_mode: selected_mode.value
            }
        },
        {
            headers: {
                Authorization: `Bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
            }
        }
    )
    .then(() => {
        store.setShopMode(selected_mode.value)
        show_mode_modal.value = false
        saving_mode.value = false
    })
    .catch((err) => {
        console.log(err)
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to save shop mode', life: 3000, group: 'br' })
        saving_mode.value = false
    })
}

const getConfigStatus = () => {
    return axios.get(`http://${import.meta.env.VITE_APP_BACKEND_HOST}/api/setup/status`, {
        headers: {
            Authorization: `bearer ${proxy.$zitadel?.oidcAuth.accessToken}`
        },
    })
}

const init = () => {
    getConfigStatus()
    .then((response) => {
        if(response.data.setup) {
            getSettings().finally(() => {
                setting_up.value = false
            })
        }else {
            router.push({ path: '/setup' }).finally(() => {
                setting_up.value = false
            })
        }
    })
    .catch((err) => {
        toast.add({severity:'error', summary: 'Error', detail: err, life: 3000, group: 'br'});
    });
}

init()

</script>

<style lang="scss">
@use '@/assets/styles.scss';

body {
    font-family: sans-serif;
    height: 100vh;
    margin:0px;
}
.my-app-dark {
    body {
        font-family: sans-serif;
        height: 100vh;
        background-color: #232327;
    }
}

/* ── First-run shop mode modal ───────────────────────────────── */
.mode-modal-wrapper {
    display: flex;
    flex-direction: column;
    background: linear-gradient(145deg, #222222 0%, #181f19 60%, #1a1a2e 100%);
    color: #fff;
    padding: 2.5rem 2.5rem 2rem;
    min-height: 500px;
}

.mode-modal-header {
    text-align: center;
    margin-bottom: 2rem;
}

.mode-modal-logo {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #1A9878, #E3D40F);
    border-radius: 1rem;
    padding: 0.4rem 1.1rem;
    margin-bottom: 1.2rem;
}

.mode-modal-logo-text {
    font-size: 1.2rem;
    font-weight: 700;
    letter-spacing: 0.08em;
    color: #fff;
}

.mode-modal-title {
    font-size: 1.6rem;
    font-weight: 700;
    margin: 0 0 0.5rem;
    background: linear-gradient(90deg, #14977B,#2E9D6D);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
}

.mode-modal-subtitle {
    font-size: 0.95rem;
    color: #dce4dc;
    margin: 0;
}

.mode-modal-cards {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.25rem;
    margin-bottom: 2rem;
}

.mode-card {
    border: 2px solid rgba(255, 255, 255, 0.08);
    border-radius: 1.2rem;
    padding: 1.75rem 1.5rem;
    cursor: pointer;
    transition: all 0.25s ease;
    background: rgba(255, 255, 255, 0.03);
    position: relative;
    overflow: hidden;

    &::before {
        content: '';
        position: absolute;
        inset: 0;
        opacity: 0;
        transition: opacity 0.25s ease;
        background: linear-gradient(135deg, rgba(99, 241, 130, 0.07), rgba(92, 246, 131, 0.07));
    }

    &:hover {
        border-color: rgba(99, 241, 163, 0.4);
        transform: translateY(-3px);
        &::before { opacity: 1; }
    }

    &--selected {
        border-color: #14977b !important;
        box-shadow: 0 0 0 3px rgba(99, 241, 130, 0.2), 0 8px 32px rgba(99, 241, 194, 0.15);
        &::before { opacity: 1; }

        .mode-card-icon {
            box-shadow: 0 0 24px rgba(99, 241, 210, 0.5);
        }
    }
}

.mode-card-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 3.5rem;
    height: 3.5rem;
    border-radius: 1rem;
    margin-bottom: 1rem;
    font-size: 1.5rem;
    transition: box-shadow 0.25s ease;

    &--kitchen {
        background: linear-gradient(135deg, #f97316, #ef4444);
        color: #fff;
    }

    &--retail {
        background: linear-gradient(135deg, #06b6d4, #3b82f6);
        color: #fff;
    }
}

.mode-card-title {
    font-size: 1.2rem;
    font-weight: 700;
    margin: 0 0 1rem;
    color: #e2e8f0;
}

.mode-card-features {
    list-style: none;
    margin: 0;
    padding: 0;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;

    li {
        display: flex;
        align-items: center;
        gap: 0.6rem;
        font-size: 0.875rem;
        color: #ffffff;

        .pi-check {
            color: #ffffff;
            font-size: 0.75rem;
        }
    }
}

.mode-modal-footer {
    display: flex;
    justify-content: center;
}
</style>