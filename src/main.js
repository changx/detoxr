import axios from 'axios'
import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './site.css'

axios.defaults.withCredentials = true
axios.defaults.baseURL = process.env.VUE_APP_API

const app = createApp(App);
app.use(router);
app.mount('#app')
