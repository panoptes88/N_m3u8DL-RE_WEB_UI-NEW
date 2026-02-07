package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:64" json:"username"`
	Password  string    `gorm:"size:256" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

func (User) TableName() string {
	return "users"
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Session 用户会话
type Session struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Token     string    `gorm:"uniqueIndex;size:64" json:"-"`
	UserID    uint      `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (Session) TableName() string {
	return "sessions"
}

// CreateSession 创建会话，生成随机 token
func CreateSession(db *gorm.DB, userID uint) (*Session, error) {
	token := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	session := &Session{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}

	if err := db.Create(session).Error; err != nil {
		return nil, err
	}

	return session, nil
}

// GetUserByToken 通过 token 获取用户（使用 JOIN 单次查询）
func GetUserByToken(db *gorm.DB, token string) (*User, error) {
	var user User
	// 使用 JOIN 在一次查询中获取用户信息
	err := db.Table("sessions").
		Select("users.*").
		Joins("JOIN users ON sessions.user_id = users.id").
		Where("sessions.token = ? AND sessions.expires_at > ?", token, time.Now()).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteSession 删除会话（登出）
func DeleteSession(db *gorm.DB, token string) error {
	return db.Where("token = ?", token).Delete(&Session{}).Error
}

// DeleteExpiredSessions 删除过期会话
func DeleteExpiredSessions(db *gorm.DB) error {
	return db.Where("expires_at < ?", time.Now()).Delete(&Session{}).Error
}
