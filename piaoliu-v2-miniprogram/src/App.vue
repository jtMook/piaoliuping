<script setup lang="ts">
import { onLaunch, onShow, onHide } from '@dcloudio/uni-app'
import { useMessageStore } from '@/stores/message'
import { useUserStore } from '@/stores/user'

// 标记是否已初始化
let appInitialized = false

onLaunch(async () => {
  console.log('App Launch')
  const userStore = useUserStore()
  const token = uni.getStorageSync('token')
  
  if (!token) {
    uni.reLaunch({ url: '/pages/login/login' })
  } else {
    // 有 token 时，加载用户数据
    await userStore.getProfile()
  }
})

onShow(async () => {
  console.log('App Show')
  const userStore = useUserStore()
  
  // 只在首次显示时检查资料和获取消息
  if (!appInitialized) {
    appInitialized = true
    
    // 检查用户是否完善了资料
    if (userStore.user) {
      const user = userStore.user
      if (!user.gender || !user.birthday || !user.cityCode) {
        const pages = getCurrentPages()
        const currentPage = pages[pages.length - 1]
        const currentRoute = currentPage?.route || ''
        
        if (currentRoute !== 'pages/complete-profile/complete-profile' && 
            currentRoute !== 'pages/login/login') {
          uni.reLaunch({ url: '/pages/complete-profile/complete-profile' })
          return
        }
      }
    }
    
    // 首次获取消息列表
    const messageStore = useMessageStore()
    try {
      await messageStore.getConversations()
      messageStore.updateTabBarBadge()
    } catch (error) {
      console.error('Failed to update badge:', error)
    }
  }
})

onHide(() => {
  console.log('App Hide')
})
</script>

<style lang="scss">
@import './styles/global.scss';
</style>
