package server

import (
	"uber_fx_init_folder_structure/internal/mw"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
)

func v1Routes(router *gin.RouterGroup, awsSession *session.Session, o *Options) {
	r := router.Group("/v1/")
	// middlewares
	r.Use(mw.ErrorHandlerX(o.Log))
	r.PUT("/user_registration", o.UserHandler.CreateUser)
	r.GET("/ws/user_chat", o.UserHandler.ChatWithBot)
	r.POST("/upload_photos", mw.AWSSessionAttach(awsSession), o.UserHandler.UserUploadPhoto)
}
