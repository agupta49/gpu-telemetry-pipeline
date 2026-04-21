package main
import ("log"; "github.com/gin-gonic/gin")
func main() {
	r := gin.Default()
	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.GET("/api/v1/gpus", func(c *gin.Context) { c.JSON(200, []string{}) })
	log.Println("api-gateway: starting on :8080")
	r.Run(":8080")
}
