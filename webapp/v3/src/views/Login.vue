<template>
    <div>
        <div v-if="formValuesLogin">
            <Form @submit="onSubmitLogin" :initial-values="formValuesLogin">
                <div class="mb-3 row">
                        <Fields :fields="fieldsLogin"></Fields>
                    <div class="col-12">
                        <router-link to="/forgot_password">Recover Password</router-link>
                    </div>
                    <div class="col-12 mt-3">
                        <button class="btn btn-primary" type="submit">Submit form</button>
                    </div>
                </div>
            </Form>
        </div>
  
        <AlertDialog ref="dialog"></AlertDialog>
    </div>
</template>
<script lang="ts" setup>

import { Form } from 'vee-validate';
import { ref, onMounted } from 'vue'
import { genResponseError } from "../utils/errorMessage";
import { getPropertyName } from "../utils/helper";
import type { PageNoContent } from "../models/page.model";
import router from "../router";
import Fields from "../components/Fields.vue";
import AlertDialog from "../components/AlertDialog.vue";
import ApiService from '../services/api.service'
import type { UserLoginForm } from "../models/user.model";


import { useRoute } from 'vue-router'
import { FormField } from '../models/app.model';
const dialog = ref()



const route = useRoute();
const formValuesLogin = ref(null as UserLoginForm);
const fieldsLogin = ref(null as FormField[]);




function onSubmitLogin(values: any, actions: any) {
    return ApiService.postPage("/user/login", values).then(
        (page: PageNoContent) => {
            if (page.status != 204) {
                dialog.value.alert("error on login");
                return;

            }
            if (route.query.redirect) {
                router.replace(route.query.redirect.toString());
                return
            }
            router.replace("/");
            return;

        },
        (err: any) => {
            const errors = genResponseError(err);
            if (errors?.fields) {
                actions.setErrors(errors.fields);
            }

            if (errors?.message) {
                dialog.value.alert(errors.message);
                return;
            }
        }
    );
};

onMounted(() => {
    formValuesLogin.value = {
        email: "",
        password: ""
    }


    fieldsLogin.value = [
        {
            label: "Email",
            name: getPropertyName(formValuesLogin.value, a => a.email),
            type: "text",
            placeholder: "max@mustermann.de"
        },
        {
            label: "Password",
            name: getPropertyName(formValuesLogin.value, a => a.password),
            type: "password"
        }
    ]





});



</script>