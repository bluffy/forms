<template>
    <div>
        <div v-if="success">
            <div class="alert alert-secondary" role="alert">
               {{ success }}
            </div>
        </div>
        <div v-else-if="formValuesNewPassword && !success">
            <Form @submit="onSubmit" :initial-values="formValuesNewPassword">
                <div class="mb-3 row">
                    <Fields :fields="fieldsNewPassword"></Fields>
                </div>
                <div class="col-12">
                    <button class="btn btn-primary" type="submit">Submit form</button>
                </div>
                <div class="col-12">
                        <router-link to="/login">login</router-link>
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
import type { PageNoContent, PageRegister } from "../models/page.model";
import Fields from "../components/Fields.vue";
import AlertDialog from "../components/AlertDialog.vue";
import ApiService from '../services/api.service'
import type { UserLoginForm, UserNewPasswordForm } from "../models/user.model";


import { FormField } from '../models/app.model';
const dialog = ref()

const success = ref()

const formValuesNewPassword = ref(null as UserNewPasswordForm);
const fieldsNewPassword = ref(null as FormField[]);


function onSubmit(values: any, actions: any) {
    return ApiService.postPage("/user/forgot_password", values).then(
        (page: PageRegister) => {
            if (page.status != 200 || !page.data.message) {
                dialog.value.alert("error on Register");
                return;
            }

            success.value = page.data.message 

       

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

    formValuesNewPassword.value = {
        email: "",
    }



    fieldsNewPassword.value = [
        {
            label: "Email",
            name: getPropertyName(formValuesNewPassword.value, a => a.email),
            type: "text",
            placeholder: "max@mustermann.de"
        }
    ]



});



</script>