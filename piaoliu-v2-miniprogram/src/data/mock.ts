import type { User, Bottle, Conversation, Message } from '@/types'

export const mockUser: User = {
  id: 1,
  openid: 'mock_openid_001',
  nickname: '海风',
  avatarUrl: '',
  provinceCode: '440000',
  provinceName: '广东省',
  cityCode: '440300',
  cityName: '深圳市',
  districtCode: '440305',
  districtName: '南山区',
  status: 1,
  addressUpdatedAt: new Date().toISOString(),
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString()
}

export const mockBottles: Bottle[] = [
  {
    id: 1,
    userId: 2,
    content: '今天心情不好，想找人聊聊。工作压力好大，感觉喘不过气来...',
    mediaUrls: [],
    bottleType: 1,
    provinceCode: '440000',
    cityCode: '440300',
    districtCode: '440305',
    cityName: '深圳市',
    status: 1,
    isPicked: false,
    pickedCount: 0,
    visibleScope: 1,
    expireAt: '',
    createdAt: new Date(Date.now() - 3600000).toISOString(),
    updatedAt: new Date(Date.now() - 3600000).toISOString(),
    canReply: true
  },
  {
    id: 2,
    userId: 3,
    content: '有没有人和我一样喜欢深夜听歌？戴上耳机，整个世界都是我的。',
    mediaUrls: [],
    bottleType: 1,
    provinceCode: '440000',
    cityCode: '440100',
    districtCode: '440105',
    cityName: '广州市',
    status: 1,
    isPicked: false,
    pickedCount: 0,
    visibleScope: 2,
    expireAt: '',
    createdAt: new Date(Date.now() - 7200000).toISOString(),
    updatedAt: new Date(Date.now() - 7200000).toISOString(),
    canReply: true
  },
  {
    id: 3,
    userId: 4,
    content: '一个人在这座城市打拼了三年，有时候真的很想家...',
    mediaUrls: [],
    bottleType: 1,
    provinceCode: '310000',
    cityCode: '310100',
    districtCode: '310115',
    cityName: '上海市',
    status: 1,
    isPicked: false,
    pickedCount: 0,
    visibleScope: 1,
    expireAt: '',
    createdAt: new Date(Date.now() - 1800000).toISOString(),
    updatedAt: new Date(Date.now() - 1800000).toISOString(),
    canReply: true
  },
  {
    id: 4,
    userId: 5,
    content: '今天阳光很好，想出去走走，但不知道去哪里...',
    mediaUrls: [],
    bottleType: 1,
    provinceCode: '330000',
    cityCode: '330100',
    districtCode: '',
    cityName: '杭州市',
    status: 1,
    isPicked: false,
    pickedCount: 0,
    visibleScope: 1,
    expireAt: '',
    createdAt: new Date(Date.now() - 900000).toISOString(),
    updatedAt: new Date(Date.now() - 900000).toISOString(),
    canReply: true
  },
  {
    id: 5,
    userId: 6,
    content: '刚刚失恋了，四年的感情说散就散。不知道该怎么面对明天...',
    mediaUrls: [],
    bottleType: 3,
    provinceCode: '440000',
    cityCode: '440300',
    districtCode: '440306',
    cityName: '深圳市',
    status: 1,
    isPicked: false,
    pickedCount: 0,
    visibleScope: 1,
    expireAt: new Date(Date.now() + 86400000).toISOString(),
    createdAt: new Date(Date.now() - 600000).toISOString(),
    updatedAt: new Date(Date.now() - 600000).toISOString(),
    canReply: true
  }
]

export const mockConversations: Conversation[] = [
  {
    id: 1,
    bottleId: 1,
    userAId: 2,
    userBId: 1,
    lastMessageId: 3,
    lastMessageAt: new Date(Date.now() - 300000).toISOString(),
    isClosed: false,
    createdAt: new Date(Date.now() - 3600000).toISOString(),
    updatedAt: new Date(Date.now() - 300000).toISOString(),
    otherUser: {
      id: 2,
      openid: 'mock_openid_002',
      nickname: '星空',
      avatarUrl: '',
      provinceCode: '440000',
      provinceName: '广东省',
      cityCode: '440300',
      cityName: '深圳市',
      districtCode: '440305',
      districtName: '南山区',
      status: 1,
      addressUpdatedAt: new Date().toISOString(),
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    },
    unreadCount: 2,
    lastMessage: {
      id: 3,
      conversationId: 1,
      senderId: 2,
      content: '谢谢你的安慰，感觉好多了',
      mediaUrl: '',
      msgType: 1,
      isRead: false,
      createdAt: new Date(Date.now() - 300000).toISOString()
    }
  },
  {
    id: 2,
    bottleId: 2,
    userAId: 3,
    userBId: 1,
    lastMessageId: 6,
    lastMessageAt: new Date(Date.now() - 1800000).toISOString(),
    isClosed: false,
    createdAt: new Date(Date.now() - 7200000).toISOString(),
    updatedAt: new Date(Date.now() - 1800000).toISOString(),
    otherUser: {
      id: 3,
      openid: 'mock_openid_003',
      nickname: '夜猫子',
      avatarUrl: '',
      provinceCode: '440000',
      provinceName: '广东省',
      cityCode: '440100',
      cityName: '广州市',
      districtCode: '440105',
      districtName: '海珠区',
      status: 1,
      addressUpdatedAt: new Date().toISOString(),
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    },
    unreadCount: 0,
    lastMessage: {
      id: 6,
      conversationId: 2,
      senderId: 1,
      content: '推荐你听一下陈粒的歌',
      mediaUrl: '',
      msgType: 1,
      isRead: true,
      createdAt: new Date(Date.now() - 1800000).toISOString()
    }
  }
]

export const mockMessages: Message[] = [
  {
    id: 1,
    conversationId: 1,
    senderId: 1,
    content: '抱抱你，我理解你的感受。有时候压力大的时候，不妨停下来休息一下。',
    mediaUrl: '',
    msgType: 1,
    isRead: true,
    createdAt: new Date(Date.now() - 3500000).toISOString()
  },
  {
    id: 2,
    conversationId: 1,
    senderId: 2,
    content: '谢谢你，你是第一个愿意听我说这些的人',
    mediaUrl: '',
    msgType: 1,
    isRead: true,
    createdAt: new Date(Date.now() - 3400000).toISOString()
  },
  {
    id: 3,
    conversationId: 1,
    senderId: 2,
    content: '谢谢你的安慰，感觉好多了',
    mediaUrl: '',
    msgType: 1,
    isRead: false,
    createdAt: new Date(Date.now() - 300000).toISOString()
  }
]
