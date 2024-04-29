package main

// this is a plugin, compile it with:
// # go build -buildmode=plugin -o <plugins-folder>/syspass.so syspass.go

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/areYouLazy/ethanol/types"
	"github.com/areYouLazy/ethanol/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// define plugin constants
const (
	name    = "syspass_search_plugin"
	label   = "SysPass"
	version = "0.1"
)

// searchPlugin structure to hold Search routine
type searchPlugin struct{}

// Name exposes plugin name
func (s *searchPlugin) Name() string {
	return name
}

// Version exposes plugin version
func (s *searchPlugin) Version() string {
	return version
}

func (s *searchPlugin) Search(getNewHTTPClient func() *http.Client, getNewHTTPGetRequest func() *http.Request, getNewHTTPPostRequest func() *http.Request, query string, resultsChan chan<- types.SearchResult) {
	backends := []backend{}

	err := viper.UnmarshalKey("plugins.syspass", &backends)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error in syspass unmarshal")
		return
	}

	logrus.WithFields(logrus.Fields{
		"configurations": backends,
	}).Debug("syspass configurations")

	var backendWG sync.WaitGroup

	for _, b := range backends {
		backendWG.Add(1)

		go func(bck backend) {
			defer backendWG.Done()

			search(getNewHTTPClient(), getNewHTTPPostRequest(), query, bck, resultsChan)
		}(b)
	}

	backendWG.Wait()
}

// accountSearch actual query to backend
func search(client *http.Client, request *http.Request, query string, backend backend, resultsChan chan<- types.SearchResult) {
	var r rawResponse

	queryURL, err := url.JoinPath(backend.Endpoint, "api.php")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error in endpoint url construction")

		return
	}

	// prepare request body for account/search
	body, err := json.Marshal(queryBody{
		JSONRPC: "2.0",
		Method:  "account/search",
		Params: queryBodyParams{
			AuthToken: backend.APIKey,
			TokenPass: backend.APIKeyPassPhrase,
			Text:      query,
			Count:     backend.Count,
		},
		ID: 1,
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error in syspass request body construction")

		return
	}

	// add content-type header
	request.Header.Del("Content-Type")
	request.Header.Add("Content-Type", "application/json")

	// format request url
	parsedQueryURL, err := url.Parse(queryURL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
			"url":   queryURL,
		}).Error("error in syspass url parse")
		return
	}

	request.URL = parsedQueryURL

	// add body to request
	request.Body = io.NopCloser(bytes.NewReader(body))

	// dump request
	utils.DumpHTTPRequest(request, "syspass request for account/search")

	// do request
	res, err := client.Do(request)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error in syspass apisearch")
		return
	}
	defer res.Body.Close()

	// dump response
	utils.DumpHTTPResponse(res, "syspass response for account/search")

	// iterate results to get password for every returned account
	json.NewDecoder(res.Body).Decode(&r)

	for _, item := range r.Result.Result {
		result := make(map[string]interface{}, 0)
		result["id"] = uuid.New().String()
		result["title"] = item.Name
		result["description"] = item.Notes
		result["client"] = item.ClientName
		result["category"] = item.CategoryName
		result["login"] = item.Login
		result["source"] = name
		result["label"] = label

		// get password for every item in response
		var p rawPasswordResponse

		// prepare request body for account/viewPass
		getPasswordBody, err := json.Marshal(queryBody{
			JSONRPC: "2.0",
			Method:  "account/viewPass",
			Params: queryBodyParams{
				AuthToken: backend.APIKey,
				TokenPass: backend.APIKeyPassPhrase,
				ID:        item.ID,
				Count:     1,
			},
			ID: 1,
		})
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("error in syspass body construction")

			return
		}

		logrus.Debug("getPasswordBody", string(getPasswordBody))

		// add body to request
		request.Body = io.NopCloser(bytes.NewReader(getPasswordBody))

		// dump request
		utils.DumpHTTPRequest(request, "syspass request for account/viewPass")

		// do request
		res, err := client.Do(request)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("error in syspass account/viewPass")
			return
		}
		defer res.Body.Close()

		// dump response
		utils.DumpHTTPResponse(res, "syspass response for account/viewPass")

		// add password to result object
		json.NewDecoder(res.Body).Decode(&p)
		result["password"] = p.Result.Result.Password
		// e.Password = p.Result.Result.Password

		resultsChan <- result
	}
}

// exposes structure as symbol
var Searcher searchPlugin
