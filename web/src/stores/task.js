import { defineStore } from 'pinia'
import { ref } from 'vue'
import { get, post, del } from '../api'

export const useTaskStore = defineStore('task', () => {
  const tasks = ref([])
  const loading = ref(false)
  let pollingTimer = null
  const pollingInterval = 5000 // 轮询间隔 5 秒

  // 统一获取任务列表
  async function fetchTasks(status = '') {
    loading.value = true
    try {
      const params = status ? `?status=${status}` : ''
      const newTasks = await get(`/tasks${params}`)
      // 只在数据变化时更新，避免不必要的重渲染
      if (JSON.stringify(tasks.value) !== JSON.stringify(newTasks)) {
        tasks.value = newTasks
      }
    } catch (e) {
      console.error('获取任务列表失败:', e)
    } finally {
      loading.value = false
    }
  }

  // 启动统一轮询（单例模式）
  function startPolling() {
    if (pollingTimer) {
      return // 已经在轮询中，不再创建新的定时器
    }
    fetchTasks()
    pollingTimer = setInterval(() => {
      fetchTasks()
    }, pollingInterval)
  }

  // 停止轮询
  function stopPolling() {
    if (pollingTimer) {
      clearInterval(pollingTimer)
      pollingTimer = null
    }
  }

  async function createTask(args) {
    const res = await post('/tasks', args)
    tasks.value.unshift(res)
    return res
  }

  async function deleteTask(id) {
    await del(`/tasks/${id}`)
    tasks.value = tasks.value.filter(t => t.id !== id)
  }

  async function getTaskProgress(id) {
    return await get(`/tasks/${id}`)
  }

  async function getTaskLog(id) {
    return await get(`/tasks/${id}/log`)
  }

  return {
    tasks,
    loading,
    fetchTasks,
    startPolling,
    stopPolling,
    createTask,
    deleteTask,
    getTaskProgress,
    getTaskLog
  }
})
