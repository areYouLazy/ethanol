package core

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// logMiddleware gin middleware to log informations about web requests
func logMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// setup logrus logging format
		logrus.Infof(
			"%s %s %s %s %d %d %s %s",
			c.RemoteIP(),
			c.Request.Method,
			c.Request.URL,
			c.Request.Proto,
			c.Writer.Status(),
			c.Request.ContentLength,
			c.Request.Referer(),
			c.Request.UserAgent(),
		)

		// call next handler
		c.Next()
	}
}

// securityMiddleware gin middleware to enhance http security
func securityMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// setup security headers
		c.Writer.Header().Add("Cache-Control", viper.GetString("Ethanol.Server.Security.Headers.CacheCOntrol"))
		c.Writer.Header().Add("X-XSS-Protection", viper.GetString("Ethanol.Server.Security.Headers.XSSProtection"))
		c.Writer.Header().Add("X-Frame-Options", viper.GetString("Ethanol.Server.Security.Headers.XFrameOptions"))
		c.Writer.Header().Add("X-Content-Type-Options", viper.GetString("Ethanol.Server.Security.Headers.XContentTypeOptions"))

		c.Writer.Header().Add("Strict-Transport-Security", viper.GetString("Ethanol.Server.Security.Headers.HSTS"))
		c.Writer.Header().Add("Content-Security-Policy", viper.GetString("Ethanol.Server.Security.Headers.CSP"))

		// call next handler
		c.Next()
	}
}

// signatureMiddleware gin middleware to add signature to http traffic
func signatureMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Server", viper.GetString("Ethanol.Server.Signature.Server"))
		c.Writer.Header().Add("X-Powered-By", viper.GetString("Ethanol.Server.Signature.XPoweredBy"))

		// call next  handler
		c.Next()
	}
}
