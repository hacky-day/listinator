<template>
  <button>
    <slot></slot>
  </button>
</template>

<style scoped>
button.inverted {
  color: var(--color-text);
  background: var(--color-surface);
}

button {
  height: 2.5em;
  min-width: 2.5em;
  border-radius: var(--border-radius);
  border: 0em;
  background: var(--color-accent);
  color: var(--color-text-inverted);
  transition: all 0.15s ease;
}

@media (hover: hover) and (pointer: fine) {
  button:hover {
    background: var(--color-accent-dark);
    cursor: pointer;
  }
  button.inverted:hover {
    background: var(--color-surface-hover);
  }
}

button:active {
  transform: scale(0.95);
  background: var(--color-primary-dark);
}

button.inverted:active {
  background: var(--color-border);
  transform: scale(0.95);
}
</style>
