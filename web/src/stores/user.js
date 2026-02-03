import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { get, post } from '../api'

export const useUserStore = defineStore('user', () => {
  const username = ref(localStorage.getItem('username') || '')
  const isLoggedIn = computed(() => !!username.value)

  async function checkLogin() {
    try {
      const res = await get('/user')
      // 只有当明确返回 username 时才认为已登录
      if (res && res.username) {
        username.value = res.username
        localStorage.setItem('username', res.username)
        return true
      }
      username.value = ''
      localStorage.removeItem('username')
      return false
    } catch (err) {
      // 401 未登录是正常情况，不做处理
      username.value = ''
      localStorage.removeItem('username')
      return false
    }
  }

  async function login(usernameVal, password) {
    const res = await post('/auth/login', { username: usernameVal, password })
    username.value = res.username
    localStorage.setItem('username', res.username)
    return res
  }

  function logout() {
    return post('/auth/logout').then(() => {
      username.value = ''
      localStorage.removeItem('username')
    })
  }

  return {
    username,
    isLoggedIn,
    checkLogin,
    login,
    logout
  }
})
