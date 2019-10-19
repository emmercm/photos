package http

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/emmercm/photos/internal/pkg/store/sqlite3"
	"github.com/emmercm/photos/internal/pkg/transport/http/handler"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func router(db *gorm.DB) *gin.Engine {
	router := gin.New()

	router.Use(static.Serve("/", static.LocalFile("./web/build", false)))

	apiV1 := router.Group("/api/v1")
	{
		apiV1.GET("/files", handler.GetFiles(db))
		apiV1.GET("/file/:id", handler.GetFile(db))
		apiV1.GET("/file/:id/image", handler.GetFileImage(db))

		apiV1.GET("/albums", handler.GetAlbums(db))
		apiV1.GET("/album/:id", handler.GetAlbum(db))
		apiV1.GET("/album/:id/files", handler.GetAlbumFiles(db))
	}

	// gin.DisableConsoleColor()
	// router.Use(gin.Logger())

	router.Use(gin.Recovery())

	return router
}

// Start starts an HTTP server and waits for OS exit signals
func Start() {
	db, err := sqlite3.Connection()
	if err != nil {
		panic(err)
	}

	router := router(db)

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error(err)
	}
	<-ctx.Done()
}
