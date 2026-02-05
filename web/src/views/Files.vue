<template>
  <div class="files-page">
    <a-card title="文件管理">
      <template #extra>
        <a-space>
          <a-button @click="fetchFiles">
            <template #icon><ReloadOutlined /></template>
            刷新
          </a-button>
        </a-space>
      </template>

      <div class="table-header" v-if="selectedRowKeys.length > 0">
        <span>已选择 {{ selectedRowKeys.length }} 个文件</span>
        <a-button type="primary" danger @click="batchDelete">
          批量删除
        </a-button>
      </div>

      <a-table
        :columns="columns"
        :data-source="files"
        :pagination="paginationConfig"
        :loading="loading"
        :row-selection="rowSelection"
        :row-key="record => record.name"
        :scroll="{ x: 600 }"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <a-space>
              <PlayCircleOutlined v-if="isVideoFile(record.name)" class="play-icon" @click="playVideo(record)" />
              <FileOutlined v-else />
              <span>{{ record.name }}</span>
            </a-space>
          </template>
          <template v-if="column.key === 'size'">
            {{ formatSize(record.size) }}
          </template>
          <template v-if="column.key === 'modTime'">
            {{ record.modTime }}
          </template>
          <template v-if="column.key === 'action'">
            <a-space>
              <a-button size="small" @click="downloadFile(record.name)">
                下载
              </a-button>
              <a-popconfirm
                title="确定删除此文件？"
                @confirm="deleteFile(record.name)"
              >
                <a-button size="small" danger>删除</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 视频播放弹窗 -->
    <a-modal
      v-model:open="videoModalVisible"
      :title="currentVideoName"
      :footer="null"
      @cancel="closeVideo"
      class="video-modal"
    >
      <div class="video-container">
        <video
          ref="videoPlayer"
          controls
          class="video-player"
          @error="handleVideoError"
        >
          您的浏览器不支持视频播放
        </video>
      </div>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { message } from 'ant-design-vue'
import { ReloadOutlined, PlayCircleOutlined, FileOutlined } from '@ant-design/icons-vue'
import { get, del } from '../api'

const loading = ref(false)
const files = ref([])
const videoModalVisible = ref(false)
const currentVideoName = ref('')
const currentVideoUrl = ref('')
const videoPlayer = ref(null)
const selectedRowKeys = ref([])
const pageSize = ref(10)

const videoExtensions = ['.mp4', '.mkv', '.avi', '.mov', '.webm', '.flv', '.wmv', '.m4v', '.3gp']

const paginationConfig = computed(() => ({
  pageSize: pageSize.value,
  showSizeChanger: true,
  pageSizeOptions: ['10', '20', '50', '100', '200'],
  showTotal: (total) => `共 ${total} 个文件`
}))

const columns = [
  {
    title: '文件名',
    dataIndex: 'name',
    key: 'name',
    ellipsis: true,
    sorter: (a, b) => a.name.localeCompare(b.name),
    width: 150
  },
  {
    title: '大小',
    dataIndex: 'size',
    key: 'size',
    width: 80,
    sorter: (a, b) => a.size - b.size
  },
  {
    title: '修改时间',
    dataIndex: 'modTime',
    key: 'modTime',
    width: 120,
    sorter: (a, b) => new Date(a.modTime) - new Date(b.modTime),
    responsive: ['lg']
  },
  { title: '操作', key: 'action', width: 100 }
]

const rowSelection = computed(() => ({
  selectedRowKeys: selectedRowKeys.value,
  onChange: (keys) => {
    selectedRowKeys.value = keys
  }
}))

function formatSize(bytes) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

function isVideoFile(filename) {
  const ext = filename.substring(filename.lastIndexOf('.')).toLowerCase()
  return videoExtensions.includes(ext)
}

async function fetchFiles() {
  loading.value = true
  try {
    files.value = await get('/files')
    // 清空选择
    selectedRowKeys.value = []
  } catch (err) {
    message.error('获取文件列表失败')
  } finally {
    loading.value = false
  }
}

function downloadFile(name) {
  window.open(`/api/files/download?name=${encodeURIComponent(name)}`, '_blank')
}

function playVideo(record) {
  currentVideoName.value = record.name
  currentVideoUrl.value = `/api/files/download?name=${encodeURIComponent(record.name)}`
  videoModalVisible.value = true

  // 等待 DOM 更新后设置视频源
  setTimeout(() => {
    if (videoPlayer.value) {
      videoPlayer.value.src = currentVideoUrl.value
      videoPlayer.value.load()
    }
  }, 100)
}

function closeVideo() {
  if (videoPlayer.value) {
    // 先移除事件监听，避免触发 error
    videoPlayer.value.onerror = null
    videoPlayer.value.pause()
    videoPlayer.value.src = ''
  }
  currentVideoName.value = ''
  currentVideoUrl.value = ''
}

function handleVideoError() {
  message.error('视频加载失败，请检查文件是否完整')
}

async function deleteFile(name) {
  try {
    await del(`/files/${encodeURIComponent(name)}`)
    message.success('删除成功')
    fetchFiles()
  } catch (err) {
    message.error(err.response?.data?.error || '删除失败')
  }
}

async function batchDelete() {
  if (selectedRowKeys.value.length === 0) {
    return
  }

  const names = selectedRowKeys.value.filter(name => {
    // 过滤掉目录
    const file = files.value.find(f => f.name === name)
    return file && !file.isDir
  })

  if (names.length === 0) {
    message.warning('没有可删除的文件')
    return
  }

  try {
    // 逐个删除
    let success = 0
    let failed = 0
    for (const name of names) {
      try {
        await del(`/files/${encodeURIComponent(name)}`)
        success++
      } catch {
        failed++
      }
    }

    if (success > 0) {
      message.success(`成功删除 ${success} 个文件`)
    }
    if (failed > 0) {
      message.warning(`删除失败 ${failed} 个文件`)
    }

    selectedRowKeys.value = []
    fetchFiles()
  } catch (err) {
    message.error('批量删除失败')
  }
}

onMounted(() => {
  fetchFiles()
})
</script>

<style scoped>
.files-page {
  display: flex;
  flex-direction: column;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  margin-bottom: 16px;
  background: #fafafa;
  border-radius: 8px;
  border: 1px solid #f0f0f0;
}

.play-icon {
  color: #1677ff;
  cursor: pointer;
  font-size: 18px;
  transition: color 0.3s;
}

.play-icon:hover {
  color: #4096ff;
}

.video-container {
  display: flex;
  justify-content: center;
  background: #000;
  border-radius: 4px;
  overflow: hidden;
}

.video-player {
  max-width: 100%;
  max-height: 450px;
  width: 100%;
}

.video-modal :deep(.ant-modal-body) {
  padding: 8px;
}

@media (max-width: 576px) {
  .files-page :deep(.ant-table-cell) {
    padding: 8px !important;
  }
}
</style>
