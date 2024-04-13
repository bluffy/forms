<template>
    <div>
        <div v-if="page == PAGE_RECOVER">
            <div v-if="success">
                <div class="alert alert-secondary" role="alert">
                {{ success }}
                </div>
            </div>
            <div v-else-if="formRecoverValues && !success">
                <Form @submit="onSubmitRecover" :initial-values="formRecoverValues">
                    <div class="mb-3 row">
                        <Fields :fields="fieldsRecover"></Fields>
                    </div>
                    <div class="col-12">
                        <button class="btn btn-primary" type="submit">Submit form</button>
                    </div>
                    <div class="col-12">
                            <router-link to="/login">login</router-link>
                        </div>                 
                </Form>

            </div>
        </div>
        <div v-else-if="page == PAGE_PASSWORD_LINK">
            <div v-if="success">
                <div class="alert alert-secondary" role="alert">
                {{ success }}
                </div>
            </div>
            <div v-else-if="formPasswordValues && !success">
                <Form @submit="onSubmitPassword" :initial-values="formPasswordValues">
                    <div class="mb-3 row">
                        <Fields :fields="fieldsPassword"></Fields>
                    </div>
                    <div class="col-12">
                        <button class="btn btn-primary" type="submit">Submit form</button>
                    </div>
                    <div class="col-12">
                            <router-link to="/login">login</router-link>
                    </div>                 
                </Form>   
            </div>             

        </div>        
        <AlertDialog ref="dialog"></AlertDialog>
    </div>
</template>
<script lang="ts" setup>

import { Form } from 'vee-validate';
import { ref, onMounted } from 'vue'
import { genResponseError } from "../../utils/errorMessage";
import { getPropertyName } from "../../utils/helper";
import type { PageRegister } from "../../models/page.model";
import Fields from "../../components/Fields.vue";
import AlertDialog from "../../components/AlertDialog.vue";
import ApiService from '../../services/api.service'
import type {  UserNewPasswordForm } from "../../models/user.model";
import { useRoute } from 'vue-router'

import router from '../../router'

const route = useRoute();

import { FormField } from '../../models/app.model';
const dialog = ref()

const PAGE_RECOVER = 1
const PAGE_PASSWORD_LINK = 2

const success = ref()
const page = ref(0);

const formRecoverValues = ref(null as UserNewPasswordForm);
const formPasswordValues = ref();
const fieldsRecover = ref(null as FormField[]);
const fieldsPassword = ref(null as FormField[]);


function onSubmitRecover(values: any, actions: any) {
    return ApiService.post("/user/forgot_password", values).then(
        (page: PageRegister) => {
            if (page.status != 200 || !page.data.message) {
                dialog.value.alert("error on Register");
                return;
            }

            dialog.value.alert(page.data.message );
            router.replace("/login")


       

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


function onSubmitPassword(values: any, actions: any) {
    return ApiService.post("/user/forgot_password/link", {password: values.password, link: route.params.link}).then(
        (page: PageRegister) => {
            if (page.status != 200 || !page.data.message) {
                dialog.value.alert("error on change");
                return;
            }

            dialog.value.alert(page.data.message );
            router.replace("/login")


       

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
    success.value = null

    if (route.name == "forgot_password_recover_form") {

        page.value = PAGE_RECOVER

        formRecoverValues.value = {
            email: "",
        }
        fieldsRecover.value = [
            {
                label: "Email",
                name: getPropertyName(formRecoverValues.value, a => a.email),
                type: "text",
                placeholder: "max@mustermann.de"
            }
        ]
        return;
    }
    
    if (route.name == "forgot_password_link") {

        page.value = PAGE_PASSWORD_LINK

        formPasswordValues.value = {
            password: "",
        }
        fieldsPassword.value = [
            {
                label: "Password",
                name: "password",
                type: "password",
            }
        ]

        return;
    }


});



</script>