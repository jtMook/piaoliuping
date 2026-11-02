package models

import "time"

type User struct {
	ID               int64       `gorm:"primaryKey;autoIncrement" json:"id"`
	OpenID           string      `gorm:"column:openid;uniqueIndex;size:64;not null" json:"-"`
	UnionID          string      `gorm:"column:unionid;size:64" json:"-"`
	Nickname         string      `gorm:"size:50;not null" json:"nickname"`
	AvatarURL        string      `gorm:"column:avatar_url;size:255" json:"avatarUrl"`
	Gender           int8        `gorm:"default:0" json:"gender"` // 0未知 1男 2女
	Birthday         string      `gorm:"column:birthday;size:20" json:"birthday"`
	ProvinceCode     string      `gorm:"column:province_code;size:20;not null" json:"provinceCode"`
	ProvinceName     string      `gorm:"column:province_name;size:50;not null" json:"provinceName"`
	CityCode         string      `gorm:"column:city_code;size:20;not null" json:"cityCode"`
	CityName         string      `gorm:"column:city_name;size:50;not null" json:"cityName"`
	DistrictCode     string      `gorm:"column:district_code;size:20" json:"districtCode"`
	DistrictName     string      `gorm:"column:district_name;size:50" json:"districtName"`
	InterestTags     StringArray `gorm:"column:interest_tags;type:json" json:"interestTags"`
	Status           int8        `gorm:"default:1" json:"status"` // 1正常 2禁言 3封号
	BottleCount      int64       `gorm:"column:bottle_count;default:0" json:"bottleCount"`
	AddressUpdatedAt *time.Time  `gorm:"column:address_updated_at" json:"addressUpdatedAt"`
	CreatedAt        time.Time   `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt        time.Time   `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (User) TableName() string {
	return "users"
}
