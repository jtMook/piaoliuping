<template>
  <view class="page">
    <!-- 顶部导航 -->
    <view class="header">
      <view class="header-content">
        <text class="back-icon" @click="goBack">‹</text>
        <text class="header-title">写瓶中信</text>
      </view>
    </view>

    <view class="content" :class="{ throwing: isThrowing }">
      <!-- 瓶子类型选择 -->
      <view class="type-tabs">
        <view class="type-tab" :class="{ active: form.bottleType === 1 }" @click="form.bottleType = 1">
          <text class="tab-icon">📋</text>
          <text class="tab-name">普通瓶</text>
        </view>
        <view class="type-tab" :class="{ active: form.bottleType === 3 }" @click="form.bottleType = 3">
          <text class="tab-icon">🔒</text>
          <text class="tab-name">秘密瓶</text>
          <text class="tab-badge" v-if="form.bottleType === 3">24h</text>
        </view>
      </view>

      <!-- 副标题 -->
      <text class="subtitle">文字 + 最多3张图片 · 同城可见</text>

      <!-- 内容输入区 -->
      <view class="input-card">
        <textarea
          class="input-textarea"
          v-model="form.content"
          placeholder="把心事装进瓶子里..."
          placeholder-class="placeholder"
          maxlength="500"
          :auto-height="false"
        />
        <text class="char-count">{{ form.content.length }}/500</text>
      </view>

      <!-- 图片上传 -->
      <view class="image-section">
        <text class="image-label">添加图片（最多3张）</text>
        <view class="image-grid">
          <view class="image-add" @click="chooseImage" v-if="imageList.length < 3">
            <text class="add-icon">+</text>
          </view>
          <view class="image-item" v-for="(img, idx) in imageList" :key="idx">
            <image :src="img" mode="aspectFill" class="image-preview" />
            <view class="image-remove" @click="removeImage(idx)">
              <text class="remove-icon">×</text>
            </view>
          </view>
        </view>
      </view>

      <!-- 标签 -->
      <view class="tag-section">
        <text class="tag-label">选择标签（最多5个）</text>
        <view class="preset-tags">
          <view
            class="preset-tag"
            v-for="tag in defaultTags"
            :key="tag"
            :class="{ active: form.tags.includes(tag) }"
            @click="toggleTag(tag)"
          >
            <text class="preset-tag-text">{{ tag }}</text>
          </view>
        </view>
        <view class="tag-input-row">
          <input
            class="tag-input"
            v-model="newTag"
            placeholder="或输入自定义标签"
            placeholder-class="tag-placeholder"
            maxlength="10"
            @confirm="addTag"
          />
          <view class="tag-add-btn" @click="addTag">
            <text class="tag-add-text">添加</text>
          </view>
        </view>
        <view class="tag-list" v-if="form.tags.length > 0">
          <view class="tag-chip" v-for="(tag, idx) in form.tags" :key="idx">
            <text class="tag-chip-text">{{ tag }}</text>
            <text class="tag-chip-remove" @click="removeTag(idx)">×</text>
          </view>
        </view>
      </view>

      <!-- 底部操作 -->
      <view class="bottom-actions">
        <view class="btn-send" @click="submitBottle">
          <text class="btn-send-icon">📤</text>
          <text class="btn-send-text">扔进海里</text>
        </view>
      </view>

      <!-- 提示 -->
      <view class="tips-row">
        <text class="tips-icon">🛡️</text>
        <text class="tips-text">同城可见 · 内容审核通过后可被捡到</text>
      </view>
    </view>

    <!-- 扔瓶子动画遮罩 -->
    <view class="throw-overlay" v-if="isThrowing">
      <view class="bottle-flying">
        <view class="flying-bottle-cap"></view>
        <view class="flying-bottle-body"></view>
      </view>
      <view class="ripple ripple-1"></view>
      <view class="ripple ripple-2"></view>
      <view class="ripple ripple-3"></view>
      <text class="throw-text">已投入海中</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useUserStore } from '@/stores/user'
import { useBottleStore } from '@/stores/bottle'

const userStore = useUserStore()
const bottleStore = useBottleStore()
const imageList = ref<string[]>([])
const isThrowing = ref(false)
const newTag = ref('')

const form = reactive({
  content: '',
  bottleType: 1,
  visibleScope: 1,
  mediaUrls: [] as string[],
  tags: [] as string[]
})

const defaultTags = ['心情', '美食', '旅行', '音乐', '电影', '运动', '阅读', '摄影', '游戏', '宠物', '科技', '艺术', '日常', '吐槽', '求助']

const toggleTag = (tag: string) => {
  const idx = form.tags.indexOf(tag)
  if (idx > -1) {
    form.tags.splice(idx, 1)
  } else {
    if (form.tags.length >= 5) {
      uni.showToast({ title: '最多5个标签', icon: 'none' })
      return
    }
    form.tags.push(tag)
  }
}

const goBack = () => {
  uni.navigateBack()
}

