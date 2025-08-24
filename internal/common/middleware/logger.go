package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StructuredLog(l *logrus.Entry) gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now()
		c.Next()
		elapsed := time.Since(now)
		l.WithFields(logrus.Fields{
			"elapsed":     elapsed.Milliseconds(),
			"request_url": c.Request.RequestURI,
			"client_ip":   c.ClientIP(),
		}).Info("Request out")
	}
}
