package service

import (
	"context"
	"hot-ai-backend/internal/access"
	"hot-ai-backend/internal/models"
	"hot-ai-backend/internal/repository"
)

// ToolService 工具服务
type ToolService struct {
	repo *repository.ToolRepository
}

// NewToolService 创建工具服务
func NewToolService(repo *repository.ToolRepository) *ToolService {
	return &ToolService{repo: repo}
}

// CategoryList 获取工具类别列表
func (s *ToolService) CategoryList(ctx context.Context) ([]models.ToolCategory, error) {
	return s.repo.GetCategories()
}

// ToolList 获取工具列表
func (s *ToolService) ToolList(ctx context.Context, params map[string]interface{}) ([]models.Tool, int64, error) {
	// 设置默认参数
	if _, ok := params["page"]; !ok {
		params["page"] = 1
	}
	if _, ok := params["page_size"]; !ok {
		params["page_size"] = 20
	}
	if _, ok := params["sort_by"]; !ok {
		params["sort_by"] = "popularity"
	}
	if _, ok := params["order"]; !ok {
		params["order"] = "desc"
	}

	return s.repo.GetTools(params)
}

// ToolDetail 获取工具详情
func (s *ToolService) ToolDetail(ctx context.Context, slug string) (*models.Tool, error) {
	tool, err := s.repo.GetToolBySlug(slug)
	if err != nil {
		return nil, err
	}
	// 转换标签为中文名称
	tool.Tags = s.convertTagIDsToNames(tool.Tags)
	return tool, nil
}

// ToolDetailByID 根据ID获取工具详情
func (s *ToolService) ToolDetailByID(ctx context.Context, id uint) (*models.Tool, error) {
	tool, err := s.repo.GetToolByID(id)
	if err != nil {
		return nil, err
	}
	// 转换标签为中文名称
	tool.Tags = s.convertTagIDsToNames(tool.Tags)
	return tool, nil
}

// convertTagIDsToNames 将标签ID转换为中文名称
func (s *ToolService) convertTagIDsToNames(tagIDs []string) []string {
	tagMap := map[string]string{
		"1": "大语言模型",
		"2": "对话",
		"3": "写作",
		"4": "图像生成",
		"5": "视频生成",
		"6": "编程辅助",
		"7": "免费",
		"8": "付费",
		"9": "开源",
		"10": "国内",
	}
	var names []string
	for _, id := range tagIDs {
		if name, ok := tagMap[id]; ok {
			names = append(names, name)
		} else {
			names = append(names, id)
		}
	}
	return names
}

// ToolView 工具详情响应 (含 access 决策)
type ToolView struct {
	*models.Tool
	IsLocked          bool                    `json:"is_locked"`
	RequiredLevel     int                     `json:"required_level,omitempty"`
	RequiredLevelName string                  `json:"required_level_name,omitempty"`
	Locked            *access.LockedContent   `json:"locked,omitempty"`
}

// ToToolView 把 Tool 包成 view，根据 userLevel 算 access
func ToToolView(t *models.Tool, userLevel int) *ToolView {
	v := &ToolView{Tool: t}
	decision := access.Decide(userLevel, t.AccessLevel)
	v.IsLocked = !decision.Allow
	if !decision.Allow {
		v.RequiredLevel = t.AccessLevel
		v.RequiredLevelName = access.LevelName(t.AccessLevel)
		preview, _ := access.TruncateContent(t.Description, access.GuestPreviewChars)
		t.Description = preview
		lp := access.LockedPlaceholder("工具", t.AccessLevel)
		v.Locked = &lp
	}
	return v
}

// ToolListView 给列表里每条工具打 is_locked 标签
func ToolListView(tools []models.Tool, userLevel int) []ToolView {
	out := make([]ToolView, 0, len(tools))
	for i := range tools {
		v := ToolView{Tool: &tools[i]}
		decision := access.Decide(userLevel, tools[i].AccessLevel)
		v.IsLocked = !decision.Allow
		if !decision.Allow {
			v.RequiredLevel = tools[i].AccessLevel
			v.RequiredLevelName = access.LevelName(tools[i].AccessLevel)
		}
		out = append(out, v)
	}
	return out
}
