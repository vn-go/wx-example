import { createApp } from 'vue';
import App from './App.vue';
import VSCodeLayout from './layouts/LayoutPanel.vue'; //<-- can i use this
import './style.css';

const app = createApp(App)
app.component('VSCodeLayout', VSCodeLayout)
app.mount('#app')