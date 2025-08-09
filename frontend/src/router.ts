import { createWebHashHistory, createRouter } from "vue-router";

import Home from "./Pages/Home.vue";
import List from "./Pages/List.vue";
import Login from "./Pages/Login.vue";

const routes = [
  { path: "/", component: Home, name: "home" },
  { path: "/list/:id", component: List, name: "list" },
  { path: "/login", component: Login, name: "login" },
];

export const router = createRouter({
  history: createWebHashHistory(),
  routes,
});
