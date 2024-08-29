package engine

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"scylla/pkg/helper"
	"syscall"
	"time"
)

func StartServerWithGracefulShutdown(server *http.Server) {
	// Channel to listen for OS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Run the server in a goroutine so that it doesn't block
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			helper.ErrorPanic(err)
		}
	}()

	// Wait for an interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		helper.ErrorPanic(err)
	}

	log.Println("Server exiting gracefully")
}
