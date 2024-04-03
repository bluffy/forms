import axios from "axios";


const host = import.meta.env.VITE_APP_API_HOST
const path = import.meta.env.VITE_APP_API_PATH
const dev = import.meta.env.DEV

const instance = axios.create({
  withCredentials: (dev)?true:false,
  baseURL: host + path,
  headers: {
    "Content-Type": "application/json",
  },
});


export default instance;


