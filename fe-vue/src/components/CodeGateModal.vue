<template>
  <div class="fixed inset-0 z-[120] flex items-center justify-center bg-ink/60 backdrop-blur-sm">
    <div
      class="p-4 sm:p-6 bg-surface border-[3px] border-ink shadow-brut-lg max-w-sm w-[calc(100%-2rem)] mx-4"
    >
      <div class="flex items-center gap-3 mb-5">
        <div
          class="flex items-center justify-center w-10 h-10 shrink-0 border-2 border-ink bg-brut-yellow text-ink"
        >
          <PhVault weight="fill" class="w-6 h-6" />
        </div>
        <div class="min-w-0">
          <h3 class="text-lg font-black leading-tight">Brankas Terkunci</h3>
          <p class="text-xs font-bold text-ink">Masukkan kode rahasia kamu dulu, ya.</p>
        </div>
      </div>

      <form @submit.prevent="submit">
        <input
          ref="inputEl"
          v-model="code"
          type="password"
          required
          autocomplete="off"
          placeholder="Kode rahasia"
          class="w-full px-3 py-2 text-ink bg-surface border-2 border-ink rounded-none outline-none focus:border-ink focus:ring-2 focus:ring-ink"
        />
        <button
          type="submit"
          class="flex items-center justify-center w-full mt-4 py-3 gap-2 text-ink border-[3px] border-ink bg-brut-cyan rounded-none shadow-brut font-black cursor-pointer transition-all hover:translate-y-[-2px] hover:translate-x-[-2px]"
        >
          <PhLockOpen weight="bold" class="w-5 h-5" />
          Buka Brankas
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { PhVault, PhLockOpen } from '@phosphor-icons/vue'

const emit = defineEmits(['submit'])

const code = ref('')
const inputEl = ref(null)

const submit = () => {
  const safeCode = code.value.toLowerCase().trim()
  if (!safeCode) return
  emit('submit', safeCode)
}

onMounted(() => inputEl.value?.focus())
</script>
