import { defineStore } from 'pinia'

export const globalStore = defineStore('global', {
  state: () => ({ 
    count: 0,
    orientation:'ltr',
    settings: null,
    colorMode: 'light',
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
    toggleDarkMode(){
        this.colorMode = this.colorMode === 'light' ? 'dark' : 'light';
    }
  },
})
