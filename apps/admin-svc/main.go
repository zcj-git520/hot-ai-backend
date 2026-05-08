package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"

	"hot-ai-backend/internal/database"
	"hot-ai-backend/internal/repository"
	"hot-ai-backend/internal/service"
	"hot-ai-backend/apps/admin-svc/handler"
)

var configFile = flag.String("f", "apps/admin-svc/etc/admin-svc.yaml", "the config file")

// AdminSvcConf 管理后台服务配置
type AdminSvcConf struct {
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

	var c AdminSvcConf
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
		fmt.Fprintf(os.Stderr, "[错误] 未配置数据库\n")
		os.Exit(1)
	}

	// 创建服务器
	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	// 初始化仓储层
	learningPathRepo := repository.NewLearningPathRepository()

	// 初始化服务层
	learningPathService := service.NewAdminService(learningPathRepo)

	// 注册路由
	registerRoutes(server, learningPathService)

	fmt.Printf("Starting admin-svc at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// registerRoutes 注册路由
func registerRoutes(server *rest.Server, learningPathService *service.AdminService) {
	// 创建处理器
	learningPathHandler := handler.NewLearningPathHandler(learningPathService)

	// ===== 学习路径管理路由 =====
	// 获取学习路径列表（管理后台）
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/learning-paths",
		Handler: learningPathHandler.GetLearningPaths,
	})

	// 获取学习路径详情
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/learning-paths/:id",
		Handler: learningPathHandler.GetLearningPathByID,
	})

	// 创建学习路径
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/learning-paths",
		Handler: learningPathHandler.CreateLearningPath,
	})

	// 更新学习路径
	server.AddRoute(rest.Route{
		Method:  http.MethodPut,
		Path:    "/api/admin/learning-paths/:id",
		Handler: learningPathHandler.UpdateLearningPath,
	})

	// 删除学习路径
	server.AddRoute(rest.Route{
		Method:  http.MethodDelete,
		Path:    "/api/admin/learning-paths/:id",
		Handler: learningPathHandler.DeleteLearningPath,
	})

	// ===== 章节管理路由 =====
	// 获取章节列表
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/learning-paths/:path_id/chapters",
		Handler: learningPathHandler.GetChapters,
	})

	// 获取章节详情
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/chapters/:id",
		Handler: learningPathHandler.GetChapterByID,
	})

	// 创建章节
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/learning-paths/:path_id/chapters",
		Handler: learningPathHandler.CreateChapter,
	})

	// 更新章节
	server.AddRoute(rest.Route{
		Method:  http.MethodPut,
		Path:    "/api/admin/chapters/:id",
		Handler: learningPathHandler.UpdateChapter,
	})

	// 删除章节
	server.AddRoute(rest.Route{
		Method:  http.MethodDelete,
		Path:    "/api/admin/chapters/:id",
		Handler: learningPathHandler.DeleteChapter,
	})

	// ===== 学习路径操作 =====
	// 发布学习路径
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/learning-paths/:id/publish",
		Handler: learningPathHandler.PublishLearningPath,
	})

	// 下架学习路径
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/learning-paths/:id/unpublish",
		Handler: learningPathHandler.UnpublishLearningPath,
	})

	// 设置推荐
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/learning-paths/:id/featured",
		Handler: learningPathHandler.SetFeatured,
	})
}