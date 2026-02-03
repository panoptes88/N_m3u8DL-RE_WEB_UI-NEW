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

      <a-table
        :columns="columns"
        :data-source="files"
        :pagination="{ pageSize: 10 }"
        :loading="loading"
      >
        <template #bodyCell="{ column, record }">
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
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import { get, del } from '../api'

const loading = ref(false)
const files = ref([])

const columns = [
  { title: '文件名', dataIndex: 'name', key: 'name', ellipsis: true },
  { title: '大小', dataIndex: 'size', key: 'size', width: 100 },
  { title: '修改时间', dataIndex: 'modTime', key: 'modTime', width: 160 },
  { title: '操作', key: 'action', width: 150 }
]

function formatSize(bytes) {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

async function fetchFiles() {
  loading.value = true
  try {
    files.value = await get('/files')
  } catch (err) {
    message.error('获取文件列表失败')
  } finally {
    loading.value = false
  }
}

function downloadFile(name) {
  window.open(`/api/files/download?name=${encodeURIComponent(name)}`, '_blank')
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

onMounted(() => {
  fetchFiles()
})
</script>

<style scoped>
.files-page {
  display: flex;
  flex-direction: column;
}
</style>
