<template>
  <view class="page">
    <!-- 筛选按钮（左侧，定位下方） -->
    <view class="filter-btn-wrap" @click="showGenderFilter">
      <view class="filter-btn">
        <view class="filter-icon-circle">
          <view class="filter-lines">
            <view class="filter-line filter-line-1"></view>
            <view class="filter-line filter-line-2"></view>
            <view class="filter-line filter-line-3"></view>
          </view>
        </view>
        <text class="filter-btn-text">{{ genderFilterText }}</text>
      </view>
    </view>

    <!-- 性别筛选弹窗 -->
    <view class="gender-filter-modal" v-if="showFilter" @click="showFilter = false">
      <view class="gender-filter-popup" @click.stop>
        <view class="filter-title">
          <text>筛选性别</text>
        </view>
        <view class="filter-options">
          <view
            class="filter-option"
            :class="{ active: genderFilter === 0 }"
            @click="selectGender(0)"
          >
            <text class="option-text">不限</text>
          </view>
          <view
            class="filter-option"
            :class="{ active: genderFilter === 1 }"
            @click="selectGender(1)"
          >
            <text class="option-text">♂ 男</text>
          </view>
          <view
            class="filter-option"
            :class="{ active: genderFilter === 2 }"
            @click="selectGender(2)"
          >
            <text class="option-text">♀ 女</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 海洋场景 -->
    <view class="ocean-scene">
      <!-- 星星 -->
      <view class="stars">
        <view class="star" v-for="i in 20" :key="i" :style="getStarStyle(i)"></view>
      </view>
      <!-- 月亮 -->
      <view class="moon"></view>

      <!-- 漂浮的瓶子 -->
      <view class="floating-bottles">
        <view class="float-bottle fb-1">
          <view class="mini-bottle"></view>
        </view>
        <view class="float-bottle fb-2">
          <view class="mini-bottle"></view>
        </view>
        <view class="float-bottle fb-3">
          <view class="mini-bottle"></view>
        </view>
        <view class="float-bottle fb-4">
          <view class="mini-bottle"></view>
        </view>
        <view class="float-bottle fb-5">
          <view class="mini-bottle"></view>
        </view>
      </view>

      <!-- 主按钮区 -->
      <view class="action-area">
        <view class="btn-pick" @click="handlePick">
          <text class="btn-pick-text">捡一个瓶子</text>
        </view>
        <view class="btn-throw" @click="goToThrow">
          <text class="btn-throw-text">扔瓶子</text>
        </view>
        <text class="nearby-text">附近有 {{ nearbyCount }} 个瓶子在漂流</text>
      </view>

      <!-- 海底光晕与海浪，填充下半屏并营造海的氛围 -->
      <view class="sea-bottom">
        <view class="sea-wave sea-wave-1"></view>
        <view class="sea-wave sea-wave-2"></view>
        <text class="sea-tagline">把心事交给大海，让陌生人听见</text>
      </view>
    </view>

    <!-- 捡瓶子加载动画 -->
    <view class="pick-loading-modal" v-if="pickLoading">
      <view class="pick-loading-content">
        <view class="loading-ocean">
          <view class="loading-bottle">
            <view class="loading-bottle-body"></view>
            <view class="loading-bottle-cap"></view>
          </view>
          <view class="ripple ripple-1"></view>
          <view class="ripple ripple-2"></view>
          <view class="ripple ripple-3"></view>
        </view>
        <text class="loading-text">正在捞瓶子...</text>
      </view>
    </view>

    <!-- 捡瓶子结果弹窗 -->
    <view class="bottle-result-modal" v-if="showBottleResult" @click="closeBottleResult">
      <view class="bottle-result-card" @click.stop>
        <!-- 关闭按钮 -->
        <view class="close-btn-wrap" @click="closeBottleResult">
          <text class="close-btn">✕</text>
        </view>

        <!-- 发送者信息 -->
        <view class="sender-header">
          <view class="sender-avatar-wrap">
            <view class="sender-avatar">
              <text class="avatar-text">{{ currentBottle?.nickname?.charAt(0) || '匿' }}</text>
            </view>
            <view class="gender-badge" :class="currentBottle?.gender === 1 ? 'male' : currentBottle?.gender === 2 ? 'female' : ''">
              <text class="gender-badge-text">{{ currentBottle?.gender === 1 ? '♂' : currentBottle?.gender === 2 ? '♀' : '?' }}</text>
            </view>
          </view>
          <view class="sender-info">
            <text class="sender-name">{{ currentBottle?.nickname || '匿名用户' }}</text>
            <text class="sender-location">
              <text class="loc-icon-small">📍</text>
              {{ currentBottle?.provinceName || '未知' }}
            </text>
          </view>
        </view>

        <!-- 瓶子内容 -->
        <view class="bottle-content-card">
          <text class="bottle-text">{{ currentBottle?.content }}</text>
          <view class="bottle-source-row">
            <view class="bottle-type-tag">
              <text class="bottle-type-text">来自交友瓶</text>
            </view>
            <view class="picked-tags">
              <view class="picked-tag" v-for="tag in currentBottle?.tags" :key="tag">
                <text class="picked-tag-text">{{ tag }}</text>
              </view>
            </view>
          </view>
        </view>

        <!-- 操作按钮 -->
        <view class="result-actions">
          <view class="btn-reply" @click="goToReply">
            <text class="btn-reply-text">回复对方</text>
          </view>
          <view class="btn-private" @click="goToPrivateChat">
            <text class="btn-private-text">私聊</text>
          </view>
          <view class="btn-pick-again" @click="handlePickAgain">
            <text class="btn-pick-again-text">再捞一个</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 回复输入弹窗（底部弹出） -->
    <view class="reply-mask" v-if="showReplyInput" @click="showReplyInput = false"></view>
    <view class="reply-bottom-panel" v-if="showReplyInput">
      <view class="reply-input-row">
        <textarea
          class="reply-input"
          v-model="replyContent"
          placeholder="请文明友善发言"
          placeholder-class="reply-input-placeholder"
          :maxlength="200"
          :auto-height="true"
          :adjust-position="true"
        />
        <view class="reply-send-btn" :class="{ disabled: !replyContent.trim() }" @click="submitReply">
          <text class="reply-send-text">发送</text>
        </view>
      </view>
    </view>

    <!-- 底部安全区 -->
    <view class="bottom-safe"></view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useUserStore } from '@/stores/user'
