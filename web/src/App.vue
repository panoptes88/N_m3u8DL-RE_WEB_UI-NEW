<template>
  <a-config-provider :theme="{ token: { colorPrimary: '#1677ff' } }">
    <div class="app-container">
      <template v-if="userStore.isLoggedIn">
        <a-layout class="layout">
          <a-layout-header class="header">
            <div class="logo">N_m3u8DL-RE</div>
            <a-menu
              v-model:selectedKeys="selectedKeys"
              mode="horizontal"
              :items="menuItems"
              @click="handleMenuClick"
            />
            <div class="user-info">
              <a-space>
                <span>{{ userStore.username }}</span>
                <a-button type="link" @click="handleLogout">退出</a-button>
              </a-space>
            </div>
          </a-layout-header>
          <a-layout-content class="content">
            <router-view />
          </a-layout-content>
        </a-layout>
      </template>
      <template v-else>
        <router-view />
      </template>
    </div>
  </a-config-provider>
</template>

<script setup>
import { ref, onMounted, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from './stores/user'
import {
  DashboardOutlined,
  CloudDownloadOutlined,
  FolderOutlined,
  UserOutlined
} from '@ant-design/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const selectedKeys = ref(['dashboard'])

const menuItems = [
  {
    key: 'dashboard',
    icon: () => h(DashboardOutlined),
    label: '首页'
  },
  {
    key: 'tasks',
    icon: () => h(CloudDownloadOutlined),
    label: '下载任务'
  },
  {
    key: 'files',
    icon: () => h(FolderOutlined),
    label: '文件管理'
  }
]

function handleMenuClick({ key }) {
  router.push({ name: key.charAt(0).toUpperCase() + key.slice(1) })
}

function handleLogout() {
  userStore.logout().then(() => {
    router.push('/login')
  })
}

onMounted(async () => {
  const isLoggedIn = await userStore.checkLogin()
  if (!isLoggedIn && route.name !== 'Login') {
    router.push('/login')
  }
})
</script>

<style scoped>
.app-container {
  min-height: 100vh;
  background: #f0f2f5;
}

.layout {
  min-height: 100vh;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.09);
}

.logo {
  font-size: 18px;
  font-weight: bold;
  color: #1677ff;
}

.user-info {
  display: flex;
  align-items: center;
}

.content {
  padding: 24px;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}
</style>
