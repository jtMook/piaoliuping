<template>
  <view class="page">
    <!-- 顶部导航 -->
    <view class="header">
      <view class="header-content">
        <text class="back-icon" @click="goBack">‹</text>
        <view class="chat-title">
          <text class="nickname">{{ otherUser?.nickname || '加载中...' }}</text>
          <text class="subtitle">{{ otherUser?.cityName || '' }} · 来自「{{ bottleInfo?.content?.substring(0, 10) || '...' }}」瓶</text>
        </view>
        <text class="header-action" @click="showMore">⋮</text>
      </view>
    </view>

    <!-- 消息列表 -->
    <scroll-view
      class="message-list"
      scroll-y
      :scroll-top="scrollTop"
      :scroll-with-animation="true"
    >
      <!-- 时间标签 -->
      <view class="time-label">
        <text>今天 23:40</text>
      </view>

      <view class="message-item"
        v-for="message in messages"
        :key="message.id"
        :class="{ self: isSelf(message.senderId) }"
      >
        <view class="avatar" v-if="!isSelf(message.senderId)">
          <view class="avatar-gradient"></view>
        </view>

        <view class="message-content">
          <view class="message-bubble" :class="{ self: isSelf(message.senderId) }">
            <text class="message-text">{{ message.content }}</text>
          </view>
          <text class="message-time" :class="{ self: isSelf(message.senderId) }">{{ formatTime(message.createdAt) }}</text>
        </view>
      </view>
    </scroll-view>

    <!-- 底部输入栏 -->
    <view class="input-bar">
      <view class="input-wrapper">
        <input
          class="chat-input"
          v-model="inputContent"
          placeholder="说点什么..."
          placeholder-class="placeholder"
          :adjust-position="true"
        />
      </view>
      <view class="btn-send" @click="sendMessage">
        <text>发送</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { useUserStore } from '@/stores/user'
import { useMessageStore } from '@/stores/message'
import type { Message, User } from '@/types'

const userStore = useUserStore()
const messageStore = useMessageStore()
const messages = ref<Message[]>([])
const inputContent = ref('')
const scrollTop = ref(0)
const conversationId = ref(0)
const otherUser = ref<{ id: number; nickname: string; avatarUrl: string; cityName: string; gender: number } | null>(null)
const bottleInfo = ref<{ id: number; content: string; bottleType: number } | null>(null)

onMounted(async () => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1] as any
  const id = currentPage?.options?.id || currentPage?.$route?.query?.id

  if (id) {
    conversationId.value = Number(id)
    // 调用接口获取会话详情
    await loadConversationDetail()
  }

  await loadMessages()
})

const loadConversationDetail = async () => {
  try {
    const detail = await messageStore.getConversation(conversationId.value)
    otherUser.value = detail.otherUser
    bottleInfo.value = detail.bottle
  } catch (error) {
    console.error('获取会话详情失败:', error)
  }
}

const loadMessages = async () => {
  try {
    if (conversationId.value) {
      messages.value = await messageStore.getMessages(conversationId.value)
    }
    scrollToBottom()
    await messageStore.getConversations()
    messageStore.updateTabBarBadge()
  } catch (error) {
    console.error(error)
    messages.value = []
  }
}

const sendMessage = async () => {
  if (!inputContent.value.trim()) return

  const content = inputContent.value.trim()
  inputContent.value = ''

  try {
    const newMessage = await messageStore.sendMessage(conversationId.value, content)
    messages.value.push(newMessage)
    scrollToBottom()
  } catch (error) {
    uni.showToast({ title: '发送失败', icon: 'none' })
  }
}

const showMore = () => {
  uni.showActionSheet({
    itemList: ['举报', '拉黑', '清空聊天记录'],
    success: (res) => {
      switch (res.tapIndex) {
        case 0:
          uni.showToast({ title: '已举报', icon: 'none' })
          break
        case 1:
          uni.showToast({ title: '已拉黑', icon: 'none' })
          break
        case 2:
          messages.value = []
          uni.showToast({ title: '已清空', icon: 'none' })
          break
      }
    }
  })
}

const scrollToBottom = () => {
  nextTick(() => {
    scrollTop.value = messages.value.length * 1000
  })
}

const isSelf = (senderId: number) => {
  return senderId === userStore.user?.id
}

const goBack = () => {
  uni.navigateBack()
}

const formatTime = (dateString: string) => {
  const date = new Date(dateString)
  return `${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}
</script>

<style lang="scss" scoped>
.page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: $bg-dark;
}

.header {
  flex-shrink: 0;
  padding-top: var(--status-bar-height, 44px);
  background: $bg-header;
  backdrop-filter: blur(10px);
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: $spacing-sm $spacing-md;
  height: $header-height;
}

.back-icon {
  width: 60rpx;
  font-size: 60rpx;
  color: $text-primary;
  line-height: 1;
}

.chat-title {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
}

.nickname {
  font-size: $font-base;
  font-weight: 600;
  color: $text-primary;
}

.subtitle {
  font-size: $font-xs;
  color: $text-muted;
}

.header-action {
  width: 60rpx;
  text-align: right;
  font-size: $font-xl;
  color: $text-primary;
}

.message-list {
  flex: 1;
  padding: $spacing-md;
  overflow-y: auto;
}

.time-label {
  text-align: center;
  margin-bottom: $spacing-md;

  text {
    font-size: $font-xs;
    color: $text-muted;
    background: $bg-card;
    padding: 4rpx 16rpx;
    border-radius: $radius-sm;
  }
}

.message-item {
  display: flex;
  align-items: flex-start;
  margin-bottom: $spacing-md;

  &.self {
    flex-direction: row-reverse;
  }
}

.avatar {
  width: 72rpx;
  height: 72rpx;
  border-radius: 50%;
  flex-shrink: 0;
  overflow: hidden;
}

.avatar-gradient {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.message-content {
  max-width: 65%;
  margin: 0 $spacing-md;
  display: flex;
  flex-direction: column;
}

.message-bubble {
  padding: $spacing-md;
  background: $bg-card;
  border-radius: 4rpx $radius-lg $radius-lg $radius-lg;
  border: 1px solid $border-color;

  &.self {
    background: linear-gradient(135deg, #d4a843 0%, #c49a38 100%);
    border-color: transparent;
    border-radius: $radius-lg 4rpx $radius-lg $radius-lg;
  }
}

.message-text {
  font-size: $font-base;
  color: $text-primary;
  line-height: 1.6;
  word-break: break-all;

  .message-bubble.self & {
    color: #0a1628;
  }
}

.message-time {
  font-size: $font-xs;
  color: $text-muted;
  margin-top: $spacing-xs;

  &.self {
    text-align: right;
  }
}

.input-bar {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  padding: $spacing-md;
  padding-bottom: calc($spacing-md + env(safe-area-inset-bottom));
  background: $bg-card;
  border-top: 1px solid $border-color;
}

.btn-send {
  background: linear-gradient(135deg, #d4a843 0%, #c49a38 100%);
  border-radius: $radius-md;
  padding: $spacing-xs $spacing-md;
  display: flex;
  align-items: center;
  justify-content: center;

  text {
    font-size: $font-sm;
    color: #0a1628;
    font-weight: 600;
  }
}

.input-wrapper {
  flex: 1;
  background: $bg-input;
  border-radius: $radius-lg;
  padding: $spacing-sm $spacing-md;
  border: 1px solid $border-color;
}

.chat-input {
  width: 100%;
  height: 60rpx;
  color: $text-primary;
  font-size: $font-base;
}

.placeholder {
  color: $text-muted;
}
</style>
