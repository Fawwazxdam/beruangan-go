<template>
  <header class="sticky top-0 z-40 bg-paper/85 backdrop-blur-sm border-b-[3px] border-ink">
    <div
      class="mx-auto flex h-14 sm:h-16 max-w-5xl items-center justify-between gap-4 px-4 sm:px-8"
    >
      <RouterLink to="/" class="flex items-center gap-2.5 no-underline text-ink">
        <span
          class="grid h-8 w-8 sm:h-9 sm:w-9 place-items-center border-2 border-ink bg-brut-yellow shadow-brut-sm"
        >
          <PhHorse weight="fill" class="h-4 w-4 sm:h-5 sm:w-5" />
        </span>
        <span class="text-lg sm:text-xl font-black tracking-tight">Beruangan</span>
      </RouterLink>

      <nav class="hidden sm:flex items-center gap-3">
        <RouterLink
          to="/"
          class="border-2 border-transparent px-3 py-1.5 text-sm font-bold text-ink no-underline transition-all hover:border-ink hover:bg-surface hover:shadow-brut-sm"
          active-class="border-ink bg-surface font-black shadow-brut-sm"
          >Dashboard</RouterLink
        >
        <RouterLink
          to="/transactions"
          class="border-2 border-transparent px-3 py-1.5 text-sm font-bold text-ink no-underline transition-all hover:border-ink hover:bg-surface hover:shadow-brut-sm"
          active-class="border-ink bg-surface font-black shadow-brut-sm"
          >Kelola Transaksi</RouterLink
        >
      </nav>

      <span
        class="sm:hidden border-2 border-ink bg-surface px-2 py-0.5 text-[10px] font-black tracking-wide shadow-brut-sm"
        >IDR</span
      >
    </div>
  </header>

  <RouterView />

  <nav
    class="sm:hidden fixed inset-x-0 bottom-0 z-40 border-t-[3px] border-ink bg-paper/95 backdrop-blur-sm"
  >
    <div class="grid grid-cols-2 gap-2 px-3 pb-[calc(0.75rem+env(safe-area-inset-bottom))] pt-3">
      <RouterLink
        to="/"
        class="flex flex-col items-center gap-1 border-2 border-transparent px-2 py-2 text-ink no-underline transition-all"
        active-class="border-ink bg-surface shadow-brut-sm"
      >
        <PhSquaresFour weight="fill" class="h-5 w-5" />
        <span class="text-[11px] font-black tracking-wide">Dashboard</span>
      </RouterLink>
      <RouterLink
        to="/transactions"
        class="flex flex-col items-center gap-1 border-2 border-transparent px-2 py-2 text-ink no-underline transition-all"
        active-class="border-ink bg-surface shadow-brut-sm"
      >
        <PhNote weight="fill" class="h-5 w-5" />
        <span class="text-[11px] font-black tracking-wide">Transaksi</span>
      </RouterLink>
    </div>
  </nav>
</template>

<script setup>
import { PhWallet, PhSquaresFour, PhNote, PhHorse } from '@phosphor-icons/vue'
</script>
