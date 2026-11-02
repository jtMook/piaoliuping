<template>
  <view class="page">
    <!-- 顶部导航 -->
    <view class="header">
      <view class="header-content">
        <text class="back-icon" @click="goBack">‹</text>
        <text class="header-title">瓶中信</text>
        <text class="header-action" @click="showMore">⋮</text>
      </view>
    </view>

    <scroll-view class="content" scroll-y v-if="bottle">
      <!-- 瓶子卡片 -->
      <view class="bottle-card">
        <!-- 发送者信息区 -->
        <view class="sender-section">
          <view class="sender-avatar-wrap">
            <view class="sender-avatar">
              <image
                v-if="bottle.user?.avatarUrl"
                class="sender-avatar-image"
                :src="bottle.user.avatarUrl"
                mode="aspectFill"
              />
              <text v-else class="avatar-text">{{ bottle.user?.nickname?.charAt(0) || (isOwner ? '我' : '匿') }}</text>
            </view>
          </view>
          <view class="sender-main">
            <view class="sender-top-row">
              <text class="sender-name">{{ bottle.user?.nickname || (isOwner ? '我' : '匿名用户') }}</text>
              <text class="sender-time">{{ formatTime(bottle.createdAt) }}</text>
            </view>
            <view class="sender-tags" v-if="!isOwner">
              <view class="tag-age" v-if="senderAge">
                <text class="tag-age-text">{{ senderAge }}岁</text>
              </view>
              <view class="tag-gender" :class="bottle.user?.gender === 1 ? 'male' : bottle.user?.gender === 2 ? 'female' : ''" v-if="bottle.user?.gender">
                <text class="tag-gender-text">{{ bottle.user?.gender === 1 ? '♂' : '♀' }}</text>
              </view>
            </view>
            <view class="sender-location">
              <text class="loc-icon">📍</text>
              <text>{{ bottle.user?.cityName || bottle.cityName || '未知' }}</text>
            </view>
          </view>
          <view class="share-btn" @click="shareBottle">
            <text class="share-icon">↗</text>
          </view>
        </view>

        <!-- 瓶子内容 -->
        <view class="bottle-content-section">
          <text class="bottle-content-text">{{ bottle.content }}</text>
          <!-- 标签显示 -->
          <view class="bottle-tags" v-if="bottle.tags && bottle.tags.length > 0">
            <view class="bottle-tag" v-for="tag in bottle.tags" :key="tag">
              <text class="bottle-tag-text">{{ tag }}</text>
            </view>
          </view>
        </view>

        <!-- 操作按钮 -->
        <view class="action-buttons">
          <view class="action-btn action-btn-reply" :class="{ disabled: isOwner }" @click="showReplyInput">
            <text class="action-btn-text">回复</text>
            <view class="comment-count-badge" v-if="comments.length > 0">
              <text class="comment-count-text">{{ comments.length }}</text>
            </view>
          </view>
          <view class="action-btn action-btn-private" @click="goToChat" v-if="hasReplied || conversationId > 0">
            <text class="action-btn-text">私聊</text>
          </view>
          <view class="action-btn action-btn-throw" @click="throwBack" v-if="!isOwner">
            <text class="action-btn-text">扔回海里</text>
          </view>
          <view class="action-btn action-btn-delete" @click="handleDelete" v-if="isOwner">
            <text class="action-btn-text">删除瓶子</text>
          </view>
        </view>
      </view>

      <!-- 统计信息 -->
      <view class="stats-bar">
        <text class="stats-text">有 <text class="stats-num">{{ comments.length }}</text> 人给 TA 留言</text>
      </view>

      <text class="stats-hint">留言通常可以更快得到回复</text>

      <!-- 留言列表 -->
      <view class="comments-section" id="comments-section">
        <view class="comments-header">
          <text class="comments-header-text">留言区</text>
        </view>

        <view class="no-comments" v-if="comments.length === 0">
          <text class="no-comments-text">暂无留言，快来抢沙发吧</text>
        </view>

        <view class="comment-item" v-for="comment in comments" :key="comment.id">
          <view class="comment-avatar">
            <image
              v-if="comment.user?.avatarUrl"
              class="comment-avatar-image"
              :src="comment.user.avatarUrl"
              mode="aspectFill"
            />
            <text v-else class="comment-avatar-text">{{ comment.user?.nickname?.charAt(0) || '?' }}</text>
          </view>
          <view class="comment-body">
            <view class="comment-user-row">
              <text class="comment-username">{{ comment.user?.nickname || '匿名用户' }}</text>
              <view class="comment-age-tag" v-if="comment.user?.birthday">
                <text class="comment-age-text">{{ getAge(comment.user.birthday) }}岁</text>
              </view>
              <text class="comment-time">{{ formatTime(comment.createdAt) }}</text>
            </view>
            <text class="comment-content">{{ comment.content }}</text>
          </view>
        </view>
      </view>

      <!-- 底部安全区 -->
      <view class="bottom-safe-area"></view>
    </scroll-view>

    <!-- 底部留言输入栏 -->
    <view class="comment-input-bar" v-if="!isOwner">
      <textarea
        class="comment-input"
        v-model="replyContent"
        placeholder="请文明友善发言"
        placeholder-class="comment-input-placeholder"
        :maxlength="200"
        :auto-height="false"
      />
      <view class="comment-send-btn" :class="{ disabled: !replyContent.trim() }" @click="sendComment">
        <text class="comment-send-text">发送</text>
      </view>
    </view>

    <!-- 自己的瓶子提示 -->
    <view class="owner-hint-bar" v-else>
      <text class="owner-hint-text">这是你扔的瓶子，无法回复</text>
    </view>

    <!-- 回复弹窗 -->
    <view class="reply-modal" v-if="showReply" @click.self="showReply = false">
      <view class="reply-modal-content">
        <text class="reply-title">回复这个瓶子</text>
        <textarea
          class="reply-textarea"
          v-model="replyContent"
          placeholder="说点什么..."
          maxlength="500"
          :auto-height="true"
        />
        <view class="reply-modal-actions">
          <view class="btn-cancel" @click="showReply = false">
            <text>取消</text>
          </view>
          <view class="btn-send-reply" @click="sendReply">
            <text>发送</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 扔回海里动画遮罩 -->
    <view class="throw-overlay" v-if="isThrowing">
      <view class="ripple ripple-1"></view>
      <view class="ripple ripple-2"></view>
      <view class="ripple ripple-3"></view>
      <text class="throw-text">已扔回海里</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useBottleStore } from '@/stores/bottle'
