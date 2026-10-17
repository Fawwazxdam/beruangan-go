<template>
  <main class="max-w-4xl mx-auto px-4 sm:px-8 py-8 sm:py-10">
    <div class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4 mb-8">
      <div>
        <h2 class="flex items-center gap-2 sm:gap-3 mb-2 text-2xl sm:text-4xl font-black text-ink">
          <PhWallet weight="fill" class="w-7 h-7 sm:w-9 sm:h-9" />
          Manajemen Dompet
        </h2>
        <p class="text-sm sm:text-base text-ink">Pantau saldo semua dompetmu dalam satu tempat.</p>
      </div>
    </div>

    <div
      v-if="isLoading"
      class="py-12 border-[3px] border-ink bg-surface shadow-brut rounded-brut text-center text-ink font-black text-lg"
    >
      Mengambil data dari server...
    </div>

    <template v-else>
      <!-- Form bikin dompet -->
      <section
        class="mb-8 border-[3px] border-ink bg-brut-yellow shadow-brut rounded-brut p-4 sm:p-6"
      >
        <div class="flex items-center justify-between gap-3 mb-4">
          <h3 class="text-lg sm:text-xl font-black uppercase tracking-wide text-ink">
            Bikin Dompet Baru
          </h3>
          <span
            class="grid shrink-0 w-8 h-8 place-items-center border-2 border-ink bg-surface shadow-brut-sm rounded-brut-sm"
          >
            <PhPlus weight="bold" class="w-4 h-4" />
          </span>
        </div>

        <form @submit.prevent="submitWallet" class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div class="sm:col-span-1">
            <label class="block mb-2 text-sm font-black uppercase text-ink"
              >Nama Dompet (BCA, Gopay, dll)</label
            >
            <input
              v-model="form.name"
              type="text"
              required
              placeholder="Contoh: Gopay"
              class="w-full px-3 py-2 text-ink bg-surface border-2 border-ink rounded-brut outline-none focus:border-ink focus:ring-2 focus:ring-ink"
            />
          </div>

          <div class="sm:col-span-1">
            <label class="block mb-2 text-sm font-black uppercase text-ink">Saldo Awal (Rp)</label>
            <input
              :value="form.balance"
              @input="form.balance = formatThousands($event.target.value)"
              type="text"
              inputmode="numeric"
              required
              placeholder="0"
              class="w-full px-3 py-2 text-ink bg-surface border-2 border-ink rounded-brut outline-none focus:border-ink focus:ring-2 focus:ring-ink"
            />
          </div>

          <div class="sm:col-span-2">
            <button
              type="submit"
              class="flex items-center justify-center w-full py-3 gap-2 text-ink border-[3px] border-ink bg-brut-cyan rounded-brut shadow-brut font-black cursor-pointer transition-all hover:translate-y-[-2px] hover:translate-x-[-2px]"
            >
              <PhFloppyDisk weight="bold" class="w-5 h-5" />
              Simpan Dompet
            </button>
          </div>
        </form>
      </section>

      <!-- Daftar dompet -->
      <h3 class="text-lg sm:text-xl font-black uppercase tracking-wide text-ink mb-4">
        Dompet Aktif
      </h3>

      <p
        v-if="wallets.length === 0"
        class="py-8 border-2 border-dashed border-ink bg-surface rounded-brut text-center font-black text-ink"
      >
        Belum ada dompet, bikin dulu yuk!
      </p>

      <div v-else class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-6">
        <!-- Kartu ringkasan total saldo -->
        <article
          class="p-5 sm:p-6 text-ink bg-brut-cyan border-[3px] border-ink shadow-brut rounded-brut transition-all hover:translate-y-[-2px] sm:col-span-2 flex flex-col"
        >
          <div class="flex items-start justify-between gap-3">
            <div>
              <h3 class="text-xs sm:text-sm font-black uppercase tracking-wide">Total Saldo</h3>
              <p class="mt-0.5 text-[11px] sm:text-xs font-bold text-ink/70">
                Akumulasi dari {{ wallets.length }} dompet
              </p>
            </div>
            <span
              class="grid shrink-0 w-8 h-8 place-items-center border-2 border-ink bg-surface shadow-brut-sm rounded-brut-sm"
            >
              <PhPiggyBank weight="fill" class="w-4 h-4" />
            </span>
          </div>

          <p class="mt-3 text-2xl sm:text-3xl font-black break-words">
            {{ formatRupiah(totalSaldo) }}
          </p>
        </article>

        <!-- Kartu dompet -->
        <article
          v-for="(wallet, index) in wallets"
          :key="wallet.id"
          class="p-4 sm:p-6 border-[3px] border-ink shadow-brut rounded-brut transition-all hover:translate-y-[-2px] flex flex-col"
          :class="walletColor(index)"
        >
          <div class="flex items-start justify-between gap-3">
            <div class="min-w-0">
              <h3 class="text-xs sm:text-sm font-black uppercase tracking-wide break-words">
                {{ wallet.name }}
              </h3>
              <p class="mt-0.5 text-[11px] sm:text-xs font-bold text-ink/70">Saldo saat ini</p>
            </div>
            <span
              class="grid shrink-0 w-8 h-8 place-items-center border-2 border-ink bg-surface shadow-brut-sm rounded-brut-sm"
            >
              <PhWallet weight="fill" class="w-4 h-4" />
            </span>
          </div>

          <p class="mt-3 text-xl sm:text-2xl font-black break-words">
            {{ formatRupiah(wallet.balance) }}
          </p>

          <p class="mt-3 text-[10px] font-black uppercase text-ink/60 break-all">
            ID: {{ wallet.id }}
          </p>
        </article>
      </div>
    </template>

    <Toast ref="toastRef" />
    <CodeGateModal v-if="showLoginGate" @submit="submitCode" />
  </main>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import Toast from '../components/Toast.vue'
