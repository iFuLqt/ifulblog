package model

import "time"

type Content struct {
	ID          int        `gorm:"id"`
	CategoryID  int64      `gorm:"category_by_id"`
	CreatedByID int64      `gorm:"created_by_id"`
	User        User       `gorm:"foreignKey:CreatedByID"`
	Category    Category   `gorm:"foreignKey:CategoryID"`
	Title       string     `gorm:"title"`
	Excerpt     string     `gorm:"excerpt"`
	Description string     `gorm:"description"`
	Image       string     `gorm:"image"`
	Tags        string     `gorm:"tags"`
	Status      string     `gorm:"status"`
	CreatedAt   time.Time  `gorm:"created_at"`
	UpdatedAt   *time.Time `gorm:"updated_at"`
}
