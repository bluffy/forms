import api from "./api";
import { useAppStore } from "../stores/app";
//const path_api = import.meta.env.VITE_APP_API_PATH
const api_path = import.meta.env.VITE_APP_API_PATH


function get(path: string, slug: string, params?: any ) {
  const appStore = useAppStore()
  appStore.startLoading()        
  return api.get(path + slug, {
    params: params
  }).then((response) => {
    return {status: response.status, data: response.data};
  }).finally(() => {
    appStore.endLoad()
  });
}

function post(path: string, slug: string, values?: any ) {
  const appStore = useAppStore()
  appStore.startLoading()        
  return api.post(path + slug, values).then((response) => {
    return {status: response.status, data: response.data};
  }).finally(() => {
    appStore.endLoad()
  });
}

class ApiService {

    get(slug: string, params?: any) {
        return get(api_path, slug, params)
    }

    post(slug: string, values?:any) {
      return post(api_path, slug, values)
    }


    /*
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
    */


}

export default new ApiService();
