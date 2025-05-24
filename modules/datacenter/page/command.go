package page

import (
	"github.com/andycai/goapi/models"
)

// CommandAddPage adds a new page
func CommandAddPage(page *models.Page) error {
	return addPage(page)
}

// CommandUpdatePage updates an existing page
func CommandUpdatePage(page *models.Page) error {
	return updatePage(page)
}

// CommandDeletePage deletes a page
func CommandDeletePage(id int64) error {
	return deletePage(id)
}
