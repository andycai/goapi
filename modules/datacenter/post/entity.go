package post

import "time"

// Post represents a blog post entity
type Post struct {
	ID        int64     `json:"id"`         // ID
	Title     string    `json:"title"`      // 标题
	Content   string    `json:"content"`    // 内容
	Slug      string    `json:"slug"`       // 别名
	Status    string    `json:"status"`     // 状态：draft, published
	AuthorID  int64     `json:"author_id"`  // 作者ID
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// PostCreate represents the data needed to create a new post
type PostCreate struct {
	Title    string `json:"title" validate:"required"`
	Content  string `json:"content" validate:"required"`
	Slug     string `json:"slug" validate:"required"`
	Status   string `json:"status" validate:"required,oneof=draft published"`
	AuthorID int64  `json:"author_id" validate:"required"`
}

// PostUpdate represents the data needed to update an existing post
type PostUpdate struct {
	Title    string `json:"title" validate:"required"`
	Content  string `json:"content" validate:"required"`
	Slug     string `json:"slug" validate:"required"`
	Status   string `json:"status" validate:"required,oneof=draft published"`
	AuthorID int64  `json:"author_id" validate:"required"`
}
