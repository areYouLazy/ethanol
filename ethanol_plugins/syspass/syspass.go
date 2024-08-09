package main

// this is a plugin, compile it with:
// # go build -buildmode=plugin -o <plugins-folder>/syspass.so syspass.go

import (
	"bytes"
	"encoding/json"
	"io"
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
	name        = "syspass_search_plugin"
	provider    = "syspass_3.2.11_apikey"
	description = "This plugins get results from a syspass installation through API Key"
	label       = "SysPass"
	raw_label   = "syspass"
	version     = "0.1"
)

// searchPlugin structure to hold Search routine
type searchPlugin struct{}

// Name exposes plugin name
func (s *searchPlugin) Name() string {
	return name
}

// Provider exposes plugin provider
func (s *searchPlugin) Provider() string {
	return provider
}

// Description exposes plugin description
func (s *searchPlugin) Description() string {
	return description
}

// Version exposes plugin version
func (s *searchPlugin) Version() string {
	return version
}

// Search prepare environment for the actual query
func (s *searchPlugin) Search(query string, resultsChan chan<- types.SearchResult) {
	// define a list of backends
	backends := []backend{}

	// load backends from configuration
	err := viper.UnmarshalKey("plugins.syspass", &backends)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error in syspass unmarshal")
		return
	}

	// log
	logrus.WithFields(logrus.Fields{
		"configurations": backends,
	}).Debug("syspass configurations")

	// define a waitgroup for backend routines
	var backendWG sync.WaitGroup

	// iterate backends
	for _, b := range backends {
		backendWG.Add(1)

		// start search in a goroutine
		go func(bck backend) {
			defer backendWG.Done()

			search(query, bck, resultsChan)
		}(b)
	}

	// wait for jobs to be done
	backendWG.Wait()
}

// search actual query to backend
func search(query string, backend backend, resultsChan chan<- types.SearchResult) {
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

	// get a new ethanol HTTP post request
	request := utils.NewEthanolHTTPClientPOSTRequest()

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

	// get a new ethanol HTTP client
	client := utils.NewEthanolHTTPClient(backend.InsecureSkipSSLVerify)

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
		result["url"] = item.URL
		result["login"] = item.Login
		result["source"] = name
		result["label"] = label
		result["raw_label"] = raw_label
		result["provider"] = provider

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
