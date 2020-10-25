package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/JackFazackerley/complete-packs/internal/database"

	"github.com/JackFazackerley/complete-packs/pkg/cache"

	"github.com/JackFazackerley/complete-packs/internal/config"
	orderController "github.com/JackFazackerley/complete-packs/internal/controllers/order"
	packsController "github.com/JackFazackerley/complete-packs/internal/controllers/packs"
	"github.com/JackFazackerley/complete-packs/internal/database/sqlite"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := config.New()
	if err != nil {
		log.WithError(err).Fatal("loading config")
	}

	db, err := database.OpenDatabase(config)
	if err != nil {
		log.WithError(err).Fatal("connecting to db")
	}
	defer db.Close()

	dbInterface := sqlite.New(db)

	sizesCache := cache.New(dbInterface)
	if err := sizesCache.Load(); err != nil {
		log.WithError(err).Fatal("loading sizes cache")
	}

	packs := packsController.New(dbInterface, sizesCache)
	order := orderController.New(sizesCache)

	router := gin.Default()
	router.Use(cors.Default())
	orderGroup := router.Group("/order")
	orderGroup.POST("/best", order.Best)
	orderGroup.POST("/fast", order.Fast)

	packsGroup := router.Group("/packs")
	packsGroup.GET("/read", packs.Read)
	packsGroup.POST("/write", packs.Write)
	packsGroup.DELETE("/delete", packs.Delete)

	srv := &http.Server{
		Addr:    ":9090",
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
