import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 30000,
  withCredentials: true
})

// 请求拦截器
api.interceptors.request.use(
  config => {
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// 响应拦截器 - 返回响应数据
api.interceptors.response.use(
  response => response.data,
  error => {
    return Promise.reject(error)
  }
)

export function get(url, params) {
  return api.get(url, { params })
}

export function post(url, data) {
  return api.post(url, data)
}

export function del(url) {
  return api.delete(url)
}

export function downloadFile(url, filename) {
  const link = document.createElement('a')
  link.href = `/api${url}&filename=${encodeURIComponent(filename)}`
  link.download = filename
  link.click()
}

export default api
