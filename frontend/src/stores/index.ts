import { defineStore } from 'pinia'

const COLOR_MODE_KEY = 'color_mode'

export const globalStore = defineStore('global', {
  state: () => ({
    count: 0,
    orientation:'ltr',
    settings: null,
    colorMode: localStorage.getItem(COLOR_MODE_KEY) || 'light',
    shopMode: '' as string,
  }),
  getters: {
    double: state => state.count * 2,
    getSettings(state) {
      return state.settings
    },
    currentOrientation(state) {
      return state.orientation;
    },
    getColorMode(state) {
      return state.colorMode;
    },
    getShopMode(state) {
      return state.shopMode;
    }
  },
  actions: {
    increment() {
      this.count++
    },
    setOrientation(orientation:string){
        this.orientation = orientation;
    },
    setSettings(settings:any){
        this.settings = settings;
    },
    applyDarkModeClass(){
        if (this.colorMode === 'dark') {
            document.documentElement.classList.add('my-app-dark');
        } else {
            document.documentElement.classList.remove('my-app-dark');
        }
    },
    toggleDarkMode(){
        this.colorMode = this.colorMode === 'light' ? 'dark' : 'light';
        localStorage.setItem(COLOR_MODE_KEY, this.colorMode);
        this.applyDarkModeClass();
    },
    setShopMode(mode: string){
        this.shopMode = mode;
    }
  },
})