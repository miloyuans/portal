<template>
  <div class="portal-page">
    <section class="glass-panel portal-section portal-auth-card">
      <div class="portal-kicker">403</div>
      <h2 class="portal-auth-headline">没有管理权限</h2>
      <p class="portal-auth-copy">管理页需要具备 portal 管理 realm role，当前 portal session 不满足该条件。</p>
      <RouterLink to="/">
        <el-button type="primary">返回门户首页</el-button>
      </RouterLink>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ElButton } from 'element-plus'
import { RouterLink } from 'vue-router'
</script>
