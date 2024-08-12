package main

// https://www.paessler.com/support/prtg/api/v2/overview/index.html

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"github.com/areYouLazy/ethanol/types"
	"github.com/areYouLazy/ethanol/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	name        = "prtg_search_plugin"
	provider    = "prtg_0.2_apitoken"
	description = "get results from a prtg installation through apitoken parameter using APIv2"
	version     = "0.2"
)

// searchPlugin structure to expose plugin methods
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

func (s *searchPlugin) Search(query string, results chan<- types.SearchResult) {
	// define a list of backends
	backends := []backend{}

	// load backends from configuration
	err := viper.UnmarshalKey("plugins.prtg", &backends)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error in check_mk unmarshal")
		return
	}

	// log
	logrus.WithFields(logrus.Fields{
		"configurations": backends,
	}).Debug("prtg configurations")

	var backendWG sync.WaitGroup

	for _, b := range backends {
		backendWG.Add(1)

		go func(bck backend) {
			defer backendWG.Done()
			search(query, bck, results)
		}(b)
	}

	backendWG.Wait()
}

func search(query string, backend backend, results chan<- types.SearchResult) {
	baseURL, err := url.JoinPath(backend.Endpoint, "api", "v2", "objects")
	if err != nil {
		logrus.Panic(err)
	}

	queryURL, err := url.Parse(baseURL)
	if err != nil {
		logrus.Panic(err)
	}

	values := queryURL.Query()
	values.Add("filter", url.QueryEscape(fmt.Sprintf("name matches %s", query)))
	queryURL.RawQuery = values.Encode()

	request := utils.NewEthanolHTTPClientPOSTRequest()
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", backend.APIToken))

	// dump request
	utils.DumpHTTPRequest(request, "prtg request")

	// get a new ethanol HTTP client
	client := utils.NewEthanolHTTPClient(backend.InsecureSkipSSLVerify)

	// do request
	res, err := client.Do(request)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error in prtg request")
		return
	}
	defer res.Body.Close()

	// dump response
	utils.DumpHTTPResponse(res, "prtg response")

	// parse response
	var r rawResponse

	json.NewDecoder(res.Body).Decode(&r)

	// iterate r and generate results
}

var Searcher searchPlugin
