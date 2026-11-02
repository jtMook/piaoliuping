<template>
  <view class="page">
    <view class="header">
      <view class="header-content">
        <text class="back-icon" @click="goBack">‹</text>
        <text class="header-title">修改资料</text>
        <view class="header-placeholder"></view>
      </view>
    </view>

    <view class="content">
      <view class="avatar-section">
        <view class="avatar" @click="chooseAvatar">
          <image
            v-if="avatarUrl"
            class="avatar-image"
            :src="avatarUrl"
            mode="aspectFill"
          />
          <text v-else class="avatar-text">{{ userStore.user?.nickname?.charAt(0) || '贝' }}</text>
          <view class="avatar-badge">
            <text class="badge-icon">📷</text>
          </view>
        </view>
      </view>

      <view class="form-card">
        <view class="form-item">
          <text class="form-label">昵称</text>
          <input
            class="form-input"
            v-model="formData.nickname"
            placeholder="请输入昵称"
            placeholder-class="placeholder"
            maxlength="20"
          />
        </view>

        <view class="form-divider"></view>

        <view class="form-item">
          <text class="form-label">性别</text>
          <view class="gender-options">
            <view
              class="gender-option"
              :class="{ active: formData.gender === 1 }"
              @click="formData.gender = 1"
            >
              <text class="gender-icon">♂</text>
              <text class="gender-text">男</text>
            </view>
            <view
              class="gender-option"
              :class="{ active: formData.gender === 2 }"
              @click="formData.gender = 2"
            >
              <text class="gender-icon">♀</text>
              <text class="gender-text">女</text>
            </view>
          </view>
        </view>

        <view class="form-divider"></view>

        <view class="form-item">
          <text class="form-label">出生日期</text>
          <view class="form-value" @click="openDatePicker">
            <text class="date-icon">📅</text>
            <text class="form-text">{{ formData.birthday || '请选择出生日期' }}</text>
            <text class="form-arrow">›</text>
          </view>
        </view>

        <view class="form-divider"></view>

        <view class="form-item">
          <text class="form-label">所在城市</text>
          <view class="form-value" @click="goToAddress">
            <text class="location-icon"></text>
            <text class="form-text">{{ cityText }}</text>
            <text class="form-arrow">›</text>
          </view>
        </view>

        <view class="form-divider"></view>

        <view class="form-item-column">
          <text class="form-label">兴趣标签（最多5个）</text>
          <view class="interest-tags">
            <view
              class="interest-tag"
              v-for="tag in allTags"
              :key="tag"
              :class="{ active: formData.interestTags.includes(tag) }"
              @click="toggleTag(tag)"
            >
              <text class="interest-tag-text">{{ tag }}</text>
            </view>
          </view>
        </view>

        <text class="hint-text">挖同城需要选择所在城市</text>
      </view>

      <view class="submit-btn" @click="submitProfile">
        <text class="submit-text">确定</text>
      </view>
    </view>

    <view class="picker-mask" v-if="showDatePicker" @click="showDatePicker = false">
      <view class="picker-popup" @click.stop>
        <view class="picker-header">
          <text class="picker-cancel" @click="showDatePicker = false">取消</text>
          <text class="picker-title">选择出生日期</text>
          <text class="picker-confirm" @click="confirmDate">确定</text>
        </view>
        <view class="picker-body">
          <view class="picker-indicator"></view>
          <picker-view
            class="picker-view"
            :value="dateIndex"
            @change="onDateChange"
            indicator-style="height: 80rpx;"
          >
            <picker-view-column>
              <view class="picker-item" v-for="year in yearList" :key="year">{{ year }}年</view>
            </picker-view-column>
            <picker-view-column>
              <view class="picker-item" v-for="month in monthList" :key="month">{{ String(month).padStart(2, '0') }}月</view>
            </picker-view-column>
            <picker-view-column>
              <view class="picker-item" v-for="day in dayList" :key="day">{{ String(day).padStart(2, '0') }}日</view>
            </picker-view-column>
          </picker-view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useUserStore } from '@/stores/user'
