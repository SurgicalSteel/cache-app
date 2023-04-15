package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SurgicalSteel/cache-app/cache"
	"github.com/SurgicalSteel/cache-app/server"
)

const (
	serverName               = "Cache-App"
	cacheExpireCheckInterval = time.Second
	timeout                  = time.Second * 5
)

func main() {
	coreCache := cache.NewCoreCache(cacheExpireCheckInterval)

	httpServer := server.InitializeServer(coreCache)
	RunApp(httpServer, serverName, coreCache)
}

// RunApp for running and shutting down server gracefully
func RunApp(httpServer *http.Server, serverName string, coreCache *cache.CoreCache) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR][httpServer] ListenAndServe: %s\n", err)
		}
	}()

	<-done
	log.Printf("Stopping %s\n", serverName)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer func() {
		// extra handling here
		coreCache.StopExpireCheck()
		cancel()
	}()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to stop http service %s %+v\n", serverName, err)
	}

	log.Printf("%s stopped successfully\n", serverName)
}
