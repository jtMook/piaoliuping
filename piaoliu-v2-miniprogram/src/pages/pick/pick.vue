<template>
  <view class="page">
    <view class="header">
      <view class="header-content">
        <text class="back-icon" @click="goBack">‹</text>
        <text class="header-title">捡瓶子</text>
        <text class="header-action" @click="resetPick">换一批</text>
      </view>
    </view>

    <view class="content">
      <view class="ocean-animation" v-if="!currentBottle">
        <view class="ripple ripple-1"></view>
        <view class="ripple ripple-2"></view>
        <view class="ripple ripple-3"></view>
        <view class="floating-bottle">🍾</view>
      </view>

      <view class="empty-state" v-if="!currentBottle && !loading">
        <text class="empty-icon">🌊</text>
        <text class="empty-title">海面上还没有瓶子</text>
        <text class="empty-desc">点击下方按钮，打捞一个漂流瓶</text>
        <view class="btn btn-primary btn-lg btn-block" @click="handlePick">
          <text class="btn-icon">🤿</text>
          <text>捡一个瓶子</text>
        </view>
      </view>

      <view class="bottle-card" v-if="currentBottle">
        <view class="bottle-header">
          <view class="bottle-type-tag" :class="getBottleTypeClass(currentBottle.bottleType)">
            {{ getBottleTypeName(currentBottle.bottleType) }}
          </view>
          <text class="bottle-time">{{ formatTime(currentBottle.createdAt) }}</text>
        </view>

        <view class="bottle-body">
          <text class="bottle-content">{{ currentBottle.content }}</text>
          <view class="voice-player" v-if="currentBottle.bottleType === 2">
            <view class="player-wave">
              <text class="player-icon">▶️</text>
              <view class="wave-bars">
                <view class="bar" v-for="i in 20" :key="i" :style="{ height: getBarHeight(i) }"></view>
              </view>
              <text class="player-duration">12"</text>
            </view>
          </view>
        </view>

        <view class="bottle-footer">
          <view class="sender-info">
            <view class="sender-avatar">
              <text class="avatar-text">{{ getGenderIcon(currentBottle.gender) }}</text>
            </view>
            <view class="sender-detail">
              <text class="sender-name">{{ currentBottle.nickname || '匿名用户' }}</text>
              <text class="sender-location">
                <text class="location-icon">📍</text>
                {{ currentBottle.cityName }}{{ currentBottle.districtName ? ' · ' + currentBottle.districtName : '' }}
              </text>
            </view>
          </view>
        </view>

        <view class="bottle-actions" v-if="currentBottle.canReply">
          <view class="action-buttons">
            <view class="btn btn-secondary action-btn" @click="showReplyInput = true">
              <text>回复</text>
            </view>
            <view class="btn btn-primary action-btn" @click="goToChat">
              <text>私聊</text>
            </view>
          </view>
        </view>

        <view class="next-actions">
          <text class="next-link" @click="handlePick">再捡一个 ›</text>
        </view>
      </view>
    </view>

    <view class="reply-modal" v-if="showReplyInput" @click="showReplyInput = false">
      <view class="reply-content" @click.stop>
        <text class="reply-title">回复瓶子</text>
        <textarea
          class="input reply-textarea"
          v-model="replyContent"
          placeholder="写下你的回复..."
          placeholder-class="placeholder"
          maxlength="300"
        />
        <view class="reply-actions">
          <view class="btn btn-ghost" @click="showReplyInput = false">取消</view>
          <view class="btn btn-primary" @click="submitReply">发送</view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { useBottleStore } from '@/stores/bottle'
import type { PickedBottle } from '@/types'

const bottleStore = useBottleStore()
const currentBottle = ref<PickedBottle | null>(null)
const loading = ref(false)
const showReplyInput = ref(false)
const replyContent = ref('')
const hasReplied = ref(false)
const conversationId = ref(0)
const genderFilter = ref(0) // 从 URL 参数获取

onLoad((options) => {
  // 从 URL 获取 gender 参数
  if (options?.gender) {
    genderFilter.value = parseInt(options.gender)
  }
  handlePick()
})

const goBack = () => {
  uni.navigateBack()
}

const resetPick = () => {
  currentBottle.value = null
  handlePick()
}

const handlePick = async () => {
  loading.value = true
  try {
    const gender = genderFilter.value > 0 ? genderFilter.value : undefined
    const data = await bottleStore.pickBottle(gender)
    currentBottle.value = data
  } catch (error) {
    uni.showToast({ title: '暂时没捡到瓶子', icon: 'none' })
  } finally {
    loading.value = false
  }
}

const submitReply = async () => {
  if (!replyContent.value.trim()) {
    uni.showToast({ title: '请输入回复内容', icon: 'none' })
    return
  }

  try {
    if (currentBottle.value) {
      const res = await bottleStore.replyBottle(currentBottle.value.id, replyContent.value)
      conversationId.value = res.conversationId
      hasReplied.value = true
    }
    uni.showToast({ title: '回复成功', icon: 'success' })
    showReplyInput.value = false
    replyContent.value = ''
  } catch (error) {
    console.error(error)
  }
}

const goToChat = () => {
  uni.navigateTo({ url: `/pages/chat/chat?id=${conversationId.value}` })
}

