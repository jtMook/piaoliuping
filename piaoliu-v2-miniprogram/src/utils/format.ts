export const formatTime = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 30) return `${days}天前`
  return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')}`
}

export const formatDateTime = (dateString: string) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}-${date.getDate().toString().padStart(2, '0')} ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`
}

export const getBottleTypeName = (type: number) => {
  const names: Record<number, string> = {
    1: '普通瓶',
    2: '语音瓶',
    3: '秘密瓶'
  }
  return names[type] || '普通瓶'
}

export const getBottleTypeClass = (type: number) => {
  const classes: Record<number, string> = {
    1: 'normal',
    2: 'voice',
    3: 'secret'
  }
  return classes[type] || 'normal'
}

export const getCategoryName = (type: number) => {
  const names: Record<number, string> = {
    1: '日常',
    2: '深夜emo',
    3: '情感'
  }
  return names[type] || '日常'
}