import { useUserStore } from '@/stores/user'
import { post } from '@/utils/request'
import type { Bottle } from '@/types'

const bottleStore = useBottleStore()
const userStore = useUserStore()
const bottle = ref<Bottle | null>(null)
const comments = ref<any[]>([])
const showReply = ref(false)
const replyContent = ref('')
const hasReplied = ref(false)
const conversationId = ref(0)
const isOwner = ref(false)
const hasOtherReply = ref(false)
const isThrowing = ref(false)

const senderAge = computed(() => {
  if (isOwner.value) return ''
  const birthday = bottle.value?.user?.birthday || bottle.value?.user?.['birthday']
  if (!birthday) return ''
  return getAge(birthday)
})

onMounted(() => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1] as any
  const id = currentPage?.options?.id || currentPage?.$route?.query?.id
  if (id) {
    loadBottleDetail(Number(id))
  }
})

const loadBottleDetail = async (id: number) => {
  try {
    const data = await bottleStore.getBottleDetail(id)
    bottle.value = data.bottle
    isOwner.value = data.isOwner || false
    hasOtherReply.value = data.hasOtherReply || false
    conversationId.value = data.conversationId || 0

    // 加载评论
    await loadComments(id)
  } catch (error) {
    console.error(error)
    uni.showToast({ title: '瓶子不存在', icon: 'none' })
  }
}

const loadComments = async (id: number) => {
  try {
    const data = await bottleStore.getBottleComments(id)
    comments.value = data || []
  } catch (error) {
    console.error('加载评论失败:', error)
  }
}