const getBottleTypeClass = (type: number) => {
  const classes: Record<number, string> = {
    1: 'normal',
    2: 'voice',
    3: 'secret'
  }
  return classes[type] || 'normal'
}

const getBottleTypeName = (type: number) => {
  const names: Record<number, string> = {
    1: '普通瓶',
    2: '语音瓶',
    3: '秘密瓶'
  }
  return names[type] || '普通瓶'
}

const formatTime = (dateString: string) => {
  const date = new Date(dateString)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const hours = Math.floor(diff / 3600000)
  if (hours < 24) return `${hours}小时前`
  return `${Math.floor(hours / 24)}天前`
}

const getBarHeight = (index: number) => {
  const heights = [8, 16, 12, 20, 14, 24, 10, 18, 12, 22, 16, 20, 8, 14, 18, 12, 24, 10, 16, 20]
  return `${heights[(index - 1) % heights.length]}rpx`
}

const getGenderIcon = (gender: number) => {
  if (gender === 1) return ''
  if (gender === 2) return '♀'
  return '?'
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: linear-gradient(180deg, $bg-dark 0%, #16213e 100%);
}

.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  padding-top: var(--status-bar-height, 44px);
  background: rgba(15, 15, 26, 0.9);
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
  font-size: $font-base;
  color: $primary-color;
  padding: $spacing-xs $spacing-md;
}

.content {
  padding-top: calc(var(--status-bar-height, 44px) + $header-height + $spacing-lg);
  padding-left: $spacing-md;
  padding-right: $spacing-md;
  min-height: 100vh;
  box-sizing: border-box;
}

.ocean-animation {
  position: relative;
  width: 300rpx;
  height: 300rpx;
  margin: 100rpx auto;
  display: flex;
  align-items: center;
  justify-content: center;
}

.ripple {
  position: absolute;
  width: 100%;
  height: 100%;
  border: 2px solid rgba(100, 255, 218, 0.3);
  border-radius: 50%;
  animation: ripple 2s ease-out infinite;
}

.ripple-2 {
  animation-delay: 0.5s;
}

.ripple-3 {
  animation-delay: 1s;
}

.floating-bottle {
  font-size: 80rpx;
  animation: float 2s ease-in-out infinite;
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

.btn-icon {
  margin-right: $spacing-sm;
}

.bottle-card {
  background: $bg-card;
  border-radius: $radius-xl;
  padding: $spacing-lg;
  border: 1px solid $border-color;
  box-shadow: $shadow-lg;
}

.bottle-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: $spacing-md;
}

.bottle-type-tag {
  padding: 4rpx 16rpx;
  border-radius: $radius-sm;
  font-size: $font-xs;

  &.normal {
    background: rgba(100, 255, 218, 0.15);
    color: $primary-color;
  }

  &.voice {
    background: rgba(237, 137, 54, 0.15);
    color: $warning-color;
  }

  &.secret {
    background: rgba(139, 92, 246, 0.15);
    color: #a78bfa;
  }
}

.bottle-time {
  font-size: $font-xs;
  color: $text-muted;
}

.bottle-body {
  margin-bottom: $spacing-lg;
}

.bottle-content {
  font-size: $font-lg;
  color: $text-primary;
  line-height: 1.8;
}

.voice-player {
  margin-top: $spacing-md;
  padding: $spacing-md;
  background: $bg-input;
  border-radius: $radius-md;
}

.player-wave {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
}

.wave-bars {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 40rpx;
}

.bar {
  width: 6rpx;
  background: $primary-color;
  border-radius: 3rpx;
}

.player-duration {
  font-size: $font-sm;
  color: $text-secondary;
}

.bottle-footer {
  margin-bottom: $spacing-lg;
  padding-bottom: $spacing-md;
  border-bottom: 1px solid $border-color;
}

.sender-info {
  display: flex;
  align-items: center;
  gap: $spacing-md;
}

.sender-avatar {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  background: linear-gradient(135deg, $primary-color 0%, $primary-dark 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-text {
  font-size: $font-lg;
  color: #fff;
  font-weight: 600;
}

.sender-detail {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: $spacing-xs;
}

.sender-name {
  font-size: $font-base;
  font-weight: 600;
  color: $text-primary;
}

.sender-location {
  font-size: $font-xs;
  color: $text-muted;
  display: flex;
  align-items: center;
  gap: $spacing-xs;
}

.location-icon {
  font-size: $font-xs;
}

.bottle-actions {
  margin-bottom: $spacing-md;
}

.action-buttons {
  display: flex;
  gap: $spacing-md;
}

.action-btn {
  flex: 1;
  padding: $spacing-md;
  text-align: center;
}

.next-actions {
  text-align: center;
}

.next-link {
  font-size: $font-sm;
  color: $primary-color;
}

.reply-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: flex-end;
  z-index: 200;
}

.reply-content {
  width: 100%;
  background: $bg-card;
  border-radius: $radius-xl $radius-xl 0 0;
  padding: $spacing-lg;
  border-top: 1px solid $border-color;
}

.reply-title {
  display: block;
  font-size: $font-lg;
  font-weight: 600;
  color: $text-primary;
  margin-bottom: $spacing-md;
}

.reply-textarea {
  min-height: 200rpx;
  margin-bottom: $spacing-md;
}

.reply-actions {
  display: flex;
  gap: $spacing-md;
  justify-content: flex-end;
}
</style>
