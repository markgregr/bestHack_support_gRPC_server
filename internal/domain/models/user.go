package models

type User struct {
	ID               int64   `gorm:"primaryKey" index:"idx_id" json:"id"`
	Email            string  `gorm:"unique" index:"idx_email" json:"email"`
	PassHash         []byte  `gorm:"not null" json:"pass_hash"`
	Role             int     `gorm:"not null" json:"role"`
	Status           int     `gorm:"not null" json:"status"`
	AvarageDuration  float32 `json:"avarage_duration"`
	TelegramUsername string  `json:"telegram_username"`
}
