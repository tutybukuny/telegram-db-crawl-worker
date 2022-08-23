package entity

import "time"

type Channel struct {
	ID        int64 `gorm:"primaryKey;autoIncrement:false"`
	Name      string
	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time `gorm:"index"`
}

func (Channel) TableName() string {
	return "channels"
}

func (Channel) IDField() string {
	return "id"
}
