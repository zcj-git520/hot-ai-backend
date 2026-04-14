package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"

	"hot-ai-backend/apps/services/learning-path-svc/handler"
	"hot-ai-backend/internal/database"
	"hot-ai-backend/internal/repository"
	"hot-ai-backend/internal/service"
)

var configFile = flag.String("f", "apps/services/learning-path-svc/etc/learning-path-svc.yaml", "the config file")

// LearningPathSvcConf 学习路径服务配置
type LearningPathSvcConf struct {
	rest.RestConf
	// 数据库配置
	DataSource struct {
		MySQL struct {
			DSN string `json:",optional"`
		}
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
	registerRoutes(server, learningPathService)

	fmt.Printf("Starting learning-path-svc at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// registerRoutes 注册路由
func registerRoutes(server *rest.Server, learningPathService *service.LearningPathService) {
	// 创建处理器
	learningPathHandler := handler.NewLearningPathHandler(learningPathService)

	// ===== 学习路径路由 =====
	// 获取路径列表
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths",
		Handler: learningPathHandler.GetLearningPaths,
	})

	// 根据ID获取路径详情
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/:id",
		Handler: learningPathHandler.GetLearningPathByID,
	})

	// 根据slug获取路径详情
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/slug/:slug",
		Handler: learningPathHandler.GetLearningPathBySlug,
	})

	// 获取推荐路径
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/featured",
		Handler: learningPathHandler.GetFeaturedPaths,
	})

	// 获取难度等级信息
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/levels",
		Handler: learningPathHandler.GetLevelInfo,
	})

	// ===== 章线路由 =====
	// 获取路径的所有章节
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/:path_id/chapters",
		Handler: learningPathHandler.GetPathChapters,
	})

	// 根据章节ID获取详情
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/chapters/:chapter_id",
		Handler: learningPathHandler.GetChapterByID,
	})

	// 根据路径slug和章节slug获取章节详情
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/:path_slug/:chapter_slug",
		Handler: learningPathHandler.GetChapterBySlug,
	})

	// ===== 学习进度路由 =====
	// 获取用户的学习进度
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/progress",
		Handler: learningPathHandler.GetPathProgress,
	})

	// 获取用户已完成的章节列表
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/completed-chapters",
		Handler: learningPathHandler.GetCompletedChapters,
	})

	// 保存学习进度
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/learning-paths/progress",
		Handler: learningPathHandler.SaveProgress,
	})

	// 获取路径学习仪表盘
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/dashboard",
		Handler: learningPathHandler.GetPathDashboard,
	})
}
