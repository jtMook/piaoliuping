<template>
  <view class="page" :class="{ 'page-locked': showCityPicker }">
    <!-- 星星背景 -->
    <view class="stars">
      <view class="star" v-for="i in 20" :key="i" :style="getStarStyle(i)"></view>
    </view>

    <view class="content">
      <!-- 标题 -->
      <view class="title-section">
        <text class="main-title">完善资料</text>
        <text class="sub-title">完善资料后才能使用漂流瓶</text>
      </view>

      <!-- 性别选择 -->
      <view class="form-section">
        <text class="form-label">性别</text>
        <view class="gender-options">
          <view
            class="gender-btn"
            :class="{ active: gender === 1 }"
            @click="gender = 1"
          >
            <text class="gender-icon">♂</text>
            <text class="gender-text">男</text>
          </view>
          <view
            class="gender-btn"
            :class="{ active: gender === 2 }"
            @click="gender = 2"
          >
            <text class="gender-icon">♀</text>
            <text class="gender-text">女</text>
          </view>
        </view>
      </view>

      <!-- 出生日期 -->
      <view class="form-section">
        <text class="form-label">出生日期</text>
        <picker mode="date" :start="minDate" :end="maxDate" @change="onDateChange">
          <view class="date-picker">
            <text class="date-text" :class="{ placeholder: !birthday }">
              {{ birthday || '请选择出生日期' }}
            </text>
            <text class="date-arrow">›</text>
          </view>
        </picker>
      </view>

      <!-- 城市选择 -->
      <view class="form-section">
        <text class="form-label">所在城市</text>
        <view class="city-picker" @click="showCityPicker = true">
          <text class="city-text" :class="{ placeholder: !cityName }">
            {{ cityName || '请选择城市' }}
          </text>
          <text class="date-arrow">›</text>
        </view>
      </view>

      <!-- 提交按钮 -->
      <view class="submit-btn" :class="{ disabled: !canSubmit }" @click="handleSubmit">
        <text class="submit-text">完成</text>
      </view>
    </view>

    <!-- 城市选择弹窗 -->
    <view class="city-modal" v-if="showCityPicker" @click="showCityPicker = false">
      <view class="city-modal-content" :class="{ 'modal-short': cityListShort }" @click.stop>
        <view class="modal-header">
          <text class="modal-title">选择城市</text>
          <text class="modal-close" @click="showCityPicker = false"></text>
        </view>

        <view class="city-tabs">
          <view
            class="city-tab"
            :class="{ active: cityTab === 'province' }"
            @click="cityTab = 'province'"
          >
            <text>{{ selectedProvince?.name || '省' }}</text>
          </view>
          <view
            class="city-tab"
            :class="{ active: cityTab === 'city' }"
            @click="cityTab = 'city'"
          >
            <text>{{ selectedCity?.name || '市' }}</text>
          </view>
        </view>

        <scroll-view class="city-list" scroll-y :style="{ maxHeight: cityListMaxHeight }">
          <view
            class="city-item"
            v-for="item in cityCurrentList"
            :key="item.code"
            :class="{ active: isCitySelected(item) }"
            @click="selectCity(item)"
          >
            <text class="city-item-name">{{ item.name }}</text>
            <text class="city-item-check" v-if="isCitySelected(item)">✓</text>
          </view>
        </scroll-view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useUserStore } from '@/stores/user'
import { provinces, getCities } from '@/data/address'

const userStore = useUserStore()

const gender = ref(0)
const birthday = ref('')
const cityName = ref('')
const provinceCode = ref('')
const provinceName = ref('')
const cityCode = ref('')

const now = new Date()
const maxDate = `${now.getFullYear() - 10}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`
const minDate = '1950-01-01'

const canSubmit = computed(() => {
  return gender.value > 0 && birthday.value && cityCode.value
})

// 城市选择
const showCityPicker = ref(false)
const cityTab = ref<'province' | 'city'>('province')
const selectedProvince = ref<{ code: string; name: string } | null>(null)
const selectedCity = ref<{ code: string; name: string } | null>(null)

const cityCurrentList = computed(() => {
  if (cityTab.value === 'province') return provinces
  if (selectedProvince.value) return getCities(selectedProvince.value.code)
  return []
})

// 城市列表是否较短（少于 8 项）
const cityListShort = computed(() => {
  return cityCurrentList.value.length < 8
})

// 根据列表长度动态计算最大高度
const cityListMaxHeight = computed(() => {
  const count = cityCurrentList.value.length
  // 每项约 88rpx 高度，最少保证 3 项高度
  const contentHeight = Math.max(count, 3) * 88
  // 最大不超过 50vh（约 1000rpx）
  return `${Math.min(contentHeight, 1000)}rpx`
})

const isCitySelected = (item: { code: string }) => {
  if (cityTab.value === 'province') return selectedProvince.value?.code === item.code
  return selectedCity.value?.code === item.code
}

