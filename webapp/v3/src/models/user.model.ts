export interface UserLoginForm {
    email: string;
    password: string;
  }
  export interface UserNewPasswordForm {
    email: string;
  }
  



  export type UserRegisterForm = {
    first_name?: string
    last_name?: string
    email?: string;
    password?: string;
    terms_agree?: boolean;
    newsletter?: boolean
  } | null

  

  export interface EmailForm {
    email: string;
    reply_email: string;
  }
  
  export interface PasswordForm {
    password: string;
    reply_password: string;
    old_password: string;
  }
  
  