const showMore = () => {
  uni.showActionSheet({
    itemList: ['举报', '拉黑'],
    success: (res) => {
      if (res.tapIndex === 0) {
        uni.showToast({ title: '已举报', icon: 'none' })
      } else {
        uni.showToast({ title: '已拉黑', icon: 'none' })
      }
    }
  })
}

const throwBack = async () => {
  if (!bottle.value) return
  try {
    await bottleStore.throwBackBottle(bottle.value.id)
    isThrowing.value = true
    setTimeout(() => {
      uni.showToast({ title: '已扔回海里', icon: 'none' })
      setTimeout(() => {
        uni.navigateBack()
      }, 800)
    }, 1200)
  } catch (error: any) {
    uni.showToast({ title: error.message || '操作失败', icon: 'none' })
  }
}

const handleDelete = () => {
  if (!bottle.value) return
  uni.showModal({
    title: '删除瓶子',
    content: '删除后不可恢复，确定删除吗？',
    success: async (res) => {
      if (!res.confirm) return
      try {
        await bottleStore.deleteBottle(bottle.value!.id)
        uni.showToast({ title: '已删除', icon: 'success' })
        setTimeout(() => {
          uni.navigateBack()
        }, 600)
      } catch (error: any) {
        uni.showToast({ title: error.message || '删除失败', icon: 'none' })
      }
    }
  })
}

const handleStamp = () => {
  uni.showToast({ title: '已盖戳', icon: 'success' })
}

const handlePickAgain = () => {
  uni.navigateBack()
}

const scrollToComment = () => {
  showReply.value = true
}

const showReplyInput = () => {
  showReply.value = true
}

const sendComment = async () => {
  if (!replyContent.value.trim()) {
    uni.showToast({ title: '请输入留言内容', icon: 'none' })
    return
  }
  if (!bottle.value) return

  try {
    await bottleStore.commentBottle(bottle.value.id, replyContent.value)
    uni.showToast({ title: '留言成功', icon: 'success' })
    replyContent.value = ''
    // 重新加载评论
    await loadComments(bottle.value.id)
  } catch (error: any) {
    uni.showToast({ title: error.message || '留言失败', icon: 'none' })
  }
}

const sendReply = async () => {
  if (!replyContent.value.trim()) {
    uni.showToast({ title: '请输入回复内容', icon: 'none' })
    return
  }

  try {
    const data = await post<any>(`/bottles/${bottle.value?.id}/reply`, {
      content: replyContent.value.trim()
    })
    hasReplied.value = true
    conversationId.value = data.conversationId
    showReply.value = false
    replyContent.value = ''
    uni.showToast({ title: '回复成功', icon: 'none' })
  } catch (error) {
    uni.showToast({ title: '回复失败', icon: 'none' })
  }
}

const goToChat = () => {
  uni.navigateTo({ url: `/pages/chat/chat?id=${conversationId.value}` })
}

const shareBottle = () => {
  uni.showToast({ title: '分享功能开发中', icon: 'none' })
}

const goBack = () => {
  uni.navigateBack()
}

const getAge = (birthday: string) => {
  if (!birthday) return 0
  const birth = new Date(birthday)
  const now = new Date()
  let age = now.getFullYear() - birth.getFullYear()
  const m = now.getMonth() - birth.getMonth()
  if (m < 0 || (m === 0 && now.getDate() < birth.getDate())) {
    age--
  }
  return age > 0 ? age : 0
}

