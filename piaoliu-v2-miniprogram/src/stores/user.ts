import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { User } from '@/types'
import { get, post, put } from '@/utils/request'

export const useUserStore = defineStore('user', () => {
  const user = ref<User | null>(null)
  const token = ref('')

  const login = async (code: string) => {
    const data = await post<{ token: string; user: User; profileComplete: boolean }>('/auth/login', { code })
    token.value = data.token
    user.value = data.user
    uni.setStorageSync('token', data.token)
    uni.setStorageSync('user', JSON.stringify(data.user))
    uni.setStorageSync('profileComplete', data.profileComplete)
    return data.profileComplete
  }

  const getProfile = async () => {
    // 如果没有 token，尝试从存储恢复
    if (!token.value) {
      const savedToken = uni.getStorageSync('token')
      const savedUser = uni.getStorageSync('user')
      if (savedToken) {
        token.value = savedToken
      } else if (savedUser) {
        user.value = JSON.parse(savedUser)
        return
      } else {
        return
      }
    }
    // 有 token 时始终从后端获取最新数据
    const data = await get<User>('/users/profile')
    user.value = data
    uni.setStorageSync('user', JSON.stringify(data))
  }

  const updateProfile = async (profile: Partial<User>) => {
    const data = await put<User>('/users/profile', profile)
    // 合并数据，保留后端可能未返回的字段
    if (user.value) {
      user.value = { ...user.value, ...data }
    } else {
      user.value = data
    }
    uni.setStorageSync('user', JSON.stringify(user.value))
  }

  const updateAddress = async (address: {
    provinceCode: string
    provinceName: string
    cityCode: string
    cityName: string
    districtCode?: string
    districtName?: string
  }) => {
    const data = await put<User>('/users/address', address)
    user.value = data
    uni.setStorageSync('user', JSON.stringify(data))
  }

  const logout = () => {
    user.value = null
    token.value = ''
    uni.removeStorageSync('token')
    uni.removeStorageSync('user')
  }

  return {
    user,
    token,
    login,
    getProfile,
    updateProfile,
    updateAddress,
    logout
  }
})
