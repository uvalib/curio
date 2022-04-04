import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'

import WaitSpinner from "@/components/WaitSpinner.vue"

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)
app.use(router)

app.component('WaitSpinner', WaitSpinner)
import '@fortawesome/fontawesome-free/css/all.css'


// actually mount to DOM
app.mount('#app')