const formatTime = (dateString: string) => {
  if (!dateString) return '刚刚'
  const date = new Date(dateString)
  if (isNaN(date.getTime())) return '刚刚'
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  return `${Math.floor(hours / 24)}天前`
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: $bg-dark;
  padding-bottom: 120rpx;
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

.back-icon {
  width: 60rpx;
  font-size: 60rpx;
  color: $text-primary;
  line-height: 1;
}

.header-title {
  font-size: $font-xl;
  font-weight: 600;
  color: $text-primary;
}

.header-action {
  width: 60rpx;
  text-align: right;
  font-size: $font-xl;
  color: $text-primary;
}

.content {
  padding-top: calc(var(--status-bar-height, 44px) + $header-height);
  min-height: 100vh;
}

// 瓶子卡片
.bottle-card {
  margin: $spacing-md;
  background: linear-gradient(180deg, #1e2d3d 0%, #162236 100%);
  border-radius: $radius-xl;
  border: 1px solid rgba(212, 168, 67, 0.15);
  padding: $spacing-lg;
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 4rpx;
    background: linear-gradient(90deg, transparent, $primary-color, transparent);
  }
}

// 发送者信息区
.sender-section {
  display: flex;
  align-items: flex-start;
  gap: $spacing-md;
  margin-bottom: $spacing-lg;
}

.sender-avatar-wrap {
  flex-shrink: 0;
}

.sender-avatar {
  width: 96rpx;
  height: 96rpx;
  background: linear-gradient(135deg, $primary-color 0%, #c49a38 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 3rpx solid rgba(212, 168, 67, 0.3);
  overflow: hidden;
}

.sender-avatar-image {
  width: 100%;
  height: 100%;
}

.avatar-text {
  font-size: $font-xl;
  font-weight: 700;
  color: #0a1628;
}

.sender-main {
  flex: 1;
  min-width: 0;
}

.sender-top-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 6rpx;
}

.sender-name {
  font-size: $font-lg;
  font-weight: 700;
  color: $text-primary;
}

.sender-time {
  font-size: $font-xs;
  color: $text-muted;
  flex-shrink: 0;
}

.sender-tags {
  display: flex;
  gap: $spacing-xs;
  margin-bottom: 6rpx;
}

.tag-age {
  padding: 2rpx 12rpx;
  background: rgba(212, 168, 67, 0.2);
  border-radius: 20rpx;
}

.tag-age-text {
  font-size: $font-xs;
  color: $primary-color;
}

.tag-gender {
  padding: 2rpx 12rpx;
  border-radius: 20rpx;

  &.male {
    background: rgba(100, 181, 246, 0.2);
  }

  &.female {
    background: rgba(240, 98, 146, 0.2);
  }
}

.tag-gender-text {
  font-size: $font-xs;
  color: $text-primary;
}

.sender-location {
  display: flex;
  align-items: center;
  gap: 4rpx;
  font-size: $font-xs;
  color: $text-muted;
}

.loc-icon {
  font-size: $font-xs;
}

.share-btn {
  flex-shrink: 0;
  width: 56rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 50%;

  &:active {
    background: rgba(255, 255, 255, 0.12);
  }
}

.share-icon {
  font-size: $font-lg;
  color: $text-secondary;
}

// 瓶子内容
.bottle-content-section {
  background: rgba(255, 255, 255, 0.04);
  border-radius: $radius-lg;
  padding: $spacing-lg;
  margin-bottom: $spacing-lg;
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.bottle-content-text {
  font-size: $font-lg;
  color: $text-primary;
  line-height: 1.8;
  word-break: break-word;
}

// 瓶子标签
.bottle-tags {
  display: flex;
  flex-wrap: wrap;
  gap: $spacing-xs;
  margin-top: $spacing-md;
}

.bottle-tag {
  padding: 4rpx 16rpx;
  background: rgba(212, 168, 67, 0.1);
  border: 1px solid rgba(212, 168, 67, 0.2);
  border-radius: $radius-sm;
}

.bottle-tag-text {
  font-size: $font-xs;
  color: $primary-color;
}

// 操作按钮
.action-buttons {
  display: flex;
  align-items: center;
  gap: $spacing-md;
}

.action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8rpx;
  padding: $spacing-sm $spacing-lg;
  border-radius: 40rpx;
  position: relative;

  &:active {
    opacity: 0.8;
    transform: scale(0.96);
  }
}

.action-btn-reply {
  flex: 1;
  background: rgba(212, 168, 67, 0.12);
  border: 1px solid rgba(212, 168, 67, 0.3);

  &.disabled {
    opacity: 0.4;
    pointer-events: none;
  }
}

.action-btn-private {
  flex: 1;
  background: rgba(100, 181, 246, 0.12);
  border: 1px solid rgba(100, 181, 246, 0.3);
}

.action-btn-throw {
  flex: 1;
  background: rgba(240, 98, 146, 0.12);
  border: 1px solid rgba(240, 98, 146, 0.3);
}

.action-btn-delete {
  flex: 1;
  background: rgba(239, 83, 80, 0.12);
  border: 1px solid rgba(239, 83, 80, 0.3);
}

.action-btn-icon {
  font-size: $font-md;
}

.action-btn-text {
  font-size: $font-base;
  color: $primary-color;
  font-weight: 600;
}

.comment-count-badge {
  position: absolute;
  top: -8rpx;
  right: -8rpx;
  min-width: 32rpx;
  height: 32rpx;
  background: $primary-color;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 8rpx;
}

.comment-count-text {
  font-size: 20rpx;
  color: #0a1628;
  font-weight: 700;
}

// 统计信息
.stats-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: $spacing-md $spacing-lg;
  margin: 0 $spacing-md;
  background: rgba(255, 255, 255, 0.04);
  border-radius: $radius-lg;
}

