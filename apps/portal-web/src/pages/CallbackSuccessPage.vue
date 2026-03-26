<template>
  <div class="portal-page">
    <div class="portal-auth-card glass-panel">
      <div class="portal-kicker">Sync Complete</div>
      <h1 class="portal-auth-headline">正在建立门户会话</h1>
      <p class="portal-auth-copy">realm、clients、当前用户与角色投影已经同步，正在进入门户。</p>
      <el-skeleton :rows="3" animated />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ElSkeleton } from 'element-plus'
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

import { useSessionStore } from '../stores/session'

const router = useRouter()
const sessionStore = useSessionStore()

onMounted(async () => {
  try {
    await sessionStore.fetchMe()
    await router.replace({ name: 'portal' })
  } catch {
    await router.replace({ name: 'login' })
  }
})
</script>