import { post } from '@/utils/request'

const userStore = useUserStore()

const avatarUrl = ref('')

const formData = ref({
  nickname: '',
  gender: 0 as number,
  birthday: '',
  cityName: '',
  districtName: '',
  interestTags: [] as string[]
})

const allTags = ['美食', '旅行', '音乐', '电影', '运动', '阅读', '摄影', '游戏', '宠物', '科技', '艺术', '舞蹈']

const toggleTag = (tag: string) => {
  const idx = formData.value.interestTags.indexOf(tag)
  if (idx > -1) {
    formData.value.interestTags.splice(idx, 1)
  } else {
    if (formData.value.interestTags.length >= 5) {
      uni.showToast({ title: '最多选择5个标签', icon: 'none' })
      return
    }
    formData.value.interestTags.push(tag)
  }
}

const cityText = computed(() => {
  if (formData.value.cityName) {
    return formData.value.cityName
  }
  return '请选择城市'
})

onMounted(async () => {
  if (!userStore.user) {
    try {
      await userStore.getProfile()
    } catch (error) {
      console.log('未登录')
    }
  }
  initFormData()
})

// 从地址选择页返回时刷新城市显示
onShow(() => {
  // 重新获取用户资料，确保数据是最新的
  userStore.getProfile().then(() => {
    initFormData()
  })
})

const initFormData = () => {
  if (userStore.user) {
    formData.value.nickname = userStore.user.nickname || ''
    formData.value.gender = userStore.user.gender || 0
    formData.value.birthday = userStore.user.birthday || ''
    formData.value.cityName = userStore.user.cityName || ''
    formData.value.districtName = userStore.user.districtName || ''
    formData.value.interestTags = userStore.user.interestTags || []
    avatarUrl.value = userStore.user.avatarUrl || ''
  }
}

const chooseAvatar = () => {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: async (res) => {
      const tempFilePath = res.tempFilePaths[0]
      uni.showLoading({ title: '上传中...' })
      try {
        const data = await uploadFile(tempFilePath)
        avatarUrl.value = data.url
        uni.hideLoading()
        uni.showToast({ title: '头像已更新', icon: 'none' })
      } catch (error) {
        uni.hideLoading()
        uni.showToast({ title: '上传失败', icon: 'none' })
      }
    }
  })
}

const uploadFile = (filePath: string) => {
  return new Promise<{ url: string }>((resolve, reject) => {
    uni.uploadFile({
      url: 'http://localhost:8080/api/upload',
      filePath,
      name: 'file',
      header: {
        Authorization: `Bearer ${uni.getStorageSync('token')}`
      },
      success: (res) => {
        const data = JSON.parse(res.data)
        if (data.code === 0) {
          resolve(data.data)
        } else {
          reject(new Error(data.message || '上传失败'))
        }
      },
      fail: reject
    })
  })
}

const goBack = () => {
  uni.navigateBack()
}

const showDatePicker = ref(false)
const dateIndex = ref([30, 0, 0])
const tempDateIndex = ref([30, 0, 0])

const yearList = computed(() => {
  const years = []
  const currentYear = new Date().getFullYear()
  for (let y = currentYear - 1; y >= 1950; y--) {
    years.push(y)
  }
  return years
})

const monthList = computed(() => {
  const months = []
  for (let m = 1; m <= 12; m++) {
    months.push(m)
  }
  return months
})

const dayList = computed(() => {
  const year = yearList.value[dateIndex.value[0]] || 2000
  const month = monthList.value[dateIndex.value[1]] || 1
  const days = new Date(year, month, 0).getDate()
  const result = []
  for (let d = 1; d <= days; d++) {
    result.push(d)
  }
  return result
})

