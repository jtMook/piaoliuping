<template>
  <view class="page">
    <!-- 顶部导航 -->
    <view class="header">
      <view class="header-content">
        <text class="back-icon" @click="goBack">‹</text>
        <text class="header-title">设置</text>
        <view class="header-placeholder"></view>
      </view>
    </view>

    <view class="content">
      <!-- 内容审核 & 聊天记录 & 举报 -->
      <view class="menu-card">
        <view class="menu-item" @click="handleContentReview">
          <text class="menu-label">内容审核</text>
          <view class="menu-right">
            <text class="menu-hint">已开启</text>
            <text class="menu-arrow">›</text>
          </view>
        </view>
        <view class="menu-divider"></view>
        <view class="menu-item" @click="handleChatHistory">
          <text class="menu-label">聊天记录管理</text>
          <view class="menu-right">
            <text class="menu-hint">30天自动清理</text>
            <text class="menu-arrow">›</text>
          </view>
        </view>
        <view class="menu-divider"></view>
        <view class="menu-item" @click="handleReportRecord">
          <text class="menu-label">举报记录</text>
          <text class="menu-arrow">›</text>
        </view>
      </view>

      <!-- 消息通知 & 清除缓存 -->
      <view class="menu-card">
        <view class="menu-item">
          <text class="menu-label">消息通知</text>
          <switch
            class="toggle-switch"
            :checked="notificationEnabled"
            @change="toggleNotification"
            color="#d4a843"
          />
        </view>
        <view class="menu-divider"></view>
        <view class="menu-item" @click="handleClearCache">
          <text class="menu-label">清除缓存</text>
          <view class="menu-right">
            <text class="menu-hint">{{ cacheSize }}</text>
            <text class="menu-arrow">›</text>
          </view>
        </view>
      </view>

      <!-- 用户协议 & 隐私政策 & 关于 -->
      <view class="menu-card">
        <view class="menu-item" @click="handleUserAgreement">
          <text class="menu-label">用户协议</text>
          <text class="menu-arrow">›</text>
        </view>
        <view class="menu-divider"></view>
        <view class="menu-item" @click="handlePrivacyPolicy">
          <text class="menu-label">隐私政策</text>
          <text class="menu-arrow">›</text>
        </view>
        <view class="menu-divider"></view>
        <view class="menu-item" @click="handleAbout">
          <text class="menu-label">关于漂流瓶</text>
          <view class="menu-right">
            <text class="menu-hint">v1.0.0</text>
            <text class="menu-arrow">›</text>
          </view>
        </view>
      </view>

      <!-- 退出登录 -->
      <view class="logout-btn" @click="handleLogout">
        <text class="logout-text">退出登录</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const notificationEnabled = ref(true)
const cacheSize = ref('12.3 MB')

onMounted(async () => {
  if (!userStore.user) {
    try {
      await userStore.getProfile()
    } catch (error) {
      console.log('未登录')
    }
  }
  getCacheSize()
})

const getCacheSize = () => {
  // #ifdef APP-PLUS
  ;(plus as any).cache.calculate((size: number) => {
    if (size < 1024) {
      cacheSize.value = size + ' B'
    } else if (size < 1024 * 1024) {
      cacheSize.value = (size / 1024).toFixed(1) + ' KB'
    } else {
      cacheSize.value = (size / 1024 / 1024).toFixed(1) + ' MB'
    }
  })
  // #endif
}

const goBack = () => {
  uni.navigateBack()
}

const handleContentReview = () => {
  uni.showToast({ title: '内容审核已开启', icon: 'none' })
}

const handleChatHistory = () => {
  uni.showModal({
    title: '聊天记录管理',
    content: '聊天记录将保留30天，过期自动清理。是否立即清理？',
    success: (res) => {
      if (res.confirm) {
        uni.showToast({ title: '已清理', icon: 'none' })
      }
    }
  })
}

const handleReportRecord = () => {
  uni.showToast({ title: '举报记录功能开发中', icon: 'none' })
}

const toggleNotification = (e: any) => {
  notificationEnabled.value = e.detail.value
  uni.showToast({
    title: notificationEnabled.value ? '消息通知已开启' : '消息通知已关闭',
    icon: 'none'
  })
}

const handleClearCache = () => {
  uni.showModal({
    title: '清除缓存',
    content: `确定要清除 ${cacheSize.value} 缓存吗？`,
    success: (res) => {
      if (res.confirm) {
        // #ifdef APP-PLUS
        ;(plus as any).cache.clear(() => {
          cacheSize.value = '0 B'
          uni.showToast({ title: '缓存已清除', icon: 'none' })
        })
        // #endif
        // #ifdef MP-WEIXIN
        uni.clearStorageSync()
        cacheSize.value = '0 B'
        uni.showToast({ title: '缓存已清除', icon: 'none' })
        // #endif
      }
    }
  })
}

const handleUserAgreement = () => {
  uni.showToast({ title: '用户协议功能开发中', icon: 'none' })
}

const handlePrivacyPolicy = () => {
  uni.showToast({ title: '隐私政策功能开发中', icon: 'none' })
}

const handleAbout = () => {
  uni.showToast({ title: '关于漂流瓶 v1.0.0', icon: 'none' })
}

const handleLogout = () => {
  uni.showModal({
    title: '退出登录',
    content: '确定要退出登录吗？',
    confirmColor: '#fc8181',
    success: (res) => {
      if (res.confirm) {
        userStore.logout()
        uni.reLaunch({ url: '/pages/login/login' })
      }
    }
  })
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: $bg-dark;
}

.header {
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
  font-size: $font-md;
  font-weight: 600;
  color: $text-primary;
}

.header-placeholder {
  width: 60rpx;
}

.content {
  padding: $spacing-md;
  padding-bottom: calc($tabbar-height + $spacing-xl);
}

.menu-card {
  background: $bg-card;
  border-radius: $radius-lg;
  border: 1px solid $border-color;
  margin-bottom: $spacing-md;
  overflow: hidden;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: $spacing-md $spacing-lg;

  &:active {
    background: $bg-card-hover;
  }
}

.menu-divider {
  height: 1px;
  background: $border-color;
  margin: 0 $spacing-lg;
}

.menu-label {
  flex: 1;
  font-size: $font-base;
  color: $text-primary;
  font-weight: 500;
}

.menu-right {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
}

.menu-hint {
  font-size: $font-sm;
  color: $text-muted;
}

.menu-arrow {
  font-size: $font-lg;
  color: $text-muted;
  flex-shrink: 0;
}

.toggle-switch {
  flex-shrink: 0;
  transform: scale(0.85);
}

.logout-btn {
  background: transparent;
  border: 2px solid $error-color;
  border-radius: $radius-lg;
  padding: $spacing-md;
  text-align: center;
  margin-top: $spacing-lg;

  &:active {
    background: rgba(252, 129, 129, 0.1);
  }
}

.logout-text {
  font-size: $font-base;
  color: $error-color;
  font-weight: 600;
}
</style>
