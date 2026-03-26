<template>
  <div class="portal-page">
    <div class="portal-shell">
      <header class="glass-panel portal-section page-header">
        <div>
          <div class="portal-kicker">{{ appTitle }}</div>
          <h1 class="page-title">Portal</h1>
          <p class="page-subtitle">Keycloak 统一认证入口，portal-api 负责会话、同步、权限解析与导航返回。</p>
        </div>
        <div class="portal-toolbar">
          <RouterLink to="/portal">
            <el-button plain>门户</el-button>
          </RouterLink>
          <RouterLink to="/profile">
            <el-button plain>Profile</el-button>
          </RouterLink>
          <RouterLink v-if="sessionStore.isAdmin" to="/admin/dashboard">
            <el-button plain type="warning">Admin</el-button>
          </RouterLink>
          <el-tag effect="dark" type="info">{{ sessionStore.displayName }}</el-tag>
          <el-button type="primary" @click="onLogout">退出登录</el-button>
        </div>
      </header>

      <nav v-if="sessionStore.isAdmin" class="glass-panel portal-section" style="margin-bottom: 20px;">
        <div class="portal-toolbar">
          <RouterLink to="/admin/dashboard"><el-button plain>Dashboard</el-button></RouterLink>
          <RouterLink to="/admin/clients"><el-button plain>Clients</el-button></RouterLink>
          <RouterLink to="/admin/users"><el-button plain>Users</el-button></RouterLink>
          <RouterLink to="/admin/settings/session"><el-button plain>Session</el-button></RouterLink>
          <RouterLink to="/admin/sync-status"><el-button plain>Sync Status</el-button></RouterLink>
        </div>
      </nav>

      <RouterView />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ElButton, ElTag } from 'element-plus'
import { RouterLink, RouterView } from 'vue-router'

import { useSessionStore } from '../stores/session'
import { getRuntimeConfig } from '../utils/runtimeConfig'

const sessionStore = useSessionStore()
const appTitle = getRuntimeConfig().appTitle

async function onLogout(): Promise<void> {
  await sessionStore.logout('manual')
}
</script>
