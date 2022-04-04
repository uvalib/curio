import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import WaitSpinner from "@/components/WaitSpinner.vue"

const app = createApp(App)

// provide store access to the rouer
store.router = router

// bind store and router to all componens as $store and $router
app.use(store)
app.use(router)

app.component('WaitSpinner', WaitSpinner)
import '@fortawesome/fontawesome-free/css/all.css'


// actually mount to DOM
app.mount('#app')
