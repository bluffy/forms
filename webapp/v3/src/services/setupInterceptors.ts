import axiosInstance from "./api";
import router from "../router";
import { useRoute } from 'vue-router'

const setup = () => {
  axiosInstance.interceptors.response.use(
    (res: any) => {

      return res;
    },
    async (err: any) => {
      
     

    

      if (!(err.config.url == "/login" || err.config.url == "/signup")) {

        console.log("FEHLER")

        if (err.response.status === 400 || err.response.status === 401) {

         const path = window.location.pathname

         if (path == "/") {
          router.push( { name: 'login'});
         }else {
          router.push( { name: 'login', query: { redirect: window.location.pathname }});

         }

        
          return Promise.reject(err);
        }else if  (err.response.status === 403){
          window.location.href = "/permission_denied";
          return Promise.reject(err);
        }
      }
      return Promise.reject(err);


    }
  );
};

export default setup;
