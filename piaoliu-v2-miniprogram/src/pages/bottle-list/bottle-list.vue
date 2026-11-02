<template>
  <view class="page">
    <view class="content">
      <!-- 统计信息 -->
      <view class="stats-row">
        <view class="stat-item">
          <text class="stat-value">{{ stats.thrown }}</text>
          <text class="stat-label">扔出的</text>
        </view>
        <view class="stat-item">
          <text class="stat-value">{{ stats.picked }}</text>
          <text class="stat-label">捡到的</text>
        </view>
        <view class="stat-item">
          <text class="stat-value">{{ stats.replied }}</text>
          <text class="stat-label">收到回复</text>
        </view>
      </view>

      <!-- 我的瓶子 -->
      <view class="section-title">我的瓶子</view>

      <!-- Tab 切换 -->
      <view class="tabs">
        <view
          class="tab-item"
          :class="{ active: activeTab === 'all' }"
          @click="switchTab('all')"
        >
          <text>全部</text>
          <view class="tab-indicator" v-if="activeTab === 'all'"></view>
        </view>
        <view
          class="tab-item"
          :class="{ active: activeTab === 'thrown' }"
          @click="switchTab('thrown')"
        >
          <text>扔出的</text>
          <view class="tab-indicator" v-if="activeTab === 'thrown'"></view>
        </view>
        <view
          class="tab-item"
          :class="{ active: activeTab === 'picked' }"
          @click="switchTab('picked')"
        >
          <text>捡到的</text>
          <view class="tab-indicator" v-if="activeTab === 'picked'"></view>
        </view>
        <view
          class="tab-item"
          :class="{ active: activeTab === 'secret' }"
          @click="switchTab('secret')"
        >
          <text>秘密瓶</text>
          <view class="tab-indicator" v-if="activeTab === 'secret'"></view>
        </view>
      </view>

      <!-- 瓶子列表 -->
      <view class="bottle-list" v-if="bottles.length > 0">
        <view
          class="bottle-item"
          v-for="bottle in bottles"
          :key="bottle.id"
          @click="goToBottleDetail(bottle.id)"
        >
          <!-- 发送者信息 -->
          <view class="sender-info">
            <view class="sender-avatar">
              <image
                v-if="bottle.user?.avatarUrl"
                class="sender-avatar-image"
                :src="bottle.user.avatarUrl"
                mode="aspectFill"
              />
              <text v-else class="avatar-text">{{ bottle.user?.nickname?.charAt(0) || '?' }}</text>
            </view>
            <view class="sender-detail">
              <view class="sender-name-row">
                <text class="sender-name">{{ bottle.user?.nickname || '匿名用户' }}</text>
                <view class="sender-tags">
                  <view class="tag-age" v-if="getSenderAge(bottle.user)">
                    <text class="tag-age-text">{{ getSenderAge(bottle.user) }}岁</text>
                  </view>
                  <view class="tag-gender" :class="bottle.user?.gender === 1 ? 'male' : bottle.user?.gender === 2 ? 'female' : ''" v-if="bottle.user?.gender">
                    <text class="tag-gender-text">{{ bottle.user?.gender === 1 ? '♂' : '♀' }}</text>
                  </view>
                </view>
              </view>
              <text class="sender-location">{{ bottle.user?.cityName || '未知' }}</text>
            </view>
          </view>

          <view class="bottle-header">
            <view class="type-tag">
              {{ getBottleTypeName(bottle.bottleType) }} · {{ getCategoryName(bottle.bottleType) }}
            </view>
            <view class="item-delete" @click.stop="handleDelete(bottle)">
              <text class="item-delete-text">🗑️</text>
            </view>
          </view>
          <text class="bottle-content">{{ bottle.content }}</text>
          <view class="bottle-footer">
            <text class="bottle-time">{{ formatTime(bottle.createdAt) }}</text>
            <text class="bottle-reply" v-if="bottle.replyCount > 0">
              {{ bottle.replyCount }}人回复
            </text>
            <text class="bottle-status" v-if="bottle.status === 2">
              已销毁
            </text>
          </view>
        </view>

        <!-- 加载状态 -->
        <view class="load-more" v-if="isLoading">
          <text class="load-more-text">加载中...</text>
        </view>
        <view class="load-more" v-else-if="!hasMore">
          <text class="load-more-text">没有更多了</text>
        </view>
      </view>

      <!-- 空状态 -->
      <view class="empty-state" v-else-if="!isLoading">
        <text class="empty-title">还没有瓶子</text>
        <text class="empty-desc">去海里扔一个或捡一个吧</text>
      </view>

    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { onShow, onReachBottom } from '@dcloudio/uni-app'
import { useUserStore } from '@/stores/user'
import { useBottleStore } from '@/stores/bottle'
import { formatTime, getBottleTypeName, getBottleTypeClass, getCategoryName } from '@/utils/format'
import type { Bottle } from '@/types'

