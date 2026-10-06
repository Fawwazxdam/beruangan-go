<template>
  <Transition
    enter-active-class="transition-all duration-300 ease-out"
    enter-from-class="opacity-0 translate-y-4"
    enter-to-class="opacity-100 translate-y-0"
    leave-active-class="transition-all duration-200 ease-in"
    leave-from-class="opacity-100 translate-y-0"
    leave-to-class="opacity-0 translate-y-4"
  >
    <div
      v-if="visible"
      :class="type === 'error' ? 'bg-brut-red text-white' : 'bg-brut-cyan text-black'"
      class="fixed bottom-4 right-4 sm:bottom-8 sm:right-8 flex items-center gap-3 px-4 py-3 border-4 border-black shadow-brut font-black z-[100]"
    >
      <PhCheckCircle v-if="type === 'success'" weight="fill" class="w-6 h-6" />
      <PhX v-else weight="bold" class="w-6 h-6" />
      <span>{{ message }}</span>
    </div>
  </Transition>
</template>

<script setup>
import { ref } from 'vue'
import { PhCheckCircle, PhX } from '@phosphor-icons/vue'

const visible = ref(false)
const message = ref('')
const type = ref('success')

// Fungsi ini yang bakal dipanggil dari luar
const show = (msg, toastType = 'success') => {
  message.value = msg
  type.value = toastType
  visible.value = true

  setTimeout(() => {
    visible.value = false
  }, 3000)
}

// Buka kunci fungsinya biar bisa diakses oleh file parent
defineExpose({ show })
</script>
