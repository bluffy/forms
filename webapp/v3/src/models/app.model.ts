
 
 export interface Field {
    [key: string]: string;
  }
  

interface ResponseError {
    fields?:  Field
    message?: string;
    code?: number;
    logout?: boolean;
}
export type AppError = ResponseError | null;  