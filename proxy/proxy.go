package proxy

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Setup() {
	var proxy_http string

	// check if we need to enable proxy
	if viper.GetBool("ethanol.client.proxy.enable") {
		// check if proxy requires authentication
		if viper.GetBool("ethanol.client.proxy.authenticate") {
			proxy_http = fmt.Sprintf(
				"http://%s:%s@%s:%d",
				viper.GetString("ethanol.client.proxy.username"),
				viper.GetString("ethanol.client.proxy.password"),
				viper.GetString("ethanol.client.proxy.address"),
				viper.GetInt("ethanol.client.proxy.port"),
			)
		} else {
			proxy_http = fmt.Sprintf(
				"http://%s:%d",
				viper.GetString("ethanol.client.proxy.address"),
				viper.GetInt("ethanol.client.proxy.port"),
			)
		}

		// setup proxy environment variables
		os.Setenv("HTTP_PROXY", proxy_http)
		os.Setenv("HTTPS_PROXY", proxy_http)
		logrus.WithFields(logrus.Fields{
			"http_proxy":  proxy_http,
			"https_proxy": proxy_http,
		}).Debug("proxy configuration")
	} else {
		logrus.Debug("no proxy configuration required")
	}
}

func GetHTTPProxyURL() func(*http.Request) (*url.URL, error) {
	var proxy_http string

	// check if we need to enable proxy
	if viper.GetBool("ethanol.client.proxy.enable") {
		// check if proxy requires authentication
		if viper.GetBool("ethanol.client.proxy.authenticate") {
			proxy_http = fmt.Sprintf(
				"http://%s:%s@%s:%d",
				viper.GetString("ethanol.client.proxy.username"),
				viper.GetString("ethanol.client.proxy.password"),
				viper.GetString("ethanol.client.proxy.address"),
				viper.GetInt("ethanol.client.proxy.port"),
			)
		} else {
			proxy_http = fmt.Sprintf(
				"http://%s:%d",
				viper.GetString("ethanol.client.proxy.address"),
				viper.GetInt("ethanol.client.proxy.port"),
			)
		}

		proxyURL, err := url.Parse(proxy_http)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error":     err.Error(),
				"proxy_url": proxy_http,
			}).Error("error parsing proxy url")
		} else {
			logrus.WithFields(logrus.Fields{
				"proxy_string": proxyURL.String(),
			}).Debug("proxy configuration loaded")
		}

		return http.ProxyURL(proxyURL)
	}

	return nil
}
