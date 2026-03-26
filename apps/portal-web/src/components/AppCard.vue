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
    <p class="portal-muted">canView={{ app.canView }} · canLaunch={{ app.canLaunch }} · canAdmin={{ app.canAdmin }}</p>
    <div style="margin-top: 20px;">
      <el-button type="primary" :disabled="!app.canLaunch" @click="openApp">进入应用</el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ElButton, ElCard, ElTag } from 'element-plus'

import type { PortalAppView } from '../api/types'

const props = defineProps<{
  app: PortalAppView
}>()

function openApp(): void {
  window.open(props.app.launchUrl, '_blank', 'noopener,noreferrer')
}
</script>
