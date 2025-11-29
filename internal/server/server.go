package server

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go-rest-api-template/internal/handler"
)

// Server holds the HTTP server dependencies
type Server struct {
	router *gin.Engine
	db     *pgxpool.Pool
	logger *slog.Logger
}

// New creates a new Server instance
func New(db *pgxpool.Pool, logger *slog.Logger, debug bool) *Server {
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	s := &Server{
		router: router,
		db:     db,
		logger: logger,
	}

	s.setupMiddleware()
	s.setupRoutes()

	return s
}

// setupMiddleware configures middleware for the server
func (s *Server) setupMiddleware() {
	// Recovery middleware
	s.router.Use(gin.Recovery())

	// Request logging middleware using slog
	s.router.Use(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		s.logger.Info("request",
			slog.String("method", c.Request.Method),
			slog.String("path", path),
			slog.String("query", query),
			slog.Int("status", status),
			slog.Duration("latency", latency),
			slog.String("client_ip", c.ClientIP()),
		)
	})
}

// setupRoutes configures all routes for the server
func (s *Server) setupRoutes() {
	// Health check handlers
	healthHandler := handler.NewHealthHandler(s.db)
	s.router.GET("/health", healthHandler.Health)
	s.router.GET("/ready", healthHandler.Ready)

	// Swagger UI
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// Router returns the gin router for use in http.Server
func (s *Server) Router() *gin.Engine {
	return s.router
}
