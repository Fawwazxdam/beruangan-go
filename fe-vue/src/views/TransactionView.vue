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
          v-for="grp in groupedTransactions"
          :key="grp.groupId"
          class="border-2 border-ink shadow-brut-sm overflow-hidden"
          :class="grp.allPaid ? 'bg-brut-lime' : 'bg-surface'"
        >
          <div
            class="p-4"
            :class="{ 'cursor-pointer': !grp.single }"
            @click="!grp.single && toggleGroup(grp.groupId)"
          >
            <div class="flex items-start justify-between gap-3">
              <h3 class="flex items-center gap-2 font-black text-base leading-snug break-words">
                <PhCaretDown
                  v-if="!grp.single"
                  weight="bold"
                  class="w-4 h-4 shrink-0 transition-transform"
                  :class="isExpanded(grp.groupId) ? 'rotate-180' : ''"
                />
                {{ grp.baseTitle }}
              </h3>
              <span
                class="shrink-0 px-2 py-1 text-[11px] font-black border-2 border-ink"
                :class="grp.allPaid ? 'bg-brut-cyan' : 'bg-brut-yellow'"
              >
                {{ statusLabel(grp) }}
              </span>
            </div>

            <div class="mt-3 flex items-center justify-between gap-3 text-sm text-ink">
              <span class="flex items-center gap-1.5">
                <PhCalendarBlank weight="bold" class="w-4 h-4" />
                {{ formatDate(grp.nextDue) }}
              </span>
              <span class="font-black text-base">{{ formatRupiah(grp.totalAmount) }}</span>
            </div>

            <p v-if="!grp.single" class="mt-1 text-xs font-black text-ink">
              {{ grp.items.length }} cicilan
            </p>

            <div class="mt-4 flex gap-2" @click.stop>
              <button
                v-if="grp.single && grp.items[0].status === 'PENDING'"
                @click="markAsPaid(grp.items[0])"
                class="flex flex-1 items-center justify-center gap-1.5 py-2 text-ink border-2 border-ink bg-brut-cyan font-black cursor-pointer shadow-brut-sm transition-all hover:translate-y-[-1px]"
              >
                <PhCheckCircle weight="fill" class="w-4 h-4" />
                Lunas
              </button>
              <button
                v-if="grp.single"
                @click="openEdit(grp.items[0])"
                class="flex flex-1 items-center justify-center gap-1.5 py-2 text-ink border-2 border-ink bg-brut-yellow font-black cursor-pointer shadow-brut-sm transition-all hover:translate-y-[-1px]"
              >
                <PhPencilSimple weight="bold" class="w-4 h-4" />
                Edit
              </button>
              <button
                v-if="grp.single"
                @click="deleteTransaction(grp.items[0].id, grp.items[0].title)"
                class="flex flex-1 items-center justify-center gap-1.5 py-2 text-white border-2 border-ink bg-brut-red font-black cursor-pointer shadow-brut-sm transition-all hover:translate-y-[-1px]"
              >
                <PhTrashSimple weight="bold" class="w-4 h-4" />
                Hapus
              </button>
              <button
                v-else
                @click="deleteGroup(grp.groupId)"
                class="flex flex-1 items-center justify-center gap-1.5 py-2 text-white border-2 border-ink bg-brut-red font-black cursor-pointer shadow-brut-sm transition-all hover:translate-y-[-1px]"
              >
                <PhTrashSimple weight="bold" class="w-4 h-4" />
                Hapus Grup
              </button>
            </div>
          </div>

          <div
            v-if="!grp.single && isExpanded(grp.groupId)"
            class="border-t-2 border-ink bg-surface p-4 space-y-3"
          >
            <div
              v-for="item in grp.items"
              :key="item.id"
              class="flex items-center justify-between gap-3"
            >
              <div class="min-w-0">
                <p class="text-sm font-black leading-snug break-words">{{ item.title }}</p>
                <p class="text-xs text-ink">
                  {{ formatDate(item.due_date) }} · {{ formatRupiah(item.amount) }}
                </p>
              </div>
              <div class="flex items-center gap-2 shrink-0">
                <span
                  class="px-2 py-1 text-[11px] font-black border-2 border-ink"
                  :class="item.status === 'PAID' ? 'bg-brut-cyan' : 'bg-brut-yellow'"
                >
                  {{ item.status }}
                </span>
                <button
                  v-if="item.status === 'PENDING'"
                  @click="markAsPaid(item)"
                  class="flex items-center gap-1 px-2 py-1 text-ink border-2 border-ink bg-brut-cyan font-black cursor-pointer shadow-brut-sm transition-all hover:translate-y-[-1px]"
                >
                  <PhCheckCircle weight="fill" class="w-4 h-4" />
                  Lunas
                </button>
                <button
                  @click="openEdit(item)"
                  title="Edit transaksi"
                  class="flex items-center justify-center w-8 h-8 text-ink border-2 border-ink bg-brut-yellow font-black cursor-pointer shadow-brut-sm transition-all hover:translate-y-[-1px]"
                >
                  <PhPencilSimple weight="bold" class="w-4 h-4" />
                </button>
                <button
                  @click="deleteTransaction(item.id, item.title)"
                  title="Hapus transaksi"
                  class="flex items-center justify-center w-8 h-8 text-white border-2 border-ink bg-brut-red font-black cursor-pointer shadow-brut-sm transition-all hover:translate-y-[-1px]"
                >
                  <PhTrashSimple weight="bold" class="w-4 h-4" />
                </button>
              </div>
            </div>
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
            <template v-for="grp in groupedTransactions" :key="grp.groupId">
              <tr
                class="cursor-pointer"
                :class="{ 'bg-brut-lime': grp.allPaid }"
                @click="!grp.single && toggleGroup(grp.groupId)"
              >
                <td class="p-2 sm:p-4 text-xs sm:text-sm border-2 border-ink">
                  <span class="flex items-center gap-2 font-bold">
                    <PhCaretDown
                      v-if="!grp.single"
                      weight="bold"
                      class="w-4 h-4 shrink-0 transition-transform"
                      :class="isExpanded(grp.groupId) ? 'rotate-180' : ''"
                    />
                    {{ grp.baseTitle }}
                    <span
                      v-if="!grp.single"
                      class="px-1.5 py-0.5 text-[10px] font-black border-2 border-ink bg-brut-cyan"
                    >
                      {{ grp.items.length }} cicilan
                    </span>
                  </span>
                </td>
                <td class="p-2 sm:p-4 text-xs sm:text-sm border-2 border-ink">
                  {{ formatDate(grp.nextDue) }}
                </td>
                <td class="p-2 sm:p-4 text-xs sm:text-sm border-2 border-ink">
                  {{ formatRupiah(grp.totalAmount) }}
                </td>
                <td class="p-2 sm:p-4 text-xs sm:text-sm border-2 border-ink">
                  <span
                    class="px-2 py-1 text-xs font-black border-2 border-ink"
                    :class="grp.allPaid ? 'bg-brut-cyan' : 'bg-brut-yellow'"
                  >
                    {{ statusLabel(grp) }}
                  </span>
                </td>
                <td class="p-2 sm:p-4 text-xs sm:text-sm border-2 border-ink">
                  <div class="flex flex-wrap gap-2" @click.stop>
                    <button
                      v-if="grp.single && grp.items[0].status === 'PENDING'"
                      @click="markAsPaid(grp.items[0])"
                      class="flex items-center gap-1 px-2 py-1 text-ink border-2 border-ink bg-brut-cyan font-black cursor-pointer shadow-brut hover:translate-y-[-1px]"
                    >
                      <PhCheckCircle weight="fill" class="w-4 h-4" />
                      Lunas
                    </button>
                    <button
                      v-if="grp.single"
                      @click="openEdit(grp.items[0])"
                      class="flex items-center gap-1 px-2 py-1 text-ink border-2 border-ink bg-brut-yellow font-black cursor-pointer shadow-brut hover:translate-y-[-1px]"
                    >
                      <PhPencilSimple weight="bold" class="w-4 h-4" />
                      Edit
                    </button>
                    <button
                      v-if="grp.single"
                      @click="deleteTransaction(grp.items[0].id, grp.items[0].title)"
                      class="flex items-center gap-1 px-2 py-1 text-white border-2 border-ink bg-brut-red font-black cursor-pointer shadow-brut hover:translate-y-[-1px]"
                    >
                      <PhTrashSimple weight="bold" class="w-4 h-4" />
                      Hapus
                    </button>
                    <button
                      v-else
                      @click="deleteGroup(grp.groupId)"
                      class="flex items-center gap-1 px-2 py-1 text-white border-2 border-ink bg-brut-red font-black cursor-pointer shadow-brut hover:translate-y-[-1px]"
                    >
                      <PhTrashSimple weight="bold" class="w-4 h-4" />
                      Hapus Grup
                    </button>
                  </div>
                </td>
              </tr>

              <tr v-if="!grp.single && isExpanded(grp.groupId)">
                <td colspan="5" class="p-0 border-x-2 border-b-2 border-ink">
                  <table class="w-full border-collapse">
                    <tbody>
                      <tr
                        v-for="item in grp.items"
                        :key="item.id"
                        :class="{ 'bg-brut-lime': item.status === 'PAID' }"
                      >
                        <td
                          class="p-2 sm:p-3 pl-8 sm:pl-12 text-xs sm:text-sm border-t-2 border-ink"
                        >
                          {{ item.title }}
                        </td>
                        <td class="p-2 sm:p-3 text-xs sm:text-sm border-t-2 border-ink">
                          {{ formatDate(item.due_date) }}
                        </td>
                        <td class="p-2 sm:p-3 text-xs sm:text-sm border-t-2 border-ink">
                          {{ formatRupiah(item.amount) }}
                        </td>
                        <td class="p-2 sm:p-3 text-xs sm:text-sm border-t-2 border-ink">
                          <span
                            class="px-2 py-1 text-xs font-black border-2 border-ink"
                            :class="item.status === 'PAID' ? 'bg-brut-cyan' : 'bg-brut-yellow'"
                          >
                            {{ item.status }}
                          </span>
                        </td>
                        <td class="p-2 sm:p-3 text-xs sm:text-sm border-t-2 border-ink">
                          <div class="flex flex-wrap gap-2">
                            <button
                              v-if="item.status === 'PENDING'"
                              @click="markAsPaid(item)"
                              class="flex items-center gap-1 px-2 py-1 text-ink border-2 border-ink bg-brut-cyan font-black cursor-pointer shadow-brut hover:translate-y-[-1px]"
                            >
                              <PhCheckCircle weight="fill" class="w-4 h-4" />
                              Lunas
                            </button>
                            <button
                              @click="openEdit(item)"
                              class="flex items-center gap-1 px-2 py-1 text-ink border-2 border-ink bg-brut-yellow font-black cursor-pointer shadow-brut hover:translate-y-[-1px]"
                            >
                              <PhPencilSimple weight="bold" class="w-4 h-4" />
                              Edit
                            </button>
                            <button
                              @click="deleteTransaction(item.id, item.title)"
                              class="flex items-center gap-1 px-2 py-1 text-white border-2 border-ink bg-brut-red font-black cursor-pointer shadow-brut hover:translate-y-[-1px]"
                            >
                              <PhTrashSimple weight="bold" class="w-4 h-4" />
                              Hapus
                            </button>
                          </div>
                        </td>
                      </tr>
                    </tbody>
                  </table>
                </td>
              </tr>
            </template>
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
                :value="form.amount"
                @input="form.amount = formatThousands($event.target.value)"
                type="text"
                inputmode="numeric"
                required
                placeholder="Contoh: 1.500.000"
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

    <div
      v-if="showEditModal"
      class="fixed inset-0 z-50 flex items-center justify-center bg-ink/40 backdrop-blur-sm"
      @click.self="showEditModal = false"
    >
      <div
        class="p-4 sm:p-6 bg-surface border-[3px] border-ink shadow-brut-lg max-w-2xl w-[calc(100%-2rem)] sm:w-full mx-auto sm:mx-4"
      >
        <div class="flex justify-between items-center mb-6">
          <h3 class="text-xl font-black">Edit Transaksi</h3>
          <button
            @click="showEditModal = false"
            class="flex items-center justify-center w-8 h-8 text-ink bg-brut-yellow border-2 border-ink rounded-none font-black cursor-pointer shadow-brut hover:translate-y-[-1px]"
          >
            <PhX weight="bold" class="w-4 h-4" />
          </button>
        </div>

        <form @submit.prevent="submitEdit" class="space-y-4">
          <div>
            <label class="block mb-2 text-sm font-black uppercase text-ink"
              >Nama Transaksi / Hutang</label
            >
            <input
              v-model="editForm.title"
              type="text"
              required
              class="w-full px-3 py-2 text-ink bg-surface border-2 border-ink rounded-none outline-none focus:border-ink focus:ring-2 focus:ring-ink"
            />
          </div>

          <div class="flex flex-col sm:flex-row gap-4">
            <div class="flex-1">
              <label class="block mb-2 text-sm font-black uppercase text-ink">Nominal (Rp)</label>
              <input
                :value="editForm.amount"
                @input="editForm.amount = formatThousands($event.target.value)"
                type="text"
                inputmode="numeric"
                required
                placeholder="Contoh: 1.500.000"
                class="w-full px-3 py-2 text-ink bg-surface border-2 border-ink rounded-none outline-none focus:border-ink focus:ring-2 focus:ring-ink"
              />
            </div>
            <div class="flex-1">
              <label class="block mb-2 text-sm font-black uppercase text-ink">Jenis</label>
              <select
                v-model="editForm.type"
                class="w-full px-3 py-2 text-ink bg-surface border-2 border-ink rounded-none outline-none focus:border-ink focus:ring-2 focus:ring-ink"
              >
                <option value="HUTANG">Hutang</option>
                <option value="PENGELUARAN">Pengeluaran Rutin</option>
              </select>
            </div>
          </div>

          <div class="flex flex-col sm:flex-row gap-4">
            <div class="flex-1">
              <label class="block mb-2 text-sm font-black uppercase text-ink">Jatuh Tempo</label>
              <input
                v-model="editForm.due_date"
                type="date"
                required
                class="w-full px-3 py-2 text-ink bg-surface border-2 border-ink rounded-none outline-none focus:border-ink focus:ring-2 focus:ring-ink"
              />
            </div>
            <div class="flex-1">
              <label class="block mb-2 text-sm font-black uppercase text-ink">Status</label>
              <select
                v-model="editForm.status"
                class="w-full px-3 py-2 text-ink bg-surface border-2 border-ink rounded-none outline-none focus:border-ink focus:ring-2 focus:ring-ink"
              >
                <option value="PENDING">Belum Lunas</option>
                <option value="PAID">Lunas</option>
              </select>
            </div>
          </div>

          <button
            type="submit"
            class="flex items-center justify-center w-full py-3 gap-2 text-ink border-[3px] border-ink bg-brut-yellow rounded-none shadow-brut font-black cursor-pointer transition-all hover:translate-y-[-2px] hover:translate-x-[-2px]"
          >
            <PhFloppyDisk weight="bold" class="w-5 h-5" />
            Simpan Perubahan
          </button>
        </form>
      </div>
    </div>

    <Toast ref="toastRef" />
    <ConfirmModal ref="confirmRef" />
  </main>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import Toast from '../components/Toast.vue'
