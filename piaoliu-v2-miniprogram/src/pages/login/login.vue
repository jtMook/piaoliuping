<template>
  <view class="page">
    <!-- 星星背景 -->
    <view class="stars">
      <view class="star" v-for="i in 30" :key="i" :style="getStarStyle(i)"></view>
    </view>

    <view class="content">
      <!-- Logo 区域 -->
      <view class="logo-section">
        <view class="bottle-icon">
          <view class="bottle-shape"></view>
          <view class="bottle-sparkles">
            <view class="sparkle s1"></view>
            <view class="sparkle s2"></view>
            <view class="sparkle s3"></view>
          </view>
        </view>
        <text class="app-title">漂流瓶</text>
        <text class="app-slogan">把心事交给大海</text>
      </view>

      <!-- 登录按钮 -->
      <view class="login-buttons">
        <view class="btn-wechat" @click="wechatLogin">
          <text class="btn-text">微信一键登录</text>
        </view>
        <view class="btn-phone" @click="phoneLogin">
          <text class="btn-text-phone">手机号登录</text>
        </view>
      </view>

      <!-- H5 开发环境：显示测试账号快速切换 -->
      <!-- #ifndef MP-WEIXIN -->
      <view class="test-accounts">
        <text class="test-label">测试账号切换</text>
        <view class="account-list">
          <view
            class="account-item"
            :class="{ active: selectedAccount === 'A' }"
            @click="switchAccount('A')"
          >
            <text>用户A (拾贝人#8975)</text>
          </view>
          <view
            class="account-item"
            :class="{ active: selectedAccount === 'B' }"
            @click="switchAccount('B')"
          >
            <text>用户B (拾贝人#1234)</text>
          </view>
        </view>
      </view>
      <!-- #endif -->

      <!-- 底部协议 -->
      <view class="footer">
        <view class="agreement-row">
          <text class="agreement-link" @click="handleUserAgreement">用户协议</text>
          <text class="agreement-dot">·</text>
          <text class="agreement-link" @click="handlePrivacyPolicy">隐私政策</text>
        </view>
        <text class="agreement-tip">登录即代表同意</text>
      </view>
    </view>

    <!-- 底部海浪 -->
    <view class="wave-bottom"></view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()

// #ifndef MP-WEIXIN
// H5 开发环境：测试账号切换
const selectedAccount = ref<'A' | 'B'>('A')

const switchAccount = (account: 'A' | 'B') => {
  selectedAccount.value = account
}

// 获取测试账号对应的 code
const getTestCode = (): string => {
  return selectedAccount.value === 'A' ? 'dev_test_code_a' : 'dev_test_code_b'
}
// #endif

const wechatLogin = async () => {
  try {
    let code: string

    // #ifdef MP-WEIXIN
    // 微信小程序环境：使用真实微信 code
    const res = await uni.login({ provider: 'weixin' })
    code = (res as any).code || ''
    if (!code) {
      uni.showToast({ title: '获取登录凭证失败', icon: 'none' })
      return
    }
    // #endif

    // #ifndef MP-WEIXIN
    // H5 开发环境：使用测试账号 code
    code = getTestCode()
    // #endif

    const profileComplete = await userStore.login(code)
    
    // 根据后端返回的标志决定跳转
    if (profileComplete) {
      uni.switchTab({ url: '/pages/index/index' })
    } else {
      uni.redirectTo({ url: '/pages/complete-profile/complete-profile' })
    }
  } catch (error) {
    uni.showToast({ title: '登录失败，请重试', icon: 'none' })
  }
}

const phoneLogin = () => {
  uni.showToast({ title: '手机号登录功能开发中', icon: 'none' })
}

const handleUserAgreement = () => {
  uni.showToast({ title: '用户协议功能开发中', icon: 'none' })
}

const handlePrivacyPolicy = () => {
  uni.showToast({ title: '隐私政策功能开发中', icon: 'none' })
}

const getStarStyle = (i: number) => {
  const left = Math.random() * 100
  const top = Math.random() * 50
  const size = Math.random() * 4 + 1
  const delay = Math.random() * 3
  const opacity = Math.random() * 0.5 + 0.2
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
  display: flex;
  flex-direction: column;
}

.stars {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 60%;
  z-index: 1;
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

.content {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: $spacing-xl $spacing-lg;
  position: relative;
  z-index: 2;
}

.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 120rpx;
}

.bottle-icon {
  position: relative;
  width: 160rpx;
  height: 160rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: $spacing-xl;
}

.bottle-shape {
  width: 80rpx;
  height: 100rpx;
  background: linear-gradient(180deg, #2a3a4a 0%, #1a2a3a 100%);
  border-radius: 12rpx 12rpx 24rpx 24rpx;
  border: 2rpx solid rgba(212, 168, 67, 0.3);
  position: relative;

  &::before {
    content: '';
    position: absolute;
    top: -20rpx;
    left: 50%;
    transform: translateX(-50%);
    width: 32rpx;
    height: 24rpx;
    background: linear-gradient(180deg, #3a4a5a 0%, #2a3a4a 100%);
    border-radius: 6rpx 6rpx 0 0;
    border: 2rpx solid rgba(212, 168, 67, 0.3);
    border-bottom: none;
  }
}

.bottle-sparkles {
  position: absolute;
  width: 100%;
  height: 100%;
}

.sparkle {
  position: absolute;
  width: 8rpx;
  height: 8rpx;
  background: #d4a843;
  border-radius: 50%;
  animation: sparkle 1.5s ease-in-out infinite;
}

.s1 {
  top: 20%;
  right: 10%;
  animation-delay: 0s;
}

.s2 {
  top: 40%;
  left: 5%;
  animation-delay: 0.5s;
}

.s3 {
  bottom: 30%;
  right: 5%;
  animation-delay: 1s;
}

@keyframes sparkle {
  0%, 100% { opacity: 0.3; transform: scale(0.8); }
  50% { opacity: 1; transform: scale(1.2); }
}

.app-title {
  font-size: 72rpx;
  font-weight: 700;
  color: $primary-color;
  margin-bottom: $spacing-md;
  letter-spacing: 8rpx;
}

.app-slogan {
  font-size: $font-base;
  color: $text-secondary;
  letter-spacing: 4rpx;
}

.login-buttons {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: $spacing-md;
  margin-bottom: auto;
}

.btn-wechat {
  width: 100%;
  height: 96rpx;
  background: linear-gradient(135deg, #d4a843 0%, #c49a38 100%);
  border-radius: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: $spacing-sm;
  box-shadow: 0 8rpx 32rpx rgba(212, 168, 67, 0.3);

  &:active {
    transform: scale(0.98);
    opacity: 0.9;
  }
}

.wechat-icon {
  font-size: $font-lg;
}

.btn-text {
  font-size: $font-lg;
  font-weight: 600;
  color: #0a1628;
}

.btn-phone {
  width: 100%;
  height: 96rpx;
  background: transparent;
  border: 2rpx solid rgba(212, 168, 67, 0.5);
  border-radius: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;

  &:active {
    background: rgba(212, 168, 67, 0.1);
    transform: scale(0.98);
  }
}

.btn-text-phone {
  font-size: $font-lg;
  font-weight: 600;
  color: $primary-color;
}

.test-accounts {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: $spacing-sm;
  margin-top: $spacing-lg;
}

.test-label {
  font-size: $font-xs;
  color: $text-muted;
  text-align: center;
}

.account-list {
  display: flex;
  gap: $spacing-sm;
}

.account-item {
  flex: 1;
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: $radius-lg;
  border: 2rpx solid $border-color;
  background: $bg-card;

  text {
    font-size: $font-sm;
    color: $text-muted;
  }

  &.active {
    border-color: $primary-color;
    background: rgba(212, 168, 67, 0.1);

    text {
      color: $primary-color;
      font-weight: 600;
    }
  }

  &:active {
    opacity: 0.8;
  }
}

.footer {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: $spacing-sm;
  margin-bottom: $spacing-xl;
}

.agreement-row {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
}

.agreement-link {
  font-size: $font-sm;
  color: $text-secondary;
}

.agreement-dot {
  font-size: $font-sm;
  color: $text-muted;
}

.agreement-tip {
  font-size: $font-xs;
  color: $text-muted;
}

.wave-bottom {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 120rpx;
  background: linear-gradient(180deg, transparent 0%, rgba(10, 22, 40, 0.8) 100%);
  border-radius: 50% 50% 0 0 / 100% 100% 0 0;
  z-index: 1;
}
</style>
