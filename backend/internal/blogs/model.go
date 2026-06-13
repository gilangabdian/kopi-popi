package blogs

import (
	"time"
)

// Blog merepresentasikan tabel blogs
type Blog struct {
	ID                    string    `json:"id" gorm:"primaryKey;column:id"`
	AuthorID              string    `json:"author_id" gorm:"column:author_id"`
	Title                 string    `json:"title" gorm:"column:title"`
	Slug                  string    `json:"slug" gorm:"column:slug"`
	Content               string    `json:"content" gorm:"column:content"`
	CoverImage            *string   `json:"cover_image" gorm:"column:cover_image"`
	EstimatedReadTimeMins int       `json:"estimated_read_time_mins" gorm:"column:estimated_read_time_mins"`
	CreatedAt             time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt             time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

type CreateBlogRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateBlogRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}
