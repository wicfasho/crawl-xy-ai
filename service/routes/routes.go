package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wicfasho/crawl-xy/handlers/ask"
	// "github.com/wicfasho/crawl-xy/scrapper"
)

// Configure all routes for the API
func SetupRoutes(r *gin.RouterGroup) {
	// @dev Only Test or Admin
	apiAsk := r.Group("/ask")
	{
		apiAsk.POST("/", ask.Send)
	}
	// r.GET("/scrapper", scrapper.Start)
}
