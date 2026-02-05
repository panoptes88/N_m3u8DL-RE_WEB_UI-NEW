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
          <!-- PC端左侧导航栏 (移动端隐藏) -->
          <a-layout-sider
            v-model:collapsed="appStore.collapsed"
            :trigger="null"
            collapsible
            class="sider pc-sider"
            :class="{ 'sider-dark': appStore.theme === 'dark' }"
            :width="200"
            :collapsedWidth="80"
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
                <!-- 移动端菜单图标 -->
                <menu-unfold-outlined
                  class="trigger mobile-trigger"
                  @click="mobileMenuVisible = true"
                />
                <!-- PC端收缩图标 -->
                <menu-unfold-outlined
                  v-if="appStore.collapsed"
                  class="trigger pc-trigger"
                  @click="appStore.toggleCollapsed"
                />
                <menu-fold-outlined
                  v-else
                  class="trigger pc-trigger"
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
                        <a-menu-item key="change-password">
                          <key-outlined />
                          <span style="margin-left: 8px;">修改密码</span>
                        </a-menu-item>
                        <a-menu-divider />
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

        <!-- 移动端抽屉式侧边栏 -->
        <a-drawer
          v-model:open="mobileMenuVisible"
          placement="left"
          :width="200"
          :closable="false"
          :show-header="false"
          class="mobile-drawer"
          :body-style="{ padding: 0 }"
        >
          <div class="mobile-menu-header">
            <span>M3U8</span>
          </div>
          <a-menu
            v-model:selectedKeys="selectedKeys"
            mode="inline"
            :theme="appStore.theme"
            :items="menuItems"
            @click="handleMobileMenuClick"
          />
        </a-drawer>
      </template>
      <template v-else>
        <router-view />
      </template>

      <!-- 修改密码弹窗 -->
      <a-modal
        v-model:open="changePasswordVisible"
        title="修改密码"
        :confirm-loading="changePasswordLoading"
        @ok="handleChangePasswordOk"
        @cancel="handleChangePasswordCancel"
      >
        <a-form :model="changePasswordForm" layout="vertical">
          <a-form-item label="新密码" name="newPassword" :rules="[{ required: true, message: '请输入新密码' }]">
            <a-input-password
              v-model:value="changePasswordForm.newPassword"
              placeholder="请输入新密码"
              @pressEnter="handleChangePasswordOk"
            />
          </a-form-item>
        </a-form>
      </a-modal>
    </div>
  </a-config-provider>
</template>

<script setup>
import { ref, reactive, onMounted, h } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import { theme } from 'ant-design-vue'
import { useUserStore } from './stores/user'
import { useAppStore } from './stores/app'
import { post } from './api'
import {
  DashboardOutlined,
  CloudDownloadOutlined,
  FolderOutlined,
  UserOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  LogoutOutlined,
  KeyOutlined,
  BulbOutlined,
  BulbFilled
} from '@ant-design/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const appStore = useAppStore()

const selectedKeys = ref(['dashboard'])
const mobileMenuVisible = ref(false)

// 修改密码相关
const changePasswordVisible = ref(false)
const changePasswordLoading = ref(false)
const changePasswordForm = reactive({
  newPassword: ''
})

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

function handleMobileMenuClick({ key }) {
  mobileMenuVisible.value = false
  router.push({ name: key.charAt(0).toUpperCase() + key.slice(1) })
}

function handleUserMenuClick({ key }) {
  if (key === 'change-password') {
    changePasswordVisible.value = true
    changePasswordForm.newPassword = ''
  } else if (key === 'logout') {
    userStore.logout().then(() => {
      router.push('/login')
    })
  }
}

async function handleChangePasswordOk() {
  if (!changePasswordForm.newPassword) {
    message.warning('请输入新密码')
    return
  }

  changePasswordLoading.value = true
  try {
    await post('/auth/change-password', {
      new_password: changePasswordForm.newPassword
    })
    message.success('密码修改成功')
    changePasswordVisible.value = false
    userStore.logout().then(() => {
      router.push('/login')
    })
  } catch (err) {
    message.error(err.response?.data?.error || '修改密码失败')
  } finally {
    changePasswordLoading.value = false
  }
}

function handleChangePasswordCancel() {
  changePasswordVisible.value = false
  changePasswordForm.newPassword = ''
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

/* PC端侧边栏 */
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

/* PC端触发器 */
.pc-trigger {
  display: inline-block;
}

/* 移动端触发器 */
.mobile-trigger {
  display: none;
}

/* 移动端抽屉样式 */
.mobile-menu-header {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: bold;
  color: #1677ff;
  border-bottom: 1px solid #f0f0f0;
}

:global(.dark) .mobile-menu-header {
  color: #fff;
  border-bottom-color: #303030;
}

/* 移动端响应式 */
@media (max-width: 768px) {
  .pc-sider {
    display: none !important;
  }

  .pc-trigger {
    display: none !important;
  }

  .mobile-trigger {
    display: inline-block !important;
  }

  .header {
    padding: 0 16px;
  }

  .content {
    padding: 16px;
  }

  .username {
    max-width: 60px;
  }
}

/* 移动端侧边栏完全隐藏 */
@media (max-width: 768px) {
  .sider.ant-layout-sider-collapsed {
    width: 0 !important;
    min-width: 0 !important;
    max-width: 0 !important;
    flex: 0 0 0 !important;
  }
}
</style>
