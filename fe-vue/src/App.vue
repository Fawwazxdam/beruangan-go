<template>
  <nav
    class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-2 sm:gap-6 px-4 sm:px-8 py-4 bg-white border-b-4 border-black"
  >
    <div
      class="flex items-center justify-center sm:justify-start gap-2 text-xl sm:text-2xl font-black"
    >
      <PhWallet weight="fill" class="w-6 h-6 sm:w-7 sm:h-7" />
      <span>Dompetku</span>
    </div>
    <div class="flex justify-center sm:justify-end gap-3 sm:gap-6">
      <RouterLink
        to="/"
        class="px-2 py-1 text-sm sm:text-base text-black font-bold no-underline transition-all border-2 border-transparent hover:border-black hover:shadow-brut"
        active-class="border-2 border-black font-black shadow-brut"
        >Dashboard</RouterLink
      >
      <RouterLink
        to="/transactions"
        class="px-2 py-1 text-sm sm:text-base text-black font-bold no-underline transition-all border-2 border-transparent hover:border-black hover:shadow-brut"
        active-class="border-2 border-black font-black shadow-brut"
        >Kelola Transaksi</RouterLink
      >
    </div>
  </nav>

  <RouterView />
</template>

<script setup>
import { PhWallet } from '@phosphor-icons/vue'
</script>
