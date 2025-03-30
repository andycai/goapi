package note

type NoteRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	CategoryID uint   `json:"category_id"`
	ParentID   uint   `json:"parent_id"`
	IsPublic   uint8  `json:"is_public"`
	RoleIDs    []uint `json:"role_ids"`
}

type CategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    uint   `json:"parent_id"`
	IsPublic    uint8  `json:"is_public"`
	RoleIDs     []uint `json:"role_ids"`
}