const openDatePicker = () => {
  if (formData.value.birthday) {
    const [y, m, d] = formData.value.birthday.split('-').map(Number)
    const yearIdx = yearList.value.indexOf(y)
    dateIndex.value = [Math.max(0, yearIdx), m - 1, d - 1]
  }
  tempDateIndex.value = [...dateIndex.value]
  showDatePicker.value = true
}

const onDateChange = (e: any) => {
  const newIndex = e.detail.value
  const newYearIdx = newIndex[0]
  const newMonthIdx = newIndex[1]
  const newDayIdx = newIndex[2]

  if (newYearIdx !== dateIndex.value[0] || newMonthIdx !== dateIndex.value[1]) {
    const year = yearList.value[newYearIdx]
    const month = newMonthIdx + 1
    const daysInMonth = new Date(year, month, 0).getDate()
    if (newDayIdx >= daysInMonth) {
      dateIndex.value = [newYearIdx, newMonthIdx, daysInMonth - 1]
      return
    }
  }
  dateIndex.value = newIndex
}

const confirmDate = () => {
  const year = yearList.value[dateIndex.value[0]]
  const month = dateIndex.value[1] + 1
  const day = dateIndex.value[2] + 1
  formData.value.birthday = `${year}-${String(month).padStart(2, '0')}-${String(day).padStart(2, '0')}`
  showDatePicker.value = false
}

const goToAddress = () => {
  uni.navigateTo({ url: '/pages/address/address' })
}

const submitProfile = async () => {
  if (!formData.value.nickname.trim()) {
    uni.showToast({ title: '请输入昵称', icon: 'none' })
    return
  }

  try {
    await userStore.updateProfile({
      nickname: formData.value.nickname.trim(),
      gender: formData.value.gender,
      birthday: formData.value.birthday,
      avatarUrl: avatarUrl.value,
      interestTags: formData.value.interestTags
    })
    uni.showToast({ title: '修改成功', icon: 'success' })
    setTimeout(() => {
      uni.navigateBack()
    }, 1000)
  } catch (error) {
    uni.showToast({ title: '修改失败', icon: 'none' })
  }
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
  padding-bottom: calc($tabbar-height + $spacing-xl);
}

.avatar-section {
  display: flex;
  justify-content: center;
  margin-bottom: $spacing-xl;
}

