package parameter

import (
	"encoding/json"
)

// 获取参数列表
func QueryParameters(limit, page int, search string) ([]ParameterResponse, int64, error) {
	parameters, total, err := getParameters(limit, page, search)
	if err != nil {
		return nil, 0, err
	}

	// 转换为响应格式
	response := make([]ParameterResponse, 0, len(parameters))
	for _, param := range parameters {
		paramResp := ParameterResponse{
			ID:        param.ID,
			Type:      param.Type,
			Name:      param.Name,
			CreatedAt: param.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: param.UpdatedAt.Format("2006-01-02 15:04:05"),
			CreatedBy: param.CreatedBy,
			UpdatedBy: param.UpdatedBy,
		}

		// 解析JSON参数
		var fields []ParameterField
		if param.Parameters != "" {
			if err := json.Unmarshal([]byte(param.Parameters), &fields); err != nil {
				// 解析失败时使用空数组
				fields = []ParameterField{}
			}
		} else {
			fields = []ParameterField{}
		}

		paramResp.Parameters = fields
		response = append(response, paramResp)
	}

	return response, total, nil
}

// 获取单个参数
func QueryParameter(id uint) (*ParameterResponse, error) {
	param, err := getParameter(id)
	if err != nil {
		return nil, err
	}

	response := &ParameterResponse{
		ID:        param.ID,
		Type:      param.Type,
		Name:      param.Name,
		CreatedAt: param.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: param.UpdatedAt.Format("2006-01-02 15:04:05"),
		CreatedBy: param.CreatedBy,
		UpdatedBy: param.UpdatedBy,
	}

	// 解析JSON参数
	var fields []ParameterField
	if param.Parameters != "" {
		if err := json.Unmarshal([]byte(param.Parameters), &fields); err != nil {
			// 解析失败时使用空数组
			fields = []ParameterField{}
		}
	} else {
		fields = []ParameterField{}
	}

	response.Parameters = fields
	return response, nil
}
