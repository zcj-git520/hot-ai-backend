package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"

	"hot-ai-backend/apps/gateway/handler"
	"hot-ai-backend/apps/gateway/middleware"
	"hot-ai-backend/internal/database"
	"hot-ai-backend/internal/repository"
	"hot-ai-backend/internal/service"
	"hot-ai-backend/internal/utils/mailutil"
	lpathHandler "hot-ai-backend/apps/services/learning-path-svc/handler"
)

var configFile = flag.String("f", "apps/gateway/etc/gateway.yaml", "the config file")

// GatewayConf 网关配置
type GatewayConf struct {
	rest.RestConf
	Auth struct {
		AccessSecret string `json:",default=your-secret-key-change-in-production"`
		AccessExpire int    `json:",default=86400"`
	}
	// 数据库配置
	DataSource struct {
		MySQL struct {
			DSN string `json:",optional"`
		}
	}
	// 邮件配置
	Mail struct {
		SMTPServer string `json:",optional, default="`
		SMTPPort   int    `json:",default=587"`
		Username   string `json:",optional"`
		Password   string `json:",optional"`
		FromName   string `json:",default=Hot AI"`
		FromEmail  string `json:",optional"`
	}
}

func main() {
	flag.Parse()

	var c GatewayConf
	if err := conf.Load(*configFile, &c); err != nil {
		fmt.Fprintf(os.Stderr, "load config error: %v\n", err)
		os.Exit(1)
	}

	// 初始化数据库
	if c.DataSource.MySQL.DSN != "" {
		// 从 DSN 字符串解析数据库配置
		dbConfig, err := database.ParseDSN(c.DataSource.MySQL.DSN)
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse DSN error: %v\n", err)
			os.Exit(1)
		}

		if err := database.InitDB(*dbConfig); err != nil {
			fmt.Fprintf(os.Stderr, "init database error: %v\n", err)
			os.Exit(1)
		}

		//// 自动迁移表结构
		//if err := database.AutoMigrate(); err != nil {
		//	fmt.Fprintf(os.Stderr, "auto migrate error: %v\n", err)
		//	os.Exit(1)
		//}

		fmt.Println("Database initialized successfully")
	} else {
		fmt.Println("[警告] 未配置数据库，将使用内存模式（仅用于测试）")
	}

	// 创建服务器
	server := rest.MustNewServer(c.RestConf, rest.WithCors())
	defer server.Stop()

	// 中间件
	authMiddleware := middleware.NewAuthMiddleware(c.Auth.AccessSecret, int64(c.Auth.AccessExpire))

	// 初始化服务层（依赖注入）
	userRepo := repository.NewUserRepository()
	favoriteRepo := repository.NewFavoriteRepository()
	authService := service.NewAuthService(userRepo, favoriteRepo, c.Auth.AccessSecret)
	userService := service.NewUserService(userRepo, favoriteRepo)
	favoriteService := service.NewFavoriteService(favoriteRepo)

	// 初始化学习路径服务
	learningPathRepo := repository.NewLearningPathRepository()
	learningPathService := service.NewLearningPathService(learningPathRepo)

	// 初始化文章服务
	articleRepo := repository.NewArticleRepository()
	articleService := service.NewArticleService(articleRepo)

	// 初始化职业服务
	professionRepo := repository.NewProfessionRepository()
	professionService := service.NewProfessionService(professionRepo)

	// 初始化工具服务
	toolDB, err := database.GetDB().DB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "get database connection error: %v\n", err)
		os.Exit(1)
	}
	toolRepo := repository.NewToolRepository(toolDB)
	toolService := service.NewToolService(toolRepo)

	// 配置邮件服务
	if c.Mail.SMTPServer != "" {
		authService.SetSMTPConfig(&mailutil.SMTPConfig{
			Server:    c.Mail.SMTPServer,
			Port:      c.Mail.SMTPPort,
			Username:  c.Mail.Username,
			Password:  c.Mail.Password,
			FromName:  c.Mail.FromName,
			FromEmail: c.Mail.FromEmail,
		})
	} else {
		fmt.Println("[提示] 未配置 SMTP，将使用日志模式发送邮件验证码")
	}

	// 注册业务路由
	registerRoutes(server, authMiddleware, authService, userService, favoriteService, learningPathService, articleService, professionService, toolService)

	fmt.Printf("Starting gateway at %s:%d...\n", c.Host, c.Port)

	server.Start()
}

