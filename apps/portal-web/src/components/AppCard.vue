<template>
  <el-card shadow="hover" class="glass-panel">
    <template #header>
      <div style="display: flex; align-items: center; justify-content: space-between; gap: 12px;">
        <div>
          <strong>{{ app.displayName }}</strong>
          <div class="portal-muted" style="margin-top: 6px;">{{ app.clientId }}</div>
        </div>
        <el-tag v-if="app.category" type="success" effect="plain">{{ app.category }}</el-tag>
      </div>
    </template>

    <p class="portal-muted">
      mode={{ app.launchMode }} | canView={{ app.canView }} | canLaunch={{ app.canLaunch }} | canAdmin={{ app.canAdmin }}
    </p>

    <div style="margin-top: 20px;">
      <el-button type="primary" :disabled="!app.canLaunch" :loading="launching" @click="openApp">
        Open App
      </el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ElButton, ElCard, ElMessage, ElTag } from 'element-plus'
import { ref } from 'vue'

import { apiClient } from '../api/client'
import type { ApiEnvelope, PortalAppView, PortalLaunchView } from '../api/types'

const props = defineProps<{
  app: PortalAppView
}>()

const launching = ref(false)

async function openApp(): Promise<void> {
  if (!props.app.canLaunch || launching.value) {
    return
  }

  launching.value = true
  try {
    const response = await apiClient.get<ApiEnvelope<PortalLaunchView>>(
      `/portal/apps/${encodeURIComponent(props.app.clientId)}/launch`,
    )
    const launchUrl = response.data.data?.launchUrl
    if (!launchUrl) {
      throw new Error('missing launchUrl')
    }
    window.open(launchUrl, '_blank', 'noopener,noreferrer')
  } catch {
    ElMessage.error('Portal could not resolve the app launch target.')
  } finally {
    launching.value = false
  }
}
</script>
