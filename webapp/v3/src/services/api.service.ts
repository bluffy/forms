import api from "./api";
import { useAppStore } from "../stores/app";
//const path_api = import.meta.env.VITE_APP_API_PATH
const path_page = import.meta.env.VITE_APP_API_PATH_PAGE

class PageService {

    get(slug: string, params?: any) {

        const appStore = useAppStore()
        appStore.startLoading()
        
        return api.get(path_page + slug, {
          params: params
        }).then((response) => {
          return response.data;
        }).finally(() => {
          appStore.endLoad()
        });


    }

    query(params?: any) {

        const appStore = useAppStore()
        appStore.startLoading()
        if (!params) {
          return api.get("/query").finally(() => {
            appStore.endLoad()
          });
        }
        return api.get("/query/", {
          params: params
        }).finally(() => {
          appStore.endLoad()
        });
    }
    
    post(params: any){
      const appStore = useAppStore()
      appStore.startLoading()  
      return api.post("/query/", params).finally(() => {
        appStore.endLoad()
      });
    }

    rotate(params: any){
      const appStore = useAppStore()
      appStore.startLoading() 
      return api.get("/rotate", {
        params: params
      }).finally(() => {
        appStore.endLoad()
      });
    }

    file(params:any ) {   
      return api.get("/query", {
        params: params,
        responseType: 'blob'        
      });      
    }

    upload(params: any, formular: any) {
      return api.post("/upload", formular, {
        params: params,
        headers: {
          'Content-Type': 'multipart/form-data'
      }
      });
    }


}

export default new PageService();
