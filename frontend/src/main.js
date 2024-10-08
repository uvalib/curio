import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import Curio from './assets/theme/curio'
import TreeTable from 'primevue/treetable'
import Column from 'primevue/column'
import Image from 'primevue/image'
import InputText from 'primevue/inputtext'
import Button from 'primevue/button'
import ToastService from 'primevue/toastservice'
import WaitSpinner from "@/components/WaitSpinner.vue"
import 'primeicons/primeicons.css'

const app = createApp(App)
app.use(router)

app.use(PrimeVue, {
   theme: {
      preset: Curio,
      options: {
         prefix: 'p',
         darkModeSelector: '.dpg-dark'
      }
   }
})
app.use(ToastService)

app.component('WaitSpinner', WaitSpinner)
app.component('TreeTable',  TreeTable)
app.component('Column',  Column)
app.component('Image',  Image)
app.component('InputText', InputText)
app.component('Button', Button)

app.use(createPinia())

app.mount('#app')
