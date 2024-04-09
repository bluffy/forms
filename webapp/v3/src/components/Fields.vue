<template>

    <template v-for="field in fields">

        <div :class="defCols(field)">
            <template v-if="field.type == 'checkbox'">
                <Field :id="field.name" :type="field.type" :value="true" :name="field.name" class="form-check-input"
                    :placeholder="field.placeholder" />
                <label :for="field.name" class="form-check-label">{{ field.label }}</label>
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
</template>
<script lang="ts" setup name="Fields">

import { FormField } from "../models/app.model";
import { Field, ErrorMessage } from 'vee-validate';



/*
export type FormField = {
    label?: string
    name: string
    type: string
    col_md?: number
    placeholder?: string
   }
   */
defineProps<{
    fields: FormField[]
}>()




function defCols(field: any) {
    var ret = "";
    if (field.col) {
        ret = field.col;
    } else {
        ret = "col-12"
    }
    if (field.col_xs) {
        ret = ret + " col-xs-" + field.col_xs;
    }
    if (field.col_sm) {
        ret = ret + " col-sm-" + field.col_sm;
    }

    if (field.col_md) {
        ret = ret + " col-md-" + field.col_md;
    }
    if (field.col_lg) {
        ret = ret + " col-lg-" + field.col_lg;
    }
    if (field.col_xl) {
        ret = ret + " col-xl-" + field.col_xl;
    }
    if (field.col_xxl) {
        ret = ret + " col-xxl-" + field.col_xxl;
    }

    return ret;
}


</script>