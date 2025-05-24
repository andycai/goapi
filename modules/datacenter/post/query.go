package post

// 获取文章列表
func QueryPostList(page, limit int, status string) ([]PostResponse, int64, error) {
	posts, total, err := getPostList(page, limit, status)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	response := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		postResp := PostResponse{
			ID:        int64(post.ID),
			Title:     post.Title,
			Content:   post.Content,
			Slug:      post.Slug,
			Status:    post.Status,
			AuthorID:  int64(post.AuthorID),
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		response = append(response, postResp)
	}

	return response, total, nil
}

// 根据ID获取文章
func QueryPostByID(id uint) (*PostResponse, error) {
	post, err := getPostByID(id)
	if err != nil {
		return nil, err
	}

	response := &PostResponse{
		ID:        int64(post.ID),
		Title:     post.Title,
		Content:   post.Content,
		Slug:      post.Slug,
		Status:    post.Status,
		AuthorID:  int64(post.AuthorID),
		CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, nil
}

// 根据别名获取文章
func QueryPostBySlug(slug string) (*PostResponse, error) {
	post, err := getPostBySlug(slug)
	if err != nil {
		return nil, err
	}

	response := &PostResponse{
		ID:        int64(post.ID),
		Title:     post.Title,
		Content:   post.Content,
		Slug:      post.Slug,
		Status:    post.Status,
		AuthorID:  int64(post.AuthorID),
		CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return response, nil
}

// 搜索文章
func QuerySearchPosts(query string, page, limit int) ([]PostResponse, int64, error) {
	posts, total, err := searchPosts(query, page, limit)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	response := make([]PostResponse, 0, len(posts))
	for _, post := range posts {
		postResp := PostResponse{
			ID:        int64(post.ID),
			Title:     post.Title,
			Content:   post.Content,
			Slug:      post.Slug,
			Status:    post.Status,
			AuthorID:  int64(post.AuthorID),
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		response = append(response, postResp)
	}

	return response, total, nil
}
