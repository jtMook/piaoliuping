<template>
  <view class="page">
    <!-- 顶部导航 -->
    <view class="header">
      <view class="header-content">
        <text class="header-title">个人中心</text>
        <text class="header-action" @click="showMore">⋮</text>
      </view>
    </view>

    <view class="content">
      <!-- 个人资料卡片 -->
      <view class="profile-card" @click="goToEditProfile">
        <view class="avatar-section">
          <view class="avatar">
            <image
              v-if="userStore.user?.avatarUrl"
              class="avatar-image"
              :src="userStore.user.avatarUrl"
              mode="aspectFill"
            />
            <text v-else class="avatar-text">{{ userStore.user?.nickname?.charAt(0) || '贝' }}</text>
          </view>
          <view class="profile-info">
            <view class="name-row">
              <text class="nickname">{{ userStore.user?.nickname || '海边拾贝人#8975' }}</text>
            </view>
            <view class="id-row">
              <text class="id-label">ID {{ userStore.user?.id || 8975 }}</text>
              <view class="copy-btn" @click.stop="copyId">
                <text class="copy-text">复制</text>
              </view>
            </view>
            <view class="tag-row">
              <view class="tag" v-if="ageText">
                <text class="tag-text">{{ ageText }}</text>
              </view>
              <view class="tag gender-tag" :class="userStore.user?.gender === 1 ? 'male' : userStore.user?.gender === 2 ? 'female' : ''" v-if="userStore.user?.gender">
                <text class="tag-text">{{ userStore.user?.gender === 1 ? '♂ 男' : '♀ 女' }}</text>
              </view>
              <view class="tag">
                <text class="tag-text">{{ userStore.user?.cityName || '杭州市' }}</text>
              </view>
              <view class="tag bottle-count-tag">
                <text class="tag-text">{{ bottleCount }}个瓶子</text>
              </view>
            </view>
          </view>
          <text class="profile-arrow">›</text>
        </view>
      </view>

      <!-- 福利中心 -->
      <view class="menu-section-title">福利中心</view>
      <view class="menu-card">
        <view class="menu-item" @click="handleCheckIn">
          <view class="menu-icon checkin-icon">
            <text class="icon-text">📅</text>
          </view>
          <text class="menu-label">每日签到</text>
          <view class="menu-right">
            <text class="menu-hint">签到领瓶子</text>
            <text class="menu-arrow">›</text>
          </view>
        </view>

        <view class="menu-item" @click="handleWatchAd">
          <view class="menu-icon ad-icon">
            <text class="icon-text">▶</text>
          </view>
          <text class="menu-label">看广告领瓶子</text>
          <view class="menu-right">
            <text class="menu-hint highlight">+10 瓶子</text>
            <text class="menu-arrow">›</text>
          </view>
        </view>

        <view class="menu-item" @click="handleVip">
          <view class="menu-icon vip-icon">
            <text class="icon-text">👑</text>
          </view>
          <text class="menu-label vip-label">VIP</text>
          <view class="menu-right">
            <text class="menu-hint">未开通</text>
            <text class="menu-arrow">›</text>
          </view>
        </view>
      </view>

      <!-- 瓶子 -->
      <view class="menu-section-title">瓶子</view>
      <view class="menu-card">
        <view class="menu-item" @click="handleStampRecharge">
          <view class="menu-icon stamp-icon">
            <text class="icon-text">🏷️</text>
          </view>
          <text class="menu-label">邮票充值</text>
          <view class="menu-right">
            <text class="menu-arrow">›</text>
          </view>
        </view>

        <view class="menu-item" @click="handleShare">
          <view class="menu-icon share-icon">
            <text class="icon-text">🔗</text>
          </view>
          <text class="menu-label">分享给好友</text>
          <view class="menu-right">
            <text class="menu-hint">可获瓶子奖励</text>
            <text class="menu-arrow">›</text>
          </view>
        </view>
      </view>

      <!-- 通用 -->
      <view class="menu-section-title">通用</view>
      <view class="menu-card">
        <view class="menu-item" @click="handleSubscribe">
          <view class="menu-icon subscribe-icon">
            <text class="icon-text">🔔</text>
          </view>
          <text class="menu-label">订阅消息</text>
          <view class="menu-right">
            <text class="menu-hint">未开启，可能错过消息</text>
            <text class="menu-arrow">›</text>
          </view>
        </view>

        <view class="menu-item" @click="handleContact">
          <view class="menu-icon contact-icon">
            <text class="icon-text">📞</text>
          </view>
          <text class="menu-label">联系我们</text>
          <view class="menu-right">
            <text class="menu-arrow">›</text>
          </view>
        </view>

        <view class="menu-item" @click="goToSettings">
          <view class="menu-icon settings-icon">
            <text class="icon-text">⚙</text>
          </view>
          <text class="menu-label">设置</text>
          <view class="menu-right">
            <text class="menu-arrow">›</text>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useUserStore } from '@/stores/user'
