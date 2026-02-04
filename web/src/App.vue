<template>
  <a-config-provider
    :theme="{
      algorithm: appStore.theme === 'dark' ? theme.darkAlgorithm : theme.defaultAlgorithm,
      token: { colorPrimary: '#1677ff' }
    }"
  >
    <div class="app-container" :class="{ dark: appStore.theme === 'dark' }">
      <template v-if="userStore.isLoggedIn">
        <a-layout class="layout">
          <!-- 左侧导航栏 -->
          <a-layout-sider
            v-model:collapsed="appStore.collapsed"
            :trigger="null"
            collapsible
            class="sider"
            :class="{ 'sider-dark': appStore.theme === 'dark' }"
          >
            <div class="logo">
              <span v-if="!appStore.collapsed">N_m3u8DL-RE</span>
              <span v-else>M3U8</span>
            </div>
            <a-menu
              v-model:selectedKeys="selectedKeys"
              mode="inline"
              :theme="appStore.theme"
              :items="menuItems"
              @click="handleMenuClick"
            />
          </a-layout-sider>

          <a-layout>
            <!-- 顶部栏 -->
            <a-layout-header class="header" :class="{ 'header-dark': appStore.theme === 'dark' }">
              <div class="header-left">
                <menu-unfold-outlined
                  v-if="appStore.collapsed"
                  class="trigger"
                  @click="appStore.toggleCollapsed"
                />
                <menu-fold-outlined
                  v-else
                  class="trigger"
                  @click="appStore.toggleCollapsed"
                />
              </div>

              <div class="header-right">
                <a-space :size="16">
                  <!-- 主题切换 -->
                  <a-tooltip :title="appStore.theme === 'dark' ? '切换到亮色模式' : '切换到暗色模式'">
                    <a-button type="text" @click="appStore.toggleTheme">
                      <template #icon>
                        <bulb-filled v-if="appStore.theme === 'dark'" />
                        <bulb-outlined v-else />
                      </template>
                    </a-button>
                  </a-tooltip>

                  <!-- 用户下拉菜单 -->
                  <a-dropdown>
                    <a-space class="user-dropdown">
                      <a-avatar :size="28">
                        <template #icon><user-outlined /></template>
                      </a-avatar>
                      <span class="username">{{ userStore.username }}</span>
                    </a-space>
                    <template #overlay>
                      <a-menu @click="handleUserMenuClick">
                        <a-menu-item key="logout">
                          <logout-outlined />
                          <span style="margin-left: 8px;">退出登录</span>
                        </a-menu-item>
                      </a-menu>
                    </template>
                  </a-dropdown>
                </a-space>
              </div>
            </a-layout-header>

            <!-- 内容区 -->
            <a-layout-content class="content">
              <router-view />
            </a-layout-content>
          </a-layout>
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
import { theme } from 'ant-design-vue'
import { useUserStore } from './stores/user'
import { useAppStore } from './stores/app'
import {
  DashboardOutlined,
  CloudDownloadOutlined,
  FolderOutlined,
  UserOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  LogoutOutlined,
  BulbOutlined,
  BulbFilled
} from '@ant-design/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const appStore = useAppStore()

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

function handleUserMenuClick({ key }) {
  if (key === 'logout') {
    userStore.logout().then(() => {
      router.push('/login')
    })
  }
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

.app-container.dark {
  background: #141414;
}

.layout {
  min-height: 100vh;
}

.sider {
  background: #fff;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.05);
}

.sider-dark {
  background: #001529;
  box-shadow: 2px 0 8px rgba(0, 0, 0, 0.2);
}

.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: bold;
  color: #1677ff;
  border-bottom: 1px solid #f0f0f0;
}

.sider-dark .logo {
  color: #fff;
  border-bottom-color: #303030;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.header-dark {
  background: #1f1f1f;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.header-left {
  display: flex;
  align-items: center;
}

.trigger {
  font-size: 18px;
  cursor: pointer;
  transition: color 0.3s;
  padding: 8px;
}

.trigger:hover {
  color: #1677ff;
}

.header-dark .trigger {
  color: rgba(255, 255, 255, 0.85);
}

.header-right {
  display: flex;
  align-items: center;
}

.user-dropdown {
  cursor: pointer;
  transition: opacity 0.3s;
}

.user-dropdown:hover {
  opacity: 0.8;
}

.username {
  max-width: 100px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.header-dark .username {
  color: rgba(255, 255, 255, 0.85);
}

.content {
  padding: 24px;
  margin: 0;
  overflow: auto;
}
</style>
