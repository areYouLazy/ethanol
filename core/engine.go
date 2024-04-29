package core

import (
	"github.com/areYouLazy/ethanol/plugins"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func getCoreEngine() *gin.Engine {
	// setup Release mode
	gin.SetMode(gin.ReleaseMode)

	// instantiate a new gin-gonic router
	r := gin.New()

	// set engine customizations
	r.ForwardedByClientIP = viper.GetBool("Ethanol.Server.ForwardedByClientIP")
	r.RemoteIPHeaders = viper.GetStringSlice("Ethanol.Server.RemoteIPHeaders")
	r.SetTrustedProxies(viper.GetStringSlice("Ethanol.Server.TrustedProxies"))

	// load htmx templates
	r.LoadHTMLGlob("ui/templates/*")

	// setup middlewares
	r.Use(gin.Recovery())
	r.Use(logMiddleware())
	r.Use(signatureMiddleware())

	// check if we need to be a secure webserver
	if viper.GetBool("Ethanol.Server.Security.Secure") {
		r.Use(securityMiddleware())
	}

	// init websocket hub
	webSocketHub = InitWebSocketHub()

	// setup handlers
	// since we use htmx we need to split JSON API from HTML API
	// /api/v1/ui/.... for HTML API
	// /api/v1/.... for JSON API
	r.GET("/status", getStatusHandler)
	r.GET("/status/json", getStatusJSONHandler)
	// r.GET("/ws", getWebSocketHandler)
	r.GET("/api/v1/search", getAPIV1SearchHandler)
	r.GET("/api/v1/ui/search", getAPIV1UISearchHandler)

	// serve static content
	r.Static("/ui/css", "./ui/css")
	r.Static("/ui/img", "./ui/img")

	// serve htmx ui
	r.GET("/ui/:template", uixHandler)

	// log engine configuration
	logrus.WithFields(logrus.Fields{
		"ethanol_forwarded_by_client_ip":    r.ForwardedByClientIP,
		"ethanol_remote_ip_headers":         r.RemoteIPHeaders,
		"ethanol_trusted_proxies":           viper.GetStringSlice("Ethanol.Server.TrustedProxies"),
		"ethanol_max_multipart_memory":      r.MaxMultipartMemory,
		"ethanol_redirect_fixed_path":       r.RedirectFixedPath,
		"ethanol_remove_extra_slash":        r.RemoveExtraSlash,
		"ethanol_trusted_platform":          r.TrustedPlatform,
		"ethanol_unescape_path_values":      r.UnescapePathValues,
		"ethanol_use_H2C":                   r.UseH2C,
		"ethanol_use_raw_path":              r.UseRawPath,
		"ethanol_base_path":                 r.BasePath(),
		"ethanol_handle_method_not_allowed": r.HandleMethodNotAllowed,
	}).Debug("gin engine configuration")

	// log engine routes
	for _, v := range r.Routes() {
		logrus.WithFields(logrus.Fields{
			"route_path":   v.Path,
			"route_mathod": v.Method,
		}).Debug("gin engine route loaded")
	}

	// init plugins manager
	plugins.Init()

	// return the structured gin engine
	return r
}