import { useBottleStore } from '@/stores/bottle'

const userStore = useUserStore()
const bottleStore = useBottleStore()
const bottleCount = ref(0)

const ageText = computed(() => {
  const birthday = userStore.user?.birthday
  if (!birthday) return ''
  const birth = new Date(birthday)
  const now = new Date()
  let age = now.getFullYear() - birth.getFullYear()
  const m = now.getMonth() - birth.getMonth()
  if (m < 0 || (m === 0 && now.getDate() < birth.getDate())) {
    age--
  }
  return age > 0 ? `${age}岁` : ''
})

const refreshData = async () => {
  try {
    await userStore.getProfile()
  } catch (error) {
    console.log('获取资料失败')
  }
  loadBottleCount()
}

onMounted(async () => {
  if (!userStore.user) {
    try {
      await userStore.getProfile()
    } catch (error) {
      console.log('未登录')
    }
  }
  loadBottleCount()
})

// 从编辑页返回时刷新数据
onShow(() => {
  refreshData()
})

const loadBottleCount = async () => {
  try {
    await userStore.getProfile()
    bottleCount.value = userStore.user?.bottleCount || 0
  } catch (error) {
    bottleCount.value = 0
  }
}

const goBack = () => {
  uni.navigateBack()
}

const showMore = () => {
  uni.showActionSheet({
    itemList: ['关于漂流瓶', '意见反馈'],
    success: (res) => {
      if (res.tapIndex === 0) {
        uni.showToast({ title: '关于漂流瓶 v1.0.0', icon: 'none' })
      } else {
        uni.showToast({ title: '意见反馈功能开发中', icon: 'none' })
      }
    }
  })
}

const copyId = () => {
  const id = String(userStore.user?.id || 8975)
  uni.setClipboardData({
    data: id,
    success: () => {
      uni.showToast({ title: 'ID 已复制', icon: 'none' })
    }
  })
}

const handleCheckIn = () => {
  uni.navigateTo({ url: '/pages/checkin/checkin' })
}

const handleStampRecharge = () => {
  uni.showToast({ title: '邮票充值功能开发中', icon: 'none' })
}

const handleWatchAd = async () => {
  try {
    await bottleStore.watchAd()
    uni.showToast({ title: '获得10个瓶子', icon: 'none' })
    await loadBottleCount()
  } catch (error) {
    uni.showToast({ title: '领取失败', icon: 'none' })
  }
}

const handleVip = () => {
  uni.showToast({ title: 'VIP功能开发中', icon: 'none' })
}

const handleShare = () => {
  uni.showToast({ title: '分享功能开发中', icon: 'none' })
}

const handleMoreApps = () => {
  uni.showToast({ title: '更多小程序功能开发中', icon: 'none' })
}

const handleContact = () => {
  uni.showToast({ title: '联系我们：support@piaoliu.com', icon: 'none' })
}

const handleSubscribe = () => {
  uni.showToast({ title: '订阅消息功能开发中', icon: 'none' })
}

const goToSettings = () => {
  uni.navigateTo({ url: '/pages/settings/settings' })
}

