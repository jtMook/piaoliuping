<template>
  <view class="empty-state">
    <text class="empty-icon">{{ icon }}</text>
    <text class="empty-title">{{ title }}</text>
    <text class="empty-desc" v-if="description">{{ description }}</text>
  </view>
</template>

<script setup lang="ts">
interface Props {
  icon: string
  title: string
  description?: string
}

defineProps<Props>()
</script>

<style lang="scss" scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 40rpx;
}

.empty-icon {
  font-size: 100rpx;
  margin-bottom: 32rpx;
}

.empty-title {
  font-size: 32rpx;
  font-weight: 600;
  color: $text-primary;
  margin-bottom: 16rpx;
}

.empty-desc {
  font-size: 26rpx;
  color: $text-muted;
  text-align: center;
}
</style>
