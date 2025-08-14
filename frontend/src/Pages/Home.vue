<script setup lang="ts">
import DefaultLayout from "@/Layouts/DefaultLayout.vue";
import VerticalMenu from "@/Components/VerticalMenu.vue";
import LinkButton from "@/Components/LinkButton.vue";
import Button from "@/Components/Button.vue";
import { apiCreateList } from "@/api/api.ts";
import { router } from "@/router.ts";
import { type List } from "@/types.ts";

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
</script>

<template>
  <DefaultLayout>
    <template v-slot:header>
      <h1>Listinator</h1>
    </template>
    <template v-slot:main>
      <VerticalMenu>
        <Button @click="listAsGuest">New List as Guest</Button>
      </VerticalMenu>
    </template>
  </DefaultLayout>
</template>
