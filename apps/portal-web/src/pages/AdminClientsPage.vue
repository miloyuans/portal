<template>
  <section class="glass-panel portal-section">
    <header class="page-header" style="margin-bottom: 12px;">
      <div>
        <h2 class="page-title" style="font-size: 28px;">Client Metadata</h2>
        <p class="page-subtitle">维护 portal_client_meta，决定 client 是否展示、展示名、跳转地址与 accessRules。</p>
      </div>
      <el-button plain @click="loadRows">刷新</el-button>
    </header>

    <el-table :data="rows" style="width: 100%;" height="520" @row-click="openEditor">
      <el-table-column prop="client.clientId" label="Client ID" min-width="160" />
      <el-table-column prop="client.name" label="Client Name" min-width="160" />
      <el-table-column prop="meta.displayName" label="Portal Name" min-width="180" />
      <el-table-column prop="meta.visible" label="Visible" width="100">
        <template #default="{ row }">
          <el-tag :type="row.meta?.visible ? 'success' : 'info'">{{ row.meta?.visible ? 'Yes' : 'No' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="meta.launchUrl" label="Launch URL" min-width="220" />
    </el-table>
  </section>

  <el-dialog v-model="dialogVisible" title="编辑 portal_client_meta" width="760px">
    <el-form label-position="top" class="portal-form-stack">
      <el-form-item label="Client ID">
        <el-input v-model="editing.clientId" disabled />
      </el-form-item>
      <el-form-item label="Display Name">
        <el-input v-model="editing.displayName" />
      </el-form-item>
      <el-form-item label="Icon">
        <el-input v-model="editing.icon" />
      </el-form-item>
      <el-form-item label="Category">
        <el-input v-model="editing.category" />
      </el-form-item>
      <el-form-item label="Sort">
        <el-input-number v-model="editing.sort" :min="0" :max="9999" />
      </el-form-item>
      <el-form-item label="Launch URL">
        <el-input v-model="editing.launchUrl" />
      </el-form-item>
      <el-form-item label="Any Realm Roles (comma separated)">
        <el-input v-model="anyRealmRolesInput" />
      </el-form-item>
      <el-form-item label="Any Client Roles (comma separated)">
        <el-input v-model="anyClientRolesInput" />
      </el-form-item>
      <el-switch v-model="editing.visible" active-text="Visible" inactive-text="Hidden" />
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="saveMeta">保存</el-button>
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
import type { AdminClientRow, ApiEnvelope, PortalClientMeta } from '../api/types'

const rows = ref<AdminClientRow[]>([])
const dialogVisible = ref(false)
const editing = reactive<PortalClientMeta>({
  realmId: '',
  clientId: '',
  displayName: '',
  icon: '',
  category: '',
  sort: 0,
  launchUrl: '',
  visible: false,
  accessRules: {
    anyRealmRoles: [],
    anyClientRoles: [],
    adminRealmRoles: ['portal_admin'],
  },
})
const anyRealmRolesInput = ref('')
const anyClientRolesInput = ref('')

async function loadRows(): Promise<void> {
  const response = await apiClient.get<ApiEnvelope<AdminClientRow[]>>('/admin/clients')
  rows.value = response.data.data ?? []
}

function openEditor(row: AdminClientRow): void {
  Object.assign(editing, {
    realmId: row.meta?.realmId ?? row.client.realmId,
    clientId: row.client.clientId,
    displayName: row.meta?.displayName ?? row.client.name ?? row.client.clientId,
    icon: row.meta?.icon ?? '',
    category: row.meta?.category ?? '',
    sort: row.meta?.sort ?? 0,
    launchUrl: row.meta?.launchUrl ?? row.client.baseUrl ?? row.client.rootUrl ?? '',
    visible: row.meta?.visible ?? false,
    accessRules: {
      anyRealmRoles: row.meta?.accessRules?.anyRealmRoles ?? [],
      anyClientRoles: row.meta?.accessRules?.anyClientRoles ?? [],
      adminRealmRoles: row.meta?.accessRules?.adminRealmRoles ?? ['portal_admin'],
    },
  })
  anyRealmRolesInput.value = editing.accessRules?.anyRealmRoles?.join(', ') ?? ''
  anyClientRolesInput.value = editing.accessRules?.anyClientRoles?.join(', ') ?? ''
  dialogVisible.value = true
}

async function saveMeta(): Promise<void> {
  await apiClient.put(`/admin/clients/${editing.clientId}/meta`, {
    ...editing,
    accessRules: {
      anyRealmRoles: splitInput(anyRealmRolesInput.value),
      anyClientRoles: splitInput(anyClientRolesInput.value),
      adminRealmRoles: ['portal_admin'],
    },
  })
  ElMessage.success('portal_client_meta 已保存')
  dialogVisible.value = false
  await loadRows()
}

function splitInput(value: string): string[] {
  return value
    .split(',')
    .map((item) => item.trim())
    .filter(Boolean)
}

onMounted(loadRows)
</script>
