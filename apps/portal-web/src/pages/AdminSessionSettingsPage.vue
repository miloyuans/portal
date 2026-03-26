<template>
  <section class="glass-panel portal-section">
    <header class="page-header" style="margin-bottom: 12px;">
      <div>
        <h2 class="page-title" style="font-size: 28px;">Session Settings</h2>
        <p class="page-subtitle">修改 portal 自己的 idle timeout 和 warning seconds，不改 Keycloak 全局超时。</p>
      </div>
      <el-button plain @click="loadSettings">刷新</el-button>
    </header>

    <el-form label-position="top" class="portal-form-stack">
      <el-form-item label="Idle Timeout Minutes">
        <el-input-number v-model="settings.idleTimeoutMinutes" :min="1" :max="240" />
      </el-form-item>
      <el-form-item label="Idle Warn Seconds">
        <el-input-number v-model="settings.idleWarnSeconds" :min="5" :max="600" />
      </el-form-item>
      <el-button type="primary" @click="saveSettings">保存设置</el-button>
    </el-form>
  </section>
</template>

<script setup lang="ts">
import { ElButton, ElForm, ElFormItem, ElInputNumber, ElMessage } from 'element-plus'
import { onMounted, reactive } from 'vue'

import { apiClient } from '../api/client'
import type { ApiEnvelope, PortalSettings } from '../api/types'

const settings = reactive<PortalSettings>({
  id: 'global',
  idleTimeoutMinutes: 15,
  idleWarnSeconds: 60,
})

async function loadSettings(): Promise<void> {
  const response = await apiClient.get<ApiEnvelope<PortalSettings>>('/admin/settings/session')
  if (response.data.data) {
    Object.assign(settings, response.data.data)
  }
}

async function saveSettings(): Promise<void> {
  await apiClient.put('/admin/settings/session', settings)
  ElMessage.success('session 设置已保存')
}

onMounted(loadSettings)
</script>
