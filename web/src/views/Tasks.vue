<template>
  <div class="tasks-page">
    <!-- 创建任务表单 -->
    <a-card title="创建下载任务">
      <a-form
        ref="formRef"
        :model="formState"
        :rules="formRules"
        layout="vertical"
      >
        <a-row :gutter="16">
          <a-col :span="24">
            <a-form-item label="m3u8 URL" name="url">
              <a-input
                v-model:value="formState.url"
                placeholder="https://example.com/video.m3u8"
                size="large"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="12" :md="8">
            <a-form-item label="输出文件名" name="outputName">
              <a-input
                v-model:value="formState.outputName"
                placeholder="output.mp4"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="8">
            <a-form-item label="线程数" name="threadCount">
              <a-input-number v-model:value="formState.threadCount" :min="1" :max="128" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="8">
            <a-form-item label="重试次数" name="retryCount">
              <a-input-number v-model:value="formState.retryCount" :min="0" :max="100" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="12">
            <a-form-item label="请求头" name="headers">
              <a-input
                v-model:value="formState.headers"
                placeholder='如: Cookie: xxx; User-Agent: xxx'
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12">
            <a-form-item label="Base URL" name="baseUrl">
              <a-input
                v-model:value="formState.baseUrl"
                placeholder="可选，用于补全相对路径"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="16">
          <a-col :xs="24" :sm="12" :md="8">
            <a-form-item label=" " name="delAfterDone">
              <a-checkbox v-model:checked="formState.delAfterDone">
                下载完成后删除临时文件
              </a-checkbox>
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="8">
            <a-form-item label=" " name="binaryMerge">
              <a-checkbox v-model:checked="formState.binaryMerge">
                启用二进制合并
              </a-checkbox>
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="8">
            <a-form-item label=" " name="autoSelect">
              <a-checkbox v-model:checked="formState.autoSelect">
                自动选择最佳轨道
              </a-checkbox>
            </a-form-item>
          </a-col>
        </a-row>

        <!-- 解密选项 -->
        <a-divider>解密选项</a-divider>
        <a-row :gutter="16">
          <a-col :xs="24" :sm="12" :md="8">
            <a-form-item label="解密密钥" name="key">
              <a-input
                v-model:value="formState.key"
                placeholder="KID:KEY 或直接 KEY"
              />
            </a-form-item>
          </a-col>
          <a-col :xs="24" :sm="12" :md="8">
            <a-form-item label="解密引擎" name="decryptionEngine">
              <a-select v-model:value="formState.decryptionEngine">
                <a-select-option value="MP4DECRYPT">MP4DECRYPT</a-select-option>
                <a-select-option value="FFMPEG">FFMPEG</a-select-option>
                <a-select-option value="SHAKA_PACKAGER">SHAKA_PACKAGER</a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>

        <!-- 代理设置 -->
        <a-divider>代理设置</a-divider>
        <a-row :gutter="16">
          <a-col :xs="24" :sm="24" :md="12">
            <a-form-item label="自定义代理" name="customProxy">
              <a-input
                v-model:value="formState.customProxy"
                placeholder="如: http://127.0.0.1:7890"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <!-- 自定义参数 -->
        <a-divider>其他参数</a-divider>
        <a-row :gutter="16">
          <a-col :span="24">
            <a-form-item label="自定义参数" name="customArgs">
              <a-input
                v-model:value="formState.customArgs"
                placeholder="其他命令行参数，如: --log-level DEBUG"
              />
            </a-form-item>
          </a-col>
        </a-row>

        <a-form-item>
          <a-space>
            <a-button type="primary" size="large" :loading="creating" @click="handleCreate">
              创建任务
            </a-button>
            <a-button size="large" @click="resetForm">
              重置
            </a-button>
            <a-checkbox v-model:checked="keepFormAfterCreate">
              创建后保留表单
            </a-checkbox>
          </a-space>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- 任务列表 -->
    <a-card title="下载任务" class="task-list">
      <template #extra>
        <a-space>
          <a-select
            v-model:value="statusFilter"
            style="width: 120px"
            @change="fetchTasks"
          >
            <a-select-option value="">全部</a-select-option>
            <a-select-option value="pending">等待中</a-select-option>
            <a-select-option value="downloading">下载中</a-select-option>
            <a-select-option value="completed">已完成</a-select-option>
            <a-select-option value="failed">失败</a-select-option>
          </a-select>
          <a-button @click="fetchTasks">
            <template #icon><ReloadOutlined /></template>
            刷新
          </a-button>
        </a-space>
      </template>

      <a-table
        :columns="columns"
        :data-source="taskStore.tasks"
        :pagination="{ pageSize: 10 }"
        :loading="taskStore.loading"
        :scroll="{ x: 800 }"
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
            />
          </template>
          <template v-if="column.key === 'createdAt'">
            {{ formatTime(record.created_at) }}
          </template>
          <template v-if="column.key === 'action'">
            <a-space>
              <a-button size="small" @click="viewLog(record)">日志</a-button>
              <a-popconfirm
                title="确定删除此任务？"
                @confirm="handleDelete(record.id)"
              >
                <a-button size="small" danger>删除</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 日志弹窗 -->
    <a-modal
      v-model:open="logModalVisible"
      title="任务日志"
      :footer="null"
      class="log-modal"
    >
      <a-spin :spinning="logLoading">
        <a-textarea
          v-model:value="logContent"
          :rows="20"
          readonly
          style="font-family: monospace"
        />
      </a-spin>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { message } from 'ant-design-vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'