import { useBottleStore } from '@/stores/bottle'
import type { PickedBottle } from '@/types'

const userStore = useUserStore()
const bottleStore = useBottleStore()
const nearbyCount = ref(0)
const showFilter = ref(false)
const genderFilter = ref(0) // 0=不限, 1=男, 2=女
const showBottleResult = ref(false)
const currentBottle = ref<PickedBottle | null>(null)
const pickLoading = ref(false)
const showReplyInput = ref(false)
const replyContent = ref('')

const genderFilterText = computed(() => {
  if (genderFilter.value === 1) return '限男'
  if (genderFilter.value === 2) return '限女'
  return '筛选'
})

const showGenderFilter = () => {
  showFilter.value = true
}

const selectGender = (gender: number) => {
  genderFilter.value = gender
  showFilter.value = false
  loadBottleCount(gender)
}

onMounted(async () => {
  if (!userStore.user) {
    try {
      await userStore.getProfile()
    } catch (error) {
      console.log('未登录，需要登录')
      uni.reLaunch({ url: '/pages/login/login' })
      return
    }
  }
  await loadBottleCount()
})

const loadBottleCount = async (gender?: number) => {
  try {
    const data = await bottleStore.getBottleCount(gender)
    nearbyCount.value = data.nearbyCount || 0
  } catch (error) {
    console.error('获取瓶子数量失败:', error)
    nearbyCount.value = 0
  }
}

const handlePick = async () => {
  if (pickLoading.value) return

  // 检查瓶子余额
  if (userStore.user && userStore.user.bottleCount <= 0) {
    uni.showToast({ title: '您的瓶子数量不足！', icon: 'none' })
    return
  }

  pickLoading.value = true
  try {
    const gender = genderFilter.value > 0 ? genderFilter.value : undefined
    const data = await bottleStore.pickBottle(gender)
    currentBottle.value = data
    showBottleResult.value = true
  } catch (error) {
    uni.showToast({ title: '暂时没捡到瓶子', icon: 'none' })
  } finally {
    pickLoading.value = false
  }
}

const handlePickAgain = async () => {
  currentBottle.value = null
  await handlePick()
}

