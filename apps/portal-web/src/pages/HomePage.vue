<template>
  <section class="glass-panel portal-section">
    <header class="page-header" style="margin-bottom: 12px;">
      <div>
        <h2 class="page-title" style="font-size: 28px;">用户导航页面</h2>
        <p class="page-subtitle">按 portal 权限规则计算后的可访问应用。</p>
      </div>
      <div class="portal-toolbar">
        <el-select v-model="selectedCategory" clearable placeholder="按分类筛选" style="width: 200px;">
          <el-option v-for="category in categories" :key="category" :label="category" :value="category" />
        </el-select>
        <el-select v-model="selectedRealm" disabled style="width: 220px;">
          <el-option v-for="realm in sessionStore.realms" :key="realm.realmId" :label="realm.displayName || realm.realmName" :value="realm.realmId" />
        </el-select>
        <el-button plain @click="reload">刷新</el-button>
      </div>
    </header>

    <el-skeleton v-if="loading" :rows="6" animated />
    <div v-else-if="!filteredApps.length" class="portal-empty">
      <h3>当前没有可见应用</h3>
      <p class="portal-muted">确认 portal_client_meta.visible、accessRules 与当前用户角色是否匹配。</p>
    </div>
    <div v-else class="portal-grid">
      <AppCard v-for="app in filteredApps" :key="app.clientId" :app="app" />
    </div>
  </section>
</template>

<script setup lang="ts">
import { ElButton, ElOption, ElSelect, ElSkeleton } from 'element-plus'
import { computed, onMounted, ref } from 'vue'

import AppCard from '../components/AppCard.vue'
import { useSessionStore } from '../stores/session'

const sessionStore = useSessionStore()
const loading = ref(false)
const selectedCategory = ref<string>()
const selectedRealm = ref('')

const categories = computed(() =>
  Array.from(new Set(sessionStore.apps.map((item) => item.category).filter(Boolean))) as string[],
)

const filteredApps = computed(() => {
  return sessionStore.apps.filter((app) => {
    if (selectedCategory.value && app.category !== selectedCategory.value) {
      return false
    }
    return true
  })
})

async function reload(): Promise<void> {
  loading.value = true
  try {
    await Promise.all([sessionStore.fetchApps(), sessionStore.fetchRealms()])
    selectedRealm.value = sessionStore.realms[0]?.realmId ?? ''
  } finally {
    loading.value = false
  }
}

onMounted(reload)
</script>
