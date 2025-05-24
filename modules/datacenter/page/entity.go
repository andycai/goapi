package page

import (
	"errors"
	"time"
)

var (
	ErrPageNotFound          = errors.New("页面不存在")
	ErrPageSlugAlreadyExists = errors.New("页面别名已存在")
)

// Page represents a static page entity
type Page struct {
	ID        int64     `json:"id"`         // ID
	Title     string    `json:"title"`      // 标题
	Content   string    `json:"content"`    // 内容
	Slug      string    `json:"slug"`       // 别名
	Status    string    `json:"status"`     // 状态：draft, published
	AuthorID  int64     `json:"author_id"`  // 作者ID
	CreatedAt time.Time `json:"created_at"` // 创建时间
	UpdatedAt time.Time `json:"updated_at"` // 更新时间
}

// PageCreate represents the data needed to create a new page
type PageCreate struct {
	Title    string `json:"title" validate:"required"`
	Content  string `json:"content" validate:"required"`
	Slug     string `json:"slug" validate:"required"`
	Status   string `json:"status" validate:"required,oneof=draft published"`
	AuthorID int64  `json:"author_id" validate:"required"`
}

// PageUpdate represents the data needed to update an existing page
type PageUpdate struct {
	Title    string `json:"title" validate:"required"`
	Content  string `json:"content" validate:"required"`
	Slug     string `json:"slug" validate:"required"`
	Status   string `json:"status" validate:"required,oneof=draft published"`
	AuthorID int64  `json:"author_id" validate:"required"`
}
