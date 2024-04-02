import type { Token } from "../models/token.model";

class TokenService {
  getToken(): Token | null {
    const tokenString = localStorage.getItem("t");
    if (tokenString != null) {
      return JSON.parse(tokenString);
    }
    return null;
  }

  
  getLocalRefreshToken() {
    const token = this.getToken();
    return token?.rt;
  }
  
  getLocalAccessToken() {
    const token = this.getToken();
    return token?.at;
  }

  updateLocalAccessToken(at: string, rt: string) {
    const token = this.getToken();
    if (token == null) {
      return;
    }
    token.at = at
    token.rt = rt

    localStorage.setItem("t", JSON.stringify(token));
  }

  setToken(token: Token) {
    localStorage.setItem("t", JSON.stringify(token));
  }

  removeToken() {
    localStorage.removeItem("t");
  }
}

export default new TokenService();
