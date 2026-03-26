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
    <p class="portal-muted">{{ app.description || '已同步到门户，点击即可进入。' }}</p>
    <div v-if="app.tags?.length">
      <el-tag v-for="tag in app.tags" :key="tag" class="portal-card-tag" effect="plain">
        {{ tag }}
      </el-tag>
    </div>
    <div style="margin-top: 20px;">
      <el-button type="primary" @click="openApp">进入应用</el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ElButton, ElCard, ElTag } from 'element-plus'

import type { PortalApp } from '../api/types'

const props = defineProps<{
  app: PortalApp
}>()

function openApp(): void {
  window.open(props.app.targetUrl, '_blank', 'noopener,noreferrer')
}
</script>
