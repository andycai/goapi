package post

import (
	"github.com/andycai/goapi/models"
)

// 获取文章列表
func QueryPostList(page, limit int, status string) ([]models.Post, int64, error) {
	return getPostList(page, limit, status)
}

// 根据ID获取文章
func QueryPostByID(id int64) (models.Post, error) {
	return getPostByID(id)
}

// 根据别名获取文章
func QueryPostBySlug(slug string) (models.Post, error) {
	return getPostBySlug(slug)
}

// 搜索文章
func QuerySearchPosts(query string, page, limit int) ([]models.Post, int64, error) {
	return searchPosts(query, page, limit)
}
