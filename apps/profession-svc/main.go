package main

import (
	"context"
	"flag"
	"fmt"
	"hot-ai-backend/apps/profession-svc/handler"
	"net/http"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"

	"hot-ai-backend/internal/access"
	"hot-ai-backend/internal/database"
	"hot-ai-backend/internal/middleware"
	"hot-ai-backend/internal/repository"
	"hot-ai-backend/internal/service"
)

var configFile = flag.String("f", "apps/profession-svc/etc/profession-svc.yaml", "the config file")

// ProfessionSvcConf 职业风险服务配置
type ProfessionSvcConf struct {
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

	var c ProfessionSvcConf
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
	professionRepo := repository.NewProfessionRepository()

	// 付费墙开关：true 按 level 拦截；false 全员可读完整内容
	access.SetPaywallEnabled(c.Paywall.Enabled)
	fmt.Printf("[config] paywall enabled=%v\n", access.PaywallEnabled())

	// 初始化服务层
	professionService := service.NewProfessionService(professionRepo)

	// 初始化速览 handler
	featuredHandler := handler.NewProfessionFeaturedHandler(professionService)

	// 注册路由
	registerRoutes(server, c.Auth.AccessSecret, professionService, featuredHandler)

	fmt.Printf("Starting profession-svc at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// registerRoutes 注册路由
func registerRoutes(server *rest.Server, jwtSecret string, professionService *service.ProfessionService, featuredHandler *handler.ProfessionFeaturedHandler) {
	// 创建处理器
	professionHandler := handler.NewProfessionHandler(professionService)

	// OptionalAuth 中间件：JWT level<2 时回查 DB 拿最新 level
	userRepo := repository.NewUserRepository()
	resolveLevel := func(ctx context.Context, userID string) int {
		u, err := userRepo.GetByID(ctx, userID)
		if err != nil || u == nil {
			return 0
		}
		return access.ComputeLevel(u)
	}
	optAuth := middleware.OptionalAuth(jwtSecret, resolveLevel)
	wrapFunc := func(h func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			h(w, r)
		}
	}

	// 职业列表
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/professions",
		Handler: optAuth(wrapFunc(professionHandler.GetProfessions)).ServeHTTP,
	})

	// 职业分类
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/professions/categories",
		Handler: optAuth(wrapFunc(professionHandler.GetCategories)).ServeHTTP,
	})

	// 风险等级
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/professions/risk-levels",
		Handler: optAuth(wrapFunc(professionHandler.GetRiskLevels)).ServeHTTP,
	})

	// 搜索职业
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/professions/search",
		Handler: optAuth(wrapFunc(professionHandler.SearchProfessions)).ServeHTTP,
	})

	// 速览首页 —— 必须在 :id 之前注册，否则 Go Zero 路由按注册顺序匹配会把 "featured" 当成 id
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/professions/featured",
		Handler: optAuth(wrapFunc(featuredHandler.GetFeatured)).ServeHTTP,
	})

	// 职业详情
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/professions/:id",
		Handler: optAuth(wrapFunc(professionHandler.GetProfessionByID)).ServeHTTP,
	})
}
