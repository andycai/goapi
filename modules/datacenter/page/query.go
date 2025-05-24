package page

// 获取页面列表
func QueryPageList(page, limit int, status string) ([]PageResponse, int64, error) {
	pages, total, err := getPageList(page, limit, status)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	response := make([]PageResponse, 0, len(pages))
	for _, page := range pages {
		pageResp := PageResponse{
			ID:        page.ID,
			Title:     page.Title,
			Content:   page.Content,
			Slug:      page.Slug,
			Status:    page.Status,
			AuthorID:  0,
			CreatedAt: page.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: page.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		response = append(response, pageResp)
	}

	return response, total, nil
}

// 根据ID获取页面
func QueryPageByID(id uint) (*PageResponse, error) {
	page, err := getPageByID(id)
	if err != nil {
		return nil, err
	}

	response := &PageResponse{
		ID:        page.ID,
		Title:     page.Title,
		Content:   page.Content,
		Slug:      page.Slug,
		Status:    page.Status,
		AuthorID:  0,
		CreatedAt: page.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: page.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, nil
}

// 根据别名获取页面
func QueryPageBySlug(slug string) (*PageResponse, error) {
	page, err := getPageBySlug(slug)
	if err != nil {
		return nil, err
	}

	response := &PageResponse{
		ID:        page.ID,
		Title:     page.Title,
		Content:   page.Content,
		Slug:      page.Slug,
		Status:    page.Status,
		AuthorID:  0,
		CreatedAt: page.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: page.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, nil
}

// 搜索页面
func QuerySearchPages(query string, page, limit int) ([]PageResponse, int64, error) {
	pages, total, err := searchPages(query, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	response := make([]PageResponse, 0, len(pages))
	for _, page := range pages {
		pageResp := PageResponse{
			ID:        page.ID,
			Title:     page.Title,
			Content:   page.Content,
			Slug:      page.Slug,
			Status:    page.Status,
			AuthorID:  0,
			CreatedAt: page.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: page.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		response = append(response, pageResp)
	}

	return response, total, nil
}
