import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import TreeTable from 'primevue/treetable'
import Column from 'primevue/column'
import Image from 'primevue/image'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button';

import WaitSpinner from "@/components/WaitSpinner.vue"

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)
app.use(router)
app.use(PrimeVue)

app.component('WaitSpinner', WaitSpinner)
app.component('TreeTable',  TreeTable)
app.component('Column',  Column)
app.component('Image',  Image)
app.component('InputText', InputText)
app.component('Button', Button)
import '@fortawesome/fontawesome-free/css/all.css'
import 'primevue/resources/themes/saga-blue/theme.css ';
import 'primevue/resources/primevue.min.css';
import 'primeicons/primeicons.css';


// actually mount to DOM
app.mount('#app')
