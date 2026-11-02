<template>
  <view class="page">
    <view class="header">
      <view class="header-content">
        <text class="back-icon" @click="goBack">‹</text>
        <text class="header-title">选择地址</text>
        <view class="header-placeholder"></view>
      </view>
    </view>

    <view class="content">
      <view class="selected-section">
        <text class="section-label">已选择</text>
        <view class="selected-address">
          <text class="address-text">{{ selectedText || '请选择省市' }}</text>
        </view>
      </view>

      <view class="picker-section">
        <view class="picker-tabs">
          <view 
            class="tab" 
            :class="{ active: activeTab === 'province' }"
            @click="activeTab = 'province'"
          >
            <text>{{ selected.province?.name || '省' }}</text>
          </view>
          <view 
            class="tab" 
            :class="{ active: activeTab === 'city' }"
            @click="activeTab = 'city'"
          >
            <text>{{ selected.city?.name || '市' }}</text>
          </view>
        </view>

        <scroll-view class="picker-list" scroll-y>
          <view 
            class="picker-item" 
            v-for="item in currentList" 
            :key="item.code"
            :class="{ active: isSelected(item) }"
            @click="selectItem(item)"
          >
            <text class="picker-name">{{ item.name }}</text>
            <text class="picker-check" v-if="isSelected(item)">✓</text>
          </view>
        </scroll-view>
      </view>

      <view class="tips-section">
        <text class="tips-item">• 地址用于匹配同城的漂流瓶</text>
        <text class="tips-item">• 修改地址后 7 天内不可再次修改</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { provinces, getCities } from '@/data/address'
import type { Address } from '@/types'

const userStore = useUserStore()
const activeTab = ref<'province' | 'city'>('province')

const selected = ref<{
  province: Address | null
  city: Address | null
}>({
  province: null,
  city: null
})

onMounted(() => {
  if (userStore.user) {
    const user = userStore.user
    selected.value.province = {
      code: user.provinceCode,
      name: user.provinceName
    }
    selected.value.city = {
      code: user.cityCode,
      name: user.cityName
    }
  }
})

const currentList = computed(() => {
  switch (activeTab.value) {
    case 'province':
      return provinces
    case 'city':
      return selected.value.province ? getCities(selected.value.province.code) : []
    default:
      return []
  }
})

const selectedText = computed(() => {
  const { province, city } = selected.value
  if (!province) return ''
  let text = province.name
  if (city) text += ` ${city.name}`
  return text
})

const isSelected = (item: Address) => {
  const current = selected.value[activeTab.value]
  return current?.code === item.code
}

const selectItem = async (item: Address) => {
  if (activeTab.value === 'province') {
    selected.value.province = item
    selected.value.city = null
    activeTab.value = 'city'
  } else if (activeTab.value === 'city') {
    selected.value.city = item
    // 选择城市后自动保存
    const province = selected.value.province
    const city = item
    if (province && city) {
      try {
        console.log('自动保存地址:', {
          provinceCode: province.code,
          provinceName: province.name,
          cityCode: city.code,
          cityName: city.name
        })
        await userStore.updateAddress({
          provinceCode: province.code,
          provinceName: province.name,
          cityCode: city.code,
          cityName: city.name,
          districtCode: '',
          districtName: ''
        })
        console.log('地址保存成功')
        uni.showToast({ title: '城市已更新', icon: 'success' })
        setTimeout(() => {
          uni.navigateBack()
        }, 800)
      } catch (error) {
        console.error('保存地址失败:', error)
        uni.showToast({ title: '保存失败', icon: 'none' })
      }
    }
  }
}

const saveAddress = async () => {
  const { province, city } = selected.value
  if (!province || !city) {
    uni.showToast({ title: '请选择完整地址', icon: 'none' })
    return
  }

  try {
    console.log('保存地址:', { provinceCode: province.code, provinceName: province.name, cityCode: city.code, cityName: city.name })
    await userStore.updateAddress({
      provinceCode: province.code,
      provinceName: province.name,
      cityCode: city.code,
      cityName: city.name,
      districtCode: '',
      districtName: ''
    })
    console.log('保存成功, userStore.user:', userStore.user)
    uni.showToast({ title: '保存成功', icon: 'success' })
    setTimeout(() => {
      uni.navigateBack()
    }, 1000)
  } catch (error) {
    console.error('保存失败:', error)
    uni.showToast({ title: '保存失败', icon: 'none' })
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
  padding-top: calc(var(--status-bar-height, 44px) + $header-height + $spacing-md);
  padding-left: $spacing-md;
  padding-right: $spacing-md;
}

.selected-section {
  background: $bg-card;
  border-radius: $radius-lg;
  padding: $spacing-md;
  margin-bottom: $spacing-md;
  border: 1px solid $border-color;
}

.section-label {
  display: block;
  font-size: $font-sm;
  color: $text-muted;
  margin-bottom: $spacing-sm;
}

.selected-address {
  padding: $spacing-md;
  background: $bg-input;
  border-radius: $radius-md;
}

.address-text {
  font-size: $font-base;
  color: $text-primary;
}

.picker-section {
  background: $bg-card;
  border-radius: $radius-lg;
  border: 1px solid $border-color;
  overflow: hidden;
  margin-bottom: $spacing-md;
}

.picker-tabs {
  display: flex;
  border-bottom: 1px solid $border-color;
}

.tab {
  flex: 1;
  padding: $spacing-md;
  text-align: center;
  border-bottom: 2px solid transparent;
  
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

.picker-list {
  max-height: 600rpx;
}

.picker-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: $spacing-md $spacing-lg;
  border-bottom: 1px solid $border-color;
  
  &:last-child {
    border-bottom: none;
  }
  
  &:active {
    background: $bg-card-hover;
  }
  
  &.active .picker-name {
    color: $primary-color;
  }
}

.picker-name {
  font-size: $font-base;
  color: $text-primary;
}

.picker-check {
  font-size: $font-md;
  color: $primary-color;
}

.tips-section {
  padding: $spacing-md;
}

.tips-item {
  display: block;
  font-size: $font-sm;
  color: $text-muted;
  line-height: 1.8;
}
</style>
