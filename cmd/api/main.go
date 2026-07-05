package main

import (
	"fmt"
	"os"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/pedro-h-castoldi/go-ddd-boilerplate/internal/config"
	"github.com/pedro-h-castoldi/go-ddd-boilerplate/internal/database"
	"github.com/pedro-h-castoldi/go-ddd-boilerplate/internal/logger"
	"github.com/pedro-h-castoldi/go-ddd-boilerplate/internal/middleware"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "application error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// 1. Config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// 2. Logger
	log, err := logger.New(cfg.Log)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	defer log.Sync()

	// opcional: global
	zap.ReplaceGlobals(log)

	log.Info("starting application")

	// 3. Database
	db := database.NewPool(cfg.Databases)
	if err := db.OpenConnections(cfg.Databases); err != nil {
		log.Error("failed to open connections to database", zap.Error(err))
		return fmt.Errorf("open connections to database: %w", err)
	}
	defer db.Close()

	log.Info("database connected")

	// 4. Gin
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	group := errgroup.Group{}
	group.Go(func() error {
		return endless.ListenAndServe(cfg.HTTP.InternalPort, internalRouter(log))
	})
	group.Go(func() error {
		return endless.ListenAndServe(cfg.HTTP.ExternalPort, externalRouter(log))
	})

	//return router.Run(fmt.Sprintf(":%d", cfg.HTTP.Port))
	return group.Wait()
}

func internalRouter(log *zap.Logger) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.GinLogger(log))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	return router
}

func externalRouter(log *zap.Logger) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.GinLogger(log))

	router.GET("/health", func(c *gin.Context) {
		zap.L().Info("health check")
		c.JSON(200, gin.H{"status": "ok"})
	})
	return router
}
