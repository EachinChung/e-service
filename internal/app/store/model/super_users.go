package model

import "time"

// SuperUsers 用户表
type SuperUsers struct {
	ID        uint      `gorm:"primaryKey;column:id" json:"-" redis:"id"`
	Username  string    `gorm:"column:username" json:"username" redis:"username"`       // 用户名
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at" redis:"created_at"` // 创建时间
}