const goToEditProfile = () => {
  uni.navigateTo({ url: '/pages/edit-profile/edit-profile' })
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

.header-action {
  width: 60rpx;
  text-align: right;
  font-size: $font-xl;
  color: $text-primary;
}

.content {
  padding: $spacing-md;
  padding-bottom: calc($tabbar-height + $spacing-md);
}

.profile-card {
  background: $bg-card;
  border-radius: $radius-lg;
  border: 1px solid $border-color;
  padding: $spacing-lg;
  margin-bottom: $spacing-md;
}

.avatar-section {
  display: flex;
  align-items: center;
  gap: $spacing-md;
}

.profile-arrow {
  font-size: $font-xl;
  color: $text-muted;
  flex-shrink: 0;
}

.avatar {
  width: 120rpx;
  height: 120rpx;
  background: linear-gradient(135deg, $primary-color 0%, $primary-dark 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 4rpx solid $primary-color;
  flex-shrink: 0;
  overflow: hidden;
}

.avatar-image {
  width: 100%;
  height: 100%;
}

.avatar-text {
  font-size: $font-lg;
  font-weight: 700;
  color: $bg-dark;
}

.profile-info {
  display: flex;
  flex-direction: column;
  gap: $spacing-xs;
  flex: 1;
}

.name-row {
  display: flex;
  align-items: center;
}

.nickname {
  font-size: $font-md;
  font-weight: 600;
  color: $text-primary;
}

.id-row {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
}

.id-label {
  font-size: $font-sm;
  color: $text-muted;
}

.copy-btn {
  background: rgba(212, 168, 67, 0.2);
  border-radius: $radius-sm;
  padding: 2rpx 12rpx;
}

.copy-text {
  font-size: $font-xs;
  color: $primary-color;
}

.tag-row {
  display: flex;
  gap: $spacing-sm;
}

.tag {
  background: rgba(255, 255, 255, 0.08);
  border-radius: $radius-sm;
  padding: 4rpx 16rpx;
}

.gender-tag {
  &.male {
    background: rgba(100, 181, 246, 0.15);
  }

  &.female {
    background: rgba(240, 98, 146, 0.15);
  }
}

.tag-text {
  font-size: $font-xs;
  color: $text-secondary;
}

.bottle-count-tag {
  background: rgba(212, 168, 67, 0.15);

  .tag-text {
    color: $primary-color;
    font-weight: 600;
  }
}

.menu-card {
  background: $bg-card;
  border-radius: $radius-lg;
  border: 1px solid $border-color;
  margin-bottom: $spacing-md;
  overflow: hidden;
}

.menu-section-title {
  font-size: $font-xs;
  color: $text-muted;
  padding: $spacing-md $spacing-md $spacing-xs;
  letter-spacing: 2rpx;
}

.menu-item {
  display: flex;
  align-items: center;
  padding: $spacing-md $spacing-lg;
  border-bottom: 1px solid $border-color;

  &:last-child {
    border-bottom: none;
  }

  &:active {
    background: $bg-card-hover;
  }
}

.menu-icon {
  width: 64rpx;
  height: 64rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: $spacing-md;
  flex-shrink: 0;

  .icon-text {
    font-size: $font-md;
  }
}

.checkin-icon {
  background: rgba(212, 168, 67, 0.15);
}

.stamp-icon {
  background: rgba(99, 179, 237, 0.15);
}

.ad-icon {
  background: rgba(237, 137, 54, 0.15);
}

.vip-icon {
  background: rgba(212, 168, 67, 0.15);
}

.share-icon {
  background: rgba(72, 187, 120, 0.15);
}

.apps-icon {
  background: rgba(99, 179, 237, 0.15);
}

.contact-icon {
  background: rgba(252, 129, 129, 0.15);
}

.subscribe-icon {
  background: rgba(237, 137, 54, 0.15);
}

.settings-icon {
  background: rgba(90, 106, 122, 0.15);
}

.menu-label {
  flex: 1;
  font-size: $font-base;
  color: $text-primary;
  font-weight: 500;
}

.vip-label {
  color: $primary-color;
  font-weight: 700;
}

.menu-right {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
}

.menu-hint {
  font-size: $font-sm;
  color: $text-muted;

  &.highlight {
    color: $primary-color;
  }
}

.menu-arrow {
  font-size: $font-lg;
  color: $text-muted;
}
</style>
