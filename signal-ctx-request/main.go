package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.GET("/long-task", func(c *gin.Context) {
		ctx := c.Request.Context()
		log.Println("Handler started, entering select")

		go func() {
			<-ctx.Done()
			log.Println("Context done in goroutine:", ctx.Err())
		}()

		select {
		case <-time.After(10 * time.Hour):
			log.Println("Timer completed")
			c.JSON(200, gin.H{"status": "completed"})
		case <-ctx.Done():
			log.Println("Request canceled in select:", ctx.Err())
			return
		}
		log.Println("Handler finished")
	})

	// Setup signal catching
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Create a context you can cancel
	requestCtx, requestCancel := context.WithCancel(context.Background())

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
		BaseContext: func(net.Listener) context.Context {
			return requestCtx
		},
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// ... later when terminating ...
	<-quit

	// cancel the base context to signal all requests
	requestCancel()

	// then shutdown the server giving time to request handlers to wrap up within a grace period
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server forced to shutdown:", err)
	}

	time.Sleep(2 * time.Second)
	log.Println("Server exiting")
}
