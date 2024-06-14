// logger.go
package mylogger

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

var Logger *log.Logger

func init() {
	Logger = log.New()
	Logger.SetFormatter(&log.JSONFormatter{})
	Logger.SetOutput(os.Stdout)
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		entry := Logger.WithFields(log.Fields{
			"method": c.Request.Method,
			"uri":    c.Request.RequestURI,
			"status": c.Writer.Status(),
			"time":   latency,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else {
			entry.Info()
		}
	}
}
