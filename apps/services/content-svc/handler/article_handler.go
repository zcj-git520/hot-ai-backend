package handler

import (
	"fmt"
	"net/http"

	"hot-ai-backend/internal/service"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// ArticleHandler 文章处理器
type ArticleHandler struct {
	articleService *service.ArticleService
}

// NewArticleHandler 创建文章处理器实例
func NewArticleHandler(articleService *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		articleService: articleService,
	}
}

// GetArticles 获取文章列表
func (h *ArticleHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageSize := 10
	category := r.URL.Query().Get("category")
	keyword := r.URL.Query().Get("keyword")

	if p := r.URL.Query().Get("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := r.URL.Query().Get("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	resp, err := h.articleService.GetArticles(&service.GetArticlesRequest{
		Page:     page,
		PageSize: pageSize,
		Category: category,
		Keyword:  keyword,
	})
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, resp)
}

// GetArticleByID 根据ID获取文章详情
func (h *ArticleHandler) GetArticleByID(w http.ResponseWriter, r *http.Request) {
	// 从URL路径参数获取ID - go-zero将路径参数存储在r.URL.Path中
	// 对于路由 /api/articles/:id,需要手动解析
	idStr := r.PathValue("id")
	
	// 如果PathValue不可用,尝试从context中获取
	if idStr == "" {
		// 尝试从URL path中手动提取
		// URL格式: /api/articles/123
		path := r.URL.Path
		// 找到最后一个/后的内容
		for i := len(path) - 1; i >= 0; i-- {
			if path[i] == '/' {
				idStr = path[i+1:]
				break
			}
		}
	}
	
	if idStr == "" {
		httpx.ErrorCtx(r.Context(), w, fmt.Errorf("缺少文章ID"))
		return
	}

	article, err := h.articleService.GetArticleByID(idStr)
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, article)
}

// GetCategories 获取文章分类列表
func (h *ArticleHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.articleService.GetCategories()
	if err != nil {
		httpx.ErrorCtx(r.Context(), w, err)
		return
	}

	httpx.OkJsonCtx(r.Context(), w, categories)
}
