package main

import (
	"log"

	_ "github.com/lukcba/club-pulse-system-api/backend/docs"
	"github.com/lukcba/club-pulse-system-api/backend/internal/bootstrap"
)

// @title           Club Pulse API
// @version         1.0
// @description     Management API for Club Pulse System.
// @termsOfService  http://swagger.io/terms/

// @contact.name    API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io

// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html

// @host            localhost:8080
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	app, err := bootstrap.NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize app: %v", err)
	}

	app.Run()
}
