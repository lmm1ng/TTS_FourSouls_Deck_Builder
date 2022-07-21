<template>
  <div class="games-wrapper">
    <add-game-modal
     :value="isAddModal"
     @close="closeModal"
     @submit="onCreate"
    />
    <header class="games-wrapper__header">
      <n-button type="primary" @click="openModal">Добавить игру</n-button>
    </header>
    <div class="games-wrapper__content">
      <entity-card
        v-for="game in games"
        :key="game.id"
        :title="game.name"
        :description="game.description"
        :image="game.image"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, Ref, toRaw, onMounted, computed } from 'vue'
import AddGameModal from '@/components/AddGameModal.vue'
import EntityCard from '@/components/EntityCard.vue'
import { IGameCreatePayload } from '@/models/Game'
import useGameStore from '@/stores/game'

const gameStore = useGameStore()

const isAddModal = ref(false)

const closeModal = (): void => {
  isAddModal.value = false
}
const openModal = (): void => {
  isAddModal.value = true
}
const onCreate = (form: Ref<IGameCreatePayload>): void => {
  gameStore.fetchCreateGame(toRaw(form.value)).then(() => {
    gameStore.fetchGames()
    isAddModal.value = false
  })
}

const games = computed(() => gameStore.games)

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
  &__content {
    display: flex;
    flex-wrap: wrap;
    gap: 30px;
  }
}
</style>

