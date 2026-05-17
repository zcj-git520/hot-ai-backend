package handler

import (
	"net/http"
	"strings"
)

// getPathValue 获取路径参数（兼容 Go 1.21）
func getPathValue(r *http.Request, name string) string {
	// 优先使用 Go 1.22+ PathValue
	val := r.PathValue(name)
	if val != "" {
		return val
	}
	// 回退：从 URL path 中提取参数值
	// Go Zero路由匹配后，:path_id 被替换为实际值
	// URL格式: /admin/learning-paths/1/chapters -> path_id=1
	// URL格式: /admin/chapters/123 -> id=123
	path := r.URL.Path
	parts := strings.Split(path, "/")

	// 定义参数名与前置路由段的映射
	paramMapping := map[string][]string{
		"path_id": {"learning-paths"},
		"id":     {"chapters", "learning-paths"},
	}

	prefixes, ok := paramMapping[name]
	if !ok {
		// 未知参数名，使用通用回退：返回最后一个非静态段
		if len(parts) > 0 {
			last := parts[len(parts)-1]
			staticSegments := map[string]bool{
				"chapters": true, "learning-paths": true, "admin": true,
			}
			if !staticSegments[last] {
				return last
			}
		}
		return ""
	}

	// 查找前置路由段
	for i, part := range parts {
		for _, prefix := range prefixes {
			if part == prefix && i+1 < len(parts) {
				return parts[i+1]
			}
		}
	}

	return ""
}