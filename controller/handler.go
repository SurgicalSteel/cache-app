package controller

import (
	"io"
	"net/http"

	"github.com/SurgicalSteel/cache-app/cache"

	"github.com/gin-gonic/gin"
)

// CacheAppController is the struct that defines app controller and its dependency
type CacheAppController struct {
	Cache *cache.CoreCache
}

// HandleInsert is the handler func for inserting a key into the cache
func (cac *CacheAppController) HandleInsert(c *gin.Context) {
	rawKey := c.Param("key")

	value := ""

	rawValue, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	}

	value = string(rawValue)

	cac.Cache.Set(rawKey, value)

	c.Header("Content-Type", "text")
	c.JSON(http.StatusOK, "success")
}

// HandleGetByKey is the handler func for getting the value of a key from the cache. If the key doesn't exist, then it will return HTTP 404
func (cac *CacheAppController) HandleGetByKey(c *gin.Context) {
	rawKey := c.Param("key")

	value, err := cac.Cache.Get(rawKey)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	}

	c.JSON(http.StatusOK, value)

}
