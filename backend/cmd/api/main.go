package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/audit"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/http/middlewares"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/logger"
	platformRedis "github.com/lukcba/club-pulse-system-api/backend/internal/platform/redis"

	// Auth Module Imports
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/application"
	authHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/http"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/repository"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/token"

	// User Module Imports
	userApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/application"
	userHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/http"
	userRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/repository"

	// Facilities Module Imports
	facilityApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/application"
	facilitiesHTTP "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/infrastructure/http"
	facilitiesRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/infrastructure/repository"

	// Membership Module Imports
	membershipApplication "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/application"
	membershipHTTP "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/infrastructure/http"
	membershipRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/infrastructure/repository"

	// Booking Module Imports
	bookingApplication "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/application"
	bookingHTTP "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/http"
	bookingRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/repository"

	// Access Module Imports
	accessApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/access/application"
	accessHTTP "github.com/lukcba/club-pulse-system-api/backend/internal/modules/access/infrastructure/http"
	accessRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/access/infrastructure/repository"

	// Attendance Module Imports
	attendanceApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/attendance/application"
	attendanceHTTP "github.com/lukcba/club-pulse-system-api/backend/internal/modules/attendance/infrastructure/http"
	attendanceRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/attendance/infrastructure/repository"

	// Notification Module
	notificationService "github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/service"

	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
)

func main() {
	// 1. Initialize Logger
	logger.InitLogger()
	logger.Info("Starting Club Pulse System API v2.0 (High Performance Edition)...")

	// 2. Initialize Database
	database.InitDB()
	db := database.GetDB()

	// 3. Initialize Redis
	platformRedis.InitRedis()
	redisClient := platformRedis.GetClient()
	logger.Info("Redis client initialized")

	// 4. Initialize Audit Queue (Zero-Impact Logging)
	auditQueue := audit.NewAuditQueue(db, 5*time.Minute)
	auditQueue.StartFlushWorker()
	logger.Info("Audit queue worker started (flush every 5 min)")

	// 5. Initialize Router
	router := gin.New()

	// 6. Register Global Middlewares
	router.Use(gin.Recovery())
	router.Use(middlewares.SecurityHeadersMiddleware())
	router.Use(middlewares.CORSMiddleware())

	// Use Redis-based Rate Limiting (distributed)
	middlewares.InitRedisRateLimiter(100, time.Minute) // 100 req/min
	router.Use(middlewares.RedisRateLimitMiddleware())
	logger.Info("Redis Rate Limiter enabled (100 req/min)")

	router.Use(gin.Logger())

	// 7. Routes
	api := router.Group("/api/v1")

	// --- Module: Auth ---
	authRepo := repository.NewPostgresAuthRepository(db)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "SECRET_KEY_DEV"
	}
	tokenService := token.NewJWTService(jwtSecret)
	authUseCase := application.NewAuthUseCases(authRepo, tokenService)
	authHandler := authHttp.NewAuthHandler(authUseCase)
	authMiddleware := authHttp.AuthMiddleware(tokenService)

	authHttp.RegisterRoutes(api, authHandler, authMiddleware)

	// --- Module: User ---
	userRepository := userRepo.NewPostgresUserRepository(db)
	userUseCase := userApp.NewUserUseCases(userRepository)
	userHandler := userHttp.NewUserHandler(userUseCase)

	userHttp.RegisterRoutes(api, userHandler, authMiddleware)

	// --- Module: Facilities ---
	facilityRepository := facilitiesRepo.NewPostgresFacilityRepository(db)
	facilityUseCase := facilityApp.NewFacilityUseCases(facilityRepository)
	facilityHandler := facilitiesHTTP.NewFacilityHandler(facilityUseCase)

	facilitiesHTTP.RegisterRoutes(api, facilityHandler, authMiddleware)

	// --- Semantic Search (Facilities) ---
	semanticSearchUseCase := facilityApp.NewSemanticSearchUseCase(facilityRepository)
	searchHandler := facilitiesHTTP.NewSearchHandler(semanticSearchUseCase)
	facilitiesHTTP.RegisterSearchRoutes(api, searchHandler)
	logger.Info("Semantic search enabled for facilities")

	// --- Module: Membership ---
	membershipRepository := membershipRepo.NewPostgresMembershipRepository(db)
	membershipUseCase := membershipApplication.NewMembershipUseCases(membershipRepository)
	membershipHandler := membershipHTTP.NewMembershipHandler(membershipUseCase)

	membershipHTTP.RegisterRoutes(api, membershipHandler, authMiddleware)

	// --- Module: Booking ---
	notifier := notificationService.NewConsoleNotificationSender()
	bookingRepository := bookingRepo.NewPostgresBookingRepository(db)
	bookingUseCase := bookingApplication.NewBookingUseCases(bookingRepository, facilityRepository, notifier)
	bookingHandler := bookingHTTP.NewBookingHandler(bookingUseCase)

	bookingHTTP.RegisterRoutes(api, bookingHandler, authMiddleware)

	// --- Module: Access (New) ---
	accessRepository := accessRepo.NewPostgresAccessRepository(db)
	accessUseCase := accessApp.NewAccessUseCases(accessRepository, userRepository, membershipRepository)
	accessHandler := accessHTTP.NewAccessHandler(accessUseCase)

	accessHTTP.RegisterRoutes(api, accessHandler, authMiddleware)

	// --- Module: Attendance (New) ---
	attendanceRepository := attendanceRepo.NewPostgresAttendanceRepository(db)
	attendanceUseCase := attendanceApp.NewAttendanceUseCases(attendanceRepository, userRepository)
	attendanceHandler := attendanceHTTP.NewAttendanceHandler(attendanceUseCase)

	attendanceHTTP.RegisterRoutes(api, attendanceHandler, authMiddleware)

	// --- Health Endpoints ---
	router.GET("/health", func(c *gin.Context) {
		// Check Redis
		redisStatus := "UP"
		if err := redisClient.Ping(context.Background()); err != nil {
			redisStatus = "DOWN"
		}

		// Check DB
		dbStatus := "UP"
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			dbStatus = "DOWN"
		}

		status := "UP"
		if redisStatus == "DOWN" || dbStatus == "DOWN" {
			status = "DEGRADED"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  status,
			"system":  "club-pulse-backend",
			"version": "2.0.0",
			"services": gin.H{
				"database": dbStatus,
				"redis":    redisStatus,
			},
		})
	})

	// 8. Start Server with Graceful Shutdown
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Run server in goroutine
	go func() {
		logger.Info(fmt.Sprintf("Server listening on port %s", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("Failed to start server: %v", err))
			os.Exit(1)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Stop audit queue worker
	auditQueue.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error(fmt.Sprintf("Server forced to shutdown: %v", err))
	}

	logger.Info("Server exited properly")
}