.stats-text {
  font-size: $font-sm;
  color: $text-muted;
}

.stats-num {
  color: $primary-color;
  font-weight: 700;
}

.stats-hint {
  display: block;
  text-align: center;
  font-size: $font-xs;
  color: $text-muted;
  padding: $spacing-sm $spacing-md;
  margin: 0 $spacing-md;
}

// 留言区
.comments-section {
  margin: $spacing-md;
  background: rgba(255, 255, 255, 0.03);
  border-radius: $radius-xl;
  padding: $spacing-lg;
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.comments-header {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  margin-bottom: $spacing-lg;
}

.comments-header-icon {
  font-size: $font-lg;
}

.comments-header-text {
  font-size: $font-base;
  font-weight: 600;
  color: $text-primary;
}

.no-comments {
  text-align: center;
  padding: $spacing-xl 0;
}

.no-comments-text {
  font-size: $font-sm;
  color: $text-muted;
}

.comment-item {
  display: flex;
  gap: $spacing-md;
  padding: $spacing-md 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);

  &:last-child {
    border-bottom: none;
  }
}

.comment-avatar {
  width: 64rpx;
  height: 64rpx;
  background: linear-gradient(135deg, $primary-color 0%, #c49a38 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  overflow: hidden;
}

.comment-avatar-image {
  width: 100%;
  height: 100%;
}

.comment-avatar-text {
  font-size: $font-base;
  font-weight: 700;
  color: #0a1628;
}

.comment-body {
  flex: 1;
  min-width: 0;
}

.comment-user-row {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  margin-bottom: 8rpx;
  flex-wrap: wrap;
}

.comment-username {
  font-size: $font-sm;
  font-weight: 600;
  color: $text-primary;
}

.comment-age-tag {
  padding: 2rpx 10rpx;
  background: rgba(100, 181, 246, 0.2);
  border-radius: 16rpx;
}

.comment-age-text {
  font-size: 20rpx;
  color: #64b5f6;
}

.comment-time {
  font-size: $font-xs;
  color: $text-muted;
  margin-left: auto;
}

.comment-content {
  font-size: $font-base;
  color: $text-secondary;
  line-height: 1.6;
  word-break: break-word;
}

// 底部留言输入栏
.comment-input-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 999;
  background: #1a2332;
  border-top: 1px solid rgba(212, 168, 67, 0.2);
  padding: $spacing-md;
  padding-bottom: calc($spacing-md + env(safe-area-inset-bottom));
  display: flex;
  align-items: flex-end;
  gap: $spacing-sm;
  box-shadow: 0 -4rpx 24rpx rgba(0, 0, 0, 0.3);
}

