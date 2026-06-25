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
	professionRepo := repository.NewProfessionRepository()

	// 初始化服务层
	adminService := service.NewAdminService(learningPathRepo)
	professionService := service.NewProfessionService(professionRepo)
	userRepo := repository.NewUserRepository()
	favoriteRepo := repository.NewFavoriteRepository()
	userService := service.NewUserService(userRepo, favoriteRepo)

	// 注册路由
	registerRoutes(server, adminService, professionService, userService)

	fmt.Printf("Starting admin-svc at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

func registerRoutes(server *rest.Server, adminService *service.AdminService, professionService *service.ProfessionService, userService *service.UserService) {
	learningPathHandler := handler.NewAdminLearningPathHandler(adminService)
	professionHandler := handler.NewProfessionHandler(professionService)
	chapterHandler := handler.NewAdminChapterHandler(adminService)
	toolReviewHandler := handler.NewAdminToolReviewHandler()
	articleReviewHandler := handler.NewAdminArticleReviewHandler()
	userHandler := handler.NewUserHandler(userService)
	accessLevelHandler := handler.NewAdminAccessLevelHandler()

	// ===== 学习路径管理 =====
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/learning-paths",
		Handler: learningPathHandler.GetLearningPaths,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/learning-paths/:id",
		Handler: learningPathHandler.GetLearningPathByID,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/learning-paths",
		Handler: learningPathHandler.CreateLearningPath,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPut,
		Path:    "/api/admin/learning-paths/:id",
		Handler: learningPathHandler.UpdateLearningPath,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodDelete,
		Path:    "/api/admin/learning-paths/:id",
		Handler: learningPathHandler.DeleteLearningPath,
	})

	// ===== 审核流程 =====
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/learning-paths/:id/submit-review",
		Handler: learningPathHandler.SubmitReview,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/learning-paths/:id/approve",
		Handler: learningPathHandler.Approve,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/learning-paths/:id/reject",
		Handler: learningPathHandler.Reject,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/learning-paths/:id/publish",
		Handler: learningPathHandler.Publish,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/learning-paths/:id/unpublish",
		Handler: learningPathHandler.Unpublish,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/learning-paths/:id/feature",
		Handler: learningPathHandler.SetFeatured,
	})

	// ===== 章节管理 =====
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/learning-paths/:path_id/chapters",
		Handler: chapterHandler.GetChapters,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/learning-paths/:path_id/chapters",
		Handler: chapterHandler.CreateChapter,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/chapters/:id",
		Handler: chapterHandler.GetChapterByID,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPut,
		Path:    "/api/admin/chapters/:id",
		Handler: chapterHandler.UpdateChapter,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodDelete,
		Path:    "/api/admin/chapters/:id",
		Handler: chapterHandler.DeleteChapter,
	})

	// ===== 工具审核 =====
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/tools/pending",
		Handler: toolReviewHandler.GetPendingTools,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/tools/:id",
		Handler: toolReviewHandler.GetToolDetailForReview,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/tools/:id/approve",
		Handler: toolReviewHandler.ApproveTool,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/tools/:id/reject",
		Handler: toolReviewHandler.RejectTool,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/tools/:id/request-revision",
		Handler: toolReviewHandler.RequestRevision,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/tools/:id/online",
		Handler: toolReviewHandler.SetToolOnline,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/tools/:id/offline",
		Handler: toolReviewHandler.SetToolOffline,
	})

	// ===== 文章审核 =====
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/articles/pending",
		Handler: articleReviewHandler.GetPendingArticles,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/articles/:id",
		Handler: articleReviewHandler.GetArticleDetailForReview,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/articles/:id/approve",
		Handler: articleReviewHandler.ApproveArticle,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/articles/:id/reject",
		Handler: articleReviewHandler.RejectArticle,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/articles/:id/publish",
		Handler: articleReviewHandler.PublishArticle,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/articles/:id/unpublish",
		Handler: articleReviewHandler.UnpublishArticle,
	})

	// ===== 职业管理 =====
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/professions",
		Handler: professionHandler.GetProfessions,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/professions/categories",
		Handler: professionHandler.GetCategories,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/professions/risk-levels",
		Handler: professionHandler.GetRiskLevels,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/professions/:id",
		Handler: professionHandler.GetProfessionByID,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/professions",
		Handler: professionHandler.CreateProfession,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPut,
		Path:    "/api/admin/professions/:id",
		Handler: professionHandler.UpdateProfession,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodDelete,
		Path:    "/api/admin/professions/:id",
		Handler: professionHandler.DeleteProfession,
	})

	// ===== 用户管理 =====
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/users",
		Handler: userHandler.GetUsers,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/users/:id",
		Handler: userHandler.GetUserByID,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPut,
		Path:    "/api/admin/users/:id/role",
		Handler: userHandler.UpdateUserRole,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/users/:id/disable",
		Handler: userHandler.DisableUser,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/users/:id/enable",
		Handler: userHandler.EnableUser,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/users/:id/logs",
		Handler: userHandler.GetAdminLogs,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/admin/users",
		Handler: userHandler.CreateUser,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPut,
		Path:    "/api/admin/users/:id/password",
		Handler: userHandler.UpdatePassword,
	})

	// ===== 访问级别管理 =====
	server.AddRoute(rest.Route{
		Method:  http.MethodPut,
		Path:    "/api/admin/content/access-level",
		Handler: accessLevelHandler.SetAccessLevel,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/admin/content/access-level",
		Handler: accessLevelHandler.GetAccessLevel,
	})
}
