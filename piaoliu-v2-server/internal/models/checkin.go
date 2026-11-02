package models

import "time"

type CheckinRecord struct {
	ID              int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          int64     `gorm:"column:user_id;not null;index" json:"userId"`
	ConsecutiveDays int64     `gorm:"column:consecutive_days;default:1" json:"consecutiveDays"`
	BottleReward    int64     `gorm:"column:bottle_reward;default:1" json:"bottleReward"`
	CreatedAt       time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (CheckinRecord) TableName() string {
	return "checkin_records"
}
