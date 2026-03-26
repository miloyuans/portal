<template>
  <section class="glass-panel portal-section">
    <header class="page-header" style="margin-bottom: 12px;">
      <div>
        <h2 class="page-title" style="font-size: 28px;">Admin Dashboard</h2>
        <p class="page-subtitle">查看当前 realm、最近同步和门户设置总览。</p>
      </div>
      <el-button plain @click="loadDashboard">刷新</el-button>
    </header>

    <div class="portal-grid">
      <el-card class="glass-panel">
        <strong>Projected Realms</strong>
        <p class="portal-muted">{{ realms.length }}</p>
      </el-card>
      <el-card class="glass-panel">
        <strong>Client Count</strong>
        <p class="portal-muted">{{ syncStatus?.clientCount ?? '-' }}</p>
      </el-card>
      <el-card class="glass-panel">
        <strong>Realm Synced At</strong>
        <p class="portal-muted">{{ syncStatus?.realmSyncedAt ?? '-' }}</p>
      </el-card>
      <el-card class="glass-panel">
        <strong>User Synced At</strong>
        <p class="portal-muted">{{ syncStatus?.userSyncedAt ?? '-' }}</p>
      </el-card>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ElButton, ElCard } from 'element-plus'
import { onMounted, ref } from 'vue'

import { apiClient } from '../api/client'
import type { ApiEnvelope, RealmProjection, SyncStatus } from '../api/types'

const realms = ref<RealmProjection[]>([])
const syncStatus = ref<SyncStatus>()

async function loadDashboard(): Promise<void> {
  const [realmsResponse, syncResponse] = await Promise.all([
    apiClient.get<ApiEnvelope<RealmProjection[]>>('/admin/realms'),
    apiClient.get<ApiEnvelope<SyncStatus>>('/admin/sync-status'),
  ])
  realms.value = realmsResponse.data.data ?? []
  syncStatus.value = syncResponse.data.data
}

onMounted(loadDashboard)
</script>
