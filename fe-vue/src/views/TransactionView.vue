<template>
  <main class="max-w-4xl mx-auto px-4 sm:px-8 py-8 sm:py-10">
    <div class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-4 mb-8">
      <div>
        <h2 class="flex items-center gap-2 sm:gap-3 mb-2 text-2xl sm:text-4xl font-black text-ink">
          <PhNote weight="fill" class="w-7 h-7 sm:w-9 sm:h-9" />
          Daftar Transaksi
        </h2>
        <p class="text-sm sm:text-base text-ink">Catat hutang atau pengeluaran barumu di sini.</p>
      </div>
      <button
        @click="showModal = true"
        class="flex items-center justify-center gap-2 w-full sm:w-auto px-4 py-2 text-ink border-[3px] border-ink bg-brut-cyan rounded-none shadow-brut font-black cursor-pointer transition-all hover:translate-y-[-2px] hover:translate-x-[-2px]"
      >
        <PhPlus weight="bold" class="w-5 h-5" />
        Tambah Transaksi
      </button>
    </div>

    <template v-if="transactions?.length > 0">
      <ul class="sm:hidden space-y-3">
        <li
          v-for="trx in transactions"
          :key="trx.id"
          class="p-4 border-2 border-ink shadow-brut-sm"
          :class="trx.status === 'PAID' ? 'bg-brut-lime' : 'bg-surface'"
        >
          <div class="flex items-start justify-between gap-3">
            <h3 class="font-black text-base leading-snug break-words">{{ trx.title }}</h3>
            <span
              class="shrink-0 px-2 py-1 text-[11px] font-black border-2 border-ink"
              :class="trx.status === 'PAID' ? 'bg-brut-cyan' : 'bg-brut-yellow'"
            >
              {{ trx.status }}
            </span>
          </div>

          <div class="mt-3 flex items-center justify-between gap-3 text-sm text-ink">
            <span class="flex items-center gap-1.5">
              <PhCalendarBlank weight="bold" class="w-4 h-4" />
              {{ trx.due_date }}
            </span>
            <span class="font-black text-base">{{ formatRupiah(trx.amount) }}</span>
          </div>

          <div class="mt-4 flex gap-2">
            <button
              v-if="trx.status === 'PENDING'"
              @click="markAsPaid(trx)"
              class="flex flex-1 items-center justify-center gap-1.5 py-2 text-ink border-2 border-ink bg-brut-cyan font-black cursor-pointer shadow-brut-sm transition-all hover:translate-y-[-1px]"
            >
              <PhCheckCircle weight="fill" class="w-4 h-4" />
              Lunas
            </button>
            <button
              @click="deleteGroup(trx.group_id)"
              class="flex flex-1 items-center justify-center gap-1.5 py-2 text-white border-2 border-ink bg-brut-red font-black cursor-pointer shadow-brut-sm transition-all hover:translate-y-[-1px]"
            >
              <PhTrashSimple weight="bold" class="w-4 h-4" />
              Hapus Grup
            </button>
          </div>
        </li>
      </ul>

      <div class="hidden sm:block overflow-x-auto border-2 border-ink bg-surface shadow-brut-sm">
        <table class="w-full border-collapse">
          <thead>
            <tr class="bg-ink text-paper">
              <th
                class="p-2 sm:p-4 text-left text-[10px] sm:text-xs font-black border-2 border-ink"
              >
                Nama
              </th>
              <th
                class="p-2 sm:p-4 text-left text-[10px] sm:text-xs font-black border-2 border-ink"
              >
                Jatuh Tempo
              </th>
              <th
                class="p-2 sm:p-4 text-left text-[10px] sm:text-xs font-black border-2 border-ink"
              >
                Nominal
              </th>
              <th
                class="p-2 sm:p-4 text-left text-[10px] sm:text-xs font-black border-2 border-ink"
              >
                Status
              </th>
              <th
                class="p-2 sm:p-4 text-left text-[10px] sm:text-xs font-black border-2 border-ink"
              >
                Aksi
              </th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="trx in transactions"
              :key="trx.id"
              :class="{ 'bg-brut-lime': trx.status === 'PAID' }"
            >
              <td class="p-2 sm:p-4 text-xs sm:text-sm border-2 border-ink">{{ trx.title }}</td>
              <td class="p-2 sm:p-4 text-xs sm:text-sm border-2 border-ink">{{ trx.due_date }}</td>
              <td class="p-2 sm:p-4 text-xs sm:text-sm border-2 border-ink">
                {{ formatRupiah(trx.amount) }}
              </td>
              <td class="p-2 sm:p-4 text-xs sm:text-sm border-2 border-ink">
                <span
                  :class="
                    trx.status === 'PAID'
                      ? 'bg-brut-cyan text-ink font-black border-2 border-ink px-2 py-1 text-xs'
                      : 'bg-brut-yellow text-ink font-black border-2 border-ink px-2 py-1 text-xs'
                  "
                  class="px-2 py-1 text-xs font-black border-2 border-ink"
                >
                  {{ trx.status }}
                </span>
              </td>
              <td class="p-2 sm:p-4 text-xs sm:text-sm border-2 border-ink">
                <div class="flex flex-wrap gap-2">
                  <button
                    v-if="trx.status === 'PENDING'"
                    @click="markAsPaid(trx)"
                    class="flex items-center gap-1 px-2 py-1 text-ink border-2 border-ink bg-brut-cyan font-black cursor-pointer shadow-brut hover:translate-y-[-1px]"
                  >
                    <PhCheckCircle weight="fill" class="w-4 h-4" />
                    Lunas
                  </button>
                  <button
                    @click="deleteGroup(trx.group_id)"
                    class="flex items-center gap-1 px-2 py-1 text-white border-2 border-ink bg-brut-red font-black cursor-pointer shadow-brut hover:translate-y-[-1px]"
                  >
                    <PhTrashSimple weight="bold" class="w-4 h-4" />
                    Hapus Grup
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </template>

    <p v-else class="py-8 text-center font-bold text-ink">Belum ada transaksi, yuk catat dulu!</p>

    <div
      v-if="showModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-ink/40 backdrop-blur-sm"
      @click.self="showModal = false"
    >
      <div
        class="p-4 sm:p-6 bg-surface border-[3px] border-ink shadow-brut-lg max-w-2xl w-[calc(100%-2rem)] sm:w-full mx-auto sm:mx-4"
      >
        <div class="flex justify-between items-center mb-6">
          <h3 class="text-xl font-black">Tambah Transaksi Baru</h3>
          <button
            @click="showModal = false"
            class="flex items-center justify-center w-8 h-8 text-ink bg-brut-yellow border-2 border-ink rounded-none font-black cursor-pointer shadow-brut hover:translate-y-[-1px]"
          >
            <PhX weight="bold" class="w-4 h-4" />
          </button>
        </div>

        <form @submit.prevent="submitTransaction" class="space-y-4">
          <div>
            <label class="block mb-2 text-sm font-black uppercase text-ink"
              >Nama Transaksi / Hutang</label
            >
            <input
              v-model="form.title"
              type="text"
              required
              placeholder="Contoh: Paylater Kredivo"
              class="w-full px-3 py-2 text-ink bg-surface border-2 border-ink rounded-none outline-none focus:border-ink focus:ring-2 focus:ring-ink"
            />
          </div>

          <div class="flex flex-col sm:flex-row gap-4">
            <div class="flex-1">
              <label class="block mb-2 text-sm font-black uppercase text-ink">Nominal (Rp)</label>
              <input
                v-model="form.amount"
                type="number"
                required
                class="w-full px-3 py-2 text-ink bg-surface border-2 border-ink rounded-none outline-none focus:border-ink focus:ring-2 focus:ring-ink"
              />
            </div>
            <div class="flex-1">
              <label class="block mb-2 text-sm font-black uppercase text-ink">Jenis</label>
              <select
                v-model="form.type"
                class="w-full px-3 py-2 text-ink bg-surface border-2 border-ink rounded-none outline-none focus:border-ink focus:ring-2 focus:ring-ink"
              >
                <option value="HUTANG">Hutang</option>
                <option value="PENGELUARAN">Pengeluaran Rutin</option>
              </select>
            </div>
          </div>

          <div class="flex flex-col sm:flex-row gap-4">
            <div class="flex-1">
              <label class="block mb-2 text-sm font-black uppercase text-ink"
                >Tanggal Mulai / Jatuh Tempo</label
              >
              <input
                v-model="form.due_date"
                type="date"
                required
                class="w-full px-3 py-2 text-ink bg-surface border-2 border-ink rounded-none outline-none focus:border-ink focus:ring-2 focus:ring-ink"
              />
            </div>
            <div class="flex-1">
              <label class="block mb-2 text-sm font-black uppercase text-ink"
                >Tenor (Berapa Kali?)</label
              >
              <input
                v-model="form.tenor"
                type="number"
                min="1"
                required
                class="w-full px-3 py-2 text-ink bg-surface border-2 border-ink rounded-none outline-none focus:border-ink focus:ring-2 focus:ring-ink"
              />
            </div>
          </div>

          <button
            type="submit"
            class="flex items-center justify-center w-full py-3 gap-2 text-ink border-[3px] border-ink bg-brut-cyan rounded-none shadow-brut font-black cursor-pointer transition-all hover:translate-y-[-2px] hover:translate-x-[-2px]"
          >
            <PhFloppyDisk weight="bold" class="w-5 h-5" />
            Simpan Transaksi
          </button>
        </form>
      </div>
    </div>
  </main>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  PhNote,
  PhFloppyDisk,
  PhCheckCircle,
  PhTrashSimple,
  PhPlus,
  PhX,
  PhCalendarBlank,
} from '@phosphor-icons/vue'

