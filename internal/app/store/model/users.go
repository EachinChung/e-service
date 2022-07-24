package model

import (
	"time"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"

	"github.com/eachinchung/component-base/auth"
	"github.com/eachinchung/errors"
)

const ctxKey = "USER"

// Users 用户表
type Users struct {
	ID           uint           `gorm:"primaryKey;column:id" json:"-" redis:"id"`
	Phone        string         `gorm:"column:phone" json:"phone" redis:"phone"`                // 手机号
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

func (u *Users) Map() map[string]interface{} {
	return map[string]interface{}{
		"id":            u.ID,
		"phone":         u.Phone,
		"username":      u.Username,
		"password_hash": u.PasswordHash,
		"avatar":        u.Avatar,
		"state":         u.State,
		"created_at":    u.CreatedAt,
		"updated_at":    u.UpdatedAt,
		"deleted_at":    u.DeletedAt,
	}
}

func (u *Users) SaveToContext(c *gin.Context) {
	c.Set(ctxKey, u)
}

func ExtractUsersFromContext(c *gin.Context) *Users {
	u, exists := c.Get(ctxKey)
	if !exists {
		return nil
	}

	return u.(*Users)
}
