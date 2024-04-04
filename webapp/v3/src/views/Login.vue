<template>
<div>

  <div  v-if="formValues">
    <Form  @submit="onSubmit"  :initial-values="formValues">
        <div class="mb-3 row">
            <div v-for="field in fields">
                <div class="col-12">
                    <label :for="field.name" class="form-label">{{ field.label }}</label>
                    <Field  :id="field.name" :type="field.type" :name="field.name"  class="form-control" :placeholder="field.placeholder"   />
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
<AlertDialog ref="dialog"></AlertDialog>

</div>
</template>
<script lang="ts" setup>

import { Form, Field, ErrorMessage} from 'vee-validate';
import { ref, onMounted } from 'vue'
import { genResponseError } from "../utils/errorMessage";
import type { PageNoContent } from "../models/page.model";
import router from "../router";
import AlertDialog from "../components/AlertDialog.vue";
import ApiService from '../services/api.service'
import type { UserLoginForm } from "../models/user.model";


import { useRoute } from 'vue-router'
const dialog = ref()

const route = useRoute();
const formValues = ref(null as UserLoginForm);

const fields = [
    {
        label: "Email",
        name: "email",
        type: "text",
        placeholder: "max@mustermann.de"
    },
    {
        label: "Password",
        name: "password",
        type: "password"
    }
]

function onSubmit(values: any, actions: any) {
    return ApiService.postPage("/login", values).then(
        (page: PageNoContent) => {
            if (page.status != 204) {
                dialog.value.alert("error on login");
                return;

            }
            if (route.query.redirect) {
                router.push(route.query.redirect.toString());
                return            
            }
            router.push("/");
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
    formValues.value = {
        "email": "",
        "password": ""
    }
});



</script>
  