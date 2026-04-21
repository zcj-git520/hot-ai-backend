package crawler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hot-ai-backend/internal/models"
	"io"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// TranslateClient 翻译客户端
type TranslateClient struct {
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
	SourceText     string `json:"source_text"`
	SourceLanguage string `json:"source_language"`
	TargetLanguage string `json:"target_language"`
	Model          string `json:"model"`
}

// NewTranslateClient 创建翻译客户端
func NewTranslateClient(baseURL string) *TranslateClient {
	return &TranslateClient{
		baseURL: baseURL,
		model:   "glm",
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

// Translate 翻译文本
func (c *TranslateClient) Translate(ctx context.Context, req TranslateRequest) (*TranslateResponse, error) {
	if req.SourceLanguage == "" {
		req.SourceLanguage = "auto"
	}
	if req.TargetLanguage == "" {
		req.TargetLanguage = ""
	}
	if req.Model == "" {
		req.Model = c.model
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/translate", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("请求翻译服务失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("翻译服务返回错误: %d - %s", resp.StatusCode, string(body))
	}

	var result TranslateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// TranslateArticle 翻译文章标题和内容
func (c *TranslateClient) TranslateArticle(ctx context.Context, article *models.Article) error {

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
			// 根据源语言决定标题存储
			if resp.SourceLanguage == "cn" {
				article.TitleEn = resp.TranslatedText
			} else {
				article.Content = resp.TranslatedText // 原文是英文，content存译文
				article.ContentEn = article.Title     // content_en存原文
			}
		}
	}

	// 翻译内容（取前2000字符作为摘要）
	//summaryContent := content
	//if len(summaryContent) > 2000 {
	//	summaryContent = summaryContent[:2000]
	//}
	//
	//if summaryContent != "" {
	//	resp, err := c.Translate(ctx, TranslateRequest{
	//		Content: summaryContent,
	//		Model:   c.model,
	//	})
	//	if err != nil {
	//		logx.Errorf("翻译摘要失败: %v", err)
	//	}
	//	// 摘要翻译结果暂不单独存储，主要用于判断源语言
	//	if resp != nil && resp.Success {
	//		result.SourceLang = resp.SourceLanguage
	//	}
	//}

	// 翻译完整内容
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
			// 根据源语言决定内容存储
			if resp.SourceLanguage == "cn" {
				// 原文是中文
				article.ContentEn = resp.TranslatedText
			} else {
				article.Content = resp.TranslatedText // 原文是英文，content存译文
				article.ContentEn = article.Content   // content_en存原文
			}
		}
	}

	return nil
}
