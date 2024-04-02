import type { UserLoginForm } from "../models/user.model";
import type { Token } from "../models/token.model";
import { defineStore } from "pinia";
import AuthService from "../services/auth.service";
import TokenService from "../services/token.service";
import { useAppStore } from './app'
import tokenService from "../services/token.service";

const token = TokenService.getToken();
const initialState = token
  ? { loggedIn: true, token }
  : { loggedIn: false, token: null };

export const useAuthStore = defineStore("auth", {
  state: () => {
    return {
    token: initialState.token,
      loggedIn: initialState.loggedIn,
    };
  },
  actions: {
    login(user: UserLoginForm) {
      return AuthService.login(user).then(
        (token: Token) => {
          this.token = token;
          this.loggedIn = true;
          return Promise.resolve(token);
        },
        (error: any) => {
          this.loggedIn = false;
          this.token = null;
          return Promise.reject(error);
        }
      );
    },
    cleanUp() {
      tokenService.removeToken();
      const appStore = useAppStore()
      appStore.cleanUp();
      this.loggedIn = false;
      this.token = null;
    },
    logout() {
      return AuthService.logout().then(
        () => {
          const appStore = useAppStore()
          appStore.cleanUp();
          this.loggedIn = false;
          this.token = null;
          return Promise.resolve();
        },
        () => {
          this.loggedIn = false;
          this.token = null;
          return Promise.resolve();
        }
      );      

    },
    
    refreshToken(at: string, rt: string) {

      if (this.token == null) {
        return;
      }
      this.loggedIn = true;
      this.token = { ...this.token, at: at, rt: rt };
    },
    
  },
});