const selectCity = (item: { code: string; name: string }) => {
  if (cityTab.value === 'province') {
    selectedProvince.value = item
    selectedCity.value = null
    cityTab.value = 'city'
  } else {
    selectedCity.value = item
    provinceCode.value = selectedProvince.value!.code
    provinceName.value = selectedProvince.value!.name
    cityCode.value = item.code
    cityName.value = item.name
    showCityPicker.value = false
  }
}

const onDateChange = (e: any) => {
  birthday.value = e.detail.value
}

const handleSubmit = async () => {
  if (!canSubmit.value) {
    uni.showToast({ title: '请完善所有信息', icon: 'none' })
    return
  }

  try {
    await userStore.updateProfile({
      gender: gender.value,
      birthday: birthday.value,
      provinceCode: provinceCode.value,
      provinceName: provinceName.value,
      cityCode: cityCode.value,
      cityName: cityName.value
    })
    uni.showToast({ title: '资料完善成功', icon: 'success' })
    setTimeout(() => {
      uni.reLaunch({ url: '/pages/index/index' })
    }, 1000)
  } catch (error) {
    uni.showToast({ title: '保存失败', icon: 'none' })
  }
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
}

// 弹窗打开时锁定页面滚动
.page-locked {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  overflow: hidden;
}

.stars {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 40%;
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
  position: relative;
  z-index: 2;
  padding: calc(var(--status-bar-height, 44px) + 40rpx) $spacing-lg $spacing-lg;
}

.title-section {
  text-align: center;
  margin-bottom: 60rpx;
}

.main-title {
  display: block;
  font-size: 56rpx;
  font-weight: 700;
  color: $primary-color;
  margin-bottom: $spacing-sm;
}

.sub-title {
  display: block;
  font-size: $font-sm;
  color: $text-muted;
}

.form-section {
  margin-bottom: $spacing-xl;
}

.form-label {
  display: block;
  font-size: $font-base;
  color: $text-secondary;
  margin-bottom: $spacing-md;
  font-weight: 500;
}

.gender-options {
  display: flex;
  gap: $spacing-md;
}

.gender-btn {
  flex: 1;
  height: 96rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: $spacing-sm;
  border-radius: $radius-lg;
  border: 2rpx solid $border-color;
  background: $bg-card;

  &:active {
    opacity: 0.8;
  }

  &.active {
    border-color: $primary-color;
    background: rgba(212, 168, 67, 0.1);
  }
}

.gender-icon {
  font-size: $font-lg;
  color: $text-secondary;

  .gender-btn.active & {
    color: $primary-color;
  }
}

.gender-text {
  font-size: $font-base;
  color: $text-primary;
  font-weight: 500;
}

.date-picker,
.city-picker {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 96rpx;
  padding: 0 $spacing-md;
  background: $bg-card;
  border: 2rpx solid $border-color;
  border-radius: $radius-lg;

  &:active {
    background: $bg-card-hover;
  }
}

.date-text,
.city-text {
  font-size: $font-base;
  color: $text-primary;

  &.placeholder {
    color: $text-muted;
  }
}

.date-arrow {
  font-size: $font-xl;
  color: $text-muted;
}

.submit-btn {
  width: 100%;
  height: 96rpx;
  background: linear-gradient(135deg, #d4a843 0%, #c49a38 100%);
  border-radius: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8rpx 32rpx rgba(212, 168, 67, 0.3);
  margin-top: 40rpx;

  &:active {
    transform: scale(0.98);
    opacity: 0.9;
  }

  &.disabled {
    opacity: 0.4;
    box-shadow: none;
  }
}

.submit-text {
  font-size: $font-lg;
  font-weight: 700;
  color: #0a1628;
}

// 城市弹窗
.city-modal {
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
}

.city-modal-content {
  width: 100%;
  max-height: 70vh;
  background: $bg-card;
  border-radius: $radius-xl;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  margin: 0 30rpx;
}

// 列表较短时，弹窗保持最小高度
.modal-short {
  min-height: 400rpx;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: $spacing-lg;
  border-bottom: 1px solid $border-color;
}

.modal-title {
  font-size: $font-lg;
  font-weight: 600;
  color: $text-primary;
}

.modal-close {
  font-size: $font-lg;
  color: $text-muted;
  padding: $spacing-xs;
}

.city-tabs {
  display: flex;
  border-bottom: 1px solid $border-color;
}

.city-tab {
  flex: 1;
  padding: $spacing-md;
  text-align: center;
  border-bottom: 4rpx solid transparent;

  text {
    font-size: $font-base;
    color: $text-secondary;
  }

  &.active {
    border-bottom-color: $primary-color;

    text {
      color: $primary-color;
      font-weight: 600;
    }
  }
}

.city-list {
  min-height: 200rpx;
}

.city-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: $spacing-md $spacing-lg;
  border-bottom: 1px solid rgba(255, 255, 255, 0.04);

  &:active {
    background: $bg-card-hover;
  }

  &.active .city-item-name {
    color: $primary-color;
  }
}

.city-item-name {
  font-size: $font-base;
  color: $text-primary;
}

.city-item-check {
  font-size: $font-md;
  color: $primary-color;
}
</style>
