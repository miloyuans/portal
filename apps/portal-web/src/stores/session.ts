import { defineStore } from 'pinia'
import { computed, ref } from 'vue'

import { apiClient, buildApiUrl } from '../api/client'
import type { ApiEnvelope, CurrentUserProfile, PortalApp } from '../api/types'

let idleTimer: number | undefined
let listenersBound = false

export const useSessionStore = defineStore('session', () => {
  const profile = ref<CurrentUserProfile | null>(null)
  const apps = ref<PortalApp[]>([])
  const ready = ref(false)

  const isAuthenticated = computed(() => profile.value !== null)
  const isAdmin = computed(() => profile.value?.isAdmin ?? false)
  const displayName = computed(() => profile.value?.user.displayName || profile.value?.user.username || 'Portal User')

  async function bootstrap(): Promise<void> {
    if (ready.value) {
      return
    }
    try {
      await fetchSession()
    } catch {
      profile.value = null
      apps.value = []
    } finally {
      ready.value = true
    }
  }

  async function fetchSession(): Promise<CurrentUserProfile> {
    const response = await apiClient.get<ApiEnvelope<CurrentUserProfile>>('/me')
    if (!response.data.success || !response.data.data) {
      throw new Error(response.data.error?.message ?? '无法获取当前会话')
    }
    profile.value = response.data.data
    bindActivityListeners()
    scheduleIdleLogout()
    return response.data.data
  }

  async function fetchApps(): Promise<void> {
    const response = await apiClient.get<ApiEnvelope<PortalApp[]>>('/apps')
    if (!response.data.success || !response.data.data) {
      throw new Error(response.data.error?.message ?? '无法获取应用列表')
    }
    apps.value = response.data.data
  }

  function goLogin(): void {
    window.location.href = buildApiUrl('/auth/login')
  }

  function logout(reason: 'manual' | 'expired' = 'manual'): void {
    const suffix = reason === 'expired' ? '?reason=expired' : ''
    window.location.href = buildApiUrl(`/auth/logout${suffix}`)
  }

  function scheduleIdleLogout(): void {
    clearIdleTimer()
    const timeoutMinutes = profile.value?.settings.idleTimeoutMinutes ?? 15
    idleTimer = window.setTimeout(() => logout('expired'), timeoutMinutes * 60 * 1000)
  }

  function touchActivity(): void {
    if (!profile.value) {
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
    fetchSession,
    goLogin,
    isAdmin,
    isAuthenticated,
    logout,
    profile,
    ready,
  }
})
