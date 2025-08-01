package router

import (
	"mindx/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

// InitRouter initializes the Echo router with middleware and routes
func InitRouter(db *gorm.DB) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	// Initialize handlers
	h := handlers.NewHandler(db)

	// Routes
	e.POST("/evaluate", h.EvaluateRisk)
	e.GET("/students", h.ListStudents)

	return e
}