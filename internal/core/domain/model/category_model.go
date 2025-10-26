package model

import "time"

type Category struct {
	ID          int        `gorm:"id"`
	CreatedByID int64      `gorm:"created_by_id"`
	User        User       `gorm:"foreignKey:CreatedByID"`
	Title       string     `gorm:"title"`
	Slug        string     `gorm:"slug"`
	CreatedAt   time.Time  `gorm:"created_at"`
	UpdatedAt   *time.Time `gorm:"updated_at"`
}