import { useTaskStore } from '../stores/task'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

const taskStore = useTaskStore()

const creating = ref(false)
const formRef = ref(null)
const statusFilter = ref('')
const keepFormAfterCreate = ref(false) // 创建后保留表单
const logModalVisible = ref(false)
const logLoading = ref(false)
const logContent = ref('')

// 表单数据
const formState = reactive({
  url: '',
  outputName: '',
  threadCount: 32,
  retryCount: 15,
  headers: '',
  baseUrl: '',
  delAfterDone: true,
  binaryMerge: false,
  autoSelect: false,
  key: '',
  decryptionEngine: 'MP4DECRYPT',
  customArgs: '',
  customProxy: ''
})

const formRules = {
  url: [{ required: true, message: '请输入m3u8链接' }],
  outputName: [{ required: true, message: '请输入输出文件名' }]
}

const columns = [
  { title: 'ID', dataIndex: 'id', key: 'id', width: 50 },
  { title: 'URL', dataIndex: 'url', key: 'url', ellipsis: true, width: 200 },
  { title: '输出', dataIndex: 'output_name', key: 'outputName', ellipsis: true, width: 80 },
  { title: '状态', dataIndex: 'status', key: 'status', width: 70 },
  { title: '进度', dataIndex: 'progress', key: 'progress', width: 100 },
  { title: '创建时间', dataIndex: 'createdAt', key: 'createdAt', width: 130, responsive: ['lg'] },
  { title: '操作', key: 'action', width: 100 }
]

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

function formatTime(time) {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

async function fetchTasks() {
  const status = statusFilter.value ? `?status=${statusFilter.value}` : ''
  await taskStore.fetchTasks(status)
}

function resetForm() {
  formState.url = ''
  formState.outputName = ''
  formState.threadCount = 32
  formState.retryCount = 15
  formState.headers = ''
  formState.baseUrl = ''
  formState.delAfterDone = true
  formState.binaryMerge = false
  formState.autoSelect = false
  formState.key = ''
  formState.decryptionEngine = 'MP4DECRYPT'
  formState.customArgs = ''
  formState.customProxy = ''
  formRef.value?.clearValidate()
}

async function handleCreate() {
  try {
    await formRef.value.validate()
  } catch {
    return
  }

  creating.value = true
  try {
    // 构建参数字符串
    const args = {
      url: formState.url,
      output_name: formState.outputName,
      thread_count: formState.threadCount,
      retry_count: formState.retryCount,
      headers: formState.headers,
      base_url: formState.baseUrl,
      del_after_done: formState.delAfterDone,
      binary_merge: formState.binaryMerge,
      auto_select: formState.autoSelect,
      key: formState.key,
      decryption_engine: formState.decryptionEngine,
      custom_args: formState.customArgs,
      custom_proxy: formState.customProxy
    }

    await taskStore.createTask(args)
    message.success('任务创建成功')
    // 根据开关决定是否保留表单
    if (!keepFormAfterCreate.value) {
      resetForm()
    }
  } catch (err) {
    message.error(err.response?.data?.error || '创建任务失败')
  } finally {
    creating.value = false
  }
}

async function handleDelete(id) {
  try {
    await taskStore.deleteTask(id)
    message.success('删除成功')
  } catch {
    message.error('删除失败')
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

onMounted(() => {
  // 使用 store 统一管理的轮询（单例模式）
  taskStore.startPolling()
})

onUnmounted(() => {
  // 不停止轮询，因为 Dashboard 页面可能还在使用
  // 轮询由 store 统一管理
})
</script>

<style scoped>
.tasks-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.task-list {
  flex: 1;
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
</style>
