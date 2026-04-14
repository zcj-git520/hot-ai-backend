package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"

	"hot-ai-backend/apps/content-svc/handler"
	"hot-ai-backend/internal/database"
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

	// 初始化服务层
	articleService := service.NewArticleService(articleRepo)

	// 注册路由
	registerRoutes(server, articleService)

	fmt.Printf("Starting content-svc at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// registerRoutes 注册路由
func registerRoutes(server *rest.Server, articleService *service.ArticleService) {
	// 创建处理器
	articleHandler := handler.NewArticleHandler(articleService)

	// ===== 文章路由 =====
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/articles",
		Handler: articleHandler.GetArticles,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/articles/categories",
		Handler: articleHandler.GetCategories,
	})

	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/articles/:id",
		Handler: articleHandler.GetArticleByID,
	})
}
