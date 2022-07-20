<template>
  <div class="app-wrapper">
    <header class="app-wrapper__header">
    </header>
    <main class="app-wrapper__content">
      <router-view/>
    </main>
  </div>
</template>

<style lang="scss">
html, body, #app {
  padding: 0;
  margin: 0;
  height: 100vh;
}
.app-wrapper {
  height: 100%;
  display: flex;
  flex-direction: column;
  &__header {
    height: 40px;
    background-color: #5BAD92;
  }
  &__content {
    flex: 1;
    padding: 0 40px;
  }
}
</style>
