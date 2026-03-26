<template>
  <div class="portal-page">
    <section class="glass-panel portal-section portal-auth-card">
      <div class="portal-kicker">Session Expired</div>
      <h2 class="portal-auth-headline">门户会话已因空闲超时退出</h2>
      <p class="portal-auth-copy">
        portal 自己控制默认 15 分钟空闲超时，不直接改 Keycloak 全局超时。重新登录后会再次执行当前用户全量同步。
      </p>
      <el-button type="primary" @click="sessionStore.goLogin()">重新登录</el-button>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ElButton } from 'element-plus'

import { useSessionStore } from '../stores/session'

const sessionStore = useSessionStore()
</script>
