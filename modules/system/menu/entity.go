package menu

import "github.com/andycai/goapi/models"

const (
	MenuIdSystem = 1000
	MenuIdTools  = 2000
	MenuIdGame   = 3000
	MenuIdWebApp = 4000
)

// MenuTree 菜单树结构
type MenuTree struct {
	Menu     *models.Menu `json:"menu"`
	Children []*MenuTree  `json:"children"`
}
