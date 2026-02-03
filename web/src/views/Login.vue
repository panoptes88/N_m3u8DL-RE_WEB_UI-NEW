<template>
  <div class="login-container">
    <a-card class="login-card">
      <template #title>
        <div class="login-title">N_m3u8DL-RE Web UI</div>
      </template>
      <a-form
        :model="formState"
        :rules="rules"
        layout="vertical"
        @finish="handleSubmit"
      >
        <a-form-item label="用户名" name="username">
          <a-input
            v-model:value="formState.username"
            placeholder="请输入用户名"
            size="large"
          >
            <template #prefix>
              <UserOutlined />
            </template>
          </a-input>
        </a-form-item>
        <a-form-item label="密码" name="password">
          <a-input-password
            v-model:value="formState.password"
            placeholder="请输入密码"
            size="large"
            @press-enter="handleSubmit"
          >
            <template #prefix>
              <LockOutlined />
            </template>
          </a-input-password>
        </a-form-item>
        <a-form-item>
          <a-button
            type="primary"
            html-type="submit"
            :loading="loading"
            block
            size="large"
          >
            登录
          </a-button>
        </a-form-item>
      </a-form>
      <div class="login-footer">
        默认用户名: admin | 默认密码: admin123
      </div>
    </a-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { message } from 'ant-design-vue'
import { UserOutlined, LockOutlined } from '@ant-design/icons-vue'
import { useUserStore } from '../stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const loading = ref(false)
const formState = reactive({
  username: '',
  password: ''
})

const rules = {
  username: [{ required: true, message: '请输入用户名' }],
  password: [{ required: true, message: '请输入密码' }]
}

async function handleSubmit() {
  loading.value = true
  try {
    const res = await userStore.login(formState.username, formState.password)
    console.log('登录响应:', res)
    message.success('登录成功')
    // 使用 window.location 强制跳转，而不是 router.push
    window.location.href = '/'
  } catch (err) {
    console.error('登录错误:', err)
    message.error(err.response?.data?.error || err.message || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.login-card {
  width: 400px;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.login-title {
  text-align: center;
  font-size: 20px;
  font-weight: bold;
  color: #1677ff;
}

.login-footer {
  text-align: center;
  color: #999;
  font-size: 12px;
  margin-top: 16px;
}
</style>
