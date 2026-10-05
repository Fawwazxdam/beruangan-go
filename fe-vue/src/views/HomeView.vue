<template>
  <main class="max-w-5xl mx-auto p-4 sm:p-8 font-sans">
    <header
      class="mb-6 sm:mb-8 flex flex-col sm:flex-row sm:justify-between sm:items-end gap-3 sm:gap-4"
    >
      <div>
        <h1 class="flex items-center gap-2 sm:gap-3 text-2xl sm:text-3xl font-black text-ink">
          <PhChartLine weight="fill" class="w-7 h-7 sm:w-9 sm:h-9" />
          Dashboard Keuangan
        </h1>
        <p class="mt-1 text-sm sm:text-base text-ink">Pantau hutang dan tagihanmu bulan ini.</p>
      </div>
      <span
        class="flex items-center gap-1.5 self-start sm:self-auto border-2 border-ink bg-surface px-2.5 py-1 text-[11px] sm:text-xs font-black uppercase tracking-wide shadow-brut-sm"
      >
        <PhCalendarBlank weight="bold" class="w-3.5 h-3.5 sm:w-4 sm:h-4" />
        {{ bulanLabel }}
      </span>
    </header>

    <div
      v-if="isLoading"
      class="py-12 border-[3px] border-ink bg-surface shadow-brut text-center text-ink font-black text-lg"
    >
      Mengambil data dari server...
    </div>

    <div v-else>
      <div class="grid grid-cols-1 gap-4 sm:gap-6 md:grid-cols-3">
        <article
          class="p-4 sm:p-6 text-ink bg-brut-cyan border-[3px] border-ink shadow-brut transition-all hover:translate-y-[-2px]"
        >
          <div class="flex items-start justify-between gap-3">
            <div>
              <h3 class="text-xs sm:text-sm font-black uppercase tracking-wide">
                Tagihan Bulan Ini
              </h3>
              <p class="mt-0.5 text-[11px] sm:text-xs font-bold text-ink/70">
                Periode {{ bulanLabel }}
              </p>
            </div>
            <span
              class="grid shrink-0 w-8 h-8 place-items-center border-2 border-ink bg-surface shadow-brut-sm"
            >
              <PhReceipt weight="fill" class="w-4 h-4" />
            </span>
          </div>

          <p class="mt-3 text-xl sm:text-2xl md:text-3xl font-black break-words">
            {{ formatRupiah(summary.total_tagihan_bulan_ini) }}
          </p>

          <div class="mt-4">
            <div class="flex items-center justify-between gap-2 mb-1.5">
              <span class="text-[11px] sm:text-xs font-black uppercase">Sudah dibayar</span>
              <span class="text-sm sm:text-base font-black leading-none">{{ pctTagihan }}%</span>
            </div>
            <div class="w-full h-4 border-2 border-ink bg-surface overflow-hidden">
              <div
                class="h-full bg-ink transition-all duration-700"
                :style="{ width: pctTagihan + '%' }"
              ></div>
            </div>
            <p class="mt-1.5 text-[11px] sm:text-xs font-bold">
              {{ formatRupiah(summary.tagihan_terbayar_bulan_ini) }}
              <span class="text-ink/70"
                >dari {{ formatRupiah(summary.total_tagihan_bulan_ini) }}</span
              >
            </p>
          </div>
        </article>

        <article
          class="p-4 sm:p-6 text-ink bg-brut-pink border-[3px] border-ink shadow-brut transition-all hover:translate-y-[-2px]"
        >
          <div class="flex items-start justify-between gap-3">
            <div>
              <h3 class="text-xs sm:text-sm font-black uppercase tracking-wide">
                Total Hutang Aktif
              </h3>
              <p class="mt-0.5 text-[11px] sm:text-xs font-bold text-ink/70">
                Sisa yang belum lunas
              </p>
            </div>
            <span
              class="grid shrink-0 w-8 h-8 place-items-center border-2 border-ink bg-surface shadow-brut-sm"
            >
              <PhWallet weight="fill" class="w-4 h-4" />
            </span>
          </div>

          <p class="mt-3 text-xl sm:text-2xl md:text-3xl font-black break-words">
            {{ formatRupiah(sisaHutang) }}
          </p>

          <div class="mt-4 grid grid-cols-2 gap-2">
            <div class="border-2 border-ink bg-surface p-2">
              <p class="text-[10px] font-black uppercase text-ink/70">Total Awal</p>
              <p class="mt-0.5 text-xs sm:text-sm font-black break-words">
                {{ formatRupiah(summary.total_hutang) }}
              </p>
            </div>
            <div class="border-2 border-ink bg-surface p-2">
              <p class="text-[10px] font-black uppercase text-ink/70">Terbayar</p>
              <p class="mt-0.5 text-xs sm:text-sm font-black break-words">
                {{ formatRupiah(summary.total_hutang_terbayar) }}
              </p>
            </div>
          </div>
        </article>

        <article
          class="p-4 sm:p-6 text-ink bg-brut-lime border-[3px] border-ink shadow-brut transition-all hover:translate-y-[-2px]"
        >
          <div class="flex items-start justify-between gap-3">
            <div>
              <h3 class="text-xs sm:text-sm font-black uppercase tracking-wide">
                Total Lunas (All Time)
              </h3>
              <p class="mt-0.5 text-[11px] sm:text-xs font-bold text-ink/70">Akumulasi selesai</p>
            </div>
            <span
              class="grid shrink-0 w-8 h-8 place-items-center border-2 border-ink bg-surface shadow-brut-sm"
            >
              <PhCheckCircle weight="fill" class="w-4 h-4" />
            </span>
          </div>

          <p class="mt-3 text-xl sm:text-2xl md:text-3xl font-black break-words">
            {{ formatRupiah(summary.total_hutang_terbayar) }}
          </p>

          <div class="mt-4">
            <div class="flex items-end justify-between gap-2 mb-1.5">
              <span class="text-[11px] sm:text-xs font-black uppercase">Progress lunas</span>
              <span class="text-lg sm:text-xl font-black leading-none">{{ pctLunas }}%</span>
            </div>
            <div class="w-full h-4 border-2 border-ink bg-surface overflow-hidden">
              <div
                class="h-full bg-ink transition-all duration-700"
                :style="{ width: pctLunas + '%' }"
              ></div>
            </div>
            <p class="mt-1.5 text-[11px] sm:text-xs font-bold">
              <template v-if="summary.total_hutang > 0">
                {{ formatRupiah(summary.total_hutang_terbayar) }}
                <span class="text-ink/70">dari {{ formatRupiah(summary.total_hutang) }}</span>
              </template>
              <template v-else>
                <span class="text-ink/70">Belum ada hutang tercatat</span>
              </template>
            </p>
          </div>
        </article>
      </div>

      <div
        v-if="summary.total_tagihan_h3 > 0 || summary.reminder_h3.length > 0"
        class="mt-4 sm:mt-6 border-[3px] border-ink bg-brut-yellow shadow-brut"
      >
        <div class="flex flex-col sm:flex-row sm:items-center gap-3 p-4">
          <PhWarningCircle weight="fill" class="w-7 h-7 mx-auto sm:mx-0 text-brut-red" />
          <div class="text-center sm:text-left">
            <div class="flex flex-wrap items-center justify-center sm:justify-start gap-2">
              <h3 class="text-lg sm:text-xl font-black uppercase tracking-wide text-ink">
                Peringatan Jatuh Tempo H-3!
              </h3>
              <span
                class="border-2 border-ink bg-surface px-2 py-0.5 text-[11px] font-black shadow-brut-sm"
              >
                {{ summary.reminder_h3.length }} tagihan
              </span>
            </div>
            <p class="mt-1 text-sm sm:text-base text-ink">
              Ada tagihan sebesar
              <span class="font-black">{{ formatRupiah(summary.total_tagihan_h3) }}</span>
              yang harus dibayar dalam 3 hari ke depan.
            </p>
          </div>
        </div>
        <div class="px-4 pb-4">
          <button
            @click="showModal = true"
            class="flex items-center justify-center w-full py-2 gap-2 text-white border-2 border-ink bg-brut-red rounded-none shadow-brut font-black cursor-pointer transition-all hover:translate-y-[-2px] hover:translate-x-[-2px]"
          >
            <PhList weight="bold" class="w-5 h-5" />
            Lihat Detail
          </button>
        </div>
      </div>
    </div>

    <div
      v-if="showModal"
      class="fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-ink/40 backdrop-blur-sm"
      @click.self="showModal = false"
    >
      <div
        class="w-full sm:max-w-lg max-h-[85vh] flex flex-col overflow-hidden bg-surface border-t-[3px] sm:border-[3px] border-ink shadow-brut-lg"
      >
        <div
          class="flex justify-between items-center p-4 border-b-[3px] border-ink bg-brut-yellow shrink-0"
        >
          <h2 class="flex items-center gap-2 text-lg sm:text-xl font-black text-ink">
            <PhWarningCircle weight="fill" class="w-6 h-6 text-brut-red" />
            Segera Jatuh Tempo
          </h2>
          <button
            @click="showModal = false"
            class="flex items-center justify-center w-8 h-8 text-ink bg-brut-lime border-2 border-ink font-black cursor-pointer shadow-brut transition-all hover:translate-y-[-1px]"
          >
            <PhX weight="bold" class="w-4 h-4" />
          </button>
        </div>

        <div class="p-4 overflow-y-auto min-h-0">
          <div v-if="summary.reminder_h3.length === 0" class="py-4 text-center text-ink font-black">
            Aman! Tidak ada tagihan mendesak.
          </div>

          <div v-else class="space-y-3">
            <div
              v-for="trx in summary.reminder_h3"
              :key="trx.id"
              class="p-4 border-2 border-ink bg-surface shadow-brut-sm"
            >
              <div class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-2">
                <div>
                  <h4 class="font-black text-ink">{{ trx.title }}</h4>
                  <p class="mt-1 text-sm font-black text-brut-red">
                    Jatuh Tempo: {{ trx.due_date }}
                  </p>
                </div>
                <div class="text-left sm:text-right">
                  <div class="text-lg font-black text-ink">
                    {{ formatRupiah(trx.amount) }}
                  </div>
                  <div
                    class="px-2 py-1 mt-1 text-xs font-black text-ink bg-brut-yellow border border-ink inline-block"
                  >
                    PENDING
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="p-4 border-t-[3px] border-ink bg-brut-yellow shrink-0">
          <RouterLink
            to="/transactions"
            @click="showModal = false"
            class="flex items-center justify-center w-full py-2 gap-2 text-ink border-[3px] border-ink bg-brut-cyan rounded-none shadow-brut font-black cursor-pointer transition-all hover:translate-y-[-2px] hover:translate-x-[-2px]"
          >
            <PhCreditCard weight="fill" class="w-5 h-5" />
            Bayar Sekarang
          </RouterLink>
        </div>
      </div>
    </div>
  </main>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import {
  PhChartLine,
  PhWarningCircle,
  PhReceipt,
  PhWallet,
  PhCheckCircle,
  PhList,
  PhX,
  PhCreditCard,
  PhCalendarBlank,
} from '@phosphor-icons/vue'

