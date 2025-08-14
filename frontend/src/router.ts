import { createWebHashHistory, createRouter } from "vue-router";

import Home from "./Pages/Home.vue";
import List from "./Pages/List.vue";

const routes = [
  { path: "/", component: Home, name: "home" },
  { path: "/list/:id", component: List, name: "list" },
];

export const router = createRouter({
  history: createWebHashHistory(),
  routes,
});
