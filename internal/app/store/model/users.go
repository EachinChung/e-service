package model

import (
	"time"

	"gorm.io/gorm"

	"github.com/eachinchung/component-base/auth"
	"github.com/eachinchung/component-base/utils/id"
	"github.com/eachinchung/errors"
)

// Users 用户表
type Users struct {
	ID           uint           `gorm:"primaryKey;column:id" json:"-" redis:"id"`
	UserID       uint64         `gorm:"column:user_id" json:"user_id" redis:"user_id"`          // 用户ID
	Phone        string         `gorm:"column:phone" json:"phone" redis:"phone"`                // 手机号
	Email        *string        `gorm:"column:email" json:"email,omitempty" redis:"email"`      // 邮箱
	Username     string         `gorm:"column:username" json:"username" redis:"username"`       // 用户名
	PasswordHash string         `gorm:"column:password_hash" json:"-" redis:"password_hash"`    // 密码
	Avatar       *string        `gorm:"column:avatar" json:"avatar,omitempty" redis:"avatar"`   // 头像
	State        Status         `gorm:"column:state" json:"state" redis:"state"`                // 状态
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at" redis:"created_at"` // 创建时间
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at" redis:"updated_at"` // 更新时间
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at" json:"-" redis:"deleted_at"`          // 删除时间
}

// ComparePasswordHash with the plain text password. Returns true if it's the same as the encrypted one (in the `Users` struct).
func (u *Users) ComparePasswordHash(pwd string) error {
	if err := auth.ComparePasswordHash(u.PasswordHash, pwd); err != nil {
		return errors.Wrap(err, "failed to compile password")
	}

	return nil
}

func (u *Users) BeforeCreate(_ *gorm.DB) (err error) {
	u.UserID = id.GenUint64ID()
	return
}
