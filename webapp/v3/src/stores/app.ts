// Utilities
import { defineStore } from 'pinia'


/*
function isObjEmpty (obj: any) {
  return Object.keys(obj).length === 0;
}
*/

type Timer = ReturnType<typeof setTimeout>

export const useAppStore = defineStore('app', {
  state: () => ({
    isLoading: false,
    isLoadingLong: false,
    isRequesting: false,
    loaderVerzoegert: {} as Timer,
    loaderVerzoegertLong: {} as Timer,  
  }),
  getters: {

  }, 
  actions: {
    cancleLoadingLong(){
      
    },
    
    cleanUp() {

    },
    startLoading() {
      this.isRequesting = true

      if (this.loaderVerzoegert) {
        clearTimeout(this.loaderVerzoegert);
      }
      if (this.loaderVerzoegertLong) {
        clearTimeout(this.loaderVerzoegertLong);
      }      
    
      this.loaderVerzoegert = setTimeout(() => {
        if (this.isRequesting) {
          this.isLoading = true;
        }
      }, 200);
    
      this.loaderVerzoegertLong = setTimeout(() => {
        if (this.isRequesting) {
          this.isLoadingLong = true;
        }
      }, 2000);
    
    },
    endLoad() {
      clearTimeout(this.loaderVerzoegert);
      clearTimeout(this.loaderVerzoegertLong);  
      this.isRequesting = false;
      this.isLoading = false;
      this.isLoadingLong = false;
    },




  } 
})
