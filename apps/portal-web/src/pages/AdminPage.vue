<template>
  <div class="portal-two-column">
    <section class="glass-panel portal-section">
      <header class="page-header" style="margin-bottom: 12px;">
        <div>
          <h2 class="page-title" style="font-size: 28px;">门户配置</h2>
          <p class="page-subtitle">管理 portal 自己的空闲超时，以及面向用户展示的 client 元数据。</p>
        </div>
      </header>

      <div class="portal-form-stack">
        <el-form label-position="top">
          <el-form-item label="Idle Timeout Minutes">
            <el-input-number v-model="settings.idleTimeoutMinutes" :min="1" :max="240" />
          </el-form-item>
          <el-button type="primary" @click="saveSettings" :loading="savingSettings">保存门户设置</el-button>
        </el-form>
      </div>
    </section>

    <section class="glass-panel portal-section">
      <div class="portal-kicker">Meta Editor</div>
      <p class="portal-muted">选中一条 client meta 进行编辑，保存后立即影响门户可见应用。</p>
      <el-table :data="clientMetas" style="width: 100%; margin-top: 16px;" height="420" @row-click="editMeta">
        <el-table-column prop="clientId" label="Client ID" min-width="160" />
        <el-table-column prop="displayName" label="Display Name" min-width="160" />
        <el-table-column prop="targetUrl" label="Target URL" min-width="180" />
        <el-table-column prop="enabled" label="Enabled" width="100">
          <template #default="{ row }">
            <el-tag :type="row.enabled ? 'success' : 'danger'">{{ row.enabled ? 'Yes' : 'No' }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
      <div style="margin-top: 16px;">
        <el-button plain type="warning" @click="newMeta">新增元数据</el-button>
      </div>
    </section>
  </div>

  <el-dialog v-model="dialogVisible" title="编辑 portal_client_meta" width="720px">
    <el-form label-position="top" class="portal-form-stack">
      <el-form-item label="Client ID">
        <el-input v-model="editing.clientId" />
      </el-form-item>
      <el-form-item label="Display Name">
        <el-input v-model="editing.displayName" />
      </el-form-item>
      <el-form-item label="Description">
        <el-input v-model="editing.description" type="textarea" :rows="3" />
      </el-form-item>
      <el-form-item label="Target URL">
        <el-input v-model="editing.targetUrl" />
      </el-form-item>
      <el-form-item label="Category">
        <el-input v-model="editing.category" />
      </el-form-item>
      <el-form-item label="Sort Order">
        <el-input-number v-model="editing.sortOrder" :min="0" :max="9999" />
      </el-form-item>
      <el-form-item label="Required Realm Roles (comma separated)">
        <el-input v-model="requiredRealmRolesInput" />
      </el-form-item>
      <el-form-item label="Required Client Roles (comma separated)">
        <el-input v-model="requiredClientRolesInput" />
      </el-form-item>
      <el-switch v-model="editing.enabled" active-text="Enabled" inactive-text="Disabled" />
      <el-switch v-model="editing.showInPortal" active-text="Visible In Portal" inactive-text="Hidden In Portal" />
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="danger" plain @click="removeMeta" :disabled="!editing.clientId">删除</el-button>
      <el-button type="primary" @click="saveMeta" :loading="savingMeta">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import {
  ElButton,
  ElDialog,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber,
  ElMessage,
  ElSwitch,
  ElTable,
  ElTableColumn,
  ElTag,
} from 'element-plus'
import { onMounted, reactive, ref } from 'vue'

import { apiClient } from '../api/client'
import type { ApiEnvelope, PortalClientMeta, PortalSettings } from '../api/types'

const dialogVisible = ref(false)
const savingMeta = ref(false)
const savingSettings = ref(false)
const clientMetas = ref<PortalClientMeta[]>([])
const settings = reactive<PortalSettings>({
  realm: '',
  idleTimeoutMinutes: 15,
})

const editing = reactive<PortalClientMeta>({
  realm: '',
  clientId: '',
  displayName: '',
  description: '',
  targetUrl: '',
  icon: '',
  category: '',
  sortOrder: 0,
  enabled: true,
  showInPortal: true,
  requiredRealmRoles: [],
  requiredClientRoles: [],
  tags: [],
})

const requiredRealmRolesInput = ref('')
const requiredClientRolesInput = ref('')

async function loadAdminData(): Promise<void> {
  const [settingsResponse, metasResponse] = await Promise.all([
    apiClient.get<ApiEnvelope<PortalSettings>>('/admin/settings'),
    apiClient.get<ApiEnvelope<PortalClientMeta[]>>('/admin/client-metas'),
  ])

  if (settingsResponse.data.data) {
    Object.assign(settings, settingsResponse.data.data)
  }
  if (metasResponse.data.data) {
    clientMetas.value = metasResponse.data.data
  }
}

function editMeta(meta: PortalClientMeta): void {
  Object.assign(editing, meta)
  requiredRealmRolesInput.value = (meta.requiredRealmRoles ?? []).join(', ')
  requiredClientRolesInput.value = (meta.requiredClientRoles ?? []).join(', ')
  dialogVisible.value = true
}

function newMeta(): void {
  Object.assign(editing, {
    realm: '',
    clientId: '',
    displayName: '',
    description: '',
    targetUrl: '',
    icon: '',
    category: '',
    sortOrder: 0,
    enabled: true,
    showInPortal: true,
    requiredRealmRoles: [],
    requiredClientRoles: [],
    tags: [],
  })
  requiredRealmRolesInput.value = ''
  requiredClientRolesInput.value = ''
  dialogVisible.value = true
}

async function saveSettings(): Promise<void> {
  savingSettings.value = true
  try {
    await apiClient.put('/admin/settings', {
      idleTimeoutMinutes: settings.idleTimeoutMinutes,
    })
    ElMessage.success('门户设置已保存')
  } finally {
    savingSettings.value = false
  }
}

async function saveMeta(): Promise<void> {
  savingMeta.value = true
  try {
    await apiClient.put(`/admin/client-metas/${editing.clientId}`, {
      ...editing,
      requiredRealmRoles: splitRoles(requiredRealmRolesInput.value),
      requiredClientRoles: splitRoles(requiredClientRolesInput.value),
    })
    ElMessage.success('client meta 已保存')
    dialogVisible.value = false
    await loadAdminData()
  } finally {
    savingMeta.value = false
  }
}

async function removeMeta(): Promise<void> {
  if (!editing.clientId) {
    return
  }
  await apiClient.delete(`/admin/client-metas/${editing.clientId}`)
  ElMessage.success('client meta 已删除')
  dialogVisible.value = false
  await loadAdminData()
}

function splitRoles(value: string): string[] {
  return value
    .split(',')
    .map((item) => item.trim())
    .filter(Boolean)
}

onMounted(loadAdminData)
</script>
