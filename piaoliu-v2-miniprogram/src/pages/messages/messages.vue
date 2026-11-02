<template>
  <view class="page">
    <!-- 顶部导航 -->
    <view class="header">
      <view class="header-content">
        <text class="header-title">消息</text>
        <view class="header-actions">
          <text class="action-icon" @click="markAllRead">✓</text>
        </view>
      </view>
    </view>

    <view class="content">
      <!-- 会话列表 -->
      <view class="conversation-list" v-if="conversations.length > 0">
        <view
          class="conversation-item"
          v-for="conversation in conversations"
          :key="conversation.id"
          @click="goToChat(conversation)"
        >
          <view class="avatar">
            <text class="avatar-text">{{ getAvatarText(conversation.otherUser.nickname) }}</text>
            <view class="unread-badge" v-if="conversation.unreadCount > 0">
              <text>{{ conversation.unreadCount }}</text>
            </view>
          </view>

          <view class="conversation-info">
            <view class="info-header">
              <text class="nickname">{{ conversation.otherUser.nickname }}</text>
              <text class="time">{{ formatTime(conversation.lastMessageAt) }}</text>
            </view>
            <text class="last-message">{{ conversation.lastMessage?.content || '暂无消息' }}</text>
          </view>

          <view class="item-delete" @click.stop="handleDelete(conversation)">
            <text class="item-delete-text">🗑️</text>
          </view>
        </view>
      </view>

      <!-- 空状态 -->
      <view class="empty-state" v-else>
        <text class="empty-icon">📭</text>
        <text class="empty-title">还没有消息</text>
        <text class="empty-desc">去捡一个瓶子，开始一段新的对话</text>
        <view class="btn btn-primary" @click="goToPick">去捡瓶子</view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useMessageStore } from '@/stores/message'
import type { Conversation } from '@/types'

const messageStore = useMessageStore()
const conversations = ref<Conversation[]>([])

onMounted(() => {
  loadConversations()
})

onShow(() => {
  loadConversations()
})

const loadConversations = async () => {
  try {
    conversations.value = await messageStore.getConversations()
    messageStore.updateTabBarBadge()
  } catch (error) {
    console.error(error)
    conversations.value = []
  }
}

const goToChat = (conversation: Conversation) => {
  messageStore.setCurrentConversation(conversation)
  uni.navigateTo({ url: `/pages/chat/chat?id=${conversation.id}` })
}

const handleDelete = (conversation: Conversation) => {
  uni.showModal({
    title: '删除会话',
    content: '删除后该对话将不再显示，对方发来新消息时会重新出现。确定删除吗？',
    success: async (res) => {
      if (!res.confirm) return
      try {
        await messageStore.deleteConversation(conversation.id)
        conversations.value = conversations.value.filter((c) => c.id !== conversation.id)
        messageStore.updateTabBarBadge()
        uni.showToast({ title: '已删除', icon: 'success' })
      } catch (error: any) {
        uni.showToast({ title: error.message || '删除失败', icon: 'none' })
      }
    }
  })
}

const goToPick = () => {
  uni.navigateTo({ url: '/pages/pick/pick' })
}

const markAllRead = () => {
  conversations.value.forEach(item => {
    item.unreadCount = 0
  })
  uni.showToast({ title: '已全部已读', icon: 'none' })
}

const getAvatarText = (nickname: string) => {
  return nickname ? nickname.charAt(0) : '?'
}

const formatTime = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  return `${date.getMonth() + 1}/${date.getDate()}`
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: $bg-dark;
}

.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
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

.header-title {
  font-size: $font-xl;
  font-weight: 600;
  color: $text-primary;
}

.header-actions {
  width: 60rpx;
  height: 60rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.action-icon {
  font-size: $font-xl;
  color: $primary-color;
}

.content {
  padding-top: calc(var(--status-bar-height, 44px) + $header-height);
}

.conversation-list {
  padding: $spacing-sm $spacing-md;
}

.conversation-item {
  display: flex;
  align-items: center;
  gap: $spacing-md;
  padding: $spacing-md;
  background: $bg-card;
  border-radius: $radius-lg;
  margin-bottom: $spacing-sm;
  border: 1px solid $border-color;

  .item-delete {
    flex-shrink: 0;
    padding: 8rpx;

    .item-delete-text {
      font-size: 28rpx;
      opacity: 0.6;
    }

    &:active {
      opacity: 0.4;
    }
  }

  &:active {
    background: $bg-card-hover;
  }
}

.avatar {
  position: relative;
  width: 96rpx;
  height: 96rpx;
  background: linear-gradient(135deg, $primary-color 0%, $primary-dark 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.avatar-text {
  font-size: $font-xl;
  font-weight: 600;
  color: $bg-dark;
}

.unread-badge {
  position: absolute;
  top: -4rpx;
  right: -4rpx;
  min-width: 32rpx;
  height: 32rpx;
  background: $error-color;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 8rpx;

  text {
    font-size: 20rpx;
    color: #fff;
  }
}

.conversation-info {
  flex: 1;
  min-width: 0;
}

.info-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: $spacing-xs;
}

.nickname {
  font-size: $font-base;
  font-weight: 600;
  color: $text-primary;
}

.time {
  font-size: $font-xs;
  color: $text-muted;
  flex-shrink: 0;
}

.last-message {
  font-size: $font-sm;
  color: $text-secondary;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding-top: 200rpx;
}

.empty-icon {
  font-size: 120rpx;
  margin-bottom: $spacing-lg;
}

.empty-title {
  font-size: $font-xl;
  font-weight: 600;
  color: $text-primary;
  margin-bottom: $spacing-sm;
}

.empty-desc {
  font-size: $font-sm;
  color: $text-muted;
  margin-bottom: $spacing-xl;
}
</style>
