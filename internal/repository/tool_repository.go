package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"hot-ai-backend/internal/models"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
)

// ToolRepository 工具仓储
type ToolRepository struct {
	db *sql.DB
}

// NewToolRepository 创建工具仓储
func NewToolRepository(db *sql.DB) *ToolRepository {
	return &ToolRepository{db: db}
}

// GetCategories 获取所有工具类别
func (r *ToolRepository) GetCategories() ([]models.ToolCategory, error) {
	query := "SELECT id, name, slug, icon, description, sort_order, featured, status, created_at, updated_at, deleted_at FROM tool_categories WHERE deleted_at IS NULL AND status = 1 ORDER BY sort_order"

	rows, err := r.db.Query(query)
	if err != nil {
		logx.Errorf("Query tool_categories error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var categories []models.ToolCategory
	for rows.Next() {
		var cat models.ToolCategory
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Slug, &cat.Icon, &cat.Description, &cat.SortOrder, &cat.Featured, &cat.Status, &cat.CreatedAt, &cat.UpdatedAt, &cat.DeletedAt)
		if err != nil {
			logx.Errorf("Scan tool_category error: %v", err)
			continue
		}
		categories = append(categories, cat)
	}

	return categories, nil
}

// GetTools 获取工具列表
func (r *ToolRepository) GetTools(params map[string]interface{}) ([]models.Tool, int64, error) {
	query := "SELECT id, name, slug, icon, description, official_url, documentation_url, pricing, pricing_description, category_id, difficulty, rating, review_count, view_count, popularity, tags, featured, status, external_id, created_by, updated_by, created_at, updated_at, deleted_at, is_online, access_level FROM tools WHERE deleted_at IS NULL AND status = 1 AND is_online = 1"
	countQuery := "SELECT COUNT(*) FROM tools WHERE deleted_at IS NULL AND status = 1 AND is_online = 1"

	whereClauses := []string{}
	args := []interface{}{}

	if categoryID, ok := params["category_id"].(int); ok && categoryID > 0 {
		whereClauses = append(whereClauses, "category_id = ?")
		args = append(args, categoryID)
	}

	if isFree, ok := params["is_free"].(bool); ok {
		whereClauses = append(whereClauses, "is_free = ?")
		args = append(args, isFree)
	}

	if difficulty, ok := params["difficulty"].(string); ok && difficulty != "" {
		whereClauses = append(whereClauses, "difficulty = ?")
		args = append(args, difficulty)
	}

	if minRating, ok := params["min_rating"].(float64); ok && minRating > 0 {
		whereClauses = append(whereClauses, "rating >= ?")
		args = append(args, minRating)
	}

	if search, ok := params["search"].(string); ok && search != "" {
		whereClauses = append(whereClauses, "(name LIKE ? OR description LIKE ?)")
		args = append(args, "%"+search+"%", "%"+search+"%")
	}

	if len(whereClauses) > 0 {
		query += " AND " + strings.Join(whereClauses, " AND ")
		countQuery += " AND " + strings.Join(whereClauses, " AND ")
	}

	var total int64
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		logx.Errorf("Query tools count error: %v", err)
		return nil, 0, err
	}

	sortBy := params["sort_by"].(string)
	order := params["order"].(string)

	if order != "asc" && order != "desc" {
		order = "desc"
	}

	if sortBy == "popularity" {
		query += " ORDER BY popularity " + order
	} else if sortBy == "rating" {
		query += " ORDER BY rating " + order
	} else if sortBy == "created_at" {
		query += " ORDER BY created_at " + order
	} else {
		query += " ORDER BY popularity " + order
	}

	page := params["page"].(int)
	pageSize := params["page_size"].(int)
	offset := (page - 1) * pageSize
	query += " LIMIT ? OFFSET ?"
	args = append(args, pageSize, offset)

	rows, err := r.db.Query(query, args...)
	if err != nil {
		logx.Errorf("Query tools error: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	var tools []models.Tool
	for rows.Next() {
		var tool models.Tool
		var tagsStr sql.NullString
		err := rows.Scan(&tool.ID, &tool.Name, &tool.Slug, &tool.Icon, &tool.Description, &tool.OfficialURL, &tool.DocumentationURL, &tool.Pricing, &tool.PricingDesc, &tool.CategoryID, &tool.Difficulty, &tool.Rating, &tool.ReviewCount, &tool.ViewCount, &tool.Popularity, &tagsStr, &tool.Featured, &tool.Status, &tool.ExternalID, &tool.CreatedBy, &tool.UpdatedBy, &tool.CreatedAt, &tool.UpdatedAt, &tool.DeletedAt, &tool.IsOnline, &tool.AccessLevel)
		if err != nil {
			logx.Errorf("Scan tool error: %v", err)
			continue
		}
		// 处理标签转换
		if tagsStr.Valid {
			var tags []string
			if err := json.Unmarshal([]byte(tagsStr.String), &tags); err == nil {
				tool.Tags = tags
			}
		}
		tools = append(tools, tool)
	}

	return tools, total, nil
}

// GetToolBySlug 获取工具详情
func (r *ToolRepository) GetToolBySlug(slug string) (*models.Tool, error) {
	query := "SELECT id, name, slug, icon, description, official_url, documentation_url, pricing, pricing_description, category_id, difficulty, rating, review_count, view_count, popularity, tags, featured, status, external_id, created_by, updated_by, created_at, updated_at, deleted_at, access_level FROM tools WHERE slug = ? AND deleted_at IS NULL AND status = 1"

	var tool models.Tool
	var tagsStr sql.NullString
	err := r.db.QueryRow(query, slug).Scan(&tool.ID, &tool.Name, &tool.Slug, &tool.Icon, &tool.Description, &tool.OfficialURL, &tool.DocumentationURL, &tool.Pricing, &tool.PricingDesc, &tool.CategoryID, &tool.Difficulty, &tool.Rating, &tool.ReviewCount, &tool.ViewCount, &tool.Popularity, &tagsStr, &tool.Featured, &tool.Status, &tool.ExternalID, &tool.CreatedBy, &tool.UpdatedBy, &tool.CreatedAt, &tool.UpdatedAt, &tool.DeletedAt, &tool.AccessLevel)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("tool not found")
		}
		logx.Errorf("Query tool error: %v", err)
		return nil, err
	}

	// 处理标签转换
	if tagsStr.Valid {
		var tags []string
		if err := json.Unmarshal([]byte(tagsStr.String), &tags); err == nil {
			tool.Tags = tags
		}
	}

	r.db.Exec("UPDATE tools SET view_count = view_count + 1 WHERE id = ?", tool.ID)

	return &tool, nil
}

