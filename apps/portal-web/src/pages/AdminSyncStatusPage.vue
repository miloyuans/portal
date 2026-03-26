<template>
  <section class="glass-panel portal-section">
    <header class="page-header" style="margin-bottom: 12px;">
      <div>
        <h2 class="page-title" style="font-size: 28px;">Sync Status</h2>
        <p class="page-subtitle">展示当前登录用户最近一次登录触发同步的投影状态。</p>
      </div>
      <el-button plain @click="loadStatus">刷新</el-button>
    </header>

    <el-descriptions v-if="status" :column="1" border>
      <el-descriptions-item label="Realm ID">{{ status.realmId }}</el-descriptions-item>
      <el-descriptions-item label="Realm Synced At">{{ status.realmSyncedAt }}</el-descriptions-item>
      <el-descriptions-item label="User Synced At">{{ status.userSyncedAt }}</el-descriptions-item>
      <el-descriptions-item label="Client Count">{{ status.clientCount }}</el-descriptions-item>
      <el-descriptions-item label="Settings Updated At">{{ status.settingsUpdatedAt }}</el-descriptions-item>
    </el-descriptions>
  </section>
</template>

<script setup lang="ts">
import { ElButton, ElDescriptions, ElDescriptionsItem } from 'element-plus'
import { onMounted, ref } from 'vue'

import { apiClient } from '../api/client'
import type { ApiEnvelope, SyncStatus } from '../api/types'

const status = ref<SyncStatus>()

async function loadStatus(): Promise<void> {
  const response = await apiClient.get<ApiEnvelope<SyncStatus>>('/admin/sync-status')
  status.value = response.data.data
}

onMounted(loadStatus)
</script>
