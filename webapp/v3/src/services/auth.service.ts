import api from "./api";
import TokenService from "./token.service";
import type { EmailForm, PasswordForm, UserLoginForm } from "../models/user.model";
import type { Token } from "../models/token.model";
//const path_api = import.meta.env.VITE_APP_API_PATH
const path_page = import.meta.env.VITE_APP_API_PATH_PAGE

class AuthService {
  login(user: UserLoginForm) {
    return api.post(path_page +"/user/login", user).then((response) => {
      

      return response.data;
    });
  }

  logout() {
    return api.delete("user").then(() => {
      TokenService.removeToken();
    }).catch(() => {
      TokenService.removeToken();
    });
  }
  kewnnwortVergessen(form: EmailForm) {
    return api.get("forgot_password",{
      params: form
    })
  }

  kennwort(form: PasswordForm) {

    const param = {
      pch1: form.old_password,
      pch2: form.password,
      pch3: form.reply_password,
    }

    return api.get("query",{
      params: param
    })
  }

  kennwortAnfordern(email: string, token: string) {

    const params = {
      email: email,
      token: token
    }

    return api.post("new_password",params)
  }
  /*
  register(user: UserLoginForm) {
    return api.post("/auth/register", user);
  }
  */
}

export default new AuthService();
