<script setup lang="ts">
import DefaultLayout from "@/Layouts/DefaultLayout.vue";
import VerticalMenu from "@/Components/VerticalMenu.vue";
import LinkButton from "@/Components/LinkButton.vue";
import Button from "@/Components/Button.vue";
import { apiCreateList, apiDeleteSession, apiGetSession } from "@/api/api.ts";
import { router } from "@/router.ts";
import { type List } from "@/types.ts";
import { onMounted, ref } from "vue";

const loggedIn = ref<boolean>(false);
const loaded = ref<boolean>(false);

async function listAsGuest() {
  try {
    const list = await apiCreateList({} as List);

    router.push({
      name: "list",
      params: { id: list.ID },
    });
  } catch (error) {
    alert("Unable to get list from server: " + error);
  }
}

async function checkSession() {
  try {
    await apiGetSession();
    loggedIn.value = true;
  } catch (error) {
    loggedIn.value = false;
  }
}

async function logout() {
  await apiDeleteSession();
  loggedIn.value = false;
}

onMounted(async () => {
  await checkSession();
  loaded.value = true;
});
</script>

<template>
  <DefaultLayout>
    <template v-slot:header>
      <h1>Listinator</h1>
    </template>
    <template v-slot:main>
      <VerticalMenu v-if="loaded">
        <Button @click="listAsGuest">New List as Guest</Button>
        <LinkButton v-if="!loggedIn" to="/login">Login</LinkButton>
        <Button v-if="loggedIn" @click="logout">Logout</Button>
      </VerticalMenu>
    </template>
  </DefaultLayout>
</template>
