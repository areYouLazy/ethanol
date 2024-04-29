package main

// this is a plugin, compile it with:
// # go build -buildmode=plugin -o <plugins-folder>/check_mk.so check_mk.go

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	name    = "check_mk_search_plugin"
	label   = "Check_MK"
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

// Search main routine exported from plugin
func (s *searchPlugin) Search(getNewHTTPClient func() *http.Client, getNewHTTPGetRequest func() *http.Request, getNewHTTPPostRequest func() *http.Request, query string, resultsChan chan<- types.SearchResult) {
	backends := []backend{}

	// unmarshal plugin configuration
	err := viper.UnmarshalKey("Plugins.CheckMK", &backends)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error in check_mk unmarshal")
		return
	}

	logrus.WithFields(logrus.Fields{
		"configurations": backends,
	}).Debug("check_mk configurations")

	// search for every backend
	var backendWG sync.WaitGroup

	// add tasks for every backend
	for _, b := range backends {
		backendWG.Add(1)

		go func(bck backend) {
			defer backendWG.Done()
			// use a channel to collect results
			search(getNewHTTPClient(), getNewHTTPGetRequest(), query, bck, resultsChan)
		}(b)
	}

	backendWG.Wait()
}

func search(client *http.Client, request *http.Request, query string, backend backend, resultsChan chan<- types.SearchResult) {
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
	// queryURL, err := url.JoinPath(backend.Endpoint, "/api/1.0/domain-types/host/collections/all", queryString)
	// if err != nil {
	// 	logrus.WithFields(logrus.Fields{
	// 		"error": err.Error(),
	// 	}).Error("error in endpoint url composition")

	// 	return
	// }

	// add content-type header
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// build authorization header
	authHeader := fmt.Sprintf(
		"Bearer %s %s",
		backend.Username,
		backend.Password,
	)

	// add authorization header
	request.Header.Add("Authorization", authHeader)

	// add accept header
	request.Header.Add("Accept", "*/*")

	// format request url
	parsedQueryURL, err := url.Parse(queryURL)
	if err != nil {
		return
	}

	// set request url
	request.URL = parsedQueryURL

	// dump http request
	utils.DumpHTTPRequest(request, "check_mk request")

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

		// send result to results channel
		resultsChan <- result
	}
}

// expose structure as symbol
var Searcher searchPlugin
