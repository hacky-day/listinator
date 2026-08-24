import { createApp } from "vue";
import { registerSW } from "virtual:pwa-register";
import { router } from "@/router";

import App from "./App.vue";

registerSW({ immediate: true });
void navigator.storage?.persist?.();

const app = createApp(App);

app.use(router);

app.mount("#app");
