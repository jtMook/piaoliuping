import type { UniAppRequestOptions } from '@dcloudio/uni-app'

// #ifdef MP-WEIXIN
// @ts-ignore
let BASE_URL = 'http://localhost:8080/api'
// #endif
// #ifdef H5
// @ts-ignore
let BASE_URL = '/api'
// #endif

interface RequestOptions extends UniAppRequestOptions {
  showLoading?: boolean
  loadingText?: string
}

export function request<T = any>(options: RequestOptions): Promise<T> {
  return new Promise((resolve, reject) => {
    const { showLoading = true, loadingText = '加载中...', url, ...rest } = options
    
    if (showLoading) {
      uni.showLoading({ title: loadingText })
    }
    
    uni.request({
      url: `${BASE_URL}${url}`,
      header: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${uni.getStorageSync('token') || ''}`
      },
      ...rest,
      success: (res) => {
        if (showLoading) {
          uni.hideLoading()
        }
        
        const data = res.data as { code: number; data: T; message: string }
        
        if (data.code === 0) {
          resolve(data.data)
        } else if (data.code === 401) {
          uni.removeStorageSync('token')
          uni.removeStorageSync('user')
          uni.reLaunch({ url: '/pages/login/login' })
          reject(new Error('登录已过期'))
        } else {
          uni.showToast({ title: data.message, icon: 'none' })
          reject(new Error(data.message))
        }
      },
      fail: (err) => {
        if (showLoading) {
          uni.hideLoading()
        }
        uni.showToast({ title: '网络请求失败', icon: 'none' })
        reject(err)
      }
    })
  })
}

export const get = <T = any>(url: string, options?: Omit<RequestOptions, 'url' | 'method'>) => {
  return request<T>({ url, method: 'GET', ...options })
}

export const post = <T = any>(url: string, data?: any, options?: Omit<RequestOptions, 'url' | 'method'>) => {
  return request<T>({ url, method: 'POST', data, ...options })
}

export const put = <T = any>(url: string, data?: any, options?: Omit<RequestOptions, 'url' | 'method'>) => {
  return request<T>({ url, method: 'PUT', data, ...options })
}

export const del = <T = any>(url: string, options?: Omit<RequestOptions, 'url' | 'method'>) => {
  return request<T>({ url, method: 'DELETE', ...options })
}
