import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import TransactionView from '@/views/TransactionView.vue'
import WalletView from '@/views/WalletView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/transactions',
      name: 'transaction',
      component: TransactionView,
    },
    {
      path: '/wallets',
      name: 'wallet',
      component: WalletView,
    },
  ],
})

export default router