import ConfirmModal from '../components/ConfirmModal.vue'
import { formatRupiah, formatDate, formatThousands, toDigits } from '../utils/format.js'
import {
  PhNote,
  PhFloppyDisk,
  PhCheckCircle,
  PhTrashSimple,
  PhPlus,
  PhX,
  PhCalendarBlank,
  PhCaretDown,
  PhPencilSimple,
} from '@phosphor-icons/vue'

const transactions = ref([])
const showModal = ref(false)
const showEditModal = ref(false)
const expandedGroups = ref([])
const toastRef = ref(null)
const confirmRef = ref(null)

const form = ref({
  title: '',
  amount: '',
  type: 'HUTANG',
  due_date: '',
  tenor: 1,
})

const editForm = ref({
  id: null,
  title: '',
  amount: '',
  type: 'HUTANG',
  due_date: '',
  status: 'PENDING',
})

const fetchTransactions = async () => {
  try {
    const res = await fetch('/api/transactions')
    const data = await res.json()

    if (data) {
      // PROSES SORTING: Tanggal terdekat / paling awal ada di atas
      data.sort((a, b) => new Date(a.due_date) - new Date(b.due_date))

      transactions.value = data
    } else {
      transactions.value = []
    }
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
      amount: parseInt(toDigits(form.value.amount), 10),
      type: form.value.type,
      due_date: form.value.due_date,
      tenor: parseInt(form.value.tenor),
    }),
  })

  form.value.title = ''
  form.value.amount = ''
  form.value.tenor = 1
  showModal.value = false
  fetchTransactions()
  toastRef.value?.show('Transaksi berhasil disimpan!', 'success')
}

