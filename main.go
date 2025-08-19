package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gbrayhan/microservices-go/docs"
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port string
}

// loadServerConfig loads server configuration from environment variables
func loadServerConfig() ServerConfig {
	return ServerConfig{
		Port: getEnvOrDefault("SERVER_PORT", "8080"),
	}
}

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/v1

//	@securityDefinitions.basic	BasicAuth

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	// swagger setting
	setSwaggerConfiguration()
	// load .env file
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
	// Initialize logger first based on environment
	env := getEnvOrDefault("GO_ENV", "development")
	var loggerInstance *logger.Logger
	var err error

	if env == "development" {
		loggerInstance, err = logger.NewDevelopmentLogger()
	} else {
		loggerInstance, err = logger.NewLogger()
	}

	if err != nil {
		panic(fmt.Errorf("error initializing logger: %w", err))
	}
	defer func() {
		if err := loggerInstance.Log.Sync(); err != nil {
			loggerInstance.Log.Error("Failed to sync logger", zap.Error(err))
		}
	}()

	loggerInstance.Info("Starting microservices application")

	// Load server configuration
	serverConfig := loadServerConfig()

	// Initialize application context with dependencies and logger
	appContext, err := di.SetupDependencies(loggerInstance)
	if err != nil {
		loggerInstance.Panic("Error initializing application context", zap.Error(err))
	}

	// Close relation resource
	defer func() {
		if err := appContext.Close(); err != nil {
			loggerInstance.Error("Error closing application context", zap.Error(err))
		}
	}()

	// Setup router
	router := setupRouter(appContext, loggerInstance)

	// Setup server
	server := setupServer(router, serverConfig.Port)

	// setup scheduler
	appContext.TaskScheduler.Start()

	// Start server
	// 在goroutine中启动服务器，以便可以捕获关闭信号
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			loggerInstance.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	loggerInstance.Info("Shutting down server...")

	// 创建一个带超时的上下文用于关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		loggerInstance.Fatal("Server forced to shutdown", zap.Error(err))
	}

	loggerInstance.Info("Server exiting")
}

func setupRouter(appContext *di.ApplicationContext, logger *logger.Logger) *gin.Engine {
	// Configurar Gin para usar el logger de Zap basado en el entorno
	env := getEnvOrDefault("GO_ENV", "development")
	if env == "development" {
		logger.SetupGinWithZapLoggerInDevelopment()
	} else {
		logger.SetupGinWithZapLogger()
	}

	// Crear el router después de configurar el logger
	router := gin.New()

	// set file upload configuration
	router.MaxMultipartMemory = 10 << 20 // 10 MB
	router.Static("/public", "./public")
	router.RedirectTrailingSlash = false

	// Agregar middlewares de recuperación y logger personalizados
	router.Use(gin.Recovery())
	router.Use(middlewares.CorsHeader())
	// Add middlewares
	router.Use(middlewares.ErrorHandler())
	router.Use(middlewares.GinBodyLogMiddleware(appContext.DB, appContext.Logger))
	router.Use(middlewares.SecurityHeaders())
	// Add logger middleware
	router.Use(logger.GinZapLogger())
	// Setup routes
	routes.ApplicationRouter(router, appContext)

	return router
}

func setupServer(router *gin.Engine, port string) *http.Server {
	return &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    18000 * time.Second,
		WriteTimeout:   18000 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

// Helper function
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// swagger some set
func setSwaggerConfiguration() {
	// programatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "petstore.swagger.io"
	docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
