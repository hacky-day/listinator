<template>
  <div class="vertical-menu"><slot></slot></div>
</template>

<style>
.vertical-menu {
  min-height: 60svh;
  flex-direction: column;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1em;
}
.vertical-menu * {
  width: 10em;
}
</style>
