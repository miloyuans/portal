import { createRouter, createWebHistory } from 'vue-router'

import AppShell from '../layouts/AppShell.vue'
import AdminPage from '../pages/AdminPage.vue'
import CallbackSuccessPage from '../pages/CallbackSuccessPage.vue'
import ForbiddenPage from '../pages/ForbiddenPage.vue'
import HomePage from '../pages/HomePage.vue'
import LoginPage from '../pages/LoginPage.vue'
import NotFoundPage from '../pages/NotFoundPage.vue'
import SessionExpiredPage from '../pages/SessionExpiredPage.vue'
import { useSessionStore } from '../stores/session'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginPage,
    },
    {
      path: '/auth/callback/success',
      name: 'callback-success',
      component: CallbackSuccessPage,
    },
    {
      path: '/403',
      name: 'forbidden',
      component: ForbiddenPage,
    },
    {
      path: '/session-expired',
      name: 'session-expired',
      component: SessionExpiredPage,
    },
    {
      path: '/',
      component: AppShell,
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          name: 'home',
          component: HomePage,
        },
        {
          path: 'admin',
          name: 'admin',
          component: AdminPage,
          meta: { requiresAuth: true, admin: true },
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: NotFoundPage,
    },
  ],
})

router.beforeEach(async (to) => {
  const store = useSessionStore()
  await store.bootstrap()

  if (to.name === 'login' && store.isAuthenticated) {
    return { name: 'home' }
  }
  if (to.meta.requiresAuth && !store.isAuthenticated) {
    return { name: 'login' }
  }
  if (to.meta.admin && !store.isAdmin) {
    return { name: 'forbidden' }
  }
  return true
})

export default router
