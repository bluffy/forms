
import App from './App.vue'
import { createApp } from 'vue'
import { registerPlugins } from './plugins'

//import './style.css'
import './style/App.scss'
import 'bootstrap/dist/js/bootstrap.js'
import setupInterceptors from "./services/setupInterceptors";
import 'bootstrap-icons/font/bootstrap-icons.css'


const app = createApp(App)

setupInterceptors();

registerPlugins(app)


app.mount('#app')
