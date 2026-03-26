<template>
  <div class="portal-page">
    <div class="portal-auth-card glass-panel">
      <div class="portal-kicker">OIDC Login</div>
      <h1 class="portal-auth-headline">Portal only accepts Keycloak login.</h1>
      <p class="portal-auth-copy">
        Portal stores no passwords and performs no local authentication. This page redirects to Keycloak, then
        portal-api completes the login sync and creates the portal session.
      </p>
      <div class="portal-toolbar">
        <el-button type="primary" size="large" @click="sessionStore.goLogin()">Continue to Keycloak</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ElButton } from 'element-plus'
import { onMounted } from 'vue'

import { useSessionStore } from '../stores/session'

const sessionStore = useSessionStore()

onMounted(() => {
  window.setTimeout(() => {
    sessionStore.goLogin()
  }, 150)
})
</script>
