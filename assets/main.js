import { createApp } from "vue";
import { createPinia } from "pinia";
import router from "./router";
import vuetify from "./plugins/vuetify";
import App from "./App.vue";
import "@mdi/font/css/materialdesignicons.css";
import "./styles/global.css";

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(vuetify);

app.mount("#app");
