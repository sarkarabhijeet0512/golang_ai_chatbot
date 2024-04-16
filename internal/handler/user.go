package handler

import (
	"context"
	"fmt"
	"log"
	"uber_fx_init_folder_structure/er"
	model "uber_fx_init_folder_structure/utils/models"

	"net/http"
	"uber_fx_init_folder_structure/pkg/cache"
	"uber_fx_init_folder_structure/pkg/user"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	log          *logrus.Logger
	userService  *user.Service
	cacheService *cache.Service
}

var messageChan = make(chan string)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func newUserHandler(
	log *logrus.Logger,
	userService *user.Service,
	cacheService *cache.Service,
) *UserHandler {
	return &UserHandler{
		log,
		userService,
		cacheService,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var (
		err  error
		res  = model.GenericRes{}
		req  = &user.User{}
		dCtx = context.Background()
	)
	defer func() {
		if err != nil {
			c.Error(err)
			h.log.WithField("span", res).Warn(err.Error())
			return
		}
	}()
	if err = c.ShouldBind(&req); err != nil {
		err = er.New(err, er.UncaughtException).SetStatus(http.StatusUnprocessableEntity)
		return
	}
	err = h.userService.UpsertUserRegistration(dCtx, req)
	if err != nil {
		err = er.New(err, er.UncaughtException).SetStatus(http.StatusUnprocessableEntity)
		return
	}
	res.Message = "Registration Sucessfully Done"
	res.Success = true
	res.Data = req
	c.JSON(http.StatusOK, res)
}
func (h *UserHandler) ChatWithBot(c *gin.Context) {
	var (
		err  error
		res  = model.GenericRes{}
		req  = &model.BotReq{}
		dCtx = context.Background()
	)
	defer func() {
		if err != nil {
			c.Error(err)
			h.log.WithField("span", res).Warn(err.Error())
			return
		}
	}()
	if err = c.ShouldBind(&req); err != nil {
		err = er.New(err, er.UncaughtException).SetStatus(http.StatusUnprocessableEntity)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	for {
		go func() {
			for msg := range messageChan {
				if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
					log.Printf("Error writing message to WebSocket: %v", err)
					break
				}
			}
		}()
		// Read message from WebSocket client
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from WebSocket: %v", err)
			break
		}
		// Process the message using OpenAI API
		response, err := h.userService.ProcessMessage(dCtx, string(msg))
		if err != nil {
			log.Printf("Error processing message: %v", err)
			continue
		}

		if err := conn.WriteMessage(websocket.TextMessage, []byte(response)); err != nil {
			log.Printf("Error writing message to WebSocket: %v", err)
			break
		}
	}
}

func (h *UserHandler) UserUploadPhoto(c *gin.Context) {
	var (
		err  error
		res  = model.GenericRes{}
		req  = user.UserImages{}
		dCtx = context.Background()
	)
	defer func() {
		if err != nil {
			c.Error(err)
			h.log.WithField("span", res).Warn(err.Error())
			return
		}
	}()

	username := ""
	err = h.cacheService.Get("username", &username)
	if err != nil {
		go h.SendMessageToSocket("Please enter username in the chatbox to proceed!!")
		err = er.New(err, er.UserNotFound).SetStatus(http.StatusBadRequest)
		return
	}

	err = c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		err = er.New(err, er.UncaughtException).SetStatus(http.StatusBadRequest)
		return
	}
	sess := c.MustGet("sess").(*session.Session)
	contentType := "image/jpeg"
	form, err := c.MultipartForm()
	if err != nil {
		err = er.New(err, er.UncaughtException).SetStatus(http.StatusBadRequest)
		return
	}
	userDetails, err := h.userService.FetchUserByUsername(dCtx, username)
	if err != nil {
		err = er.New(err, er.UncaughtException).SetStatus(http.StatusBadRequest)
		return
	}
	req = user.UserImages{
		UserID: userDetails.ID,
	}
	files := form.File["images"]
	for _, file := range files {
		// Open the uploaded file
		f, err := file.Open()
		if err != nil {
			err = er.New(err, er.UncaughtException).SetStatus(http.StatusBadRequest)
			return
		}
		defer f.Close()
		// Generate a unique filename
		filename := fmt.Sprintf("%s-%s", file.Filename, uuid.New())
		err = h.userService.UserUploadPhoto(dCtx, req, f, filename, contentType, sess)
		if err != nil {
			err = er.New(err, er.UncaughtException).SetStatus(http.StatusUnprocessableEntity)
			return
		}
		fmt.Printf("Uploaded file: %s\n", filename)
		go h.SendMessageToSocket("uploaded successfully", filename)
	}

	res.Message = "uploaded successfully"
	res.Success = true
	res.Data = req
	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) SendMessageToSocket(message string, file ...string) {
	messageChan <- fmt.Sprint(file, " ", message)
}
