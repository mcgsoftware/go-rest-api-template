package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"go-rest-api-template/internal/config"
	"go-rest-api-template/internal/database"
	"go-rest-api-template/internal/server"
)

var (
	cfgPort   int
	cfgDebug  bool
	cfgDBURL  string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "server",
	Short: "Go REST API Template Server",
	Long:  `A production-ready REST API server built with Gin, PostgreSQL, and modern Go tooling.`,
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Long:  `Start the HTTP server with the specified configuration.`,
	RunE:  runServe,
}

func init() {
	// Add serve command to root
	rootCmd.AddCommand(serveCmd)

	// Define flags for serve command
	serveCmd.Flags().IntVar(&cfgPort, "port", 8080, "Server port")
	serveCmd.Flags().BoolVar(&cfgDebug, "debug", false, "Enable debug mode")
	serveCmd.Flags().StringVar(&cfgDBURL, "db", "", "Database connection URL")

	// Bind flags to viper
	viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
	viper.BindPFlag("debug", serveCmd.Flags().Lookup("debug"))
	viper.BindPFlag("database_url", serveCmd.Flags().Lookup("db"))
}

func runServe(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Setup logger
	logLevel := slog.LevelInfo
	if cfg.Debug {
		logLevel = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	logger.Info("starting server",
		slog.Int("port", cfg.Port),
		slog.Bool("debug", cfg.Debug),
	)

	// Connect to database
	ctx := context.Background()
	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()
	logger.Info("connected to database")

	// Create server
	srv := server.New(db, logger, cfg.Debug)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      srv.Router(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info("server listening", slog.String("addr", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	// Give outstanding requests 30 seconds to complete
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	logger.Info("server stopped")
	return nil
}
