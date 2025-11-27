import PrimeVue from 'primevue/config';

import emitter from '@core/eventBus';
import libs from '@core/lib';
import Aura from '@primeuix/themes/aura';
import { createApp } from 'vue';
import 'vue-select/dist/vue-select.css';
import App from './App.vue';
import './primevue-fix.css';
import './style.css';
import './style.form.css';
import './style.grid.form.css';
import './style.primvue.table.hack.css';
const app = createApp(App);
app.provide('emitter', emitter);

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
    },
    theme: {
        preset: Aura,
        options: {
            darkModeSelector: null,   // tắt dark mode
            cssLayer: false           // không inject color-scheme vào CSS
        }
    }
    //unstyled: true,
    //pt: Aura
});
app.mount('#app');
libs.currentApp = app;