import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

import { apiClient, buildApiUrl } from '../api/client'
import type { ApiEnvelope, CurrentUserProfile, PortalAppView, RealmProjection, SessionView } from '../api/types'

let idleTimer: number | undefined
let listenersBound = false

export const useSessionStore = defineStore('session', () => {
  const me = ref<SessionView | null>(null)
  const profile = ref<CurrentUserProfile | null>(null)
  const apps = ref<PortalAppView[]>([])
  const realms = ref<RealmProjection[]>([])
  const ready = ref(false)

  const isAuthenticated = computed(() => me.value !== null)
  const isAdmin = computed(() => me.value?.realmRoles.includes('portal_admin') ?? false)
  const displayName = computed(() => me.value?.displayName || me.value?.username || 'Portal User')

  async function bootstrap(): Promise<void> {
    if (ready.value) {
      return
    }
    try {
      await fetchMe()
    } catch {
      me.value = null
      profile.value = null
      apps.value = []
      realms.value = []
    } finally {
      ready.value = true
    }
  }

  async function fetchMe(): Promise<SessionView> {
    const response = await apiClient.get<ApiEnvelope<SessionView>>('/auth/me')
    if (!response.data.success || !response.data.data) {
      throw new Error(response.data.error?.message ?? '无法获取当前会话')
    }
    me.value = response.data.data
    bindActivityListeners()
    scheduleIdleLogout()
    return response.data.data
  }

  async function fetchProfile(): Promise<void> {
    const response = await apiClient.get<ApiEnvelope<CurrentUserProfile>>('/portal/profile')
    if (!response.data.success || !response.data.data) {
      throw new Error(response.data.error?.message ?? '无法获取当前用户资料')
    }
    profile.value = response.data.data
  }

  async function fetchApps(): Promise<void> {
    const response = await apiClient.get<ApiEnvelope<PortalAppView[]>>('/portal/apps')
    if (!response.data.success || !response.data.data) {
      throw new Error(response.data.error?.message ?? '无法获取应用列表')
    }
    apps.value = response.data.data
  }

  async function fetchRealms(): Promise<void> {
    const response = await apiClient.get<ApiEnvelope<RealmProjection[]>>('/portal/realms')
    if (!response.data.success || !response.data.data) {
      throw new Error(response.data.error?.message ?? '无法获取 realm 列表')
    }
    realms.value = response.data.data
  }

  function goLogin(): void {
    window.location.href = buildApiUrl('/auth/login')
  }

  async function logout(reason: 'manual' | 'expired' = 'manual'): Promise<void> {
    try {
      const suffix = reason === 'expired' ? '?reason=expired' : ''
      const response = await apiClient.post<ApiEnvelope<{ logoutUrl: string }>>(`/auth/logout${suffix}`)
      const logoutUrl = response.data.data?.logoutUrl
      window.location.href = logoutUrl ?? (reason === 'expired' ? '/session-expired' : '/login')
    } catch {
      window.location.href = reason === 'expired' ? '/session-expired' : '/login'
    }
  }

  function scheduleIdleLogout(): void {
    clearIdleTimer()
    const timeoutMinutes = me.value?.idleTimeoutMinutes ?? 15
    idleTimer = window.setTimeout(() => {
      void logout('expired')
    }, timeoutMinutes * 60 * 1000)
  }

  function touchActivity(): void {
    if (!me.value) {
      return
    }
    scheduleIdleLogout()
  }

  function bindActivityListeners(): void {
    if (listenersBound) {
      return
    }
    ;['click', 'keydown', 'mousemove', 'touchstart', 'scroll'].forEach((eventName) => {
      window.addEventListener(eventName, touchActivity, { passive: true })
    })
    listenersBound = true
  }

  function clearIdleTimer(): void {
    if (idleTimer !== undefined) {
      window.clearTimeout(idleTimer)
      idleTimer = undefined
    }
  }

  return {
    apps,
    bootstrap,
    displayName,
    fetchApps,
    fetchMe,
    fetchProfile,
    fetchRealms,
    goLogin,
    isAdmin,
    isAuthenticated,
    logout,
    me,
    profile,
    realms,
    ready,
  }
})
