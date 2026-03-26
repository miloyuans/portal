<template>
  <section class="glass-panel portal-section">
    <header class="page-header" style="margin-bottom: 12px;">
      <div>
        <h2 class="page-title" style="font-size: 28px;">应用导航</h2>
        <p class="page-subtitle">仅展示当前 portal session 权限可见的应用。</p>
      </div>
      <el-button plain @click="reload">刷新</el-button>
    </header>

    <el-skeleton v-if="loading" :rows="6" animated />
    <div v-else-if="!sessionStore.apps.length" class="portal-empty">
      <h3>当前没有可见应用</h3>
      <p class="portal-muted">检查 Keycloak client roles、realm roles 或 portal_client_meta 的启用状态。</p>
    </div>
    <div v-else class="portal-grid">
      <AppCard v-for="app in sessionStore.apps" :key="app.clientId" :app="app" />
    </div>
  </section>
</template>

<script setup lang="ts">
import { ElButton, ElSkeleton } from 'element-plus'
import { onMounted, ref } from 'vue'

import AppCard from '../components/AppCard.vue'
import { useSessionStore } from '../stores/session'

const sessionStore = useSessionStore()
const loading = ref(false)

async function reload(): Promise<void> {
  loading.value = true
  try {
    await sessionStore.fetchApps()
  } finally {
    loading.value = false
  }
}

onMounted(reload)
</script>