// GetToolByID 获取工具详情（含ID）
func (r *ToolRepository) GetToolByID(id uint) (*models.Tool, error) {
	query := "SELECT id, name, slug, icon, description, official_url, documentation_url, pricing, pricing_description, category_id, difficulty, rating, review_count, view_count, popularity, tags, featured, status, external_id, created_by, updated_by, created_at, updated_at, deleted_at, access_level FROM tools WHERE id = ? AND deleted_at IS NULL AND status = 1"

	var tool models.Tool
	var tagsStr sql.NullString
	err := r.db.QueryRow(query, id).Scan(&tool.ID, &tool.Name, &tool.Slug, &tool.Icon, &tool.Description, &tool.OfficialURL, &tool.DocumentationURL, &tool.Pricing, &tool.PricingDesc, &tool.CategoryID, &tool.Difficulty, &tool.Rating, &tool.ReviewCount, &tool.ViewCount, &tool.Popularity, &tagsStr, &tool.Featured, &tool.Status, &tool.ExternalID, &tool.CreatedBy, &tool.UpdatedBy, &tool.CreatedAt, &tool.UpdatedAt, &tool.DeletedAt, &tool.AccessLevel)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("tool not found")
		}
		logx.Errorf("Query tool error: %v", err)
		return nil, err
	}

	// 处理标签转换
	if tagsStr.Valid {
		var tags []string
		if err := json.Unmarshal([]byte(tagsStr.String), &tags); err == nil {
			tool.Tags = tags
		}
	}

	return &tool, nil
}
