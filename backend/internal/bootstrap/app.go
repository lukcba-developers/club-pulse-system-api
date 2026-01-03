package bootstrap

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/logger"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/middleware"

	// Module Imports
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/application"
	authHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/http"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/repository"
	authService "github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/service"
	"github.com/lukcba/club-pulse-system-api/backend/internal/modules/auth/infrastructure/token"

	userApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/application"
	userHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/http"
	userRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/user/infrastructure/repository"

	facilityApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/application"
	facilitiesHTTP "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/infrastructure/http"
	facilitiesRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/facilities/infrastructure/repository"

	membershipApplication "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/application"
	membershipHTTP "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/infrastructure/http"
	membershipRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/membership/infrastructure/repository"

	bookingApplication "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/application"
	bookingHTTP "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/http"
	bookingRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/booking/infrastructure/repository"

	accessApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/access/application"
	accessHTTP "github.com/lukcba/club-pulse-system-api/backend/internal/modules/access/infrastructure/http"
	accessRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/access/infrastructure/repository"

	attendanceApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/attendance/application"
	attendanceHTTP "github.com/lukcba/club-pulse-system-api/backend/internal/modules/attendance/infrastructure/http"
	attendanceRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/attendance/infrastructure/repository"

	disciplineApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/application"
	disciplineHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/infrastructure/http"
	disciplineRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/disciplines/infrastructure/repository"

	paymentGateway "github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/infrastructure/gateways"
	paymentHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/infrastructure/http"
	paymentRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/payment/infrastructure/repository"

	notificationProviders "github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/infrastructure/providers"
	notificationService "github.com/lukcba/club-pulse-system-api/backend/internal/modules/notification/service"

	clubApp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/application"
	clubHttp "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/infrastructure/http"
	clubRepo "github.com/lukcba/club-pulse-system-api/backend/internal/modules/club/infrastructure/repository"
)

type App struct {
	Infrastructure *Infrastructure
	Server         *Server
}

func NewApp() (*App, error) {
	// 1. Init Infrastructure (DB, Redis, Logs)
	logger.InitLogger()
	logger.Info("Starting Club Pulse System API v2.0 (High Performance Edition - Refactored)...")

	infra, err := InitInfrastructure()
	if err != nil {
		return nil, err
	}

	// 2. Init Server
	server := NewServer()

	// 3. Setup Health Check (Basic) - We can move this to a module if it grows
	server.Engine.GET("/health", healthCheckHandler(infra))

	// 4. Register Modules
	v1 := server.Engine.Group("/api/v1")

	// We need clubRepo for Middleware
	// Initialize it early or restructure.
	clubRepository := clubRepo.NewPostgresClubRepository(infra.DB)

	// Middleware Instantiation
	// REMOVED GLOBAL USAGE: v1.Use(middleware.TenantMiddleware(clubRepository))
	// Instead, we will pass it to modules that need it, to be applied AFTER Auth.
	tenantMiddleware := middleware.TenantMiddleware(clubRepository)

	registerModules(v1, infra, tenantMiddleware)

	return &App{
		Infrastructure: infra,
		Server:         server,
	}, nil
}

func (app *App) Run() {
	app.Server.Start()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	app.Infrastructure.Shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.Server.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: " + err.Error())
	}

	logger.Info("Server exited properly")
}

