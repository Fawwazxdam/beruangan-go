<template>
  <main class="max-w-5xl mx-auto p-6 font-sans">
    <header
      class="mb-8 flex flex-col sm:flex-row sm:justify-between sm:items-end gap-4"
    >
      <div>
        <h1
          class="flex items-center gap-2 sm:gap-3 text-2xl sm:text-3xl font-black text-black"
        >
          <PhChartLine weight="fill" class="w-7 h-7 sm:w-9 sm:h-9" />
          Dashboard Keuangan
        </h1>
        <p class="mt-1 text-sm sm:text-base text-black">
          Pantau hutang dan tagihanmu bulan ini.
        </p>
      </div>
    </header>

    <div
      v-if="isLoading"
      class="py-12 text-center text-black font-black text-lg"
    >
      Mengambil data dari server...
    </div>

    <div v-else>
      <div class="grid grid-cols-1 gap-6 md:grid-cols-3">
        <div
          class="relative p-6 text-black transition-all bg-brut-cyan border-4 border-black shadow-brut hover:translate-y-[-2px]"
        >
          <PhReceipt weight="fill" class="w-6 h-6" />
          <h3 class="mt-4 mb-2 text-sm font-black uppercase">
            Tagihan Bulan Ini
          </h3>
          <h2 class="mb-4 text-2xl md:text-3xl font-black break-words">
            {{ formatRupiah(summary.total_tagihan_bulan_ini) }}
          </h2>
          <div
            class="p-3 border-2 border-black bg-black text-white rounded-none text-sm md:text-base"
          >
            Sudah dibayar:
            <span class="font-black">{{
              formatRupiah(summary.tagihan_terbayar_bulan_ini)
            }}</span>
          </div>
        </div>

        <div
          class="relative p-6 text-white transition-all bg-brut-pink border-4 border-black shadow-brut hover:translate-y-[-2px]"
        >
          <PhWallet weight="fill" class="w-6 h-6" />
          <h3 class="mt-4 mb-2 text-sm font-black uppercase">
            Total Hutang Aktif
          </h3>
          <h2 class="mb-4 text-2xl md:text-3xl font-black break-words">
            {{
              formatRupiah(summary.total_hutang - summary.total_hutang_terbayar)
            }}
          </h2>
          <div
            class="p-3 border-2 border-black bg-white text-black rounded-none text-sm md:text-base"
          >
            Total Awal:
            <span class="font-black">{{
              formatRupiah(summary.total_hutang)
            }}</span>
          </div>
        </div>

        <div
          class="relative p-6 text-black transition-all bg-brut-lime border-4 border-black shadow-brut hover:translate-y-[-2px]"
        >
          <PhCheckCircle weight="fill" class="w-6 h-6" />
          <h3 class="mt-4 mb-2 text-sm font-black uppercase">
            Total Lunas (All Time)
          </h3>
          <h2 class="mb-4 text-2xl md:text-3xl font-black break-words">
            {{ formatRupiah(summary.total_hutang_terbayar) }}
          </h2>
          <div
            class="p-3 border-2 border-black bg-black text-white rounded-none flex justify-between items-center text-sm md:text-base"
          >
            <span>Progress:</span>
            <span class="text-base md:text-lg font-black"
              >{{ summary.persentase_terbayar }}%</span
            >
          </div>
        </div>
      </div>

      <div
        v-if="summary.total_tagihan_h3 > 0 || summary.reminder_h3.length > 0"
        class="mt-6 border-4 border-black bg-brut-yellow shadow-brut"
      >
        <div class="flex flex-col sm:flex-row sm:items-center gap-3 p-4">
          <PhWarningCircle
            weight="fill"
            class="w-7 h-7 mx-auto sm:mx-0 text-brut-red"
          />
          <div class="text-center sm:text-left">
            <h3
              class="text-lg sm:text-xl font-black uppercase tracking-wide text-black"
            >
              Peringatan Jatuh Tempo H-3!
            </h3>
            <p class="mt-1 text-sm sm:text-base text-black">
              Ada tagihan sebesar
              <span class="font-black">{{
                formatRupiah(summary.total_tagihan_h3)
              }}</span>
              yang harus dibayar dalam 3 hari ke depan.
            </p>
          </div>
        </div>
        <div class="px-4 pb-4">
          <button
            @click="showModal = true"
            class="flex items-center justify-center w-full py-2 gap-2 text-white border-2 border-black bg-brut-red rounded-none shadow-brut font-black cursor-pointer transition-all hover:translate-y-[-2px] hover:translate-x-[-2px]"
          >
            <PhList weight="bold" class="w-5 h-5" />
            Lihat Detail
          </button>
        </div>
      </div>
    </div>

    <div
      v-if="showModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/75"
      @click.self="showModal = false"
    >
      <div
        class="w-[calc(100%-2rem)] sm:w-full sm:max-w-lg overflow-hidden bg-white border-4 border-black shadow-brut-lg"
      >
        <div
          class="flex justify-between items-center p-4 border-b-4 border-black bg-brut-yellow"
        >
          <h2
            class="flex items-center gap-2 text-lg sm:text-xl font-black text-black"
          >
            <PhWarningCircle weight="fill" class="w-6 h-6 text-brut-red" />
            Segera Jatuh Tempo
          </h2>
          <button
            @click="showModal = false"
            class="flex items-center justify-center w-8 h-8 text-black bg-brut-lime border-2 border-black font-black cursor-pointer shadow-brut transition-all hover:translate-y-[-1px]"
          >
            <PhX weight="bold" class="w-4 h-4" />
          </button>
        </div>

        <div class="p-4 max-h-[60vh] overflow-y-auto">
          <div
            v-if="summary.reminder_h3.length === 0"
            class="py-4 text-center text-black font-black"
          >
            Aman! Tidak ada tagihan mendesak.
          </div>

          <div v-else class="space-y-3">
            <div
              v-for="trx in summary.reminder_h3"
              :key="trx.id"
              class="p-4 border-2 border-black bg-white"
            >
              <div
                class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-2"
              >
                <div>
                  <h4 class="font-black text-black">{{ trx.title }}</h4>
                  <p class="mt-1 text-sm font-black text-brut-red">
                    Jatuh Tempo: {{ trx.due_date }}
                  </p>
                </div>
                <div class="text-right">
                  <div class="text-lg font-black text-black">
                    {{ formatRupiah(trx.amount) }}
                  </div>
                  <div
                    class="px-2 py-1 mt-1 text-xs font-black text-black bg-brut-yellow border border-black inline-block"
                  >
                    PENDING
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="p-4 border-t-4 border-black bg-brut-yellow">
          <RouterLink
            to="/transactions"
            @click="showModal = false"
            class="flex items-center justify-center w-full py-2 gap-2 text-black border-4 border-black bg-brut-cyan rounded-none shadow-brut font-black cursor-pointer transition-all hover:translate-y-[-2px] hover:translate-x-[-2px]"
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
import { ref, onMounted } from "vue";
import {
  PhChartLine,
  PhWarningCircle,
  PhReceipt,
  PhWallet,
  PhCheckCircle,
  PhList,
  PhX,
  PhCreditCard,
} from "@phosphor-icons/vue";

const summary = ref({
  total_hutang: 0,
  total_hutang_terbayar: 0,
  total_tagihan_bulan_ini: 0,
  tagihan_terbayar_bulan_ini: 0,
  persentase_terbayar: 0,
  total_tagihan_h3: 0,
  reminder_h3: [],
});

const isLoading = ref(true);
const showModal = ref(false);

const fetchSummary = async () => {
  try {
    const response = await fetch("/api/summary");
    const data = await response.json();

    if (!data.reminder_h3) data.reminder_h3 = [];
    if (data.total_tagihan_h3 === undefined) data.total_tagihan_h3 = 0;

    summary.value = data;
  } catch (error) {
    console.error("Gagal mengambil data backend:", error);
    alert("Pastikan server Golang sudah nyala!");
  } finally {
    isLoading.value = false;
  }
};

const formatRupiah = (angka) => {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    minimumFractionDigits: 0,
  }).format(angka);
};

onMounted(() => {
  fetchSummary();
});
</script>
