package post

import (
	"github.com/andycai/goapi/models"
)

// CommandAddPost adds a new post
func CommandAddPost(post *models.Post) error {
	return addPost(post)
}

// CommandUpdatePost updates an existing post
func CommandUpdatePost(post *models.Post) error {
	return updatePost(post)
}

// CommandDeletePost deletes a post
func CommandDeletePost(id int64) error {
	return deletePost(id)
}