const summary = ref({
  total_hutang: 0,
  total_hutang_terbayar: 0,
  total_tagihan_bulan_ini: 0,
  tagihan_terbayar_bulan_ini: 0,
  persentase_terbayar: 0,
  total_tagihan_h3: 0,
  reminder_h3: [],
})

const isLoading = ref(true)
const showModal = ref(false)

const clampPercent = (nilai) => Math.min(100, Math.max(0, Math.round(Number(nilai) || 0)))

const pctLunas = computed(() => clampPercent(summary.value.persentase_terbayar))

const pctTagihan = computed(() => {
  const total = Number(summary.value.total_tagihan_bulan_ini) || 0
  if (total <= 0) return 0
  return clampPercent(((Number(summary.value.tagihan_terbayar_bulan_ini) || 0) / total) * 100)
})

const sisaHutang = computed(() => {
  const total = Number(summary.value.total_hutang) || 0
  const lunas = Number(summary.value.total_hutang_terbayar) || 0
  return Math.max(0, total - lunas)
})

const bulanLabel = new Intl.DateTimeFormat('id-ID', {
  month: 'long',
  year: 'numeric',
}).format(new Date())

const fetchSummary = async () => {
  try {
    const response = await fetch('/api/summary')
    const data = await response.json()

    if (!data.reminder_h3) data.reminder_h3 = []
    if (data.total_tagihan_h3 === undefined) data.total_tagihan_h3 = 0

    summary.value = data
  } catch (error) {
    console.error('Gagal mengambil data backend:', error)
    alert('Pastikan server Golang sudah nyala!')
  } finally {
    isLoading.value = false
  }
}

const formatRupiah = (angka) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(angka)
}

onMounted(() => {
  fetchSummary()
})
</script>
