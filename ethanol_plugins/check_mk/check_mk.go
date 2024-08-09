package main

// this is a plugin, compile it with:
// # go build -buildmode=plugin -o <plugins-folder>/check_mk.so check_mk.go

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/areYouLazy/ethanol/types"
	"github.com/areYouLazy/ethanol/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// define plugin constants
const (
	name        = "check_mk_search_plugin"
	provider    = "check_mk_0.1_username_password"
	description = "This plugins get results from a check_mk installation through username/password authentication"
	label       = "Check_MK"
	raw_label   = "check_mk"
	version     = "0.1"
)

// searchPlugin structure to hold Search routine
type searchPlugin struct{}

// Name exposes plugin name
func (s *searchPlugin) Name() string {
	return name
}

// Provider exposes plugin name
func (s *searchPlugin) Provider() string {
	return provider
}

// Description exposes plugin name
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
	err := viper.UnmarshalKey("plugins.checkmk", &backends)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error in check_mk unmarshal")
		return
	}

	// log
	logrus.WithFields(logrus.Fields{
		"configurations": backends,
	}).Debug("check_mk configurations")

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

	// get site name from endpoint url
	//	[
	//		"0": "http:",
	//		"1": "",
	//		"2", "<host-address>",
	//		"3": "site", <-- Money (That's What I Want!)
	//		"4": "check_mk"
	//	]
	endpointSplit := strings.Split(backend.Endpoint, "/")
	site := endpointSplit[3]

	objectURLFormat := "/index.py?start_url="
	objectQueryFormat := fmt.Sprintf(
		"/%s/%s",
		site,
		"check_mk/view.py?host=%s&view_name=host",
	)

	// search with case-insensitive regexp match (~~) on 'name' and 'address' fields for every host
	// return a match if name or address matches
	queryFormat := "query={\"op\":\"or\",\"expr\":[{\"op\":\"~~\",\"left\":\"name\",\"right\":\"%s\"},{\"op\":\"~~\",\"left\":\"address\",\"right\":\"%s\"}]}&columns=address"
	queryString := fmt.Sprintf(queryFormat, query, query)

	// format request url
	queryURL := fmt.Sprintf(
		"%s/%s?%s",
		backend.Endpoint,
		"/api/1.0/domain-types/host/collections/all",
		queryString,
	)

	// get a new ethanol GET request
	request := utils.NewEthanolHTTPClientGETRequest()

	// add content-type header
	// request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// build authorization header
	authHeader := fmt.Sprintf(
		"Bearer %s %s",
		backend.Username,
		backend.Password,
	)

	// add authorization header
	request.Header.Add("Authorization", authHeader)

	// add accept header
	request.Header.Add("Accept", "application/json")

	// format request url
	parsedQueryURL, err := url.Parse(queryURL)
	if err != nil {
		return
	}

	// set request url
	request.URL = parsedQueryURL

	// dump http request
	utils.DumpHTTPRequest(request, "check_mk request")

	// get a new ethanol HTTP client
	client := utils.NewEthanolHTTPClient(backend.InsecureSkipSSLVerify)

	// do request
	res, err := client.Do(request)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error in check_mk hostsearch")
		return
	}
	// close response body on exit
	defer res.Body.Close()

	// dump http response
	utils.DumpHTTPResponse(res, "check_mk response")

	// decode response
	json.NewDecoder(res.Body).Decode(&r)

	// extract values from response
	for _, item := range r.Value {
		result := types.SearchResult{}
		result["id"] = uuid.New().String()
		result["name"] = item.Extensions.Name
		result["address"] = item.Extensions.Address
		result["site"] = site
		result["url"] = fmt.Sprintf("%s%s%s", backend.Endpoint, objectURLFormat, url.QueryEscape(fmt.Sprintf(objectQueryFormat, item.Extensions.Name)))
		result["source"] = name
		result["label"] = label
		result["raw_label"] = raw_label

		// send result to results channel
		resultsChan <- result
	}
}

// expose structure as symbol
var Searcher searchPlugin
