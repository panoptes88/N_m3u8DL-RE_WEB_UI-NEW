import { defineStore } from 'pinia'
import { ref, watch } from 'vue'

export const useAppStore = defineStore('app', () => {
  // 侧边栏折叠状态
  const collapsed = ref(localStorage.getItem('sidebarCollapsed') === 'true')

  // 主题模式：light / dark
  const theme = ref(localStorage.getItem('theme') || 'light')

  // 监听变化并持久化
  watch(collapsed, (val) => {
    localStorage.setItem('sidebarCollapsed', val)
  })

  watch(theme, (val) => {
    localStorage.setItem('theme', val)
    applyTheme(val)
  })

  // 应用主题到 DOM
  function applyTheme(mode) {
    if (mode === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  }

  // 切换侧边栏
  function toggleCollapsed() {
    collapsed.value = !collapsed.value
  }

  // 切换主题
  function toggleTheme() {
    theme.value = theme.value === 'light' ? 'dark' : 'light'
  }

  // 初始化时应用主题
  applyTheme(theme.value)

  return {
    collapsed,
    theme,
    toggleCollapsed,
    toggleTheme
  }
})
