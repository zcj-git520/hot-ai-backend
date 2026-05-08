package crawler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"hot-ai-backend/internal/models"
)

// AIClient AI 服务统一客户端
type AIClient struct {
	baseURL string
	model   string
	client  *http.Client
}

// TranslateRequest 翻译请求
type TranslateRequest struct {
	Content        string `json:"content"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
	Model          string `json:"model"`
}

// TranslateResponse 翻译响应
type TranslateResponse struct {
	Success        bool   `json:"success"`
	TranslatedText string `json:"translated_text"`
}

// AIArticleRequest AI文章分析请求
type AIArticleRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Model   string `json:"model"`
}

// AIArticleResponse AI文章分析响应
type AIArticleResponse struct {
	Category     string   `json:"category"`      // 分类结果：动态/学习/工具/职业
	Confidence   float64  `json:"confidence"`     // 置信度 (0-1)
	Reason       string   `json:"reason"`         // 分类理由
	IsAIRelated  bool     `json:"is_ai_related"` // 是否AI相关内容
	KeyPoints    []string `json:"key_points"`    // 核心要点列表
	Summary      string   `json:"summary"`       // 文章摘要
	Keywords     []string `json:"keywords"`       // 关键词列表
	Model        string   `json:"model"`
}

// NewAIClient 创建 AI 服务统一客户端
func NewAIClient(baseURL string) *AIClient {
	return &AIClient{
		baseURL: baseURL,
		model:   "glm",
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// SetModel 设置模型
func (c *AIClient) SetModel(model string) {
	c.model = model
}

// doRequest 统一请求方法
func (c *AIClient) doRequest(ctx context.Context, endpoint string, reqBody interface{}, respBody interface{}) error {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("请求 AI 服务失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("AI 服务返回错误: %d - %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(respBody); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	return nil
}

// Translate 翻译文本
func (c *AIClient) Translate(ctx context.Context, req TranslateRequest) (*TranslateResponse, error) {
	if req.SourceLanguage == "" {
		req.SourceLanguage = "auto"
	}
	if req.TargetLanguage == "" {
		req.TargetLanguage = ""
	}
	if req.Model == "" {
		req.Model = c.model
	}

	var resp TranslateResponse
	if err := c.doRequest(ctx, "/api/translate", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// TranslateArticle 翻译文章标题和内容
func (c *AIClient) TranslateArticle(ctx context.Context, article *models.Article) error {
	// 翻译标题
	if article.Title != "" {
		resp, err := c.Translate(ctx, TranslateRequest{
			Content: article.Title,
			Model:   c.model,
		})
		if err != nil {
			logx.Errorf("翻译标题失败: %v", err)
			return err
		}
		if resp.Success {
			// 直接使用翻译结果存储到 Title 字段
			article.Title = resp.TranslatedText
		}
	}

	// 翻译摘要
	if article.Summary != "" {
		resp, err := c.Translate(ctx, TranslateRequest{
			Content: article.Summary,
			Model:   c.model,
		})
		if err != nil {
			logx.Errorf("翻译摘要失败: %v", err)
		} else if resp.Success {
			// 直接使用翻译结果存储到 Summary 字段
			article.Summary = resp.TranslatedText
		}
	}

	// 翻译内容
	if article.Content != "" {
		resp, err := c.Translate(ctx, TranslateRequest{
			Content: article.Content,
			Model:   c.model,
		})
		if err != nil {
			logx.Errorf("翻译内容失败: %v", err)
			return nil
		}
		if resp.Success {
			// 直接使用翻译结果存储到 Content 字段
			article.Content = resp.TranslatedText
		}
	}

	return nil
}

// AnalyzeArticle 分析文章
func (c *AIClient) AnalyzeArticle(ctx context.Context, title, content string) (*AIArticleResponse, error) {
	if title == "" && content == "" {
		return nil, fmt.Errorf("标题和内容不能同时为空")
	}

	// 截取内容前2000字符进行分析
	analyzeContent := content
	if len([]rune(analyzeContent)) > 2000 {
		runes := []rune(analyzeContent)
		analyzeContent = string(runes[:2000])
	}

	req := AIArticleRequest{
		Title:   title,
		Content: analyzeContent,
		Model:   c.model,
	}

	var resp AIArticleResponse
	if err := c.doRequest(ctx, "/api/ai-article", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AnalyzeArticleAndSkip 如果文章非AI相关则返回true（应跳过）
func (c *AIClient) AnalyzeArticleAndSkip(ctx context.Context, title, content string) (*AIArticleResponse, bool, error) {
	result, err := c.AnalyzeArticle(ctx, title, content)
	if err != nil {
		logx.Errorf("AI文章分析失败: %v", err)
		return nil, false, err // 分析失败时不跳过
	}

	if !result.IsAIRelated {
		logx.Infof("文章非AI相关内容，跳过处理 | 标题: %s | 原因: %s", title, result.Reason)
		return result, true, nil
	}

	logx.Infof("AI分析完成 | 标题: %s | 分类: %s | 置信度: %.2f | 摘要: %s",
		title, result.Category, result.Confidence, result.Summary)

	return result, false, nil
}