func registerModules(api *gin.RouterGroup, infra *Infrastructure, tenantMiddleware gin.HandlerFunc) {
	db := infra.DB

	// --- Module: Auth ---
	authRepo := repository.NewPostgresAuthRepository(db)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "SECRET_KEY_DEV"
	}
	tokenService := token.NewJWTService(jwtSecret)
	googleAuthService := authService.NewGoogleAuthService()
	authUseCase := application.NewAuthUseCases(authRepo, tokenService, googleAuthService)
	authHandler := authHttp.NewAuthHandler(authUseCase)
	authMiddleware := authHttp.AuthMiddleware(tokenService)

	authHttp.RegisterRoutes(api, authHandler, authMiddleware)

	// --- Module: User ---
	userRepository := userRepo.NewPostgresUserRepository(db)
	userUseCase := userApp.NewUserUseCases(userRepository)
	userHandler := userHttp.NewUserHandler(userUseCase)

	userHttp.RegisterRoutes(api, userHandler, authMiddleware, tenantMiddleware)

	// --- Module: Facilities ---
	facilityRepository := facilitiesRepo.NewPostgresFacilityRepository(db)
	facilityUseCase := facilityApp.NewFacilityUseCases(facilityRepository)
	facilityHandler := facilitiesHTTP.NewFacilityHandler(facilityUseCase)

	facilitiesHTTP.RegisterRoutes(api, facilityHandler, authMiddleware, tenantMiddleware)

	// --- Semantic Search (Facilities) ---
	semanticSearchUseCase := facilityApp.NewSemanticSearchUseCase(facilityRepository)
	searchHandler := facilitiesHTTP.NewSearchHandler(semanticSearchUseCase)
	facilitiesHTTP.RegisterSearchRoutes(api, searchHandler)
	logger.Info("Semantic search enabled for facilities")

	// --- Module: Membership ---
	membershipRepository := membershipRepo.NewPostgresMembershipRepository(db)
	membershipUseCase := membershipApplication.NewMembershipUseCases(membershipRepository)
	membershipHandler := membershipHTTP.NewMembershipHandler(membershipUseCase)

	membershipHTTP.RegisterRoutes(api, membershipHandler, authMiddleware, tenantMiddleware)

	// --- Module: Notification (Real Providers) ---
	var emailProvider notificationService.EmailProvider
	if mt := os.Getenv("SENDGRID_API_KEY"); mt != "" {
		emailProvider = notificationProviders.NewSendGridProvider(mt, os.Getenv("SENDGRID_FROM_NAME"), os.Getenv("SENDGRID_FROM_EMAIL"))
		logger.Info("SendGrid Email Provider enabled")
	}

	var smsProvider notificationService.SMSProvider
	if ts := os.Getenv("TWILIO_ACCOUNT_SID"); ts != "" {
		smsProvider = notificationProviders.NewTwilioProvider(ts, os.Getenv("TWILIO_AUTH_TOKEN"), os.Getenv("TWILIO_FROM_NUMBER"))
		logger.Info("Twilio SMS Provider enabled")
	}

	notifier := notificationService.NewNotificationService(emailProvider, smsProvider)

	// --- Module: Booking ---
	bookingRepository := bookingRepo.NewPostgresBookingRepository(db)
	recurringRepository := bookingRepo.NewPostgresRecurringRepository(db)
	bookingUseCase := bookingApplication.NewBookingUseCases(bookingRepository, recurringRepository, facilityRepository, notifier)
	bookingHandler := bookingHTTP.NewBookingHandler(bookingUseCase)

	bookingHTTP.RegisterRoutes(api, bookingHandler, authMiddleware, tenantMiddleware)

	// --- Module: Access (New) ---
	accessRepository := accessRepo.NewPostgresAccessRepository(db)
	accessUseCase := accessApp.NewAccessUseCases(accessRepository, userRepository, membershipRepository)
	accessHandler := accessHTTP.NewAccessHandler(accessUseCase)

	accessHTTP.RegisterRoutes(api, accessHandler, authMiddleware, tenantMiddleware)

	// --- Module: Attendance (New) ---
	attendanceRepository := attendanceRepo.NewPostgresAttendanceRepository(db)
	attendanceUseCase := attendanceApp.NewAttendanceUseCases(attendanceRepository, userRepository, membershipRepository)
	attendanceHandler := attendanceHTTP.NewAttendanceHandler(attendanceUseCase)

	attendanceHTTP.RegisterRoutes(api, attendanceHandler, authMiddleware, tenantMiddleware)

	// --- Module: Disciplines (New) ---
	dRepo := disciplineRepo.NewPostgresDisciplineRepository(db)
	tRepo := disciplineRepo.NewPostgresTournamentRepository(db)
	dUseCase := disciplineApp.NewDisciplineUseCases(dRepo, tRepo, userRepository)
	dHandler := disciplineHttp.NewDisciplineHandler(dUseCase)

	disciplineHttp.RegisterRoutes(api, dHandler, authMiddleware, tenantMiddleware)

	// --- Module: Payment ---
	paymentRepo := paymentRepo.NewPostgresPaymentRepository(db)
	paymentProc := paymentGateway.NewMercadoPagoGateway()
	paymentHandler := paymentHttp.NewPaymentHandler(paymentRepo, paymentProc)
	paymentHttp.RegisterRoutes(api, paymentHandler, authMiddleware, tenantMiddleware)

	// --- Module: Club (Super Admin) ---
	clubRepository := clubRepo.NewPostgresClubRepository(db)
	clubUseCase := clubApp.NewClubUseCases(clubRepository)
	clubHandler := clubHttp.NewClubHandler(clubUseCase)

	// Register Club Routes
	clubHttp.RegisterRoutes(api, clubHandler, authMiddleware, tenantMiddleware)
}

func healthCheckHandler(infra *Infrastructure) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check Redis
		redisStatus := "UP"
		if err := infra.Redis.Ping(context.Background()); err != nil {
			redisStatus = "DOWN"
		}

		// Check DB
		dbStatus := "UP"
		sqlDB, err := infra.DB.DB()
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
	}
}
