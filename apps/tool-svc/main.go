package main

import (
	"flag"
	"fmt"
	"hot-ai-backend/apps/tool-svc/handler"
	"net/http"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"hot-ai-backend/internal/database"
	"hot-ai-backend/internal/middleware"
	"hot-ai-backend/internal/repository"
	"hot-ai-backend/internal/service"
)

var configFile = flag.String("f", "apps/tool-svc/etc/tool-svc.yaml", "the config file")

// ToolSvcConf 工具服务配置
type ToolSvcConf struct {
	rest.RestConf
	// 数据库配置
	DataSource struct {
		MySQL struct {
			DSN string `json:",optional"`
		}
	}
	// JWT 配置
	Auth struct {
		AccessSecret string `json:",default=your-secret-key-change-in-production"`
		AccessExpire int    `json:",default=86400"`
	}
}

func main() {
	flag.Parse()

	var c ToolSvcConf
	if err := conf.Load(*configFile, &c); err != nil {
		fmt.Fprintf(os.Stderr, "load config error: %v\n", err)
		os.Exit(1)
	}

	// 初始化数据库
	if c.DataSource.MySQL.DSN != "" {
		dbConfig, err := database.ParseDSN(c.DataSource.MySQL.DSN)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse DSN error: %v\n", err)
			os.Exit(1)
		}

		if err := database.InitDB(*dbConfig); err != nil {
			fmt.Fprintf(os.Stderr, "init database error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Database initialized successfully")
	} else {
		fmt.Println("[警告] 未配置数据库")
		os.Exit(1)
	}

	// 创建服务器
	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	// 初始化仓储层
	sqlDB, err := database.GetDB().DB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "get sql.DB error: %v\n", err)
		os.Exit(1)
	}
	toolRepo := repository.NewToolRepository(sqlDB)

	// 初始化服务层
	toolService := service.NewToolService(toolRepo)

	// 注册路由
	registerRoutes(server, c.Auth.AccessSecret, toolService)

	fmt.Printf("Starting tool-svc at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// registerRoutes 注册路由
func registerRoutes(server *rest.Server, jwtSecret string, toolService *service.ToolService) {
	// 创建处理器
	toolHandler := handler.NewToolServiceHandler(toolService)

	// OptionalAuth 中间件
	optAuth := middleware.OptionalAuth(jwtSecret)
	wrapFunc := func(h func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			h(w, r)
		}
	}

	// ===== 工具路由 =====
	// 工具类别
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/tools/categories",
		Handler: optAuth(wrapFunc(toolHandler.ToolCategories)).ServeHTTP,
	})

	// 工具列表
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/tools",
		Handler: optAuth(wrapFunc(toolHandler.ToolList)).ServeHTTP,
	})

	// 工具详情（通过slug）
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/tools/:slug",
		Handler: optAuth(wrapFunc(toolHandler.ToolDetail)).ServeHTTP,
	})

	// 工具详情（通过id）
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/tools/id/:id",
		Handler: optAuth(wrapFunc(toolHandler.ToolDetailByID)).ServeHTTP,
	})
}
