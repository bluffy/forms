import type { AppError } from "../models/app.model";

export function genResponseError(error: any): AppError {
  if (error.response && error.response.data && error.response.data.error) {
    return <AppError>error.response.data.error;
  }
  const err = <AppError>{};
  err!.message = error.message || error.toString();
  return err;
}
