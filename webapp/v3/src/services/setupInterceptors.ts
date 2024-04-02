import axiosInstance from "./api";
import TokenService from "./token.service";
import { useAuthStore } from "../stores/auth";

const setup = () => {
  axiosInstance.interceptors.request.use(
    (config: any) => {
      const token = TokenService.getLocalAccessToken();
      if (token) {
        config.headers["Authorization"] = "BEARER " + token; // for Spring Boot back-end
        //config.headers["x-access-token"] = token; // for Node.js Express back-end
      }
      return config;
    },
    (error: any) => {
      return Promise.reject(error);
    }
  );

  axiosInstance.interceptors.response.use(
    (res: any) => {
      return res;
    },
    async (err: any) => {
      
      const originalConfig = err.config;
 
      if (originalConfig.url !== "/user" && err.response) {
        // Access Token was expired

        if (err.response.status === 401 && !originalConfig?._retry) {
          originalConfig._retry = true;

          try {

            const rs = await axiosInstance.post("/login/refresh", {
              rt: TokenService.getLocalRefreshToken(),
            });


            const { at, rt } = rs.data;
   
            const authStore = useAuthStore();

            authStore.refreshToken(at, rt);
            TokenService.updateLocalAccessToken(at, rt);


            return axiosInstance(originalConfig);
          } catch (_error: any) {
            if (_error.response.status == 403) {
              const authStore = useAuthStore();
              authStore.cleanUp()
              window.location.href = "/login";
              //return Promise.reject(_error);
              return Promise.reject();

            }
            return Promise.reject(_error);
          }
        }
      }
      
   

      return Promise.reject(err);
    }
  );
};

export default setup;