const closeBottleResult = () => {
  showBottleResult.value = false
  currentBottle.value = null
}

const goToReply = () => {
  replyContent.value = ''
  showReplyInput.value = true
}

const submitReply = async () => {
  if (!replyContent.value.trim()) {
    uni.showToast({ title: '请输入回复内容', icon: 'none' })
    return
  }
  if (!currentBottle.value) return

  try {
    console.log('提交评论:', { id: currentBottle.value.id, content: replyContent.value })
    await bottleStore.commentBottle(currentBottle.value.id, replyContent.value)
    uni.showToast({ title: '评论成功', icon: 'success' })
    showReplyInput.value = false
    replyContent.value = ''
    showBottleResult.value = false
    currentBottle.value = null
  } catch (error: any) {
    console.error('评论失败:', error)
    uni.showToast({ title: error.message || '评论失败', icon: 'none' })
  }
}

const conversationId = ref(0)

const goToPrivateChat = async () => {
  if (!currentBottle.value) return
  try {
    const res = await bottleStore.replyBottle(currentBottle.value.id, '你好')
    conversationId.value = res.conversationId
    showBottleResult.value = false
    currentBottle.value = null
    uni.navigateTo({ url: `/pages/chat/chat?id=${conversationId.value}` })
  } catch (error) {
    uni.showToast({ title: '操作失败', icon: 'none' })
  }
}

const goToThrow = () => {
  uni.navigateTo({ url: '/pages/throw/throw' })
}

