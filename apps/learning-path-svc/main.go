package main

import (
	"flag"
	"fmt"
	"hot-ai-backend/apps/learning-path-svc/handler"
	"net/http"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"

	"hot-ai-backend/internal/database"
	"hot-ai-backend/internal/middleware"
	"hot-ai-backend/internal/repository"
	"hot-ai-backend/internal/service"
)

var configFile = flag.String("f", "apps/learning-path-svc/etc/learning-path-svc.yaml", "the config file")

// LearningPathSvcConf 学习路径服务配置
type LearningPathSvcConf struct {
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

	var c LearningPathSvcConf
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
	learningPathRepo := repository.NewLearningPathRepository()

	// 初始化服务层
	learningPathService := service.NewLearningPathService(learningPathRepo)

	// 注册路由
	registerRoutes(server, c.Auth.AccessSecret, learningPathService)

	fmt.Printf("Starting learning-path-svc at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// registerRoutes 注册路由
func registerRoutes(server *rest.Server, jwtSecret string, learningPathService *service.LearningPathService) {
	// 创建处理器
	learningPathHandler := handler.NewLearningPathHandler(learningPathService)

	// OptionalAuth 中间件
	optAuth := middleware.OptionalAuth(jwtSecret)
	wrapFunc := func(h func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			h(w, r)
		}
	}

	// ===== 学习路径路由 =====
	// 获取路径列表
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths",
		Handler: optAuth(wrapFunc(learningPathHandler.GetLearningPaths)).ServeHTTP,
	})

	// 根据ID获取路径详情
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/:id",
		Handler: optAuth(wrapFunc(learningPathHandler.GetLearningPathByID)).ServeHTTP,
	})

	// 根据slug获取路径详情
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/slug/:slug",
		Handler: optAuth(wrapFunc(learningPathHandler.GetLearningPathBySlug)).ServeHTTP,
	})

	// 获取推荐路径
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/featured",
		Handler: optAuth(wrapFunc(learningPathHandler.GetFeaturedPaths)).ServeHTTP,
	})

	// 获取难度等级信息
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/levels",
		Handler: optAuth(wrapFunc(learningPathHandler.GetLevelInfo)).ServeHTTP,
	})

	// ===== 章线路由 =====
	// 获取路径的所有章节
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/:path_id/chapters",
		Handler: optAuth(wrapFunc(learningPathHandler.GetPathChapters)).ServeHTTP,
	})

	// 根据章节ID获取详情
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/chapters/:chapter_id",
		Handler: optAuth(wrapFunc(learningPathHandler.GetChapterByID)).ServeHTTP,
	})

	// 根据路径slug和章节slug获取章节详情
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/:path_slug/:chapter_slug",
		Handler: optAuth(wrapFunc(learningPathHandler.GetChapterBySlug)).ServeHTTP,
	})

	// ===== 学习进度路由 (需要登录) =====
	// 获取用户的学习进度
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/progress",
		Handler: optAuth(wrapFunc(learningPathHandler.GetPathProgress)).ServeHTTP,
	})

	// 获取用户已完成的章节列表
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/completed-chapters",
		Handler: optAuth(wrapFunc(learningPathHandler.GetCompletedChapters)).ServeHTTP,
	})

	// 保存学习进度
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/learning-paths/progress",
		Handler: optAuth(wrapFunc(learningPathHandler.SaveProgress)).ServeHTTP,
	})

	// 获取路径学习仪表盘
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/dashboard",
		Handler: optAuth(wrapFunc(learningPathHandler.GetPathDashboard)).ServeHTTP,
	})
}