const addTag = () => {
  const tag = newTag.value.trim()
  if (!tag) return
  if (form.tags.length >= 5) {
    uni.showToast({ title: '最多5个标签', icon: 'none' })
    return
  }
  if (form.tags.includes(tag)) {
    uni.showToast({ title: '标签已存在', icon: 'none' })
    return
  }
  form.tags.push(tag)
  newTag.value = ''
}

const removeTag = (idx: number) => {
  form.tags.splice(idx, 1)
}

const chooseImage = () => {
  uni.chooseImage({
    count: 3 - imageList.value.length,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      imageList.value.push(...res.tempFilePaths)
      form.mediaUrls = [...imageList.value]
    }
  })
}

const removeImage = (idx: number) => {
  imageList.value.splice(idx, 1)
  form.mediaUrls = [...imageList.value]
}

const submitBottle = async () => {
  if (!form.content.trim() && form.mediaUrls.length === 0) {
    uni.showToast({ title: '请输入内容', icon: 'none' })
    return
  }

  if (!userStore.user) {
    uni.showToast({ title: '请先登录', icon: 'none' })
    return
  }

  try {
    await bottleStore.createBottle({
      content: form.content,
      mediaUrls: form.mediaUrls,
      tags: form.tags,
      bottleType: form.bottleType,
      visibleScope: form.visibleScope
    })

    // 触发扔瓶子动画
    isThrowing.value = true
    setTimeout(() => {
      uni.navigateBack()
    }, 2000)
  } catch (error) {
    uni.showToast({ title: '发送失败，请重试', icon: 'none' })
  }
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
  font-size: $font-xl;
  font-weight: 600;
  color: $text-primary;
}

.header-action {
  font-size: $font-base;
  color: $primary-color;
  padding: $spacing-xs $spacing-md;
}

.header-title {
  flex: 1;
  text-align: center;
  margin-right: 60rpx;
}

