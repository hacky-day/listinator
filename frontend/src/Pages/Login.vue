<script setup lang="ts">
import { ref } from "vue";

import { router } from "@/router.ts";

import DefaultLayout from "@/Layouts/DefaultLayout.vue"
import Button from "@/Components/Button.vue"
import { apiCreateSession } from "@/api/api";

const username = ref<string>("");
const password = ref<string>("");

async function login(event: Event) {
    event.preventDefault();
    try {
      await apiCreateSession(username.value, password.value)
      router.push({name: "home"})
    } catch (error) {
      alert("login failed:" + error);
    }
  }

</script>

<template>
  <DefaultLayout>
    <template v-slot:header>
      <h1>Listinator - Login</h1>
    </template>
    <template v-slot:main>
      <form id="loginForm">
        <input type="text" id="nameInput" placeholder="Name" v-model="username"></input>
        <input type="password" id="passwordInput" placeholder="Password" v-model="password"></input>
        <Button @click="login" type="submit">Login</Button>
      </form>
    </template>
  </DefaultLayout>
</template>

<style>
#loginForm {
  background: var(--color-surface);
  margin: 1svh 1svh 0svh 1svh;
  border-radius: var(--border-radius);
  min-height: 60svh;
  flex-direction: column;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1em;
}

#loginForm input {
  padding: 0em 0.5em;
}
</style>