const getStarStyle = (i: number) => {
  const left = Math.random() * 100
  const top = Math.random() * 60
  const size = Math.random() * 3 + 1
  const delay = Math.random() * 3
  const opacity = Math.random() * 0.6 + 0.2
  return {
    left: `${left}%`,
    top: `${top}%`,
    width: `${size}rpx`,
    height: `${size}rpx`,
    animationDelay: `${delay}s`,
    opacity
  }
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: linear-gradient(180deg, #0d1b2a 0%, #1b2838 40%, #1a2a3a 70%, #162236 100%);
  position: relative;
  overflow: hidden;
}

// 筛选按钮
.filter-btn-wrap {
  position: fixed;
  top: calc(var(--status-bar-height, 44px) + 160rpx);
  left: $spacing-lg;
  z-index: 100;
}

.filter-btn {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  background: rgba(212, 168, 67, 0.15);
  border: 1px solid rgba(212, 168, 67, 0.3);
  border-radius: 48rpx;
  padding: $spacing-sm $spacing-lg;
  backdrop-filter: blur(10px);

  &:active {
    background: rgba(212, 168, 67, 0.25);
  }
}

.filter-icon-circle {
  width: 44rpx;
  height: 44rpx;
  background: linear-gradient(135deg, $primary-color 0%, #c49a38 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2rpx 8rpx rgba(212, 168, 67, 0.3);
}

.filter-lines {
  display: flex;
  flex-direction: column;
  gap: 3rpx;
  align-items: center;
}

.filter-line {
  height: 2rpx;
  background: #0a1628;
  border-radius: 1rpx;
}

.filter-line-1 {
  width: 20rpx;
}

.filter-line-2 {
  width: 14rpx;
}

.filter-line-3 {
  width: 8rpx;
}

.filter-btn-text {
  font-size: $font-base;
  color: $primary-color;
  font-weight: 600;
}

.ocean-scene {
  position: relative;
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding-top: calc(var(--status-bar-height, 44px) + 80rpx);
  overflow: hidden;
}

.stars {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 60%;
}

.star {
  position: absolute;
  background: #d4a843;
  border-radius: 50%;
  animation: twinkle 2s ease-in-out infinite alternate;
}

@keyframes twinkle {
  0% { opacity: 0.2; }
  100% { opacity: 0.8; }
}

.moon {
  position: absolute;
  top: 120rpx;
  right: 80rpx;
  width: 80rpx;
  height: 80rpx;
  background: radial-gradient(circle, #e8d5a3 0%, #c4a96a 100%);
  border-radius: 50%;
  box-shadow: 0 0 40rpx rgba(212, 168, 67, 0.3);
}

.floating-bottles {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.float-bottle {
  position: absolute;
  animation: floatBottle 4s ease-in-out infinite;

  .mini-bottle {
    width: 64rpx;
    height: 84rpx;
    background: rgba(212, 168, 67, 0.22);
    border-radius: 8rpx 8rpx 16rpx 16rpx;
    border: 1px solid rgba(212, 168, 67, 0.5);
    box-shadow: 0 0 20rpx rgba(212, 168, 67, 0.25);
    position: relative;

    &::before {
      content: '';
      position: absolute;
      top: -12rpx;
      left: 50%;
      transform: translateX(-50%);
      width: 24rpx;
      height: 14rpx;
      background: rgba(212, 168, 67, 0.9);
      box-shadow: 0 0 10rpx rgba(212, 168, 67, 0.4);
      border-radius: 4rpx 4rpx 0 0;
    }
  }
}

.fb-1 { top: 25%; left: 10%; animation-delay: 0s; animation-duration: 5.5s; }
.fb-2 { top: 30%; left: 70%; animation-delay: 1s; animation-duration: 6.5s; }
.fb-3 { top: 20%; left: 40%; animation-delay: 0.5s; animation-duration: 4.5s; }
.fb-4 { top: 35%; left: 25%; animation-delay: 1.5s; animation-duration: 7s; }
.fb-5 { top: 22%; left: 80%; animation-delay: 2s; animation-duration: 5s; }

// 远近层次：靠后的瓶子更小更淡
.fb-3 .mini-bottle { transform: scale(0.78); opacity: 0.7; }
.fb-5 .mini-bottle { transform: scale(0.85); opacity: 0.8; }

@keyframes floatBottle {
  0%   { transform: translate(0, 0) rotate(-4deg); }
  25%  { transform: translate(8rpx, -10rpx) rotate(0deg); }
  50%  { transform: translate(0, -22rpx) rotate(4deg); }
  75%  { transform: translate(-8rpx, -10rpx) rotate(0deg); }
  100% { transform: translate(0, 0) rotate(-4deg); }
}

.action-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: $spacing-lg;
  z-index: 10;
  margin-top: 100rpx;
}

.btn-pick {
  width: 420rpx;
  height: 100rpx;
  background: linear-gradient(135deg, #d4a843 0%, #c49a38 100%);
  border-radius: 50rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8rpx 32rpx rgba(212, 168, 67, 0.4);

  &:active {
    transform: scale(0.96);
    opacity: 0.9;
  }
}

.btn-pick-text {
  font-size: $font-xl;
  font-weight: 700;
  color: #0a1628;
}

.btn-throw {
  width: 320rpx;
  height: 88rpx;
  background: transparent;
  border: 2rpx solid rgba(212, 168, 67, 0.6);
  border-radius: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;

  &:active {
    background: rgba(212, 168, 67, 0.1);
    transform: scale(0.96);
  }
}

.btn-throw-text {
  font-size: $font-lg;
  font-weight: 600;
  color: $primary-color;
}

.nearby-text {
  font-size: $font-sm;
  color: $text-muted;
  margin-top: $spacing-sm;
}

// 海底光晕与海浪
.sea-bottom {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  height: 40%;
  pointer-events: none;
  z-index: 1;
  background: linear-gradient(180deg, rgba(10, 22, 40, 0) 0%, rgba(13, 27, 42, 0.35) 55%, rgba(20, 34, 54, 0.95) 100%);
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.sea-wave {
  position: absolute;
  left: -50%;
  bottom: 0;
  width: 200%;
  border-radius: 45% 45% 0 0 / 100% 100% 0 0;
  animation: seaScroll 12s linear infinite;
}

.sea-wave-1 {
  height: 120rpx;
  background: rgba(212, 168, 67, 0.06);
}

.sea-wave-2 {
  height: 90rpx;
  bottom: -10rpx;
  background: rgba(99, 179, 237, 0.05);
  animation-duration: 18s;
  animation-direction: reverse;
}

@keyframes seaScroll {
  0% { transform: translateX(0); }
  100% { transform: translateX(-25%); }
}

.sea-tagline {
  position: relative;
  z-index: 2;
  font-size: $font-sm;
  color: $text-secondary;
  letter-spacing: 2rpx;
  padding-bottom: calc(#{$tabbar-height} + 40rpx);
  text-align: center;
}

.bottom-safe {
  height: $tabbar-height;
}

// 性别筛选弹窗
.gender-filter-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 300;
}

.gender-filter-popup {
  width: 560rpx;
  background: #1a2332;
  border-radius: $radius-xl;
  padding: $spacing-lg;
  border: 1px solid rgba(212, 168, 67, 0.2);
}

.filter-title {
  text-align: center;
  font-size: $font-lg;
  font-weight: 600;
  color: $text-primary;
  margin-bottom: $spacing-lg;
}

.filter-options {
  display: flex;
  flex-direction: column;
  gap: $spacing-sm;
}

.filter-option {
  padding: $spacing-md $spacing-lg;
  border-radius: $radius-md;
  background: rgba(255, 255, 255, 0.05);
  text-align: center;
  transition: all 0.2s;

  &.active {
    background: rgba(212, 168, 67, 0.2);
    border: 1px solid rgba(212, 168, 67, 0.4);
  }
}

.option-text {
  font-size: $font-base;
  color: $text-primary;
}

// 捡瓶子加载动画
.pick-loading-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(13, 27, 42, 0.85);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 400;
}

.pick-loading-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: $spacing-lg;
}

.loading-ocean {
  position: relative;
  width: 200rpx;
  height: 200rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-bottle {
  position: relative;
  z-index: 2;
  animation: bottleBob 1.5s ease-in-out infinite;
}

.loading-bottle-body {
  width: 48rpx;
  height: 72rpx;
  background: linear-gradient(180deg, rgba(212, 168, 67, 0.3) 0%, rgba(212, 168, 67, 0.15) 100%);
  border-radius: 8rpx 8rpx 20rpx 20rpx;
  border: 2rpx solid rgba(212, 168, 67, 0.5);
}

.loading-bottle-cap {
  position: absolute;
  top: -10rpx;
  left: 50%;
  transform: translateX(-50%);
  width: 20rpx;
  height: 12rpx;
  background: #d4a843;
  border-radius: 4rpx 4rpx 0 0;
}

@keyframes bottleBob {
  0%, 100% { transform: translateY(0) rotate(-5deg); }
  50% { transform: translateY(-16rpx) rotate(5deg); }
}

.ripple {
  position: absolute;
  border: 2rpx solid rgba(212, 168, 67, 0.3);
  border-radius: 50%;
  animation: rippleExpand 2s ease-out infinite;
}

.ripple-1 {
  width: 60rpx;
  height: 60rpx;
  animation-delay: 0s;
}

.ripple-2 {
  width: 60rpx;
  height: 60rpx;
  animation-delay: 0.6s;
}

.ripple-3 {
  width: 60rpx;
  height: 60rpx;
  animation-delay: 1.2s;
}

@keyframes rippleExpand {
  0% {
    transform: scale(1);
    opacity: 0.6;
  }
  100% {
    transform: scale(4);
    opacity: 0;
  }
}

.loading-text {
  font-size: $font-base;
  color: $primary-color;
  font-weight: 500;
  letter-spacing: 2rpx;
}

// 捡瓶子结果弹窗
.bottle-result-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 400;
}

.bottle-result-card {
  width: 620rpx;
  background: linear-gradient(180deg, #1a2332 0%, #162236 100%);
  border-radius: $radius-xl;
  padding: $spacing-lg;
  position: relative;
  border: 1px solid rgba(212, 168, 67, 0.2);
  box-shadow: 0 16rpx 48rpx rgba(0, 0, 0, 0.5), 0 0 60rpx rgba(212, 168, 67, 0.1);
}

.close-btn-wrap {
  position: absolute;
  top: $spacing-md;
  right: $spacing-md;
  z-index: 10;
}

.close-btn {
  font-size: 36rpx;
  color: rgba(255, 255, 255, 0.4);
  width: 56rpx;
  height: 56rpx;
  display: flex;
  align-items: center;
  justify-content: center;

  &:active {
    color: rgba(255, 255, 255, 0.7);
  }
}

.sender-header {
  display: flex;
  align-items: center;
  gap: $spacing-md;
  margin-bottom: $spacing-lg;
}

.sender-avatar-wrap {
  position: relative;
}

.sender-avatar {
  width: 88rpx;
  height: 88rpx;
  border-radius: 50%;
  background: linear-gradient(135deg, #d4a843 0%, #c49a38 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  border: 3rpx solid rgba(212, 168, 67, 0.5);
}

.avatar-text {
  font-size: $font-lg;
  color: #0a1628;
  font-weight: 700;
}

.gender-badge {
  position: absolute;
  bottom: -4rpx;
  right: -4rpx;
  width: 36rpx;
  height: 36rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2rpx solid #1a2332;

  &.male {
    background: #4a90d9;
  }

  &.female {
    background: #e8637a;
  }
}

.gender-badge-text {
  font-size: 18rpx;
  color: #fff;
}

.sender-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4rpx;
}

.sender-name {
  font-size: $font-lg;
  font-weight: 600;
  color: $text-primary;
}

.sender-location {
  font-size: $font-xs;
  color: $text-muted;
  display: flex;
  align-items: center;
  gap: 4rpx;
}

.loc-icon-small {
  font-size: 20rpx;
}

.bottle-content-card {
  background: rgba(255, 255, 255, 0.06);
  border-radius: $radius-md;
  padding: $spacing-lg;
  margin-bottom: $spacing-lg;
  min-height: 160rpx;
  border: 1px solid rgba(255, 255, 255, 0.08);
}

.bottle-text {
  font-size: $font-lg;
  color: $text-primary;
  line-height: 1.8;
  display: block;
  margin-bottom: $spacing-md;
}

.bottle-source-row {
  display: flex;
  align-items: center;
}

.bottle-type-tag {
  background: rgba(212, 168, 67, 0.15);
  border-radius: $radius-sm;
  padding: 4rpx 16rpx;
  border: 1px solid rgba(212, 168, 67, 0.3);
}

.bottle-type-text {
  font-size: $font-xs;
  color: $primary-color;
}

.picked-tags {
  display: flex;
  gap: $spacing-xs;
  margin-left: $spacing-sm;
  flex-wrap: wrap;
}

.picked-tag {
  background: rgba(212, 168, 67, 0.1);
  border-radius: $radius-sm;
  padding: 2rpx 12rpx;
  border: 1px solid rgba(212, 168, 67, 0.2);
}

.picked-tag-text {
  font-size: $font-xs;
  color: $primary-color;
}

.result-actions {
  display: flex;
  gap: $spacing-sm;
}

.btn-reply {
  flex: 1;
  height: 80rpx;
  background: transparent;
  border-radius: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(212, 168, 67, 0.5);

  &:active {
    background: rgba(212, 168, 67, 0.1);
  }
}

.btn-reply-text {
  font-size: $font-base;
  color: $primary-color;
  font-weight: 600;
}

.btn-private {
  flex: 1;
  height: 80rpx;
  background: rgba(212, 168, 67, 0.15);
  border-radius: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid rgba(212, 168, 67, 0.3);

  &:active {
    background: rgba(212, 168, 67, 0.25);
  }
}

.btn-private-text {
  font-size: $font-base;
  color: $primary-color;
  font-weight: 600;
}

.btn-pick-again {
  flex: 1;
  height: 80rpx;
  background: linear-gradient(135deg, #d4a843 0%, #c49a38 100%);
  border-radius: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4rpx 16rpx rgba(212, 168, 67, 0.3);

  &:active {
    opacity: 0.9;
    transform: scale(0.96);
  }
}

.btn-pick-again-text {
  font-size: $font-base;
  color: #0a1628;
  font-weight: 700;
}

// 底部回复输入面板
.reply-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.4);
  z-index: 499;
}

.reply-bottom-panel {
  position: fixed;
  left: 0;
  right: 0;
  bottom: $tabbar-height;
  z-index: 999;
  background: #1a2332;
  border-top: 1px solid rgba(212, 168, 67, 0.2);
  padding: $spacing-md;
  box-shadow: 0 -4rpx 24rpx rgba(0, 0, 0, 0.3);
}

.reply-input-row {
  display: flex;
  align-items: flex-end;
  gap: $spacing-sm;
}

.reply-input {
  flex: 1;
  min-height: 72rpx;
  max-height: 200rpx;
  background: rgba(255, 255, 255, 0.08);
  border-radius: 36rpx;
  padding: $spacing-sm $spacing-md;
  font-size: $font-base;
  color: $text-primary;
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-sizing: border-box;
  line-height: 1.5;
}

.reply-input-placeholder {
  color: $text-muted;
}

.reply-send-btn {
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

.reply-send-text {
  font-size: $font-base;
  color: #0a1628;
  font-weight: 700;
  white-space: nowrap;
}
</style>
