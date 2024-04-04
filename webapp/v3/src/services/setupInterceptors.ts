import axiosInstance from "./api";
import router from "../router";


const setup = () => {
  axiosInstance.interceptors.response.use(
    (res: any) => {

      return res;
    },
    async (err: any) => {


      if (!(err.config.url == "/login" || err.config.url == "/signup")) {
        console.log("FEHLER")

        if (err.response.status === 400 || err.response.status === 401) {
         // window.location.href = "/login?redir=" + err.config.url;
        
//         router.push( { name: 'login', params: {redir: err.config.url, error: err }});
         router.push( { name: 'login', query: { redirect: err.config.url }});
        
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
