import { createApp } from 'vue';
import App from './App.vue';


import PrimeVue from 'primevue/config';

import emitter from '@core/eventBus';
import libs from '@core/lib';


import Aura from '@primevue/themes/aura';
import 'vue-select/dist/vue-select.css';
import './primevue-fix.css';
import './style.aura.hack.css';
import './style.css';
import './style.form.css';
import './style.grid.form.css';
import './style.primevue-override.css';
import './style.primvue.table.hack.css';
import './style.v-select.hack.css';
const app = createApp(App);

app.provide('emitter', emitter);

app.use(PrimeVue, {
    ripple: true,
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
        },
        firstDayOfWeek: 1,
    },
    theme: {
        preset: Aura,
        options: {
            darkModeSelector: null,
            cssLayer: false
        }
    }
});


//app.component('DatePicker', DatePicker);

app.mount('#app');

libs.currentApp = app;
