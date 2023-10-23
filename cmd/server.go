package cmd

import (
	"doodocsProg/api/routes"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	router := gin.Default()
	routes.SetupRoutes(router)

	router.Run(":8080")
}
