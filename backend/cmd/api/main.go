package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/http/middlewares"
	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/logger"

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

	"github.com/lukcba/club-pulse-system-api/backend/internal/platform/database"
)

func main() {
	// 1. Initialize Logger
	logger.InitLogger()
	logger.Info("Starting Club Pulse System API...")

	// 2. Initialize Database
	database.InitDB()
	db := database.GetDB()

	// 3. Initialize Router
	router := gin.New() // Empty router, we add middlewares manually

	// 4. Register Global Middlewares
	router.Use(gin.Recovery()) // Panics recovery
	router.Use(middlewares.SecurityHeadersMiddleware())
	router.Use(middlewares.CORSMiddleware())
	router.Use(middlewares.RateLimitMiddleware())

	// Custom logger middleware could be added here to use our slog wrapper
	router.Use(gin.Logger())

	// 5. Routes
	api := router.Group("/api/v1")

	// --- Module: Auth ---
	// Wiring dependencies manually (Dependency Injection)
	authRepo := repository.NewPostgresAuthRepository(db)
	tokenService := token.NewJWTService("SECRET_KEY_DEV") // Use env var in prod
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
	// --- Module: Facilities ---
	facilityRepository := facilitiesRepo.NewPostgresFacilityRepository(db)
	facilityUseCase := facilityApp.NewFacilityUseCases(facilityRepository)
	facilityHandler := facilitiesHTTP.NewFacilityHandler(facilityUseCase)

	facilitiesHTTP.RegisterRoutes(api, facilityHandler, authMiddleware)

	// --- Module: Membership ---
	membershipRepository := membershipRepo.NewPostgresMembershipRepository(db)
	membershipUseCase := membershipApplication.NewMembershipUseCases(membershipRepository)
	membershipHandler := membershipHTTP.NewMembershipHandler(membershipUseCase)

	membershipHTTP.RegisterRoutes(api, membershipHandler, authMiddleware)

	// --- Module: Booking ---
	bookingRepository := bookingRepo.NewPostgresBookingRepository(db)
	bookingUseCase := bookingApplication.NewBookingUseCases(bookingRepository, facilityRepository)
	bookingHandler := bookingHTTP.NewBookingHandler(bookingUseCase)

	bookingHTTP.RegisterRoutes(api, bookingHandler, authMiddleware)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "UP",
			"system":  "club-pulse-backend",
			"version": "1.0.0",
		})
	})

	// 5. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	logger.Info(fmt.Sprintf("Server listening on port %s", port))
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(fmt.Sprintf("Failed to start server: %v", err))
		os.Exit(1)
	}
}
