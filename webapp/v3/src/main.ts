
import App from './App.vue'
import { createApp } from 'vue'
import { registerPlugins } from './plugins'

//import './style.css'
import './style/App.scss'
import 'bootstrap/dist/js/bootstrap.js'


const app = createApp(App)


registerPlugins(app)

app.mount('#app')
