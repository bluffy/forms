<template>
<div>

    <form @submit="onSubmit">
        <div class="mb-3 row">
            <div class="col-12" v-if="pageError">
                <div class="alert alert-primary" role="alert">
                    {{ pageError }}
                </div>
            </div>
            <div class="col-12">
                <label for="email" class="form-label">Email</label>
                <input  id="email" type="text" class="form-control" v-model="email"  v-bind="emailProps" placeholder="max@mustermann.de" />
                <div class="text-danger">
                     {{ errors.email }}&nbsp;
                </div>
            </div>

            <div class="col-12">
                <label for="password" class="form-label">Password</label>
                <input  id="password" type="password" class="form-control" v-model="password"  v-bind="passwordProps" placeholder="secret123" />
                <div  class="text-danger" >
                     {{ errors.password }}&nbsp;
                </div>
            </div>
            <div class="col-12">
                <button class="btn btn-primary" type="submit">Submit form</button>
            </div>
   
        </div>
    

  </form>
  <!--
<form ref="form" @submit.prevent="submit()">
  <div class="mb-3 row">
    <label for="staticEmail" class="col-sm-2 col-form-label">Email</label>
    <div class="col-sm-10">
      <input type="text"  v-model="loginForm.email" placeholder="max@mustermann.de" :rules="ruleEmail" class="form-control-plaintext" id="staticEmail" >
    </div>
  </div>
  <div class="mb-3 row">
    <label for="inputPassword" class="col-sm-2 col-form-label">Password</label>
    <div class="col-sm-10">
      <input type="password"  v-model="loginForm.password" placeholder="Geheim123"  :rules="rulePasswort" class="form-control" id="inputPassword">
    </div>
  </div>
  <div class="mb-3 row">
    <div class="col-sm-10">
    <button type="submit" class="btn btn-primary mb-3">Anmelden</button>
    <button type="submit" class="btn btn-primary mb-3"  to="/kennwort-vergessen">Kennwort vergessen</button>
    </div>
  </div>
  </form>
     </div>
     -->
</div>
</template>
<script lang="ts" setup>
import { useForm } from 'vee-validate';
import { useAuthStore } from "../stores/auth";
import { ref, onMounted } from 'vue'
import { genResponseError } from "../utils/errorMessage";

const pageError = ref("")
const store = useAuthStore()
const { handleSubmit, setFieldValue, defineField, errors, setErrors } = useForm()
import type { UserLoginForm } from "../models/user.model";

const [email,emailProps] = defineField('email');
const [password,passwordProps] = defineField('password');

//setFieldValue('email', 'test@test.de');
//setFieldValue('password', 'test');


const onSubmit = handleSubmit(async values => {
    console.log(values)
  // Send data to the API
  const loginForm  = values as UserLoginForm;
      console.log("errors", errors)
       setErrors({
        email: null,
        password: null,
       });
       pageError.value = "";
      store.login(loginForm).then(

        (token) => {
            //loading.value = false
            //router.push("/");
            console.log(token)
            return;
        },
        (err: any) => {
            const errors = genResponseError(err);
            if (errors.fields) {
                  setErrors(errors.fields);
            }

            if (errors.message) {
                pageError.value = errors.message
            }
            console.log(errors)
            /*
            loading.value = false
            formError.value = genResponseError(err);
            if (formError.value?.message) {
                dialog.value?.alert(formError.value?.message);
                return;
            }
            dialog.value?.alert("Das Formular ist fehlerhaft!");
            */
        }
    );
    /*
  const response = await client.post('/users/', values);
  // all good
  if (!response.errors) {
    return;
  }
  // set single field error
  if (response.errors.email) {
    setFieldError('email', response.errors.email);
  }
  // set multiple errors, assuming the keys are the names of the fields
  // and the key's value is the error message
  setErrors(response.errors);
  */
});