import { formatRupiah, formatThousands, toDigits } from '../utils/format.js'
import { getSavedCode, saveCode, clearCode, authHeaders } from '../utils/userCode.js'
import CodeGateModal from '../components/CodeGateModal.vue'
import { PhWallet, PhPiggyBank, PhPlus, PhFloppyDisk } from '@phosphor-icons/vue'

// State
const wallets = ref([])
const isLoading = ref(true)
const toastRef = ref(null)

const form = ref({
  name: '',
  balance: '',
})

const userCode = ref('')
const showLoginGate = ref(false)

const lockVault = () => {
  clearCode()
  userCode.value = ''
  wallets.value = []
  showLoginGate.value = true
}

const isUnauthorized = (res) => {
  if (res.status !== 401) return false
  lockVault()
  toastRef.value?.show('Kode tidak valid, masukkan ulang.', 'error')
  return true
}

const submitCode = (code) => {
  saveCode(code)
  userCode.value = code
  showLoginGate.value = false
  fetchWallets()
}

// Pakai path relatif, di-proxy Vite ke backend (lihat vite.config.js).
const API_URL = '/api'

// Palet warna brutal buat kartu dompet, biar konsisten sama dashboard.
const WALLET_COLORS = ['bg-brut-pink', 'bg-brut-lime', 'bg-brut-blue', 'bg-brut-yellow']

const walletColor = (index) => WALLET_COLORS[index % WALLET_COLORS.length]

const totalSaldo = computed(() =>
  wallets.value.reduce((sum, wallet) => sum + (Number(wallet.balance) || 0), 0),
)

// Fetch data dompet dari Golang
const fetchWallets = async () => {
  if (!userCode.value) return
  try {
    const response = await fetch(`${API_URL}/wallets`, { headers: authHeaders(userCode.value) })
    if (isUnauthorized(response)) return
    if (!response.ok) throw new Error(`HTTP ${response.status}`)

    const data = await response.json()
    wallets.value = Array.isArray(data) ? data : []
  } catch (error) {
    console.error('Gagal ambil data dompet:', error)
    toastRef.value?.show('Pastikan server Golang sudah nyala!', 'error')
  } finally {
    isLoading.value = false
  }
}

// Kirim data dompet baru ke Golang
const submitWallet = async () => {
  const payload = {
    name: form.value.name.trim(),
    balance: parseInt(toDigits(form.value.balance) || '0', 10),
  }

  try {
    const response = await fetch(`${API_URL}/wallets`, {
      method: 'POST',
      headers: authHeaders(userCode.value),
      body: JSON.stringify(payload),
    })

    if (isUnauthorized(response)) return
    if (!response.ok) {
      toastRef.value?.show('Gagal menyimpan dompet.', 'error')
      return
    }

    // Reset form
    form.value.name = ''
    form.value.balance = ''

    // Ambil ulang data terbaru dari backend
    await fetchWallets()
    toastRef.value?.show('Dompet berhasil disimpan!', 'success')
  } catch (error) {
    console.error('Error nyimpen dompet:', error)
    toastRef.value?.show('Pastikan server Golang sudah nyala!', 'error')
  }
}

// Jalankan fetch saat komponen pertama kali dimuat
onMounted(() => {
  const savedCode = getSavedCode()
  if (savedCode) {
    userCode.value = savedCode
    fetchWallets()
  } else {
    showLoginGate.value = true
    isLoading.value = false
  }
})
</script>