.comment-input {
  flex: 1;
  min-height: 72rpx;
  max-height: 160rpx;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 36rpx;
  padding: $spacing-sm $spacing-md;
  font-size: $font-base;
  color: $text-primary;
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-sizing: border-box;
  line-height: 1.5;
}

.comment-input-placeholder {
  color: $text-muted;
}

.comment-send-btn {
  flex-shrink: 0;
  height: 72rpx;
  padding: 0 $spacing-lg;
  background: linear-gradient(135deg, #d4a843 0%, #c49a38 100%);
  border-radius: 36rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4rpx 12rpx rgba(212, 168, 67, 0.3);

  &:active {
    opacity: 0.85;
    transform: scale(0.96);
  }

  &.disabled {
    opacity: 0.4;
    box-shadow: none;
  }
}

.comment-send-text {
  font-size: $font-base;
  color: #0a1628;
  font-weight: 700;
  white-space: nowrap;
}

// 自己的瓶子提示栏
.owner-hint-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 999;
  background: #1a2332;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  padding: $spacing-lg $spacing-md;
  padding-bottom: calc($spacing-lg + env(safe-area-inset-bottom));
  text-align: center;
}

.owner-hint-text {
  font-size: $font-sm;
  color: $text-muted;
}

.bottom-safe-area {
  height: 140rpx;
}

// 回复弹窗
.reply-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  z-index: 200;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: $spacing-lg;
}

.reply-modal-content {
  width: 100%;
  background: $bg-card;
  border-radius: $radius-xl;
  padding: $spacing-lg;
  border: 1px solid $border-color;
}

.reply-title {
  font-size: $font-lg;
  font-weight: 600;
  color: $text-primary;
  margin-bottom: $spacing-md;
  display: block;
}

.reply-textarea {
  width: 100%;
  min-height: 200rpx;
  background: $bg-input;
  border-radius: $radius-lg;
  padding: $spacing-md;
  border: 1px solid $border-color;
  color: $text-primary;
  font-size: $font-base;
  line-height: 1.6;
  box-sizing: border-box;
}

.reply-modal-actions {
  display: flex;
  gap: $spacing-md;
  margin-top: $spacing-lg;
}

.btn-cancel {
  flex: 1;
  height: 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: $bg-input;
  border-radius: 40rpx;
  border: 1px solid $border-color;

  text {
    font-size: $font-base;
    color: $text-secondary;
  }

  &:active {
    background: $bg-card-hover;
  }
}

.btn-send-reply {
  flex: 1;
  height: 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #d4a843 0%, #c49a38 100%);
  border-radius: 40rpx;
  box-shadow: 0 4rpx 16rpx rgba(212, 168, 67, 0.3);

  text {
    font-size: $font-base;
    font-weight: 600;
    color: #0a1628;
  }

  &:active {
    transform: scale(0.96);
    opacity: 0.9;
  }
}

// 扔回海里动画
.throw-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(10, 22, 40, 0.85);
  z-index: 300;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.ripple {
  position: absolute;
  border-radius: 50%;
  border: 2rpx solid rgba(212, 168, 67, 0.4);
  animation: ripple-expand 1.2s ease-out forwards;
}

.ripple-1 {
  width: 200rpx;
  height: 200rpx;
  animation-delay: 0s;
}

.ripple-2 {
  width: 200rpx;
  height: 200rpx;
  animation-delay: 0.2s;
}

.ripple-3 {
  width: 200rpx;
  height: 200rpx;
  animation-delay: 0.4s;
}

@keyframes ripple-expand {
  0% {
    width: 0;
    height: 0;
    opacity: 1;
    border-width: 4rpx;
  }
  100% {
    width: 600rpx;
    height: 600rpx;
    opacity: 0;
    border-width: 1rpx;
  }
}

.throw-text {
  font-size: $font-lg;
  font-weight: 600;
  color: $primary-color;
  z-index: 301;
  animation: throw-text-fade 1.2s ease-out forwards;
}

@keyframes throw-text-fade {
  0% {
    opacity: 0;
    transform: translateY(20rpx);
  }
  30% {
    opacity: 1;
    transform: translateY(0);
  }
  100% {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
