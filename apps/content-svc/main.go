package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"

	"hot-ai-backend/apps/content-svc/handler"
	"hot-ai-backend/internal/access"
	"hot-ai-backend/internal/database"
	"hot-ai-backend/internal/middleware"
	"hot-ai-backend/internal/repository"
	"hot-ai-backend/internal/service"
)

var configFile = flag.String("f", "apps/content-svc/etc/content-svc.yaml", "the config file")

// ContentSvcConf 内容服务配置
type ContentSvcConf struct {
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
	// 付费墙开关：true=按级别拦截（默认），false=全员可读完整内容
	Paywall struct {
		Enabled bool `json:",default=false"`
	}
}

func main() {
	flag.Parse()

	var c ContentSvcConf
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
	articleRepo := repository.NewArticleRepository()

	// 付费墙开关：true 按 level 拦截；false 全员可读完整内容
	access.SetPaywallEnabled(c.Paywall.Enabled)
	fmt.Printf("[config] paywall enabled=%v\n", access.PaywallEnabled())

	// 初始化服务层
	articleService := service.NewArticleService(articleRepo)

	// 注册路由 (带 OptionalAuth 中间件，注入 user level 到 ctx)
	registerRoutes(server, c.Auth.AccessSecret, articleService)

	fmt.Printf("Starting content-svc at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// registerRoutes 注册路由
func registerRoutes(server *rest.Server, jwtSecret string, articleService *service.ArticleService) {
	// 创建处理器
	articleHandler := handler.NewArticleHandler(articleService)

	// OptionalAuth 中间件：有 token 就解析注入 level，没 token 当 level=0 (游客)
	// resolve: JWT level<2 时回查 DB 拿最新 level（处理「用户升级会员但 token 还没刷新」的场景）
	userRepo := repository.NewUserRepository()
	resolveLevel := func(ctx context.Context, userID string) int {
		u, err := userRepo.GetByID(ctx, userID)
		if err != nil || u == nil {
			return 0
		}
		return access.ComputeLevel(u)
	}
	optAuth := middleware.OptionalAuth(jwtSecret, resolveLevel)

	// wrapFunc 把 handler func 包成 http.Handler，中间件再包成 http.Handler
	wrapFunc := func(h func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			h(w, r)
		}
	}

	// ===== 文章路由 (公开，OptionalAuth 解析身份用于 access 决策) =====
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/articles",
		Handler: optAuth(wrapFunc(articleHandler.GetArticles)).ServeHTTP,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/articles/categories",
		Handler: optAuth(wrapFunc(articleHandler.GetCategories)).ServeHTTP,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/articles/:id",
		Handler: optAuth(wrapFunc(articleHandler.GetArticleByID)).ServeHTTP,
	})
}
