import { defineStore } from 'pinia'
import { ref } from 'vue'
import { get, post, del } from '../api'

export const useTaskStore = defineStore('task', () => {
  const tasks = ref([])
  const loading = ref(false)

  async function fetchTasks(status = '') {
    loading.value = true
    try {
      const params = status ? `?status=${status}` : ''
      tasks.value = await get(`/tasks${params}`)
    } finally {
      loading.value = false
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
    createTask,
    deleteTask,
    getTaskProgress,
    getTaskLog
  }
})
