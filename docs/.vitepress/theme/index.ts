// https://vitepress.dev/guide/custom-theme
import { h } from 'vue'
import type { Theme } from 'vitepress'
import DefaultTheme from 'vitepress/theme'
import './style.css'


import PrimeVue from 'primevue/config';
import 'primeicons/primeicons.css'
import Aura from '@primevue/themes/aura';
import 'primeflex/primeflex.css'
import { Button, Image, Card,Badge } from 'primevue';
import Layout from './Layout.vue'


export default {
  extends: DefaultTheme,
  Layout,
  enhanceApp({ app, router, siteData }) {
    app.use(PrimeVue, {
        theme: {
            preset: Aura
        }
    })
    app.component('Image',Image)
    app.component('Button',Button)
    app.component('Card',Card)
    app.component('Badge',Badge)
  }
} satisfies Theme
