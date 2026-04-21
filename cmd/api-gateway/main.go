package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) { c.String(200, "ok") })
	v1 := r.Group("/api/v1")
	{
		v1.GET("/gpus", ListGPUs)
		v1.GET("/gpus/:id/telemetry", GetTelemetry)
	}
	r.StaticFile("/api/openapi.yaml", "api/openapi.yaml")
	r.Run(":8080")
}

func ListGPUs(c *gin.Context) {
	c.JSON(http.StatusOK, []string{"GPU-bc7a12ab-4998-fdc5-0785-2678a929a142"})
}

func GetTelemetry(c *gin.Context) {
	c.JSON(http.StatusOK, []interface{}{})
}