onMounted(() => {
    //setFieldValue('password', 'test');

    console.log(email)
    //console.log(email)

    /*
    if (store.loggedIn) {
        router.push("/")
        return;
    }   
    if (process.env.VUE_APP_OIDC  && process.env.VUE_APP_OIDC == "true"){
        oidc_enabled.value = true
    }
    if (process.env.VUE_APP_LOGIN_FORM  && process.env.VUE_APP_LOGIN_FORM == "true"){
      login_form_enabled.value = true
    }else {
      if (route.query && route.query.formular_login && route.query.formular_login && route.query.formular_login == '1') {
          login_form_enabled.value = true
      }
    }    

    checkMeta();
    */
});


/*
import { ref, onMounted } from 'vue'

import type { UserLoginForm } from "../models/user.model";
import type { AppError } from "../models/app.model";
import { genResponseError } from "../utils/errorMessage";
import { useRoute } from 'vue-router'

//import Header from '../layouts/Header.vue'
import { useAuthStore } from "../stores/auth";
import router from "../router";
import authService from '../services/auth.service';

const store = useAuthStore()

if (store.loggedIn) {
    router.push("/")
}

const route = useRoute();

const maxEmail = 100;
const maxPassword = 20;



const dialog = ref()
const loginForm = ref(<UserLoginForm>{})
const form = ref()
const loading = ref(false)
const formError = ref<AppError>(null)

//const oidc = ref(false)
const oidc_enabled = ref(false)
const login_form_enabled = ref(false)      

/*
const ruleEmail = [
    (v: string) => !!v || 'Feld muß gefüllt sein!',
    (v: string) => /^[a-z.-]+@[a-z.-]+\.[a-z]+$/i.test(v) || 'Feld muß eine E-Mail sein!',
    (v: string) => (v && v.length <= maxEmail) || 'Feld ist zu groß!',
]
const rulePasswort = [
    (v: string) => !!v || 'Feld muß gefüllt sein!',
    (v: string) => (v && v.length <= maxPassword) || 'Feld ist zu groß!',
]


 function submit() {

    if (!form.value.validate()) {
        dialog.value?.alert("Das Formular ist fehlerhaft!");
        return;
    }

    store.login(loginForm.value).then(
        () => {
            loading.value = false
            router.push("/");
            return;
        },
        (err: any) => {
            loading.value = false
            formError.value = genResponseError(err);
            if (formError.value?.message) {
                dialog.value?.alert(formError.value?.message);
                return;
            }
            dialog.value?.alert("Das Formular ist fehlerhaft!");
        }
    );
}



function checkMeta() {
    if (route.meta) {
        if (route.meta.kennwortVersenden) {

            if (!route.params.email) {
                dialog.value?.alert("Anfrage fehlgeschlagen, Parameter stimmten nicht");
                return;
            }
            if (!route.params.token) {
                dialog.value?.alert("Anfrage fehlgeschlagen, Parameter stimmten nicht");
                return;
            }

            return authService.kennwortAnfordern(<string>route.params.email, <string> route.params.token).then((resp: any)=> {
                console.log(resp)
                dialog.value?.alert("Ihnen wurde eine E-Mail mit dem neuen Kennwort versendet!");
                return 

            }).catch((err: Error) => {
                const appErr = genResponseError(err);
                dialog.value?.alert(appErr?.message);
                return;

            });
        }
        if (oidc_enabled && route.meta.oidc) {
            //  load = true
            //  this.$store.dispatch(OIDC_CALLBACK, this.$route.query).then(() => {

            //   this.$router.push({ name: "home" })

            //   }).catch(()=> {
            //     this.load = false;
            //   });

        }
    }
}
onMounted(() => {
    if (store.loggedIn) {
        router.push("/")
        return;
    }   
    if (process.env.VUE_APP_OIDC  && process.env.VUE_APP_OIDC == "true"){
        oidc_enabled.value = true
    }
    if (process.env.VUE_APP_LOGIN_FORM  && process.env.VUE_APP_LOGIN_FORM == "true"){
      login_form_enabled.value = true
    }else {
      if (route.query && route.query.formular_login && route.query.formular_login && route.query.formular_login == '1') {
          login_form_enabled.value = true
      }
    }    

    checkMeta();
});
*/


</script>
  