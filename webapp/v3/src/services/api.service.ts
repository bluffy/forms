import api from "./api";
import { useAppStore } from "../stores/app";
//const path_api = import.meta.env.VITE_APP_API_PATH
const path_page = import.meta.env.VITE_APP_API_PATH_PAGE


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

    getPage(slug: string, params?: any) {
        return get(path_page, slug, params)
    }

    postPage(slug: string, values?:any) {
      return post(path_page, slug, values)
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

export default new ApiService();
