package server

import (
	"net/http"
	"time"

	"github.com/SurgicalSteel/cache-app/cache"
	"github.com/SurgicalSteel/cache-app/controller"

	"github.com/gin-gonic/gin"
)

var (
	applicationRoute *route
)

const (
	ginMode = "release"
)

// InitializeServer for running service (initializes an HTTP server)
func InitializeServer(coreCache *cache.CoreCache) *http.Server {

	httpServer := &http.Server{
		Addr:         ":8080",
		Handler:      InitializeRouter().Routing(coreCache),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return httpServer
}

// ServiceInject is where we inject dependencies on all levels (repository, service, controller) and return the controller
func ServiceInject(coreCache *cache.CoreCache) *controller.CacheAppController {
	return &controller.CacheAppController{Cache: coreCache}

}

// IRouter is interface for routing
type IRouter interface {
	Routing(coreCache *cache.CoreCache) *gin.Engine
}

// route is a struct for router
type route struct {
	IRouter
}

// Router is a func to initialize route struct
func InitializeRouter() IRouter {
	if applicationRoute == nil {
		applicationRoute = &route{}
	}
	return applicationRoute
}

// Routing is a function for injecting dependencies and define http routing
func (route *route) Routing(coreCache *cache.CoreCache) *gin.Engine {
	gin.SetMode(ginMode)
	defaultEngine := gin.Default()

	cacheAppController := ServiceInject(coreCache)

	apiV1 := defaultEngine.Group("/")
	{
		apiV1.POST("/:key", cacheAppController.HandleInsert)
		apiV1.GET("/:key", cacheAppController.HandleGetByKey)
	}
	return defaultEngine
}
