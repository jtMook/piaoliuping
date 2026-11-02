package models

import "time"

type Report struct {
	ID             int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ReporterID     int64      `gorm:"column:reporter_id;index;not null" json:"reporterId"`
	TargetUserID   int64      `gorm:"column:target_user_id;index;not null" json:"targetUserId"`
	TargetBottleID *int64     `gorm:"column:target_bottle_id" json:"targetBottleId"`
	TargetMessageID *int64    `gorm:"column:target_message_id" json:"targetMessageId"`
	Reason         int8       `gorm:"not null" json:"reason"` // 1色情 2广告诈骗 3人身攻击 4违法违规 5其他
	Description    string     `gorm:"size:500" json:"description"`
	Status         int8       `gorm:"default:0;index" json:"status"` // 0待处理 1通过 2删除 3封号
	HandledBy      *int64     `gorm:"column:handled_by" json:"handledBy"`
	HandledAt      *time.Time `gorm:"column:handled_at" json:"handledAt"`
	CreatedAt      time.Time  `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (Report) TableName() string {
	return "reports"
}

type Blacklist struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID        int64     `gorm:"column:user_id;index;not null" json:"userId"`
	BlockedUserID int64     `gorm:"column:blocked_user_id;not null" json:"blockedUserId"`
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (Blacklist) TableName() string {
	return "blacklists"
}
