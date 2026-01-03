package bootstrap

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/http/middlewares"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/logger"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/middleware"
	platformRedis "github.com/lukcba/club-pulse-system-api/backend/internal/platform/redis"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/tracing"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/sdk/trace"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	Engine         *gin.Engine
	srv            *http.Server
	Port           string
	TracerProvider *trace.TracerProvider
}

func NewServer() *Server {
	// Init Tracing
	tp, err := tracing.InitTracer("club-pulse-backend")
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to init tracer: %v", err))
		// We can continue without tracing or fatal exit. Let's log and continue for now.
	}

	// Initialize Router
	router := gin.New()

	// Register Global Middlewares
	router.Use(gin.Recovery())

	// OpenTelemetry Middleware (Must be first to capture everything)
	router.Use(otelgin.Middleware("club-pulse-backend"))

	router.Use(middlewares.SecurityHeadersMiddleware())
	router.Use(middlewares.CORSMiddleware())

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Redis Rate Limiting
	middlewares.InitRedisRateLimiter(100, time.Minute) // 100 req/min
	router.Use(middlewares.RedisRateLimitMiddleware())

	// Health Check
	router.GET("/healthz", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		status := gin.H{
			"status":    "ok",
			"timestamp": time.Now(),
			"services": gin.H{
				"database": "ok",
				"redis":    "ok",
			},
		}
		httpStatus := http.StatusOK

		// Check DB
		db := database.GetDB()
		if db == nil {
			status["status"] = "error"
			status["services"].(gin.H)["database"] = "connection_is_nil"
			httpStatus = http.StatusServiceUnavailable
		} else {
			sqlDB, err := db.DB()
			if err != nil || sqlDB.Ping() != nil {
				status["status"] = "error"
				status["services"].(gin.H)["database"] = "unreachable"
				httpStatus = http.StatusServiceUnavailable
			}
		}

		// Check Redis
		rdb := platformRedis.GetClient()
		if rdb == nil {
			status["services"].(gin.H)["redis"] = "connection_is_nil"
			// Redis failure might be non-critical for basic API (partial degradation), but for "Healthz" (Liveness) it usually means unhealthy.
			// Let's mark as error.
			status["status"] = "error"
			httpStatus = http.StatusServiceUnavailable
		} else {
			if err := rdb.Ping(ctx); err != nil {
				status["services"].(gin.H)["redis"] = "unreachable"
				status["status"] = "error"
				httpStatus = http.StatusServiceUnavailable
			}
		}

		c.JSON(httpStatus, status)
	})

	router.Use(middleware.StructuredLogger())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Server{
		Engine:         router,
		Port:           port,
		TracerProvider: tp,
	}
}

func (s *Server) Start() {
	s.srv = &http.Server{
		Addr:    ":" + s.Port,
		Handler: s.Engine,
	}

	go func() {
		logger.Info(fmt.Sprintf("Server listening on port %s", s.Port))
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("Failed to start server: %v", err))
			os.Exit(1)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) error {
	logger.Info("Shutting down server...")

	// Flush Traces
	if s.TracerProvider != nil {
		if err := s.TracerProvider.Shutdown(ctx); err != nil {
			logger.Error(fmt.Sprintf("Error shutting down tracer provider: %v", err))
		}
	}

	return s.srv.Shutdown(ctx)
}
