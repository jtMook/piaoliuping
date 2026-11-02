import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Conversation, ConversationDetail, Message } from '@/types'
import { get, post, del } from '@/utils/request'

export const useMessageStore = defineStore('message', () => {
  const conversations = ref<Conversation[]>([])
  const currentConversation = ref<Conversation | null>(null)
  const messages = ref<Message[]>([])

  const totalUnreadCount = computed(() => {
    return conversations.value.reduce((sum, conv) => sum + (conv.unreadCount || 0), 0)
  })

  const getConversations = async () => {
    const data = await get<Conversation[]>('/conversations')
    conversations.value = data || []
    return conversations.value
  }

  const getConversation = async (conversationId: number) => {
    return await get<ConversationDetail>(`/conversations/${conversationId}`)
  }

  const getMessages = async (conversationId: number) => {
    const data = await get<Message[]>(`/conversations/${conversationId}/messages`)
    messages.value = data || []
    return messages.value
  }

  const sendMessage = async (conversationId: number, content: string, msgType: number = 1) => {
    return await post<Message>(`/conversations/${conversationId}/messages`, { content, msgType })
  }

  // 删除会话（仅对自己隐藏，对方不受影响）
  const deleteConversation = async (conversationId: number) => {
    return await del(`/conversations/${conversationId}`)
  }

  const setCurrentConversation = (conversation: Conversation) => {
    currentConversation.value = conversation
  }

  const addMessage = (message: Message) => {
    messages.value.push(message)
  }

  const updateTabBarBadge = () => {
    const count = totalUnreadCount.value
    if (count > 0) {
      uni.setTabBarBadge({
        index: 2,
        text: count > 99 ? '99+' : String(count)
      })
    } else {
      uni.removeTabBarBadge({
        index: 2
      })
    }
  }

  return {
    conversations,
    currentConversation,
    messages,
    totalUnreadCount,
    getConversations,
    getConversation,
    getMessages,
    sendMessage,
    deleteConversation,
    setCurrentConversation,
    addMessage,
    updateTabBarBadge
  }
})
