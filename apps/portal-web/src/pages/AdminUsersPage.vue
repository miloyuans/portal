<template>
  <section class="glass-panel portal-section">
    <header class="page-header" style="margin-bottom: 12px;">
      <div>
        <h2 class="page-title" style="font-size: 28px;">Projected Users</h2>
        <p class="page-subtitle">按 userId 查看 Mongo 中保存的当前用户投影。</p>
      </div>
    </header>

    <div class="portal-toolbar" style="margin-bottom: 16px;">
      <el-input v-model="userId" placeholder="输入 userId" />
      <el-button type="primary" @click="loadUser">查询</el-button>
    </div>

    <el-empty v-if="!user && !loading" description="输入 userId 后查询" />
    <el-skeleton v-else-if="loading" :rows="6" animated />
    <el-descriptions v-else-if="user" :column="1" border>
      <el-descriptions-item label="User ID">{{ user.userId }}</el-descriptions-item>
      <el-descriptions-item label="Username">{{ user.username }}</el-descriptions-item>
      <el-descriptions-item label="Email">{{ user.email || '-' }}</el-descriptions-item>
      <el-descriptions-item label="Realm Roles">{{ user.realmRoles.join(', ') || '-' }}</el-descriptions-item>
      <el-descriptions-item label="Client Roles">{{ JSON.stringify(user.clientRoles, null, 2) }}</el-descriptions-item>
    </el-descriptions>
  </section>
</template>

<script setup lang="ts">
import { ElButton, ElDescriptions, ElDescriptionsItem, ElEmpty, ElInput, ElSkeleton } from 'element-plus'
import { ref } from 'vue'

import { apiClient } from '../api/client'
import type { ApiEnvelope, UserProjection } from '../api/types'

const userId = ref('')
const user = ref<UserProjection>()
const loading = ref(false)

async function loadUser(): Promise<void> {
  if (!userId.value) {
    return
  }
  loading.value = true
  try {
    const response = await apiClient.get<ApiEnvelope<UserProjection>>(`/admin/users/${userId.value}`)
    user.value = response.data.data
  } finally {
    loading.value = false
  }
}
</script>
