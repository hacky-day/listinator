<script setup lang="ts">
import DefaultLayout from "@/Layouts/DefaultLayout.vue";
import VerticalMenu from "@/Components/VerticalMenu.vue";
import LinkButton from "@/Components/LinkButton.vue";
import Button from "@/Components/Button.vue";
import { apiCreateList, apiDeleteSession, apiGetSession } from "@/api/api.ts";
import { router } from "@/router.ts";
import { type List } from "@/types.ts";
import { onMounted, ref } from "vue";
import { useNotificationManager } from "@/composables/useNotificationManager";
import { useConnectivity } from "@/composables/useConnectivity";
import { t } from "@/i18n";

const { show } = useNotificationManager();
const { isOnline } = useConnectivity();

const loggedIn = ref<boolean>(false);
const loaded = ref<boolean>(false);

async function listAsGuest() {
  // Creating a list always hits the server (there's no offline-created-list
  // flow), so tell the user upfront instead of letting the fetch fail.
  if (!isOnline.value) {
    show("error", t("home.offlineCreateList"));
    return;
  }
  try {
    const list = await apiCreateList({} as List);

    router.push({
      name: "list",
      params: { id: list.ID },
    });
  } catch (error) {
    show("error", t("home.errorCreateList"), { logMessage: error });
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
        <Button @click="listAsGuest">{{ t("home.newListGuest") }}</Button>
        <LinkButton v-if="!loggedIn" to="/login">{{
          t("common.login")
        }}</LinkButton>
        <Button v-if="loggedIn" @click="logout">{{ t("common.logout") }}</Button>
      </VerticalMenu>
    </template>
  </DefaultLayout>
</template>
