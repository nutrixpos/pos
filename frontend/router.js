import { createMemoryHistory, createRouter } from 'vue-router'

import HomeView from './HomeView.vue'
import AboutView from './AboutView.vue'

const routes = [
  { path: '/', component: HomeView },
  { path: '/queue', component: AboutView },
]

const router = createRouter({
  history: createMemoryHistory(),
  routes,
})