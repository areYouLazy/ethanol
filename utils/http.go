package utils

import (
	"crypto/tls"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"time"

	"github.com/areYouLazy/ethanol/proxy"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Using NewEthanolHTTPClient, NewEthanolHTTPClientGETRequest and NewEthanolHTTPClientPOSTRequest
// ensures consistency across various http calls for which ethanol is the client

// NewEthanolHTTPClient returns a new http client that can be used by plugins
func NewEthanolHTTPClient(SkipSSLValidation bool) *http.Client {
	// setup cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil
	}

	// setup transport with proxy (if required)
	transport := http.Transport{
		Proxy: proxy.GetEthanolHTTPClientProxyURL(),
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: SkipSSLValidation,
		},
	}

	// setup client
	client := http.Client{
		Timeout:   5 * time.Second,
		Transport: &transport,
		Jar:       jar,
	}

	// return client
	return &client
}

// NewEthanolHTTPClientGETRequest returns a new get request that can be used by plugins
func NewEthanolHTTPClientGETRequest() *http.Request {
	// setup GET http request
	getRequest, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		return nil
	}

	// setup user-agent header
	getRequest.Header.Add("User-Agent", viper.GetString("ethanol.client.useragent"))

	return getRequest
}

// NewHTTPClientPOSTRequest returns a new post request that can be used by plugins
func NewEthanolHTTPClientPOSTRequest() *http.Request {
	// setup post HTTP request
	postRequest, err := http.NewRequest(http.MethodPost, "", nil)
	if err != nil {
		return nil
	}
	// setup user-agent header
	postRequest.Header.Add("User-Agent", viper.GetString("ethanol.client.useragent"))

	// return requests
	return postRequest
}

// DumpHTTPRequest dumps given http request if -debug is enabled
func DumpHTTPRequest(req *http.Request, description string) {
	// is there's a body?
	hasBody := false
	if req.Method == http.MethodPost {
		hasBody = true
	}

	// dump request
	dump, err := httputil.DumpRequestOut(req, hasBody)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":       err.Error(),
			"description": description,
		}).Fatal("error dumping request")
	}

	if req.Header.Get("Transfer-Encoding") == "chunked" || !hasBody {
		logrus.Debug(description, "\n\n", string(dump))
		return
	}

	logrus.Debug(description, "\n\n", string(dump), "\n\n")
}

// DumpHTTPResponse dumps given http response if -debug is enabled
func DumpHTTPResponse(res *http.Response, description string) {
	dump, err := httputil.DumpResponse(res, true)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":       err.Error(),
			"description": description,
		}).Fatal("error dumping response")
	}

	logrus.Debug(description, "\n\n", string(dump), "\n\n")
}
