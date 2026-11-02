package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// StringArray 自定义类型，用于 JSON 数组存储
type StringArray []string

func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return nil, nil
	}
	b, err := json.Marshal(s)
	return string(b), err
}

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = nil
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case string:
		bytes = []byte(v)
	case []byte:
		bytes = v
	}
	return json.Unmarshal(bytes, s)
}

type Bottle struct {
	ID           int64       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       int64       `gorm:"column:user_id;index;not null" json:"userId"`
	Content      string      `gorm:"type:text" json:"content"`
	MediaUrls    StringArray `gorm:"column:media_urls;type:json" json:"mediaUrls"`
	Tags         StringArray `gorm:"column:tags;type:json" json:"tags"`
	BottleType   int8        `gorm:"column:bottle_type;default:1" json:"bottleType"` // 1普通瓶 2语音瓶 3秘密瓶
	ProvinceCode string      `gorm:"column:province_code;size:20;not null" json:"provinceCode"`
	CityCode     string      `gorm:"column:city_code;size:20;not null" json:"cityCode"`
	DistrictCode string      `gorm:"column:district_code;size:20" json:"districtCode"`
	CityName     string      `gorm:"-" json:"cityName"`       // 从 User 获取
	Status       int8        `gorm:"default:0" json:"status"` // 0待审核 1已通过 2未通过 3已删除
	IsPicked     bool        `gorm:"column:is_picked;default:false" json:"isPicked"`
	PickedCount  int         `gorm:"column:picked_count;default:0" json:"pickedCount"`
	VisibleScope int8        `gorm:"column:visible_scope;default:1" json:"visibleScope"` // 1同城 2同省
	ExpireAt     *time.Time  `gorm:"column:expire_at" json:"expireAt"`
	CreatedAt    time.Time   `gorm:"column:created_at;autoCreateTime;index" json:"createdAt"`
	UpdatedAt    time.Time   `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (Bottle) TableName() string {
	return "bottles"
}

type BottlePickedRecord struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BottleID  int64     `gorm:"column:bottle_id;index;not null" json:"bottleId"`
	UserID    int64     `gorm:"column:user_id;index;not null" json:"userId"`
	IsReplied bool      `gorm:"column:is_replied;default:false" json:"isReplied"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (BottlePickedRecord) TableName() string {
	return "bottle_picked_records"
}

type BottleReply struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	BottleID  int64     `gorm:"column:bottle_id;index;not null" json:"bottleId"`
	UserID    int64     `gorm:"column:user_id;index;not null" json:"userId"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (BottleReply) TableName() string {
	return "bottle_replies"
}
