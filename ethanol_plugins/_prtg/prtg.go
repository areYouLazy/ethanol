package main

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
	provider    = "prtg_0.1_apitoken"
	description = "get results from a prtg installation through apitoken parameter using APIv1"
	version     = "0.1"
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
	baseURL, _ := url.JoinPath(backend.Endpoint, "api", "table.json")

	queryURL, _ := url.Parse(baseURL)

	// according to the documentation here https://www.paessler.com/manuals/prtg/multiple_object_property_or_status
	// the query should looks something like this:
	// /api/table.json?apitoken=<api-token>&content=devices&filter_name=@sub(<query>)
	// where <api-token> is the API Token and <query> is the query input from the user
	// cannot test right now as I don't have a PRTG installation
	// TODO(areYouLazy) : verify this url format
	values := queryURL.Query()
	values.Add("apitoken", backend.APIToken)
	values.Add("content", "devices")
	values.Add("filter_name", fmt.Sprintf("@sub(%s)", query))
	queryURL.RawQuery = values.Encode()

	request := utils.NewEthanolHTTPClientGETRequest()
	request.URL = queryURL

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
