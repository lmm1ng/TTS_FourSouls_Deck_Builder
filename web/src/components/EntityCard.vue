<template>
  <n-card
    :closable="true"
    :title="title"
    style="width: 400px"
    @close="onClose"
  >
    <div style="margin-bottom: 20px;">{{ description }}</div>
    <n-image
      :src="image"
      width="300"
    />
  </n-card>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue'

defineProps({
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    required: true
  },
  image: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['close', 'click'])

const onClose = (): void => {
  emit('close')
}

const onClick = (): void => {
  console.log('clicked')
}
</script>
