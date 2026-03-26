<template>
  <section class="glass-panel portal-section">
    <header class="page-header" style="margin-bottom: 12px;">
      <div>
        <h2 class="page-title" style="font-size: 28px;">当前用户资料</h2>
        <p class="page-subtitle">当前 portal session 与 Mongo 投影中的用户快照。</p>
      </div>
      <el-button plain @click="loadProfile">刷新</el-button>
    </header>

    <el-skeleton v-if="loading" :rows="6" animated />
    <el-descriptions v-else-if="sessionStore.profile" :column="1" border>
      <el-descriptions-item label="Realm">{{ sessionStore.profile.realm.displayName || sessionStore.profile.realm.realmName }}</el-descriptions-item>
      <el-descriptions-item label="Username">{{ sessionStore.profile.user.username }}</el-descriptions-item>
      <el-descriptions-item label="Email">{{ sessionStore.profile.user.email || '-' }}</el-descriptions-item>
      <el-descriptions-item label="First Name">{{ sessionStore.profile.user.firstName || '-' }}</el-descriptions-item>
      <el-descriptions-item label="Last Name">{{ sessionStore.profile.user.lastName || '-' }}</el-descriptions-item>
      <el-descriptions-item label="Realm Roles">{{ sessionStore.profile.user.realmRoles.join(', ') || '-' }}</el-descriptions-item>
      <el-descriptions-item label="Client Roles">{{ JSON.stringify(sessionStore.profile.user.clientRoles, null, 2) }}</el-descriptions-item>
      <el-descriptions-item label="Session">{{ JSON.stringify(sessionStore.profile.session, null, 2) }}</el-descriptions-item>
    </el-descriptions>
  </section>
</template>

<script setup lang="ts">
import { ElButton, ElDescriptions, ElDescriptionsItem, ElSkeleton } from 'element-plus'
import { onMounted, ref } from 'vue'

import { useSessionStore } from '../stores/session'

const sessionStore = useSessionStore()
const loading = ref(false)

async function loadProfile(): Promise<void> {
  loading.value = true
  try {
    await sessionStore.fetchProfile()
  } finally {
    loading.value = false
  }
}

onMounted(loadProfile)
</script>
