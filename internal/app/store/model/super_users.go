package model

import "time"

// SuperUsers 用户表
type SuperUsers struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"-" redis:"id"`
	EID       string    `gorm:"column:eid" json:"eid" redis:"eid"`                      // 用户名
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at" redis:"created_at"` // 创建时间
}
