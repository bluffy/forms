export interface PageIndex {
    status: number
    data?: any
}

export interface PageNoContent {
    status: number
}
export interface PageRegister {
    status: number
    data: {
        message?: string
    }
}



