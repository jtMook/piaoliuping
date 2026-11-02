package models

import "time"

type Conversation struct {
	ID            int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	BottleID      int64      `gorm:"column:bottle_id;not null" json:"bottleId"`
	UserAID       int64      `gorm:"column:user_a_id;index;not null" json:"userAId"` // 扔瓶人
	UserBID       int64      `gorm:"column:user_b_id;index;not null" json:"userBId"` // 捡瓶/回复人
	LastMessageID *int64     `gorm:"column:last_message_id" json:"lastMessageId"`
	LastMessageAt *time.Time `gorm:"column:last_message_at;index" json:"lastMessageAt"`
	IsClosed      bool       `gorm:"column:is_closed;default:false" json:"isClosed"`
	DeletedByA    bool       `gorm:"column:deleted_by_a;not null;default:false" json:"-"` // userA 侧已删除（仅对自己隐藏）
	DeletedByB    bool       `gorm:"column:deleted_by_b;not null;default:false" json:"-"` // userB 侧已删除（仅对自己隐藏）
	CreatedAt     time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// 关联
	UserA    *User    `gorm:"foreignKey:UserAID" json:"userA,omitempty"`
	UserB    *User    `gorm:"foreignKey:UserBID" json:"userB,omitempty"`
	Bottle   *Bottle  `gorm:"foreignKey:BottleID" json:"bottle,omitempty"`
	LastMsg  *Message `gorm:"foreignKey:LastMessageID" json:"lastMessage,omitempty"`
}

func (Conversation) TableName() string {
	return "conversations"
}

type Message struct {
	ID             int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ConversationID int64     `gorm:"column:conversation_id;index;not null" json:"conversationId"`
	SenderID       int64     `gorm:"column:sender_id;index;not null" json:"senderId"`
	Content        string    `gorm:"type:text" json:"content"`
	MediaURL       string    `gorm:"column:media_url;size:255" json:"mediaUrl"`
	MsgType        int8      `gorm:"column:msg_type;default:1" json:"msgType"` // 1文字 2图片 3语音
	IsRead         bool      `gorm:"column:is_read;default:false" json:"isRead"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime;index" json:"createdAt"`

	// 关联
	Sender *User `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}

func (Message) TableName() string {
	return "messages"
}