.content {
  padding-top: calc(var(--status-bar-height, 44px) + $header-height);
  padding-bottom: $spacing-xl;
  transition: all 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

.content.throwing {
  opacity: 0;
  transform: scale(0.95);
}

.type-tabs {
  display: flex;
  gap: $spacing-sm;
  padding: $spacing-md;
}

.type-tab {
  display: flex;
  align-items: center;
  gap: $spacing-xs;
  padding: $spacing-sm $spacing-md;
  background: $bg-card;
  border-radius: $radius-lg;
  border: 1px solid $border-color;

  &.active {
    background: rgba(212, 168, 67, 0.15);
    border-color: $primary-color;
  }

  &:active {
    opacity: 0.8;
  }
}

.tab-icon {
  font-size: $font-base;
}

.tab-name {
  font-size: $font-sm;
  color: $text-secondary;

  .type-tab.active & {
    color: $primary-color;
    font-weight: 600;
  }
}

.tab-badge {
  font-size: $font-xs;
  color: $primary-color;
  background: rgba(212, 168, 67, 0.2);
  padding: 2rpx 8rpx;
  border-radius: $radius-sm;
}

.subtitle {
  display: block;
  font-size: $font-xs;
  color: $text-muted;
  padding: 0 $spacing-md;
  margin-bottom: $spacing-md;
}

.input-card {
  margin: 0 $spacing-md;
  background: $bg-card;
  border-radius: $radius-lg;
  border: 1px solid $border-color;
  padding: $spacing-md;
  position: relative;
}

.input-textarea {
  width: 100%;
  min-height: 300rpx;
  color: $text-primary;
  font-size: $font-base;
  line-height: 1.6;
}

.placeholder {
  color: $text-muted;
}

.char-count {
  position: absolute;
  right: $spacing-md;
  bottom: $spacing-md;
  font-size: $font-xs;
  color: $text-muted;
}

.image-section {
  padding: $spacing-md;
}

.image-label {
  display: block;
  font-size: $font-sm;
  color: $text-secondary;
  margin-bottom: $spacing-md;
}

.image-grid {
  display: flex;
  gap: $spacing-sm;
  flex-wrap: wrap;
}

.image-add {
  width: 180rpx;
  height: 180rpx;
  background: $bg-card;
  border-radius: $radius-md;
  border: 1px dashed $border-light;
  display: flex;
  align-items: center;
  justify-content: center;

  &:active {
    background: $bg-card-hover;
  }
}

.add-icon {
  font-size: 60rpx;
  color: $text-muted;
}

.image-item {
  width: 180rpx;
  height: 180rpx;
  position: relative;
  border-radius: $radius-md;
  overflow: hidden;
}

.image-preview {
  width: 100%;
  height: 100%;
}

.image-remove {
  position: absolute;
  top: 4rpx;
  right: 4rpx;
  width: 36rpx;
  height: 36rpx;
  background: rgba(0, 0, 0, 0.6);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.remove-icon {
  font-size: $font-md;
  color: #fff;
}

// 标签区域
.tag-section {
  padding: $spacing-md;
}

.tag-label {
  display: block;
  font-size: $font-sm;
  color: $text-secondary;
  margin-bottom: $spacing-md;
}

.preset-tags {
  display: flex;
  flex-wrap: wrap;
  gap: $spacing-sm;
  margin-bottom: $spacing-md;
}

.preset-tag {
  padding: $spacing-xs $spacing-md;
  background: rgba(255, 255, 255, 0.08);
  border-radius: $radius-lg;
  border: 1px solid transparent;

  &.active {
    background: rgba(212, 168, 67, 0.2);
    border-color: $primary-color;
  }

  &:active {
    opacity: 0.8;
  }
}

.preset-tag-text {
  font-size: $font-sm;
  color: $text-primary;
}

.tag-input-row {
  display: flex;
  gap: $spacing-sm;
  margin-bottom: $spacing-md;
}

.tag-input {
  flex: 1;
  height: 72rpx;
  background: $bg-card;
  border-radius: $radius-md;
  padding: 0 $spacing-md;
  font-size: $font-sm;
  color: $text-primary;
  border: 1px solid $border-color;
}

.tag-placeholder {
  color: $text-muted;
}

.tag-add-btn {
  flex-shrink: 0;
  width: 120rpx;
  height: 72rpx;
  background: linear-gradient(135deg, #d4a843 0%, #c49a38 100%);
  border-radius: $radius-md;
  display: flex;
  align-items: center;
  justify-content: center;

  &:active {
    opacity: 0.85;
  }
}

.tag-add-text {
  font-size: $font-sm;
  color: #0a1628;
  font-weight: 600;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: $spacing-sm;
}

.tag-chip {
  display: flex;
  align-items: center;
  gap: $spacing-xs;
  padding: $spacing-xs $spacing-sm;
  background: rgba(212, 168, 67, 0.15);
  border-radius: $radius-lg;
  border: 1px solid rgba(212, 168, 67, 0.3);
}

.tag-chip-text {
  font-size: $font-xs;
  color: $primary-color;
}

.tag-chip-remove {
  font-size: $font-sm;
  color: $text-muted;
  width: 32rpx;
  height: 32rpx;
  display: flex;
  align-items: center;
  justify-content: center;

  &:active {
    color: $primary-color;
  }
}

.bottom-actions {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: $spacing-md;
}

.btn-send {
  display: flex;
  align-items: center;
  justify-content: center;
  flex: 1;
  gap: $spacing-sm;
  padding: $spacing-sm $spacing-xl;
  background: linear-gradient(135deg, #d4a843 0%, #c49a38 100%);
  border-radius: 44rpx;
  box-shadow: 0 4rpx 16rpx rgba(212, 168, 67, 0.3);

  &:active {
    transform: scale(0.96);
    opacity: 0.9;
  }
}

.btn-send-icon {
  font-size: $font-lg;
}

.btn-send-text {
  font-size: $font-lg;
  font-weight: 600;
  color: #0a1628;
}

.tips-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: $spacing-xs;
  padding: $spacing-md;
}

.tips-icon {
  font-size: $font-sm;
}

.tips-text {
  font-size: $font-xs;
  color: $text-muted;
}

.throw-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(10, 22, 40, 0.9);
  z-index: 300;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.bottle-flying {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  animation: bottle-fly 2s ease-in-out forwards;
}

.flying-bottle-cap {
  width: 40rpx;
  height: 24rpx;
  background: #d4a843;
  border-radius: 4rpx 4rpx 0 0;
}

.flying-bottle-body {
  width: 80rpx;
  height: 120rpx;
  background: linear-gradient(180deg, rgba(212, 168, 67, 0.4) 0%, rgba(212, 168, 67, 0.15) 100%);
  border-radius: 8rpx 8rpx 20rpx 20rpx;
  border: 1px solid rgba(212, 168, 67, 0.6);
  box-shadow: 0 0 30rpx rgba(212, 168, 67, 0.3);
}

@keyframes bottle-fly {
  0% {
    transform: translateY(0) rotate(0deg);
    opacity: 1;
  }
  50% {
    transform: translateY(-200rpx) rotate(180deg);
    opacity: 1;
  }
  100% {
    transform: translateY(100rpx) rotate(360deg);
    opacity: 0;
  }
}

.ripple {
  position: absolute;
  border-radius: 50%;
  border: 2rpx solid rgba(212, 168, 67, 0.4);
  animation: ripple-expand 2s ease-out forwards;
}

.ripple-1 {
  animation-delay: 0.8s;
}

.ripple-2 {
  animation-delay: 1s;
}

.ripple-3 {
  animation-delay: 1.2s;
}

@keyframes ripple-expand {
  0% {
    width: 0;
    height: 0;
    opacity: 1;
    border-width: 4rpx;
  }
  100% {
    width: 600rpx;
    height: 600rpx;
    opacity: 0;
    border-width: 1rpx;
  }
}

.throw-text {
  font-size: $font-lg;
  font-weight: 600;
  color: $primary-color;
  z-index: 301;
  animation: throw-text-fade 2s ease-out forwards;
  margin-top: 100rpx;
}

@keyframes throw-text-fade {
  0% {
    opacity: 0;
    transform: translateY(20rpx);
  }
  40% {
    opacity: 0;
    transform: translateY(20rpx);
  }
  60% {
    opacity: 1;
    transform: translateY(0);
  }
  100% {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
