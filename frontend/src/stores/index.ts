import { defineStore } from 'pinia'

export const globalStore = defineStore('global', {
  state: () => ({ 
    count: 0,
    orientation:'ltr'
  }),
  getters: {
    double: state => state.count * 2,
    currentOrientation(state) {
      return state.orientation;
    }
  },
  actions: {
    increment() {
      this.count++
    },
    setOrientation(orientation:string){
        this.orientation = orientation;
    }
  },
})
