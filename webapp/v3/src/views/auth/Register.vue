<template>
    <div class="container">
        <template v-if="success">
            <div class="alert alert-secondary" role="alert">
                <MarkdownRenderer :source="success" ></MarkdownRenderer>
            </div>
        </template>
        <template v-else-if="formValues && !success">
            <Form @submit="onSubmit" :initial-values="formValues">
                <div class="mb-3 row">
                    <Fields :fields="fields"></Fields>
                </div>
                <div class="col-12">
                    <button class="btn btn-primary" type="submit">Submit form</button>
                </div>                
            </Form>
        </template>
        <AlertDialog ref="dialog"></AlertDialog>
    </div>
</template>
<script lang="ts" setup>

import { Form } from 'vee-validate';
import { ref, onMounted } from 'vue'
import { genResponseError } from "../../utils/errorMessage";
import type { PageMessage } from "../../models/page.model";
import MarkdownRenderer from "../../components/MarkdownRenderer.vue";
import AlertDialog from "../../components/AlertDialog.vue";
import Fields from "../../components/Fields.vue";
import ApiService from '../../services/api.service'
import type { UserRegisterForm } from "../../models/user.model";
import type { FormField } from "../../models/app.model";

const dialog = ref()
const success = ref()

const formValues = ref(null as UserRegisterForm);
const fields =ref(null as FormField[]);

function onSubmit(values: any, actions: any) {
    return ApiService.postPage("/user/register", values).then(
        (page: PageMessage) => {

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