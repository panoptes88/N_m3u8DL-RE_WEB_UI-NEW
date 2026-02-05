<template>
  <div class="dashboard">
    <a-row :gutter="16" class="stats-row">
      <a-col :xs="8" :sm="8">
        <a-card>
          <a-statistic
            title="等待中任务"
            :value="pendingCount"
            :value-style="{ color: '#faad14' }"
          >
            <template #prefix>
              <ClockCircleOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card>
          <a-statistic
            title="下载中任务"
            :value="downloadingCount"
            :value-style="{ color: '#1677ff' }"
          >
            <template #prefix>
              <CloudDownloadOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
      <a-col :span="8">
        <a-card>
          <a-statistic
            title="已完成任务"
            :value="completedCount"
            :value-style="{ color: '#52c41a' }"
          >
            <template #prefix>
              <CheckCircleOutlined />
            </template>
          </a-statistic>
        </a-card>
      </a-col>
    </a-row>

    <a-card title="快速下载" class="quick-download">
      <a-form layout="vertical" :model="quickDownloadForm">
        <a-form-item
          label="m3u8 URL"
          :rules="[{ required: true, message: '请输入m3u8链接' }]"
        >
          <a-input
            v-model:value="quickDownloadForm.url"
            placeholder="https://example.com/video.m3u8"
            size="large"
          />
        </a-form-item>
        <a-row :gutter="12">
          <a-col :xs="24" :sm="12">
            <a-form-item label="输出名称（可选）">
              <a-input
                v-model:value="quickDownloadForm.outputName"
                placeholder="output.mp4"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" class="btn-col">
            <a-form-item label=" ">
              <a-button
                type="primary"
                size="large"
                block
                :loading="quickDownloadLoading"
                @click="handleQuickDownload"
              >
                开始下载
              </a-button>
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-card>

    <a-card title="最近任务" class="recent-tasks">
      <template #extra>
        <a-button type="link" @click="$router.push('/tasks')">查看全部</a-button>
      </template>
      <a-table
        :columns="columns"
        :data-source="recentTasks"
        :pagination="false"
        size="small"
        :loading="loading"
        :scroll="{ x: 500 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="getStatusColor(record.status)">
              {{ getStatusText(record.status) }}
            </a-tag>
          </template>
          <template v-if="column.key === 'progress'">
            <a-progress
              :percent="record.progress"
              :status="record.status === 'failed' ? 'exception' : 'active'"
              size="small"
            />
          </template>
          <template v-if="column.key === 'action'">
            <a-space>
              <a-button
                size="small"
                @click="viewLog(record)"
              >
                日志
              </a-button>
              <a-popconfirm
                title="确定删除此任务？"
                @confirm="deleteTask(record.id)"
              >
                <a-button size="small" danger>删除</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal
      v-model:open="logModalVisible"
      title="任务日志"
      :footer="null"
      class="log-modal"
    >
      <a-spin :spinning="logLoading">
        <a-textarea
          v-model:value="logContent"
          :rows="15"
          readonly
          style="font-family: monospace"
        />
      </a-spin>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { message } from 'ant-design-vue'
import {
  ClockCircleOutlined,
  CloudDownloadOutlined,
  CheckCircleOutlined
} from '@ant-design/icons-vue'
import { useTaskStore } from '../stores/task'

const taskStore = useTaskStore()

const loading = ref(false)
const quickDownloadLoading = ref(false)
const quickDownloadForm = ref({ url: '', outputName: '' })
const logModalVisible = ref(false)
const logLoading = ref(false)
const logContent = ref('')

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 50 },
  { title: 'URL', dataIndex: 'url', key: 'url', ellipsis: true },
  { title: '状态', dataIndex: 'status', key: 'status', width: 70 },
  { title: '进度', dataIndex: 'progress', key: 'progress', width: 100 },
  { title: '操作', key: 'action', width: 100 }
]

const recentTasks = computed(() => taskStore.tasks.slice(0, 5))
const pendingCount = computed(() => taskStore.tasks.filter(t => t.status === 'pending').length)
const downloadingCount = computed(() => taskStore.tasks.filter(t => t.status === 'downloading').length)
const completedCount = computed(() => taskStore.tasks.filter(t => t.status === 'completed').length)

function getStatusColor(status) {
  const colors = {
    pending: 'orange',
    downloading: 'blue',
    completed: 'green',
    failed: 'red',
    interrupted: 'gold'
  }
  return colors[status] || 'default'
}

function getStatusText(status) {
  const texts = {
    pending: '等待中',
    downloading: '下载中',
    completed: '已完成',
    failed: '下载失败',
    interrupted: '已中断'
  }
  return texts[status] || status
}

async function handleQuickDownload() {
  if (!quickDownloadForm.value.url) {
    message.error('请输入m3u8链接')
    return
  }

  quickDownloadLoading.value = true
  try {
    await taskStore.createTask({
      url: quickDownloadForm.value.url,
      output_name: quickDownloadForm.value.outputName || ''
    })
    message.success('任务已创建')
    quickDownloadForm.value = { url: '', outputName: '' }
  } catch (err) {
    message.error(err.response?.data?.error || '创建任务失败')
  } finally {
    quickDownloadLoading.value = false
  }
}

async function viewLog(task) {
  logModalVisible.value = true
  logLoading.value = true
  logContent.value = ''

  try {
    const res = await taskStore.getTaskLog(task.id)
    logContent.value = res.log || '暂无日志'
  } catch {
    message.error('获取日志失败')
  } finally {
    logLoading.value = false
  }
}

async function deleteTask(id) {
  try {
    await taskStore.deleteTask(id)
    message.success('删除成功')
  } catch {
    message.error('删除失败')
  }
}

onMounted(() => {
  // 使用 store 统一管理的轮询（单例模式）
  taskStore.startPolling()
})

onUnmounted(() => {
  // 不停止轮询，因为 Tasks 页面可能还在使用
  // 轮询由 store 统一管理
})
</script>

<style scoped>
.dashboard {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.stats-row {
  margin-bottom: 16px;
}

.btn-col {
  display: flex;
  align-items: flex-end;
}

@media (max-width: 576px) {
  .stats-row :deep(.ant-card-body) {
    padding: 12px;
  }
  .stats-row :deep(.ant-statistic-title) {
    font-size: 12px;
  }
  .stats-row :deep(.ant-statistic-content) {
    font-size: 20px;
  }
}

.quick-download {
  margin-bottom: 16px;
}

.recent-tasks {
  flex: 1;
}

.size-text {
  font-size: 12px;
  color: #999;
  margin-left: 4px;
}

.log-modal {
  max-width: calc(100vw - 32px);
}

.log-modal :deep(.ant-modal-body) {
  max-height: 70vh;
  overflow-y: auto;
}

.log-modal :deep(.ant-modal-content) {
  max-width: 800px;
  margin: 0 auto;
}

@media (max-width: 576px) {
  .recent-tasks :deep(.ant-table-cell) {
    padding: 8px !important;
  }
}
</style>