const transactions = ref([])
const showModal = ref(false)

const form = ref({
  title: '',
  amount: 0,
  type: 'HUTANG',
  due_date: '',
  tenor: 1,
})

const fetchTransactions = async () => {
  try {
    const res = await fetch('/api/transactions')
    const data = await res.json()
    transactions.value = data || []
  } catch (error) {
    console.error('Gagal mengambil data:', error)
    transactions.value = []
  }
}

const submitTransaction = async () => {
  await fetch('/api/transactions', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      title: form.value.title,
      amount: parseInt(form.value.amount),
      type: form.value.type,
      due_date: form.value.due_date,
      tenor: parseInt(form.value.tenor),
    }),
  })

  form.value.title = ''
  form.value.amount = 0
  form.value.tenor = 1
  showModal.value = false
  fetchTransactions()
  alert('Transaksi berhasil disimpan!')
}

const markAsPaid = async (trx) => {
  if (!confirm(`Tandai "${trx.title}" lunas?`)) return

  await fetch(`/api/transactions/${trx.id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      title: trx.title,
      amount: trx.amount,
      type: trx.type,
      due_date: trx.due_date,
      status: 'PAID',
    }),
  })
  fetchTransactions()
}

const deleteGroup = async (groupId) => {
  if (!confirm('Yakin mau hapus SEMUA cicilan di grup ini?')) return

  await fetch(`/api/transactions/group/${groupId}`, {
    method: 'DELETE',
  })
  fetchTransactions()
}

const formatRupiah = (angka) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(angka)
}

onMounted(() => {
  fetchTransactions()
})
</script>
