import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Bottle, PickedBottle } from '@/types'
import { get, post, del } from '@/utils/request'

export const useBottleStore = defineStore('bottle', () => {
  const currentBottle = ref<Bottle | null>(null)
  const myBottles = ref<Bottle[]>([])
  const pickedBottles = ref<Bottle[]>([])

  const createBottle = async (data: {
    content: string
    mediaUrls?: string[]
    tags?: string[]
    bottleType: number
    visibleScope: number
  }) => {
    return await post<{ id: number; status: number; message: string }>('/bottles', data)
  }

  const pickBottle = async (gender?: number) => {
    let url = '/bottles/pick'
    if (gender) {
      url += `?gender=${gender}`
    }
    const data = await post<PickedBottle>(url)
    currentBottle.value = data as any
    return data
  }

  const getBottleDetail = async (id: number) => {
    const data = await get<any>(`/bottles/${id}`)
    currentBottle.value = data.bottle
    return data
  }

  const replyBottle = async (id: number, content: string) => {
    return await post<{ conversationId: number; messageId: number }>(`/bottles/${id}/reply`, { content })
  }

  const commentBottle = async (id: number, content: string) => {
    return await post<{ id: number }>(`/bottles/${id}/comment`, { content })
  }

  const getBottleComments = async (id: number) => {
    return await get<any[]>(`/bottles/${id}/comments`)
  }

  const getMyBottles = async (page: number = 1, pageSize: number = 10) => {
    const data = await get<{ list: Bottle[]; total: number; page: number }>(`/bottles/my?page=${page}&pageSize=${pageSize}`)
    return data
  }

  const getPickedBottles = async (page: number = 1, pageSize: number = 10) => {
    const data = await get<{ list: any[]; total: number; page: number }>(`/bottles/picked?page=${page}&pageSize=${pageSize}`)
    return data
  }

  const getBottleCount = async (gender?: number) => {
    let url = '/bottles/count'
    if (gender) {
      url += `?gender=${gender}`
    }
    return await get<{ nearbyCount: number }>(url)
  }

  const watchAd = async () => {
    return await post<{ reward: number }>('/bottles/watch-ad')
  }

  // 把捡到的瓶子扔回海里（从“捞到的”列表移除）
  const throwBackBottle = async (id: number) => {
    return await post(`/bottles/${id}/throw-back`)
  }

  // 删除自己扔出的瓶子
  const deleteBottle = async (id: number) => {
    return await del(`/bottles/${id}`)
  }

  return {
    currentBottle,
    myBottles,
    pickedBottles,
    createBottle,
    pickBottle,
    getBottleDetail,
    replyBottle,
    commentBottle,
    getBottleComments,
    getMyBottles,
    getPickedBottles,
    getBottleCount,
    watchAd,
    throwBackBottle,
    deleteBottle
  }
})
