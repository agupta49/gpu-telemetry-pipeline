package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/agupta49/gpu-telemetry-pipeline/internal/api"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "github.com/agupta49/gpu-telemetry-pipeline/api" // docs
)

// @title GPU Telemetry API
// @version 1.0
// @description API for querying GPU telemetry data
// @BasePath /api/v1
func main() {
	dbDSN := flag.String("db", "", "Postgres DSN")
	flag.Parse()

	repo, err := api.NewRepo(*dbDSN)
	if err != nil {
		log.Fatalf("db connect failed: %v", err)
	}
	defer repo.Close()

	r := gin.Default()
	
	// @Summary Health check
	// @Success 200 {object} map[string]string
	// @Router /healthz [get]
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	v1 := r.Group("/api/v1")
	{
		// @Summary List all GPUs
		// @Produce json
		// @Success 200 {array} string
		// @Router /gpus [get]
		v1.GET("/gpus", func(c *gin.Context) {
			gpus, err := repo.ListGPUs(c.Request.Context())
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, gpus)
		})

		// @Summary Get telemetry for GPU
		// @Param id path string true "GPU ID"
		// @Param start_time query string false "Start time RFC3339" Format(date-time)
		// @Param end_time query string false "End time RFC3339" Format(date-time)
		// @Produce json
		// @Success 200 {array} api.TelemetryPoint
		// @Router /gpus/{id}/telemetry [get]
		v1.GET("/gpus/:id/telemetry", func(c *gin.Context) {
			id := c.Param("id")
			var start, end *time.Time
			if s := c.Query("start_time"); s != "" {
				t, err := time.Parse(time.RFC3339, s)
				if err != nil {
					c.JSON(400, gin.H{"error": "invalid start_time"})
					return
				}
				start = &t
			}
			if e := c.Query("end_time"); e != "" {
				t, err := time.Parse(time.RFC3339, e)
				if err != nil {
					c.JSON(400, gin.H{"error": "invalid end_time"})
					return
				}
				end = &t
			}
			points, err := repo.GetTelemetry(c.Request.Context(), id, start, end)
			if err != nil {
				c.JSON(500, gin.H{"error": err.Error()})
				return
			}
			c.JSON(200, points)
		})
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	log.Println("api-gateway: starting on :8080")
	r.Run(":8080")
}
