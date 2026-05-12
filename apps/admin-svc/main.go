package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"hot-ai-backend/apps/admin-svc/handler"
	"hot-ai-backend/internal/database"
	"hot-ai-backend/internal/repository"
	"hot-ai-backend/internal/service"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "apps/admin-svc/etc/admin-svc.yaml", "the config file")

// AdminSvcConf admin 服务配置
type AdminSvcConf struct {
	rest.RestConf
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
	}

	// 创建服务器
	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	// 初始化仓储层
	learningPathRepo := repository.NewLearningPathRepository()

	// 初始化服务层
	adminService := service.NewAdminService(learningPathRepo)

	// 注册路由
	registerRoutes(server, adminService)

	fmt.Printf("Starting admin-svc at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

func registerRoutes(server *rest.Server, adminService *service.AdminService) {
	learningPathHandler := handler.NewAdminLearningPathHandler(adminService)
	chapterHandler := handler.NewAdminChapterHandler(adminService)

	// ===== 学习路径管理 =====
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/admin/learning-paths",
		Handler: learningPathHandler.GetLearningPaths,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/admin/learning-paths/:id",
		Handler: learningPathHandler.GetLearningPathByID,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/admin/learning-paths",
		Handler: learningPathHandler.CreateLearningPath,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPut,
		Path:    "/admin/learning-paths/:id",
		Handler: learningPathHandler.UpdateLearningPath,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodDelete,
		Path:    "/admin/learning-paths/:id",
		Handler: learningPathHandler.DeleteLearningPath,
	})

	// ===== 审核流程 =====
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/admin/learning-paths/:id/submit-review",
		Handler: learningPathHandler.SubmitReview,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/admin/learning-paths/:id/approve",
		Handler: learningPathHandler.Approve,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/admin/learning-paths/:id/reject",
		Handler: learningPathHandler.Reject,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/admin/learning-paths/:id/publish",
		Handler: learningPathHandler.Publish,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/admin/learning-paths/:id/unpublish",
		Handler: learningPathHandler.Unpublish,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/admin/learning-paths/:id/feature",
		Handler: learningPathHandler.SetFeatured,
	})

	// ===== 章节管理 =====
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/admin/learning-paths/:path_id/chapters",
		Handler: chapterHandler.GetChapters,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/admin/learning-paths/:path_id/chapters",
		Handler: chapterHandler.CreateChapter,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/admin/chapters/:id",
		Handler: chapterHandler.GetChapterByID,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPut,
		Path:    "/admin/chapters/:id",
		Handler: chapterHandler.UpdateChapter,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodDelete,
		Path:    "/admin/chapters/:id",
		Handler: chapterHandler.DeleteChapter,
	})
}
