<template>
  <Transition
    enter-active-class="transition-all duration-200 ease-out"
    enter-from-class="opacity-0"
    enter-to-class="opacity-100"
    leave-active-class="transition-all duration-150 ease-in"
    leave-from-class="opacity-100"
    leave-to-class="opacity-0"
  >
    <div
      v-if="visible"
      class="fixed inset-0 z-[110] flex items-center justify-center bg-ink/40 backdrop-blur-sm"
      @click.self="cancel"
    >
      <div
        class="p-4 sm:p-6 bg-surface border-[3px] border-ink shadow-brut-lg max-w-md w-[calc(100%-2rem)] sm:w-full mx-auto sm:mx-4"
      >
        <div class="flex items-start gap-3">
          <div
            class="flex items-center justify-center w-10 h-10 shrink-0 border-2 border-ink"
            :class="danger ? 'bg-brut-red text-white' : 'bg-brut-cyan text-ink'"
          >
            <PhWarningCircle v-if="danger" weight="fill" class="w-6 h-6" />
            <PhQuestion v-else weight="bold" class="w-6 h-6" />
          </div>
          <div class="min-w-0">
            <h3 class="text-lg font-black leading-snug">{{ title }}</h3>
            <p class="mt-1 text-sm font-bold text-ink break-words">{{ message }}</p>
          </div>
        </div>

        <div class="mt-6 flex flex-col-reverse sm:flex-row gap-3">
          <button
            type="button"
            @click="cancel"
            class="flex-1 py-2 text-ink border-[3px] border-ink bg-surface rounded-none shadow-brut-sm font-black cursor-pointer transition-all hover:translate-y-[-2px] hover:translate-x-[-2px]"
          >
            {{ cancelText }}
          </button>
          <button
            type="button"
            @click="confirm"
            class="flex-1 py-2 border-[3px] border-ink rounded-none shadow-brut-sm font-black cursor-pointer transition-all hover:translate-y-[-2px] hover:translate-x-[-2px]"
            :class="danger ? 'bg-brut-red text-white' : 'bg-brut-cyan text-ink'"
          >
            {{ confirmText }}
          </button>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { PhQuestion, PhWarningCircle } from '@phosphor-icons/vue'

const visible = ref(false)
const title = ref('Konfirmasi')
const message = ref('')
const confirmText = ref('Ya, Lanjut')
const cancelText = ref('Batal')
const danger = ref(false)

let resolver = null

const open = (options = {}) => {
  if (resolver) settle(false)

  title.value = options.title || 'Konfirmasi'
  message.value = options.message || ''
  confirmText.value = options.confirmText || 'Ya, Lanjut'
  cancelText.value = options.cancelText || 'Batal'
  danger.value = Boolean(options.danger)
  visible.value = true

  return new Promise((resolve) => {
    resolver = resolve
  })
}

const settle = (value) => {
  visible.value = false
  const resolve = resolver
  resolver = null
  resolve?.(value)
}

const confirm = () => settle(true)
const cancel = () => settle(false)

const onKeydown = (event) => {
  if (!visible.value) return
  if (event.key === 'Escape') cancel()
}

onMounted(() => window.addEventListener('keydown', onKeydown))
onUnmounted(() => window.removeEventListener('keydown', onKeydown))

defineExpose({ open })
</script>
