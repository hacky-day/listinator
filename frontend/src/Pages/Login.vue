<script setup lang="ts">
import { ref } from "vue";

import { router } from "@/router.ts";

import DefaultLayout from "@/Layouts/DefaultLayout.vue";
import Button from "@/Components/Button.vue";
import { apiCreateSession } from "@/api/api";
import { useNotificationManager } from "@/composables/useNotificationManager";
import { t } from "@/i18n";

const { show } = useNotificationManager();

const username = ref<string>("");
const password = ref<string>("");

async function login(event: Event) {
  event.preventDefault();
  try {
    await apiCreateSession(username.value, password.value);
    router.push({ name: "home" });
  } catch (error) {
    show("error", t("login.failed"), { logMessage: error });
  }
}
</script>

<template>
  <DefaultLayout>
    <template v-slot:header>
      <h1>{{ t("login.title") }}</h1>
    </template>
    <template v-slot:main>
      <form id="loginForm">
        <input
          type="text"
          id="nameInput"
          :placeholder="t('common.name')"
          v-model="username"
        />
        <input
          type="password"
          id="passwordInput"
          :placeholder="t('login.password')"
          v-model="password"
        />
        <Button @click="login" type="submit">{{ t("common.login") }}</Button>
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
