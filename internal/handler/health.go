package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db *pgxpool.Pool
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(db *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{db: db}
}

// HealthResponse represents the response for health endpoints
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// ReadyResponse represents the response for the readiness endpoint
type ReadyResponse struct {
	Status    string            `json:"status"`
	Timestamp string            `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

// Health handles GET /health - liveness check
// Always returns 200 if the server is running
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, HealthResponse{
		Status:    "ok",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

// Ready handles GET /ready - readiness check with database connectivity
// Returns 200 if all dependencies are healthy, 503 otherwise
func (h *HealthHandler) Ready(c *gin.Context) {
	checks := make(map[string]string)
	status := http.StatusOK
	overallStatus := "ok"

	// Check database connectivity using a simple query
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var result int
	err := h.db.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		checks["database"] = "unhealthy: " + err.Error()
		status = http.StatusServiceUnavailable
		overallStatus = "unhealthy"
	} else {
		checks["database"] = "healthy"
	}

	c.JSON(status, ReadyResponse{
		Status:    overallStatus,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    checks,
	})
}
