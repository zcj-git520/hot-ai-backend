package crawler

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/zeromicro/go-zero/core/logx"
)

// ArticleLink 文章链接信息
type ArticleLink struct {
	URL         string
	Title       string
	PublishedAt *time.Time
}

// ParseArticleLinks 从列表页解析出文章链接
func ParseArticleLinks(htmlContent string, baseURL string, rules map[string]interface{}) ([]ArticleLink, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("解析 HTML 失败: %w", err)
	}

	var links []ArticleLink

	// 获取选择器配置
	selector, ok := rules["article_selector"].(string)
	if !ok || selector == "" {
		// 默认尝试常见的文章链接选择器
		selector = "a[href]"
	}

	// 查找所有文章链接
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if !exists || href == "" {
			return
		}

		// 转换为绝对 URL
		absoluteURL := resolveURL(href, baseURL)
		if absoluteURL == "" {
			return
		}

		// 提取标题
		title := strings.TrimSpace(s.Text())
		if title == "" {
			// 尝试从子元素获取标题
			title = strings.TrimSpace(s.Find("h1, h2, h3, .title, .headline").First().Text())
		}

		// 过滤：只保留看起来像文章的链接
		if isArticleLink(absoluteURL, title) {
			link := ArticleLink{
				URL:   absoluteURL,
				Title: title,
			}

			// 尝试提取发布时间
			if timeSelector, ok := rules["time_selector"].(string); ok && timeSelector != "" {
				timeStr := doc.Find(timeSelector).First().Text()
				if parsedTime, err := parseDateTime(timeStr); err == nil {
					link.PublishedAt = &parsedTime
				}
			}

			links = append(links, link)
		}
	})

	logx.Infof("从列表页解析到 %d 个文章链接", len(links))
	return links, nil
}

// ParseArticleDetail 解析文章详情页，提取标题、内容等信息
func ParseArticleDetail(htmlContent string, articleURL string, rules map[string]interface{}) (map[string]interface{}, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		return nil, fmt.Errorf("解析 HTML 失败: %w", err)
	}

	article := make(map[string]interface{})

	// 1. 提取标题
	titleSelector := ".title"
	if rules != nil {
		if sel, ok := rules["title_selector"].(string); ok && sel != "" {
			titleSelector = sel
		}
	}

	title := strings.TrimSpace(doc.Find(titleSelector).First().Text())
	if title == "" {
		//  fallback: 尝试 <title> 标签或 h1
		title = strings.TrimSpace(doc.Find("title").First().Text())
		if title == "" {
			title = strings.TrimSpace(doc.Find("h1").First().Text())
		}
	}
	article["title"] = title

	// 2. 提取作者
	authorSelector := ".author, .byline, [rel='author']"
	if rules != nil {
		if sel, ok := rules["author_selector"].(string); ok && sel != "" {
			authorSelector = sel
		}
	}
	author := strings.TrimSpace(doc.Find(authorSelector).First().Text())
	if author == "" {
		author = "未知作者"
	}
	article["author"] = author

	// 3. 提取发布时间
	publishedAtSelector := ".date, .published, time, .publish-date"
	if rules != nil {
		if sel, ok := rules["date_selector"].(string); ok && sel != "" {
			publishedAtSelector = sel
		}
	}
	dateStr := strings.TrimSpace(doc.Find(publishedAtSelector).First().Text())
	if parsedTime, err := parseDateTime(dateStr); err == nil {
		article["published_at"] = parsedTime.Format(time.RFC3339)
	} else {
		article["published_at"] = time.Now().Format(time.RFC3339)
	}

	// 4. 提取正文内容
	contentSelector := ".content, .article-content, .post-content, article, .entry-content"
	if rules != nil {
		if sel, ok := rules["content_selector"].(string); ok && sel != "" {
			contentSelector = sel
		}
	}

	contentHTML := extractCleanContent(doc, contentSelector)
	article["content"] = contentHTML

	// 5. 提取摘要
	summary := generateSummary(contentHTML)
	article["summary"] = summary

	// 6. 提取标签/分类
	var tags []string
	tagSelector := ".tags a, .categories a, .tag a"
	if rules != nil {
		if sel, ok := rules["tag_selector"].(string); ok && sel != "" {
			tagSelector = sel
		}
	}
	doc.Find(tagSelector).Each(func(i int, s *goquery.Selection) {
		tag := strings.TrimSpace(s.Text())
		if tag != "" {
			tags = append(tags, tag)
		}
	})
	if len(tags) > 0 {
		article["tags"] = strings.Join(tags, ",")
	}

	// 7. 提取封面图片
	imageSelector := ".featured-image img, .cover-image img, article img:first-of-type"
	if rules != nil {
		if sel, ok := rules["image_selector"].(string); ok && sel != "" {
			imageSelector = sel
		}
	}
	if imgSrc, exists := doc.Find(imageSelector).First().Attr("src"); exists {
		article["cover_image"] = resolveURL(imgSrc, articleURL)
	}

	logx.Infof("成功解析文章详情: %s", title)
	return article, nil
}

