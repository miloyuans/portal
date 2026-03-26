<template>
  <div class="portal-page">
    <div class="portal-auth-card glass-panel">
      <div class="portal-kicker">OIDC Login</div>
      <h1 class="portal-auth-headline">Portal 只接受 Keycloak 登录。</h1>
      <p class="portal-auth-copy">
        这里不保存用户名密码，也不做本地认证。点击后将跳转到 Keycloak，成功后由 portal-api 完成全量同步并建立 portal session。
      </p>
      <div class="portal-toolbar">
        <el-button type="primary" size="large" @click="sessionStore.goLogin()">前往 Keycloak 登录</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ElButton } from 'element-plus'

import { useSessionStore } from '../stores/session'

const sessionStore = useSessionStore()
</script>
