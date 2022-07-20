<template>
  <n-modal :show="value">
    <n-card
      title="Добавить игру"
      style="width: 500px"
    >
      <n-form ref="formRef" :model="form" :rules="rules">
        <n-form-item path="name" label="Имя">
          <n-input v-model:value="form.name" placeholder="The Binding Of Isaac: Four Souls"/>
        </n-form-item>
        <n-form-item path="description" label="Описание">
          <n-input v-model:value="form.description" placeholder="bla bla"/>
        </n-form-item>
        <n-form-item path="image" label="Ссылка на изображение">
          <n-input v-model:value="form.image" placeholder="https://i.imgur.com/BAie9V7.jpeg"/>
        </n-form-item>
        <div class="add-game-modal buttons-block">
          <n-button @click="onClose">Отменить</n-button>
          <n-button type="primary" @click="onSubmit">Создать</n-button>
        </div>
      </n-form>
    </n-card>
  </n-modal>
</template>

<script lang="ts" setup>
import { defineProps, defineEmits, ref } from 'vue'
import { FormInst, FormRules, FormValidationError } from 'naive-ui'
import { IGameCreatePayload } from '@/models/Game'

defineProps({
  value: {
    type: Boolean,
    required: true
  }
})
const emit = defineEmits(['close', 'submit'])

const formRef = ref<FormInst | null>(null)

const form = ref<IGameCreatePayload>({
  name: '',
  description: '',
  image: ''
})

const rules: FormRules = {
  name: [
    {
      required: true,
      trigger: ['input', 'blur']
    }
  ],
  description: [
    {
      required: true,
      trigger: ['input', 'blur']
    }
  ],
  image: [
    {
      required: true,
      trigger: ['input', 'blur']
    }
  ]
}

const onSubmit = (e: MouseEvent): void => {
  e.preventDefault()
  formRef.value?.validate((errors: Array<FormValidationError> | undefined) => {
    if (!errors) {
      emit('submit', form)
    }
  })
}
const onClose = (): void => {
  emit('close')
}
</script>

<style lang="scss">
.add-game-modal.buttons-block {
  margin-top: 10px;
  display: flex;
  justify-content: space-between;
}
</style>