.avatar {
  width: 160rpx;
  height: 160rpx;
  background: linear-gradient(135deg, $primary-color 0%, $primary-dark 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 4rpx solid $primary-color;
  position: relative;
  overflow: hidden;

  &:active {
    opacity: 0.8;
  }
}

.avatar-image {
  width: 100%;
  height: 100%;
}

.avatar-text {
  font-size: $font-xl;
  font-weight: 700;
  color: $bg-dark;
}

.avatar-badge {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 52rpx;
  height: 52rpx;
  background: $primary-color;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 3rpx solid $bg-dark;
}

.badge-icon {
  font-size: 26rpx;
}

.form-card {
  background: $bg-card;
  border-radius: $radius-lg;
  border: 1px solid $border-color;
  padding: $spacing-lg;
  margin-bottom: $spacing-lg;
}

.form-item {
  display: flex;
  align-items: center;
  padding: $spacing-md 0;
}

.form-divider {
  height: 1px;
  background: $border-color;
}

.form-label {
  width: 140rpx;
  font-size: $font-base;
  color: $text-primary;
  font-weight: 500;
  flex-shrink: 0;
}

.form-input {
  flex: 1;
  font-size: $font-base;
  color: $text-primary;
  text-align: right;
}

.placeholder {
  color: $text-muted;
}

.gender-options {
  flex: 1;
  display: flex;
  gap: $spacing-md;
  justify-content: flex-end;
}

.gender-option {
  display: flex;
  align-items: center;
  gap: $spacing-xs;
  padding: $spacing-sm $spacing-lg;
  background: rgba(255, 255, 255, 0.08);
  border-radius: $radius-lg;
  border: 2px solid transparent;

  &.active {
    background: rgba(100, 255, 218, 0.15);
    border-color: $primary-color;
  }
}

.gender-icon {
  font-size: $font-lg;
}

.gender-text {
  font-size: $font-base;
  color: $text-primary;
}

.form-value {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: $spacing-sm;
}

.form-text {
  font-size: $font-base;
  color: $text-primary;
}

.form-arrow {
  font-size: $font-lg;
  color: $text-muted;
}

.date-icon,
.location-icon {
  font-size: $font-md;
}

.hint-text {
  display: block;
  font-size: $font-xs;
  color: $text-muted;
  text-align: center;
  margin-top: $spacing-md;
}

.form-item-column {
  padding: $spacing-md 0;
}

.interest-tags {
  display: flex;
  flex-wrap: wrap;
  gap: $spacing-sm;
  margin-top: $spacing-sm;
}

.interest-tag {
  padding: $spacing-sm $spacing-md;
  background: rgba(255, 255, 255, 0.08);
  border-radius: $radius-lg;
  border: 2px solid transparent;

  &.active {
    background: rgba(212, 168, 67, 0.2);
    border-color: $primary-color;
  }
}

.interest-tag-text {
  font-size: $font-sm;
  color: $text-primary;
}

.submit-btn {
  background: linear-gradient(135deg, #d4a843 0%, #c49a38 100%);
  border-radius: $radius-lg;
  padding: $spacing-md;
  text-align: center;
  box-shadow: 0 8rpx 32rpx rgba(212, 168, 67, 0.4);

  &:active {
    opacity: 0.9;
    transform: scale(0.98);
  }
}

.submit-text {
  font-size: $font-lg;
  font-weight: 700;
  color: #0a1628;
}

.picker-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.65);
  z-index: 1000;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
}

.picker-popup {
  background: #1a2332;
  border-radius: $radius-xl $radius-xl 0 0;
  padding-bottom: env(safe-area-inset-bottom);
  box-shadow: 0 -8rpx 40rpx rgba(0, 0, 0, 0.3);
}

.picker-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: $spacing-md $spacing-lg;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.picker-cancel {
  font-size: $font-base;
  color: $text-muted;
  padding: $spacing-xs $spacing-sm;
}

.picker-title {
  font-size: $font-base;
  font-weight: 600;
  color: $text-primary;
}

.picker-confirm {
  font-size: $font-base;
  color: #d4a843;
  font-weight: 600;
  padding: $spacing-xs $spacing-sm;
}

.picker-body {
  position: relative;
  padding: $spacing-sm 0;
}

.picker-indicator {
  position: absolute;
  top: 50%;
  left: $spacing-md;
  right: $spacing-md;
  height: 80rpx;
  margin-top: -40rpx;
  background: rgba(212, 168, 67, 0.08);
  border-radius: $radius-md;
  border: 1px solid rgba(212, 168, 67, 0.2);
  pointer-events: none;
  z-index: 1;
}

.picker-view {
  width: 100%;
  height: 440rpx;
  background: transparent !important;

  // H5 模式下穿透覆盖 picker-view 内部默认白/灰渐变背景
  :deep(.uni-picker-view),
  :deep(.picker-view),
  :deep(.uni-picker-view-wrapper),
  :deep(.picker-view-column),
  :deep(.picker-view-indicator) {
    background: transparent !important;
  }

  :deep(.uni-picker-view-mask) {
    background-image:
      linear-gradient(180deg, #1a2332 0%, rgba(26, 35, 50, 0) 100%),
      linear-gradient(0deg, #1a2332 0%, rgba(26, 35, 50, 0) 100%) !important;
  }

  :deep(.picker-view-indicator) {
    background: transparent !important;
    border: none !important;
  }
}

.picker-item {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: $font-lg;
  color: rgba(255, 255, 255, 0.4);
  line-height: 80rpx;
}
</style>
