<script setup lang="ts">
import DefaultLayout from "@/Layouts/DefaultLayout.vue";
import Button from "@/Components/Button.vue";
import { onMounted, ref } from "vue";
import { type Type } from "@/types.ts";
import { apiGetTypes } from "@/api/api";

const types = ref<Type[]>([]);

async function getTypes() {
  try {
    types.value = await apiGetTypes();
  } catch (error) {
    alert("unable to get types:" + error);
  }
}

onMounted(() => {
  getTypes();
});
</script>

<template>
  <DefaultLayout>
    <template v-slot:header>
      <input placeholder="New Type" />
      <Button class="inverted">+</Button>
    </template>
    <template v-slot:main>
      <ul>
        <li v-for="type in types">
          <div>
            <input
              type="text"
              class="icon-input"
              v-model="type.Icon"
              placeholder="Icon"
            />
            <input
              type="text"
              class="name-input"
              v-model="type.Name"
              placeholder="Name"
            />
          </div>
          <div>
            <Button>↑</Button>
            <Button>↓</Button>
          </div>
        </li>
      </ul>
    </template>
  </DefaultLayout>
</template>

<style scoped>
.icon-input {
  text-align: center;
  width: 2.5em;
}

.name-input {
  width: 10em;
}

div {
  display: flex;
  align-items: center;
  gap: 0.5em;
}

li {
  padding: 0.1em 0em;
  margin: 0.2em 0em;
  border-radius: var(--border-radius);
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: var(--color-surface);
}

@media (hover: hover) and (pointer: fine) {
  li:hover {
    background: var(--color-surface-hover);
  }
}
</style>
