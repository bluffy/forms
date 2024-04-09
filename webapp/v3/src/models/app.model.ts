
 
 export interface Field {
    [key: string]: string;
  }
  

  export type FormField = {
    label?: string
    name: string
    type: string
    col?: number
    placeholder?: string
    col_md?: number
    col_xs?: number
    col_xl?: number
    col_sm?: number
    col_lg?:number
    col_xxl?: number
   }
 

interface ResponseError {
    fields?:  Field
    message?: string;
    code?: number;
    logout?: boolean;
}
export type AppError = ResponseError | null;  