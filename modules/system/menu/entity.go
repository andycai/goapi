package menu

import "github.com/andycai/goapi/models"

// MenuTree 菜单树结构
type MenuTree struct {
	Menu     *models.Menu `json:"menu"`
	Children []*MenuTree  `json:"children"`
}
