package routes

import (
	"doodocsProg/api/handler"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/api/archive/information", handler.GetArchiveInfo)
	router.POST("/api/archive/files", handler.CreateArchieve)
	router.POST("/api/mail/file", handler.SendFileToEmailsHandler)

}
