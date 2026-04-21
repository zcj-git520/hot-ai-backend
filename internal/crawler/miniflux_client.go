package crawler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// MinifluxConfig Miniflux 配置
type MinifluxConfig struct {
	BaseURL string
	APIKey  string
}

// MinifluxClient Miniflux API 客户端
type MinifluxClient struct {
	config  MinifluxConfig
	headers map[string]string
}

// NewMinifluxClient 创建 Miniflux 客户端
func NewMinifluxClient(baseURL, apiKey string) *MinifluxClient {
	return &MinifluxClient{
		config: MinifluxConfig{
			BaseURL: baseURL,
			APIKey:  apiKey,
		},
		headers: map[string]string{
			"X-Auth-Token": apiKey,
			"Content-Type": "application/json",
		},
	}
}

// Entry Miniflux 条目
type Entry struct {
	ID          uint64    `json:"id"`
	UserID      int       `json:"user_id"`
	FeedID      uint64    `json:"feed_id"`
	Status      string    `json:"status"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	CommentsURL string    `json:"comments_url"`
	Date        time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	ChangedAt   time.Time `json:"changed_at"`
	Content     string    `json:"content"`
	Author      string    `json:"author"`
	Tags        []string  `json:"tags"`
	Hash        string    `json:"hash"`
	Starred     bool      `json:"starred"`
	ReadingTime int       `json:"reading_time"`
	Enclosures  []struct {
		URL      string `json:"url"`
		MimeType string `json:"mime_type"`
		Size     int    `json:"size"`
	} `json:"enclosures"`
	Feed *Feed `json:"feed,omitempty"`
}

// Feed Miniflux 订阅源
type Feed struct {
	ID                          uint64    `json:"id"`
	UserID                      int       `json:"user_id"`
	FeedURL                     string    `json:"feed_url"`
	SiteURL                     string    `json:"site_url"`
	Title                       string    `json:"title"`
	CheckedAt                   time.Time `json:"checked_at"`
	NextCheckAt                 time.Time `json:"next_check_at"`
	ErrorCount                  int       `json:"error_count"`
	ParsingError                string    `json:"parsing_error"`
	LastHTTPStatus              string    `json:"last_http_status"`
	Crawler                     bool      `json:"crawler"`
	BlocklistRules              string    `json:"blocklist_rules"`
	KeypasslistRules            string    `json:"keepasslist_rules"`
	UserAgent                   string    `json:"user_agent"`
	Username                    string    `json:"username"`
	Password                    string    `json:"password"`
	ClientTimeout               int       `json:"client_timeout"`
	MaxEntries                  int       `json:"max_entries"`
	AllowSelfSignedCertificates bool      `json:"allow_self_signed_certificates"`
	FetchViaProxy               bool      `json:"fetch_via_proxy"`
	ScraperRules                string    `json:"scraper_rules"`
	RewriteRules                string    `json:"rewrite_rules"`
	Category                    *Category `json:"category,omitempty"`
}

// Category Miniflux 分类
type Category struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	UserID int    `json:"user_id"`
}

// EntriesResponse 条目列表响应
type EntriesResponse struct {
	Total   int     `json:"total"`
	Entries []Entry `json:"entries"`
}

// GetLatestEntries 获取最新条目
func (c *MinifluxClient) GetLatestEntries(limit int, feedID ...uint64) ([]Entry, error) {
	url := fmt.Sprintf("%s/v1/entries", c.config.BaseURL)

	params := map[string]interface{}{
		"limit":     limit,
		"order":     "published_at",
		"direction": "desc",
		"status":    "read,unread",
	}

	if len(feedID) > 0 {
		params["feed_id"] = feedID[0]
	}

	resp, err := c.doRequest("GET", url, params)
	if err != nil {
		return nil, fmt.Errorf("获取条目失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var entriesResp EntriesResponse
	if err := json.Unmarshal(body, &entriesResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	logx.Infof("从 Miniflux 获取到 %d 条条目", len(entriesResp.Entries))
	return entriesResp.Entries, nil
}

// FetchEntryContent 获取条目完整内容（使用 fetch-content 接口）
func (c *MinifluxClient) FetchEntryContent(entryID uint64) (string, error) {
	url := fmt.Sprintf("%s/v1/entries/%d/fetch-content", c.config.BaseURL, entryID)

	// GET 请求，update_content=true
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	q := req.URL.Query()
	q.Add("update_content", "true")
	req.URL.RawQuery = q.Encode()

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("HTTP 错误: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析 Miniflux API 返回的条目内容
	var entry map[string]interface{}
	if err := json.Unmarshal(body, &entry); err != nil {
		return "", fmt.Errorf("解析响应失败: %w", err)
	}

	// 安全获取 content 字段
	content, ok := entry["content"].(string)
	if !ok {
		return "", fmt.Errorf("无法获取内容字段")
	}

	logx.Infof("成功获取条目完整内容，长度: %d", len(content))

	return content, nil
}

// GetFeeds 获取所有订阅源
func (c *MinifluxClient) GetFeeds() ([]Feed, error) {
	url := fmt.Sprintf("%s/v1/feeds", c.config.BaseURL)

	resp, err := c.doRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("获取订阅源失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var feeds []Feed
	if err := json.Unmarshal(body, &feeds); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	logx.Infof("从 Miniflux 获取到 %d 个订阅源", len(feeds))
	return feeds, nil
}

// GetFeedEntries 获取指定 Feed 的未读条目
func (c *MinifluxClient) GetFeedEntries(feedID uint64, limit int) ([]Entry, error) {
	// 使用 /v1/entries 端点，通过 feed_id 参数过滤
	url := fmt.Sprintf("%s/v1/entries", c.config.BaseURL)

	// 只获取未读的条目
	params := map[string]interface{}{
		"limit":     limit,
		"order":     "published_at",
		"direction": "desc",
		"feed_id":   feedID,
		//"status":    "unread", // 只获取未读条目
	}

	logx.Infof("请求 Miniflux API: %s, 参数: %+v", url, params)

	resp, err := c.doRequest("GET", url, params)
	if err != nil {
		return nil, fmt.Errorf("获取 Feed 条目失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	logx.Debugf("API 响应: %s", string(body))

	var entriesResp EntriesResponse
	if err := json.Unmarshal(body, &entriesResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	logx.Infof("从 Feed %d 获取到 %d 条未读条目", feedID, len(entriesResp.Entries))
	return entriesResp.Entries, nil
}

// MarkFeedAsRead 标记指定 Feed 的所有条目为已读
func (c *MinifluxClient) MarkFeedAsRead(feedID uint64) error {
	url := fmt.Sprintf("%s/v1/feeds/%d/mark-all-as-read", c.config.BaseURL, feedID)

	logx.Infof("标记 Feed %d 为已读", feedID)

	// PUT 请求，不需要 body
	req, err := http.NewRequest("PUT", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP 错误: %d, 响应: %s", resp.StatusCode, string(body))
	}

	logx.Infof("成功标记 Feed %d 为已读", feedID)
	return nil
}

// doRequest 执行 HTTP 请求
func (c *MinifluxClient) doRequest(method string, url string, params map[string]interface{}) (*http.Response, error) {
	var req *http.Request
	var err error

	if method == "GET" && params != nil {
		// GET 请求，参数放在 URL query string
		q := make(map[string]string)
		for k, v := range params {
			q[k] = fmt.Sprintf("%v", v)
		}

		// 构建带参数的URL
		if len(q) > 0 {
			url += "?"
			i := 0
			for k, v := range q {
				if i > 0 {
					url += "&"
				}
				url += fmt.Sprintf("%s=%s", k, v)
				i++
			}
		}
		req, err = http.NewRequest(method, url, nil)
	} else if (method == "POST" || method == "PUT") && params != nil {
		// POST/PUT 请求，参数放在 JSON body
		bodyBytes, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("序列化参数失败: %w", err)
		}
		req, err = http.NewRequest(method, url, nil)
		if err == nil {
			req.Body = io.NopCloser(nil)
			req.ContentLength = 0
		}
		_ = bodyBytes // 暂时不使用，Miniflux API 的 POST/PUT 可能需要调整
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("HTTP 错误: %d, 响应: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}