const userStore = useUserStore()
const bottleStore = useBottleStore()
const bottles = ref<Bottle[]>([])
const activeTab = ref<'all' | 'thrown' | 'picked' | 'secret'>('all')
const allBottles = ref<Bottle[]>([])
const stats = ref({
  thrown: 0,
  picked: 0,
  replied: 0
})

// 分页相关 - 分别管理两个列表
const myBottlesPage = ref(1)
const pickedBottlesPage = ref(1)
const hasMoreMyBottles = ref(true)
const hasMorePickedBottles = ref(true)
const pageSize = 10
const isLoading = ref(false)

// 计算是否还有更多数据
const hasMore = computed(() => hasMoreMyBottles.value || hasMorePickedBottles.value)

onMounted(async () => {
  if (!userStore.user) {
    try {
      await userStore.getProfile()
    } catch (error) {
      console.log('未登录')
    }
  }
})

// 每次进入页面都刷新瓶子列表
onShow(() => {
  resetAndLoad()
})

// 页面滚动到底部时加载更多
onReachBottom(() => {
  loadMore()
})

// 删除/扔回海里：捡到的瓶子扔回海里，扔出的瓶子删除
const handleDelete = (bottle: any) => {
  const isPicked = !!bottle.pickedAt
  uni.showModal({
    title: isPicked ? '扔回海里' : '删除瓶子',
    content: isPicked ? '扔回后将从“捡到的”列表移除，确定吗？' : '删除后不可恢复，确定删除吗？',
    success: async (res) => {
      if (!res.confirm) return
      try {
        if (isPicked) {
          await bottleStore.throwBackBottle(bottle.id)
        } else {
          await bottleStore.deleteBottle(bottle.id)
        }
        allBottles.value = allBottles.value.filter((b: any) => b.id !== bottle.id)
        filterBottles()
        const key = isPicked ? 'picked' : 'thrown'
        stats.value[key] = Math.max(0, stats.value[key] - 1)
        uni.showToast({ title: isPicked ? '已扔回海里' : '已删除', icon: 'success' })
      } catch (error: any) {
        uni.showToast({ title: error.message || '操作失败', icon: 'none' })
      }
    }
  })
}

// 重置并加载
const resetAndLoad = async () => {
  myBottlesPage.value = 1
  pickedBottlesPage.value = 1
  hasMoreMyBottles.value = true
  hasMorePickedBottles.value = true
  allBottles.value = []
  bottles.value = []
  await loadBottles()
}

// 加载瓶子列表
const loadBottles = async () => {
  if (isLoading.value || !hasMore.value) return

  isLoading.value = true
  try {
    const isFirstPage = myBottlesPage.value === 1 && pickedBottlesPage.value === 1

    // 并行请求两个接口
    const [myBottlesRes, pickedRes] = await Promise.all([
      hasMoreMyBottles.value ? bottleStore.getMyBottles(myBottlesPage.value, pageSize) : Promise.resolve({ list: [], total: 0, page: 0 }),
      hasMorePickedBottles.value ? bottleStore.getPickedBottles(pickedBottlesPage.value, pageSize) : Promise.resolve({ list: [], total: 0, page: 0 })
    ])

    const myBottles = myBottlesRes.list || []
    const picked = pickedRes.list || []

    // 更新分页状态
    if (hasMoreMyBottles.value) {
      hasMoreMyBottles.value = myBottles.length >= pageSize
      myBottlesPage.value++
    }
    if (hasMorePickedBottles.value) {
      hasMorePickedBottles.value = picked.length >= pageSize
      pickedBottlesPage.value++
    }

    // 合并并排序
    const newBottles = [...myBottles, ...picked].sort((a, b) => {
      const timeA = (a as any).pickedAt || a.createdAt
      const timeB = (b as any).pickedAt || b.createdAt
      return new Date(timeB).getTime() - new Date(timeA).getTime()
    })

    // 追加到现有列表
    allBottles.value = [...allBottles.value, ...newBottles]

    // 更新统计（只在第一页时更新）
    if (isFirstPage) {
      stats.value.thrown = myBottlesRes.total || 0
      stats.value.picked = pickedRes.total || 0
    }

    filterBottles()
  } catch (error) {
    console.error(error)
    if (myBottlesPage.value === 1 && pickedBottlesPage.value === 1) {
      bottles.value = []
    }
  } finally {
    isLoading.value = false
  }
}

// 加载更多
const loadMore = async () => {
  if (!hasMore.value || isLoading.value) return
  await loadBottles()
}

const filterBottles = () => {
  switch (activeTab.value) {
    case 'all':
      bottles.value = allBottles.value
      break
    case 'thrown':
      bottles.value = allBottles.value.filter(b => !b.isPicked)
      break
    case 'picked':
      bottles.value = allBottles.value.filter(b => b.isPicked)
      break
    case 'secret':
      bottles.value = allBottles.value.filter(b => b.bottleType === 3)
      break
  }
}

const switchTab = (tab: 'all' | 'thrown' | 'picked' | 'secret') => {
  activeTab.value = tab
  filterBottles()
}