const groupedTransactions = computed(() => {
  const map = new Map()

  for (const trx of transactions.value) {
    const key = trx.group_id || `single-${trx.id}`
    if (!map.has(key)) map.set(key, [])
    map.get(key).push(trx)
  }

  return Array.from(map, ([groupId, items]) => {
    items.sort((a, b) => (a.due_date || '').localeCompare(b.due_date || ''))

    const paidCount = items.filter((item) => item.status === 'PAID').length
    const nextPending = items.find((item) => item.status === 'PENDING')

    return {
      groupId,
      items,
      single: items.length === 1,
      baseTitle: items[0].title.replace(/\s*\(\d+\/\d+\)\s*$/, ''),
      totalAmount: items.reduce((sum, item) => sum + item.amount, 0),
      nextDue: nextPending ? nextPending.due_date : items[items.length - 1].due_date,
      paidCount,
      allPaid: paidCount === items.length,
    }
  })
})

const toggleGroup = (groupId) => {
  const index = expandedGroups.value.indexOf(groupId)
  if (index === -1) expandedGroups.value.push(groupId)
  else expandedGroups.value.splice(index, 1)
}

const isExpanded = (groupId) => expandedGroups.value.includes(groupId)

const statusLabel = (grp) => {
  if (grp.single) return grp.items[0].status
  if (grp.allPaid) return 'LUNAS'
  return `${grp.paidCount}/${grp.items.length}`
}

