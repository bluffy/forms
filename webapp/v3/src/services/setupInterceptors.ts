import axiosInstance from "./api";
import router from "../router";
const path_page = import.meta.env.VITE_APP_API_PATH_PAGE

const setup = () => {
  axiosInstance.interceptors.response.use(
    (res: any) => {
      return res;
    },
    async (err: any) => {
    

      if (err.config.url != path_page + "/login" && err.config.url != path_page +"/register") {
        
        if (err && err.code == "ERR_NETWORK") {
          return Promise.reject(err);
        }
        if (!err || !err.response || !err.response.status) {
          return Promise.reject(err);
        }
        if (err.response.status === 400 || err.response.status === 401) {
          const path = window.location.pathname

          if (path == "/") {
            router.push({ name: 'login' });
          } else {
            router.push({ name: 'login', query: { redirect: window.location.pathname } });

          }
          return Promise.reject(err);
        } else if (err.response.status === 403) {
          router.push({ name: '/permission_denied' });
          return Promise.reject(err);
        }
      }

      return Promise.reject(err);
    }
  );
};

export default setup;
