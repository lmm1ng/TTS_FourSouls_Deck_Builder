<template>
  <div class="games-wrapper">
    <add-game-modal
     :value="isAddModal"
     @close="closeModal"
     @submit="onCreate"
    />
    <header class="games-wrapper__header">
      <n-button type="primary">Добавить игру</n-button>
    </header>
    <div class="games-wrapper__content">

    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, Ref, toRaw, onMounted } from 'vue'
import AddGameModal from '@/components/AddGameModal.vue'
import { IGameCreatePayload } from '@/models/Game'
import useGameStore from '@/stores/game'

const gameStore = useGameStore()

const isAddModal = ref(true)

const closeModal = (): void => {
  isAddModal.value = false
}
const onCreate = (form: Ref<IGameCreatePayload>): void => {
  gameStore.fetchCreateGame(toRaw(form.value))
}

onMounted(() => {
  gameStore.fetchGames()
})
</script>

<style lang="scss">
.games-wrapper {
  height: 100%;
  &__header {
    padding-top: 10px;
    display: flex;
    flex-direction: row-reverse;
  }
}
</style>

