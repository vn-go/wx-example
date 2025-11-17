import PrimeVue from 'primevue/config';
import { createApp } from 'vue';
import App from './App.vue';
import './style.css';
const app = createApp(App);

// createApp(App).mount('#app')
app.use(PrimeVue, {
    ripple: true,  // tùy chọn
    locale: {
        startsWith: 'Starts with',
        contains: 'Contains',
        notContains: 'Not contains',
        endsWith: 'Ends with',
        equals: 'Equals',
        notEquals: 'Not equals',
        noFilter: 'No Filter',
        aria: {
            trueLabel: 'True',
            falseLabel: 'False'
        }
    }
});
app.mount('#app');