<template>
    <div>
        <div v-if="success">
            <div class="alert alert-secondary" role="alert">
               {{ success }}
            </div>
        </div>
        <div v-else-if="formValues && !success">
            <Form @submit="onSubmit" :initial-values="formValues">
                <div class="mb-3 row">
                    <template v-for="field in fields">
                        
                        <div :class="defCols(field)">
                            <template v-if="field.type == 'checkbox'">
                                <Field :id="field.name" :type="field.type" :value="true" :name="field.name" class="form-check-input"
                                    :placeholder="field.placeholder" />
                                <label :for="field.name"  class="form-check-label">{{ field.label }}</label>
                            </template>
                            <template v-else>
                                <label :for="field.name" class="form-label">{{ field.label }}</label>
                                <Field :id="field.name" :type="field.type" :name="field.name" class="form-control"
                                    :placeholder="field.placeholder" />
                            </template>


                            <div class="text-danger">
                                <ErrorMessage :name="field.name" />&nbsp;
                            </div>
                        </div>
                    </template>
                    <div class="col-12">
                        <button class="btn btn-primary" type="submit">Submit form</button>
                    </div>
                </div>
            </Form>
        </div>
        <AlertDialog ref="dialog"></AlertDialog>
    </div>
</template>
<script lang="ts" setup>

import { Form, Field, ErrorMessage } from 'vee-validate';
import { ref, onMounted } from 'vue'
import { genResponseError } from "../utils/errorMessage";
import type { PageRegister } from "../models/page.model";
import AlertDialog from "../components/AlertDialog.vue";
import ApiService from '../services/api.service'
import type { UserRegisterForm } from "../models/user.model";


const dialog = ref()
const success = ref()

const formValues = ref(null as UserRegisterForm);
const fields =ref();

function defCols(field: any) {
    var ret = "";
    
    if (field.col){
        ret = field.col;
    }else {
        ret = "col-12"
    }

    if (field.col_xs){
        ret =  ret + " col-xs-" + field.col_xs;
    }
    if (field.col_sm){
        ret =  ret + " col-sm-" + field.col_sm;
    }

    if (field.col_md){
        ret =  ret + " col-md-" + field.col_md;
    }
    if (field.col_lg){
        ret = ret + " col-lg-" + field.col_lg;
    }
    if (field.col_xl){
        ret =  ret + " col-xl-" + field.col_xl;
    }
    if (field.col_xxl){
        ret =  ret + " col-xxl-" + field.col_xxl;
    }



    return ret;
}
function onSubmit(values: any, actions: any) {
    return ApiService.postPage("/register", values).then(
        (page: PageRegister) => {
            if (page.status != 200 || !page.data.message) {
                dialog.value.alert("error on Register");
                return;
            }
            success.value = page.data.message
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
    fields.value = [
        {
            label: "First Name",
            name:  "first_name",
            type: "text",
            col_md: 6
        },
        {
            label: "Last Name",
            name:  "last_name",
            type: "text",
            col_md: 6
        },
        {
            label: "Email",
            name:  "email",
            type: "text",
            placeholder: "max@mustermann.de",
        },
        {
            label: "Password",
            name:  "password",
            type: "password",
        },
        {
            label: "I understand and agree to the terms of use and privacy policy",
            name:  "terms_agree",
            type: "checkbox",
        },
        {
            label: "Newsletter",
            name:  "newsletter",
            type: "checkbox",
        },
    ]
    formValues.value = {
        email: "",
        terms_agree: false,
        newsletter: false
        
    }

});



</script>