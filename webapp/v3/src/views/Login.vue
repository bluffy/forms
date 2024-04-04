<template>
<div>


    

       
  


    <div class="mb-3 row">
        <div class="col-12" v-if="pageError">
            <div class="alert alert-danger" role="alert">
                {{ pageError }}
            </div>
        </div>
    </div>

    <div  v-if="formValues">
    <Form  @submit="onSubmit"  :initial-values="formValues">
        <div class="mb-3 row">
            <div class="col-12" v-if="pageError">
                <div class="alert alert-danger" role="alert">
                    {{ pageError }}
                </div>
            </div>
            <div v-for="field in fields">
                <div class="col-12">
                    <label :for="field.name" class="form-label">Email</label>
                    <Field  :id="field.name" :type="field.type" :name="field.name"  class="form-control" placeholder="max@mustermann.de"   />
                    <div  class="text-danger" >
                        <ErrorMessage :name="field.name" />&nbsp;
                    </div>
                </div>                
            </div>
            <div class="col-12">
                <button class="btn btn-primary" type="submit">Submit form</button>
            </div>
        </div>
    </Form>

</div>
<AlertDialog ref="dialog1"></AlertDialog>

</div>
</template>
<script lang="ts" setup>

import { Form, Field, ErrorMessage} from 'vee-validate';
import { useAuthStore } from "../stores/auth";
import { ref, onMounted } from 'vue'
import { genResponseError } from "../utils/errorMessage";
import type { UserLoginForm } from "../models/user.model";
import router from "../router";
import AlertDialog from "../components/AlertDialog.vue";

import { useRoute } from 'vue-router'
const dialog1 = ref()

const route = useRoute();

const pageError = ref()
const store = useAuthStore()


const formValues = ref()

const fields = [
    {
        label: "Email",
        name: "email",
        type: "text"
    },
    {
        label: "Password",
        name: "password",
        type: "password"
    }
]



function onSubmit(values: any, actions: any) {
    const loginForm  = values as UserLoginForm;
    store.login(loginForm).then(
    (token) => {
        //loading.value = false
        if (route.query.redirect) {
            router.push(route.query.redirect.toString());
            return            
        }
        router.push("/");
        console.log(token)
        return;
    },
    (err: any) => {



        const errors = genResponseError(err);
        if (errors?.fields) {
            actions.setErrors(errors.fields);
        }

        if (errors?.message) {
                dialog1.value.alert(errors.message);
                return;
            }

            /*
        if (errors?.message) {
            pageError.value = errors.message
        }
        */

    }
    );
};
  // Submit the values...
  // set single field value
  //actions.setFieldValue('email', 'ummm@example.com');
  // set multiple values
  //actions.setValues({
  //  email: 'ummm@example.com',
  //  password: 'P@$$w0Rd',
  //});

/*
import { useForm } from 'vee-validate';

const { handleSubmit, defineField, errors, setErrors } = useForm()
const [email,emailProps] = defineField('email');
const [password,passwordProps] = defineField('password');

//setFieldValue('email', 'test@test.de');
//setFieldValue('password', 'test');


*/

/*a
ubmit = handleSubmit(async values => {
    console.log(values)
  // Send data to the API
  const loginForm  = values as UserLoginForm;
      console.log("errors", errors)
       setErrors({
        email: "",
        password: "",
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
            if (errors?.fields) {
                  setErrors(errors.fields);
            }

            if (errors?.message) {
                pageError.value = errors.message
            }
            console.log(errors)

        }
    );

});
*/

onMounted(() => {
    console.log(dialog1.value)
    dialog1.value?.alert("test");
    console.log(route.query)

    //setFieldValue('password', 'test');
    formValues.value = {
        "email": "dev@bluffy.de",
        "password": "mgr"
    }

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
  