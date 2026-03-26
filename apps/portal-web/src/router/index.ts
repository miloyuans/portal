import { createRouter, createWebHistory } from 'vue-router'

import AppShell from '../layouts/AppShell.vue'
import AdminClientsPage from '../pages/AdminClientsPage.vue'
import AdminDashboardPage from '../pages/AdminDashboardPage.vue'
import AdminSessionSettingsPage from '../pages/AdminSessionSettingsPage.vue'
import AdminSyncStatusPage from '../pages/AdminSyncStatusPage.vue'
import AdminUsersPage from '../pages/AdminUsersPage.vue'
import CallbackSuccessPage from '../pages/CallbackSuccessPage.vue'
import ForbiddenPage from '../pages/ForbiddenPage.vue'
import HomePage from '../pages/HomePage.vue'
import LoginPage from '../pages/LoginPage.vue'
import NotFoundPage from '../pages/NotFoundPage.vue'
import ProfilePage from '../pages/ProfilePage.vue'
import SessionExpiredPage from '../pages/SessionExpiredPage.vue'
import { useSessionStore } from '../stores/session'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/portal',
    },
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
      path: '/session-expired',
      name: 'session-expired',
      component: SessionExpiredPage,
    },
    {
      path: '/403',
      name: 'forbidden',
      component: ForbiddenPage,
    },
    {
      path: '/',
      component: AppShell,
      meta: { requiresAuth: true },
      children: [
        {
          path: 'portal',
          name: 'portal',
          component: HomePage,
          meta: { requiresAuth: true },
        },
        {
          path: 'profile',
          name: 'profile',
          component: ProfilePage,
          meta: { requiresAuth: true },
        },
        {
          path: 'admin/dashboard',
          name: 'admin-dashboard',
          component: AdminDashboardPage,
          meta: { requiresAuth: true, admin: true },
        },
        {
          path: 'admin/clients',
          name: 'admin-clients',
          component: AdminClientsPage,
          meta: { requiresAuth: true, admin: true },
        },
        {
          path: 'admin/users',
          name: 'admin-users',
          component: AdminUsersPage,
          meta: { requiresAuth: true, admin: true },
        },
        {
          path: 'admin/settings/session',
          name: 'admin-session-settings',
          component: AdminSessionSettingsPage,
          meta: { requiresAuth: true, admin: true },
        },
        {
          path: 'admin/sync-status',
          name: 'admin-sync-status',
          component: AdminSyncStatusPage,
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
    return { name: 'portal' }
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
