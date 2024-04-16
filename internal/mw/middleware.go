// Package mw is user Middleware package
package mw

import (
	"net/http"
	"uber_fx_init_folder_structure/er"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ErrorHandlerX(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			err := c.Errors.Last()
			if err == nil {
				// no errors, abort with success
				return
			}

			e := er.From(err.Err)

			if !e.NOP {
				sentry.CaptureException(e)
			}

			httpStatus := http.StatusInternalServerError
			if e.Status > 0 {
				httpStatus = e.Status
			}

			c.JSON(httpStatus, e)
		}()

		c.Next()
	}
}

func AWSSessionAttach(sess *session.Session) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("sess", sess)
		c.Next()
	}
}