// registerRoutes 注册所有业务路由
func registerRoutes(server *rest.Server, authMiddleware *middleware.AuthMiddleware, authService *service.AuthService, userService *service.UserService, favoriteService *service.FavoriteService, learningPathService *service.LearningPathService, articleService *service.ArticleService, professionService *service.ProfessionService, toolService *service.ToolService) {
	// 创建处理器并注入依赖
	authHandler := handler.NewAuthHandler()
	authHandler.SetAuthService(authService)

	// 创建用户处理器并注入依赖
	userHandler := handler.NewUserHandler(userService, favoriteService)

	// 创建学习路径处理器
	learningPathHandler := lpathHandler.NewLearningPathHandler(learningPathService)

	// 创建首页处理器
	homeHandler := handler.NewHomeHandler(articleService, professionService, learningPathService, userService, toolService)

	// ===== 认证路由 (公开) =====
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/auth/register",
		Handler: authHandler.Register,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/auth/send-registration-code",
		Handler: authHandler.SendRegistrationCode,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/auth/login",
		Handler: authHandler.Login,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/auth/refresh",
		Handler: authHandler.RefreshToken,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/auth/send-verification-email",
		Handler: authHandler.SendVerificationEmail,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/auth/verify-email",
		Handler: authHandler.VerifyEmail,
	})

	// ===== 认证路由 (需要认证) =====
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/auth/logout",
		Handler: authMiddleware.Handle(authHandler.Logout),
	})

	// ===== 用户路由 (需要认证) =====
	// 获取/更新用户资料
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/user/profile",
		Handler: authMiddleware.Handle(userHandler.GetProfile),
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPut,
		Path:    "/api/user/profile",
		Handler: authMiddleware.Handle(userHandler.UpdateProfile),
	})

	// 更新用户偏好
	server.AddRoute(rest.Route{
		Method:  http.MethodPut,
		Path:    "/api/user/preferences",
		Handler: authMiddleware.Handle(userHandler.UpdatePreferences),
	})

	// 修改密码
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/user/change-password",
		Handler: authMiddleware.Handle(userHandler.ChangePassword),
	})

	// 收藏管理
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/user/favorites",
		Handler: authMiddleware.Handle(userHandler.GetFavorites),
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/user/favorites",
		Handler: authMiddleware.Handle(userHandler.AddFavorite),
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodDelete,
		Path:    "/api/user/favorites/{id}",
		Handler: authMiddleware.Handle(userHandler.DeleteFavorite),
	})

	// ===== 首页路由 (公开) =====
	// 获取首页统计数据
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/home/stats",
		Handler: homeHandler.GetHomeStats,
	})

	// 获取热门话题
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/home/topics",
		Handler: homeHandler.GetHomeTopics,
	})

	// 获取热门排行
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/home/rankings",
		Handler: homeHandler.GetHomeRankings,
	})

	// ===== 学习路径路由 (公开) =====
	// 获取学习路径列表
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths",
		Handler: learningPathHandler.GetLearningPaths,
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

	// 根据 ID 获取学习路径详情
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/id/{id}",
		Handler: learningPathHandler.GetLearningPathByID,
	})

	// 根据 slug 获取学习路径详情
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/slug/{slug}",
		Handler: learningPathHandler.GetLearningPathBySlug,
	})

	// 获取路径的所有章节
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/{path_id}/chapters",
		Handler: learningPathHandler.GetPathChapters,
	})

	// 根据章节 ID 获取详情
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/chapters/{chapter_id}",
		Handler: learningPathHandler.GetChapterByID,
	})

	// 根据路径 slug 和章节 slug 获取章节详情
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/{path_slug}/chapters/{chapter_slug}",
		Handler: learningPathHandler.GetChapterBySlug,
	})

	// 获取路径学习仪表盘（综合统计）
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/dashboard",
		Handler: learningPathHandler.GetPathDashboard,
	})

	// ===== 学习路径路由 (需要认证) =====
	// 获取用户的学习进度
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/progress",
		Handler: authMiddleware.Handle(learningPathHandler.GetPathProgress),
	})

	// 获取用户已完成的章节列表
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/learning-paths/completed-chapters",
		Handler: authMiddleware.Handle(learningPathHandler.GetCompletedChapters),
	})

	// 保存学习进度
	server.AddRoute(rest.Route{
		Method:  http.MethodPost,
		Path:    "/api/learning-paths/save-progress",
		Handler: authMiddleware.Handle(learningPathHandler.SaveProgress),
	})
}
