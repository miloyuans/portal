<template>
  <div class="portal-page">
    <section class="glass-panel portal-section portal-auth-card">
      <div class="portal-kicker">404</div>
      <h2 class="portal-auth-headline">页面不存在</h2>
      <p class="portal-auth-copy">当前路由没有对应页面，返回门户首页继续访问。</p>
      <RouterLink to="/portal">
        <el-button type="primary">回到首页</el-button>
      </RouterLink>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ElButton } from 'element-plus'
import { RouterLink } from 'vue-router'
</script>
