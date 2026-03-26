<template>
  <div class="portal-page">
    <div class="portal-shell">
      <header class="glass-panel portal-section page-header">
        <div>
          <div class="portal-kicker">Portal</div>
          <h1 class="page-title">{{ appTitle }}</h1>
          <p class="page-subtitle">Keycloak 统一认证，Mongo 仅做门户投影，所有页面统一走 portal-api。</p>
        </div>
        <div class="portal-toolbar">
          <RouterLink to="/">
            <el-button plain>用户门户</el-button>
          </RouterLink>
          <RouterLink v-if="sessionStore.isAdmin" to="/admin">
            <el-button plain type="warning">管理页</el-button>
          </RouterLink>
          <el-tag effect="dark" type="info">{{ sessionStore.displayName }}</el-tag>
          <el-button type="primary" @click="sessionStore.logout()">退出登录</el-button>
        </div>
      </header>
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
</script>
