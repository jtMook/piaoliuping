export interface User {
  id: number
  openid: string
  nickname: string
  avatarUrl: string
  gender: number
  birthday: string
  provinceCode: string
  provinceName: string
  cityCode: string
  cityName: string
  districtCode: string
  districtName: string
  interestTags: string[]
  status: number
  bottleCount: number
  addressUpdatedAt: string
  createdAt: string
  updatedAt: string
}

export interface Bottle {
  id: number
  userId: number
  content: string
  mediaUrls: string[]
  tags: string[]
  bottleType: number
  provinceCode: string
  cityCode: string
  districtCode: string
  cityName: string
  status: number
  isPicked: boolean
  pickedCount: number
  visibleScope: number
  expireAt: string
  createdAt: string
  updatedAt: string
  canReply: boolean
  replyCount: number
  user?: {
    id: number
    nickname: string
    avatarUrl: string
    gender: number
    birthday: string
    cityName: string
    districtName: string
  }
}

// 捡瓶子返回的数据（扁平结构）
export interface PickedBottle {
  id: number
  content: string
  mediaUrls: string[]
  tags: string[]
  bottleType: number
  provinceName: string
  cityName: string
  districtName: string
  nickname: string
  gender: number
  createdAt: string
  canReply: boolean
}

export interface Conversation {
  id: number
  bottleId: number
  userAId: number
  userBId: number
  lastMessageId: number
  lastMessageAt: string
  isClosed: boolean
  createdAt: string
  updatedAt: string
  otherUser: User
  unreadCount: number
  lastMessage: Message
}

export interface ConversationDetail {
  id: number
  bottleId: number
  otherUser: {
    id: number
    nickname: string
    avatarUrl: string
    cityName: string
    gender: number
  }
  bottle: {
    id: number
    content: string
    bottleType: number
  }
  isClosed: boolean
}

export interface Message {
  id: number
  conversationId: number
  senderId: number
  content: string
  mediaUrl: string
  msgType: number
  isRead: boolean
  createdAt: string
}

export interface Report {
  id: number
  reporterId: number
  targetUserId: number
  targetBottleId: number
  targetMessageId: number
  reason: number
  description: string
  status: number
  handledBy: number
  handledAt: string
  createdAt: string
}

export interface Blacklist {
  id: number
  userId: number
  blockedUserId: number
  createdAt: string
}

export interface Address {
  code: string
  name: string
  children?: Address[]
}

export type BottleType = 1 | 2 | 3
export type VisibleScope = 1 | 2
export type MessageType = 1 | 2 | 3
