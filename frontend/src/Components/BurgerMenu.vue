<script setup lang="ts">
import { ref } from "vue";

const visible = ref<boolean>(false);
</script>

<template>
  <div>
    <svg
      @click="visible = true"
      xmlns="http://www.w3.org/2000/svg"
      width="20"
      height="20"
      fill="white"
      viewBox="0 0 24 24"
    >
      <path d="M3 6h18v2H3V6zm0 5h18v2H3v-2zm0 5h18v2H3v-2z" />
    </svg>
    <transition>
      <div v-if="visible" class="menu">
        <div class="close" @click="visible = false">✕</div>
        <a class="item" href="home">Home</a>
        <a class="item" href="home">Google</a>
        <a class="item" href="home">Google</a>
      </div>
    </transition>
  </div>
</template>

<style scoped>
.v-enter-active,
.v-leave-active {
  transition: transform 0.4s ease;
}

.v-enter-from,
.v-leave-to {
  transform: translateX(-100%);
}

.v-enter-to,
.v-leave-from {
  transform: translateX(0);
}
.menu {
  position: absolute;
  left: 0em;
  width: 50%;
  height: 100%;
  top: 0em;
  box-sizing: border-box;
  padding: 2em;
  border: 0.1em solid var(--color-border);
  z-index: 9999;
  display: flex;
  flex-direction: column;
  background: var(--color-background);
  box-shadow: 0.1em 0.1em 0.1em var(--color-border);
}

.item {
  text-align: center;
  padding: 1em 0em;
  margin: 0.1em;
  text-decoration: none;
  color: var(--color-text);
}

.close {
  position: absolute;
  top: 0.3em;
  left: 0.3em;
  background: none;
  border: none;
  font-size: 1.5em;
  cursor: pointer;
  color: var(--color-text);
}
</style>
