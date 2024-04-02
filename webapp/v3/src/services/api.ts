import axios from "axios";


const host = import.meta.env.VITE_APP_API_HOST
const path = import.meta.env.VITE_APP_API_PATH

const instance = axios.create({
  baseURL: host + path,
  headers: {
    "Content-Type": "application/json",
  },
});


export default instance;


