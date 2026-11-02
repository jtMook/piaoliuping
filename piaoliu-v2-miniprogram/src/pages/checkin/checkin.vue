<template>
  <view class="page">
    <!-- 顶部导航 -->
    <view class="header">
      <view class="header-content">
        <text class="back-icon" @click="goBack">‹</text>
        <text class="header-title">每日签到</text>
        <view class="header-placeholder"></view>
      </view>
    </view>

    <view class="content">
      <!-- 签到日历卡片 -->
      <view class="checkin-card">
        <view class="days-grid">
          <view
            class="day-item"
            v-for="(day, index) in days"
            :key="index"
            :class="{
              'checked': day.checked,
              'today': day.isToday,
              'reward-day': day.isRewardDay
            }"
          >
            <view class="day-reward">
              <text class="reward-text">+{{ day.reward }}</text>
            </view>
            <text class="day-label">第{{ index + 1 }}天</text>
          </view>
        </view>
      </view>

      <!-- 签到按钮 -->
      <view class="btn-container">
        <view
          class="btn-checkin"
          :class="{ disabled: checkedInToday }"
          @click="handleCheckin"
        >
          <text class="btn-text">
            {{ checkedInToday ? '今日已签到' : '立刻签到' }}
          </text>
        </view>
      </view>

      <!-- 连续签到信息 -->
      <view class="info-card" v-if="consecutiveDays > 0">
        <text class="info-text">已连续签到 {{ consecutiveDays }} 天</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { get, post } from '@/utils/request'

const checkedInToday = ref(false)
const consecutiveDays = ref(0)

const days = ref([
  { reward: 10, checked: false, isToday: false, isRewardDay: false },
  { reward: 10, checked: false, isToday: false, isRewardDay: false },
  { reward: 30, checked: false, isToday: false, isRewardDay: true },
  { reward: 10, checked: false, isToday: false, isRewardDay: false },
  { reward: 10, checked: false, isToday: false, isRewardDay: false },
  { reward: 10, checked: false, isToday: false, isRewardDay: false },
  { reward: 100, checked: false, isToday: false, isRewardDay: true }
])

onMounted(() => {
  loadCheckinStatus()
})

const loadCheckinStatus = async () => {
  try {
    const data = await get<any>('/checkin/status')
    checkedInToday.value = data.checkedInToday
    consecutiveDays.value = data.consecutiveDays || 0

    // 更新签到状态（7天循环）
    const displayDay = consecutiveDays.value % 7
    for (let i = 0; i < displayDay; i++) {
      days.value[i].checked = true
    }

    // 标记今天
    if (!checkedInToday.value) {
      days.value[displayDay].isToday = true
    }
  } catch (error) {
    console.error(error)
  }
}

const handleCheckin = async () => {
  if (checkedInToday.value) return

  try {
    const data = await post<any>('/checkin')
    checkedInToday.value = true
    consecutiveDays.value = data.consecutiveDays || consecutiveDays.value + 1

    // 更新签到状态（7天循环）
    const dayIndex = (consecutiveDays.value - 1) % 7
    if (dayIndex >= 0 && dayIndex < 7) {
      days.value[dayIndex].checked = true
    }

    uni.showToast({
      title: `签到成功，获得${data.bottleReward || 1}个瓶子`,
      icon: 'none'
    })
  } catch (error) {
    console.error(error)
  }
}

const goBack = () => {
  uni.navigateBack()
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
  padding: $spacing-lg $spacing-md;
}

.checkin-card {
  background: $bg-card;
  border-radius: $radius-lg;
  border: 1px solid $border-color;
  padding: $spacing-lg;
  margin-bottom: $spacing-lg;
}

.days-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: $spacing-sm;
}

.day-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: $spacing-md $spacing-xs;
  background: rgba(255, 255, 255, 0.05);
  border-radius: $radius-md;
  border: 1px solid $border-color;

  &.checked {
    background: rgba(212, 168, 67, 0.15);
    border-color: $primary-color;

    .reward-text {
      color: $primary-color;
    }
  }

  &.today {
    border-color: $primary-color;
    border-width: 2rpx;
  }

  &.reward-day {
    grid-column: span 2;

    .day-reward {
      width: 100%;
    }
  }
}

.day-reward {
  width: 80rpx;
  height: 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: $spacing-xs;
}

.reward-text {
  font-size: $font-lg;
  font-weight: 700;
  color: $text-primary;
}

.day-label {
  font-size: $font-xs;
  color: $text-muted;
}

.btn-container {
  margin-bottom: $spacing-lg;
}

.btn-checkin {
  width: 100%;
  height: 96rpx;
  background: linear-gradient(135deg, $primary-color 0%, $primary-dark 100%);
  border-radius: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8rpx 32rpx rgba(212, 168, 67, 0.3);

  &:active:not(.disabled) {
    transform: scale(0.98);
    opacity: 0.9;
  }

  &.disabled {
    background: rgba(255, 255, 255, 0.1);
    box-shadow: none;
  }
}

.btn-text {
  font-size: $font-lg;
  font-weight: 600;
  color: $bg-dark;

  .btn-checkin.disabled & {
    color: $text-muted;
  }
}

.info-card {
  background: $bg-card;
  border-radius: $radius-lg;
  border: 1px solid $border-color;
  padding: $spacing-md $spacing-lg;
  text-align: center;
}

.info-text {
  font-size: $font-sm;
  color: $text-secondary;
}
</style>