const markAsPaid = async (trx) => {
  const ok = await confirmRef.value?.open({
    title: 'Tandai Lunas?',
    message: `Tandai "${trx.title}" lunas?`,
    confirmText: 'Ya, Lunas',
  })
  if (!ok) return

  const res = await fetch(`/api/transactions/${trx.id}`, {
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
  if (!res.ok) {
    toastRef.value?.show('Gagal menandai transaksi lunas.', 'error')
    return
  }
  fetchTransactions()
  toastRef.value?.show(`"${trx.title}" berhasil ditandai lunas!`, 'success')
}

const openEdit = (trx) => {
  editForm.value = {
    id: trx.id,
    title: trx.title,
    amount: formatThousands(trx.amount),
    type: trx.type,
    due_date: trx.due_date,
    status: trx.status || 'PENDING',
  }
  showEditModal.value = true
}

const submitEdit = async () => {
  await fetch(`/api/transactions/${editForm.value.id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      title: editForm.value.title,
      amount: parseInt(toDigits(editForm.value.amount), 10),
      type: editForm.value.type,
      due_date: editForm.value.due_date,
      status: editForm.value.status,
    }),
  })

  showEditModal.value = false
  fetchTransactions()
  toastRef.value?.show('Transaksi berhasil diupdate!', 'success')
}

const deleteGroup = async (groupId) => {
  const ok = await confirmRef.value?.open({
    title: 'Hapus Grup Cicilan?',
    message: 'Yakin mau hapus SEMUA cicilan di grup ini?',
    confirmText: 'Hapus Semua',
    danger: true,
  })
  if (!ok) return

  const res = await fetch(`/api/transactions/group/${groupId}`, {
    method: 'DELETE',
  })
  const data = await res.json().catch(() => null)
  if (!res.ok) {
    toastRef.value?.show(data?.message || 'Gagal menghapus grup cicilan.', 'error')
    return
  }
  fetchTransactions()
  toastRef.value?.show(data?.message || 'Grup cicilan berhasil dihapus!', 'success')
}

const deleteTransaction = async (id, title) => {
  const ok = await confirmRef.value?.open({
    title: 'Hapus Transaksi?',
    message: `Yakin mau hapus transaksi "${title}" ini aja?`,
    confirmText: 'Hapus',
    danger: true,
  })
  if (!ok) return

  try {
    const res = await fetch(`/api/transactions/${id}`, {
      method: 'DELETE',
    })
    const data = await res.json().catch(() => null)
    if (!res.ok) throw new Error(data?.message || `HTTP ${res.status}`)
    fetchTransactions() // Refresh data setelah hapus
    toastRef.value?.show(data?.message || 'Transaksi berhasil dihapus!', 'success')
  } catch (error) {
    console.error('Gagal menghapus transaksi:', error)
    toastRef.value?.show('Waduh, gagal menghapus transaksi nih.', 'error')
  }
}

onMounted(() => {
  fetchTransactions()
})
</script>
