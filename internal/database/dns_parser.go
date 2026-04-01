package database

import (
	"fmt"
	"net/url"
	"strings"
)

// ParseDSN 从 DSN 字符串解析数据库配置
// DSN 格式：user:password@tcp(host:port)/dbname?params
func ParseDSN(dsn string) (*Config, error) {
	if dsn == "" {
		return nil, fmt.Errorf("DSN is empty")
	}

	// 提取 ? 前面的部分
	parts := strings.SplitN(dsn, "?", 2)
	mainPart := parts[0]

	// 解析 user:password@tcp(host:port)/dbname
	// 1. 分割 @ 符号获取用户信息和地址信息
	atIndex := strings.LastIndex(mainPart, "@")
	if atIndex == -1 {
		return nil, fmt.Errorf("invalid DSN format: missing @")
	}

	userInfo := mainPart[:atIndex]
	addressInfo := mainPart[atIndex+1:]

	// 2. 解析用户信息 (user:password)
	colonIndex := strings.Index(userInfo, ":")
	if colonIndex == -1 {
		return nil, fmt.Errorf("invalid DSN format: missing password")
	}

	user := userInfo[:colonIndex]
	password := userInfo[colonIndex+1:]

	// 3. 解析地址信息 (tcp(host:port)/dbname)
	tcpPrefix := "tcp("
	if !strings.HasPrefix(addressInfo, tcpPrefix) {
		return nil, fmt.Errorf("invalid DSN format: missing tcp()")
	}

	// 移除 tcp( 和最后的 )
	addressWithDB := strings.TrimSuffix(addressInfo[len(tcpPrefix):], ")")

	// 分割 )/ 获取地址和数据库名
	slashIndex := strings.Index(addressWithDB, "/")
	if slashIndex == -1 {
		return nil, fmt.Errorf("invalid DSN format: missing database name")
	}

	hostPort := addressWithDB[:slashIndex]
	dbName := addressWithDB[slashIndex+1:]

	// 4. 解析 host:port
	hostParts := strings.Split(hostPort, ":")
	if len(hostParts) != 2 {
		return nil, fmt.Errorf("invalid DSN format: invalid host:port")
	}

	host := hostParts[0]
	port := 0
	fmt.Sscanf(hostParts[1], "%d", &port)

	// 5. 解析查询参数（可选）
	queryParams := make(map[string]string)
	if len(parts) > 1 {
		values, err := url.ParseQuery(parts[1])
		if err == nil {
			for key, value := range values {
				if len(value) > 0 {
					queryParams[key] = value[0]
				}
			}
		}
	}

	return &Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		DBName:   dbName,
	}, nil
}