const goToBottleDetail = (id: number) => {
  uni.navigateTo({ url: `/pages/bottle-detail/bottle-detail?id=${id}` })
}

const getSenderAge = (user: any) => {
  if (!user?.birthday) return 0
  const birth = new Date(user.birthday)
  if (isNaN(birth.getTime())) return 0
  const now = new Date()
  let age = now.getFullYear() - birth.getFullYear()
  const m = now.getMonth() - birth.getMonth()
  if (m < 0 || (m === 0 && now.getDate() < birth.getDate())) {
    age--
  }
  return age > 0 ? age : 0
}
</script>

<style lang="scss" scoped>
.page {
  min-height: 100vh;
  background: $bg-dark;
}

.content {
  padding-top: calc(var(--status-bar-height, 44px) + 20rpx);
  padding-bottom: $tabbar-height;
}

.stats-row {
  display: flex;
  justify-content: space-around;
  gap: $spacing-sm;
  padding: $spacing-md;
}

.stat-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  background: $bg-card;
  border-radius: $radius-lg;
  padding: $spacing-md $spacing-sm;
  border: 1px solid $border-color;
}

.stat-value {
  font-size: $font-xl;
  font-weight: 700;
  color: $primary-color;
  margin-bottom: $spacing-xs;
}

.stat-label {
  font-size: $font-xs;
  color: $text-muted;
}

.section-title {
  font-size: $font-lg;
  font-weight: 700;
  color: $text-primary;
  padding: $spacing-md;
  padding-left: $spacing-md;
}

.tabs {
  display: flex;
  padding: 0 $spacing-md;
  margin-bottom: $spacing-md;
  border-bottom: 1px solid $border-color;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: $spacing-sm 0;
  position: relative;

  text {
    font-size: $font-sm;
    color: $text-muted;
  }

  &.active text {
    color: $primary-color;
    font-weight: 600;
  }
}

.tab-indicator {
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 40rpx;
  height: 4rpx;
  background: $primary-color;
  border-radius: 2rpx;
}

.bottle-list {
  padding: 0 $spacing-md;
  display: flex;
  flex-direction: column;
  gap: $spacing-sm;
}

.bottle-item {
  background: $bg-card;
  border-radius: $radius-lg;
  padding: $spacing-md;
  border: 1px solid $border-color;

  &:active {
    background: $bg-card-hover;
  }
}

// 发送者信息
.sender-info {
  display: flex;
  align-items: center;
  gap: $spacing-sm;
  margin-bottom: $spacing-sm;
}

.sender-avatar {
  width: 64rpx;
  height: 64rpx;
  background: linear-gradient(135deg, $primary-color 0%, #c49a38 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  overflow: hidden;
}

.sender-avatar-image {
  width: 100%;
  height: 100%;
}

.avatar-text {
  font-size: $font-base;
  font-weight: 700;
  color: #0a1628;
}

.sender-detail {
  flex: 1;
  min-width: 0;
}

.sender-name-row {
  display: flex;
  align-items: center;
  gap: $spacing-xs;
  margin-bottom: 4rpx;
}

.sender-name {
  font-size: $font-sm;
  font-weight: 600;
  color: $text-primary;
}

.sender-tags {
  display: flex;
  gap: $spacing-xs;
}

.tag-age {
  padding: 2rpx 10rpx;
  background: rgba(212, 168, 67, 0.15);
  border-radius: 20rpx;
}

.tag-age-text {
  font-size: $font-xs;
  color: $primary-color;
}

.tag-gender {
  padding: 2rpx 10rpx;
  border-radius: 20rpx;

  &.male {
    background: rgba(100, 181, 246, 0.15);
  }

  &.female {
    background: rgba(240, 98, 146, 0.15);
  }
}

.tag-gender-text {
  font-size: $font-xs;
  color: $text-primary;
}

.sender-location {
  font-size: $font-xs;
  color: $text-muted;
}

.bottle-header {
  margin-bottom: $spacing-sm;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.item-delete {
  padding: 4rpx 8rpx;

  .item-delete-text {
    font-size: 28rpx;
    opacity: 0.6;
  }

  &:active {
    opacity: 0.4;
  }
}

.bottle-tags {
  display: flex;
  gap: $spacing-xs;
}

.type-tag {
  font-size: $font-xs;
  color: $text-muted;
  letter-spacing: 1rpx;
}

.bottle-content {
  font-size: $font-base;
  color: $text-primary;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: $spacing-sm;
}

.bottle-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.bottle-time {
  font-size: $font-xs;
  color: $text-muted;
}

.bottle-reply {
  font-size: $font-xs;
  color: $primary-color;
}

.bottle-status {
  font-size: $font-xs;
  color: $text-muted;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding-top: 100rpx;
}

.empty-title {
  font-size: $font-lg;
  font-weight: 600;
  color: $text-primary;
  margin-bottom: $spacing-sm;
}

.empty-desc {
  font-size: $font-sm;
  color: $text-muted;
}

</style>