// extractCleanContent 提取干净的正文内容
func extractCleanContent(doc *goquery.Document, selector string) string {
	var contentBuilder strings.Builder

	doc.Find(selector).First().Find("p, h1, h2, h3, h4, h5, h6, ul, ol, blockquote, pre, code").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			// 根据元素类型添加格式
			tagName := strings.ToLower(s.Get(0).Data)
			switch tagName {
			case "h1", "h2", "h3":
				contentBuilder.WriteString("\n\n# " + text + "\n")
			case "h4", "h5", "h6":
				contentBuilder.WriteString("\n\n## " + text + "\n")
			case "p":
				contentBuilder.WriteString("\n\n" + text)
			case "blockquote":
				contentBuilder.WriteString("\n\n> " + text + "\n")
			case "ul", "ol":
				s.Find("li").Each(func(j int, li *goquery.Selection) {
					liText := strings.TrimSpace(li.Text())
					if liText != "" {
						contentBuilder.WriteString("\n- " + liText)
					}
				})
			case "pre", "code":
				contentBuilder.WriteString("\n\n```\n" + text + "\n```\n")
			}
		}
	})

	content := contentBuilder.String()
	content = strings.TrimSpace(content)

	// 如果没找到结构化内容，返回纯文本
	if content == "" {
		content = strings.TrimSpace(doc.Find(selector).First().Text())
		// 清理多余空白
		for strings.Contains(content, "\n\n\n") {
			content = strings.ReplaceAll(content, "\n\n\n", "\n\n")
		}
	}

	return content
}

// generateSummary 生成摘要
func generateSummary(content string) string {
	// 取前 500 个字符作为摘要
	maxLen := 500
	if len(content) <= maxLen {
		return content
	}

	summary := content[:maxLen]
	// 在句子边界截断
	lastPeriod := strings.LastIndex(summary, "。")
	lastSentence := strings.LastIndex(summary, ". ")

	cutPoint := lastPeriod
	if lastSentence > cutPoint {
		cutPoint = lastSentence
	}

	if cutPoint > maxLen/2 {
		summary = summary[:cutPoint+1]
	}

	return summary + "..."
}

// resolveURL 解析相对 URL 为绝对 URL
func resolveURL(href string, baseURL string) string {
	// 跳过锚点链接和 JavaScript 链接
	if strings.HasPrefix(href, "#") ||
		strings.HasPrefix(href, "javascript:") ||
		strings.HasPrefix(href, "mailto:") ||
		strings.HasPrefix(href, "tel:") {
		return ""
	}

	// 如果已经是绝对 URL，直接返回
	if strings.HasPrefix(href, "http://") || strings.HasPrefix(href, "https://") {
		return href
	}

	// 解析基础 URL
	base, err := url.Parse(baseURL)
	if err != nil {
		return ""
	}

	// 处理以 / 开头的相对路径
	if strings.HasPrefix(href, "/") {
		return fmt.Sprintf("%s://%s%s", base.Scheme, base.Host, href)
	}

	// 处理相对路径
	resolved, err := base.Parse(href)
	if err != nil {
		return ""
	}

	return resolved.String()
}

// isArticleLink 判断是否是文章链接
func isArticleLink(url string, title string) bool {
	// 检查 URL 是否包含文章相关关键词
	articleKeywords := []string{
		"/article/", "/blog/", "/post/", "/news/",
		"/story/", "/entry/", "/item/", "/detail/",
		"/p/", "/a/", "/n/",
	}

	urlLower := strings.ToLower(url)
	for _, keyword := range articleKeywords {
		if strings.Contains(urlLower, keyword) {
			return true
		}
	}

	// 检查标题是否有意义（至少 5 个字符）
	if len(strings.TrimSpace(title)) < 5 {
		return false
	}

	// 排除一些明显不是文章的链接
	excludeKeywords := []string{
		"/login", "/register", "/signup", "/logout",
		"/search", "/tag/", "/category/", "/archive/",
		".pdf", ".doc", ".zip", ".jpg", ".png", ".gif",
	}

	for _, keyword := range excludeKeywords {
		if strings.Contains(urlLower, keyword) {
			return false
		}
	}

	return true
}

// parseDateTime 解析日期时间字符串
func parseDateTime(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return time.Time{}, fmt.Errorf("日期字符串为空")
	}

	// 尝试多种日期格式
	formats := []string{
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"Jan 2, 2006",
		"January 2, 2006",
		"2006年01月02日",
		"2006年1月2日",
		"01/02/2006",
		"02 Jan 2006",
		time.RFC1123,
		time.RFC1123Z,
		time.RFC822,
		time.RFC822Z,
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("无法解析日期: %s", dateStr)
}
