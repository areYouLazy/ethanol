package main

//
// Should we use https://github.com/andygrunwald/go-jira ?
//

import (
	"encoding/json"
	"fmt"
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
	name        = "jira_search_plugin"
	description = "get results from a jira installation through username/password authentication"
	provider    = "jira_1.0_username_password"
	label       = "Jira"
	raw_label   = "jira"
	version     = "0.1"
)

// searchPlugin structure to hold Search routine
type searchPlugin struct{}

// Name exposes plugin name
func (s *searchPlugin) Name() string {
	return name
}

func (s *searchPlugin) Provider() string {
	return provider
}

func (s *searchPlugin) Description() string {
	return description
}

// Version exposes plugin version
func (s *searchPlugin) Version() string {
	return version
}

// Search main routine exported from plugin
func (s *searchPlugin) Search(query string, resultsChan chan<- types.SearchResult) {
	backends := []backend{}

	err := viper.UnmarshalKey("plugins.jira", &backends)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error in jira unmarshal")
		return
	}

	logrus.WithFields(logrus.Fields{
		"configurations": backends,
	}).Debug("jira configurations")

	// do thigs

	var backendWG sync.WaitGroup

	for _, b := range backends {
		backendWG.Add(1)

		go func(bck backend) {
			defer backendWG.Done()

			search(query, bck, resultsChan)
		}(b)
	}

	backendWG.Wait()
}

func search(query string, backend backend, resultsChan chan<- types.SearchResult) {
	client := utils.NewEthanolHTTPClient(backend.InsecureSkipSSLVerify)

	request := utils.NewEthanolHTTPClientGETRequest()
	request.Header.Del("Content-Type")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	request.SetBasicAuth(backend.Username, backend.Password)

	rawURL, err := url.JoinPath(backend.Endpoint, "/rest/api/latest/search")
	if err != nil {
		logrus.Panic(err)
	}

	queryURL, err := url.Parse(rawURL)
	queryURL.RawQuery = fmt.Sprintf("jql=summary~%s", query)

	logrus.WithFields(logrus.Fields{
		"url_rawpath":  queryURL.RawPath,
		"url_rawquery": queryURL.RawQuery,
	}).Debug("query raw")

	if err != nil {
		logrus.Panic(err)
	}

	// set request URL
	request.URL = queryURL

	// dump request
	utils.DumpHTTPRequest(request, "jira request")

	// do request
	res, err := client.Do(request)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error in otrs ticketSearch")
		return
	}
	defer res.Body.Close()

	// dump response
	utils.DumpHTTPResponse(res, "jira response")

	var rawResponse RawResponse
	json.NewDecoder(res.Body).Decode(&rawResponse)

	for _, issue := range rawResponse.Issues {
		issueURL, _ := url.JoinPath(backend.Endpoint, "browse", issue.Key)

		result := make(map[string]interface{})
		result["id"] = uuid.New().String()
		result["description"] = issue.Fields.Description
		result["summary"] = issue.Fields.Summary
		result["creator_name"] = issue.Fields.Creator.Name
		result["creator_email"] = issue.Fields.Creator.EmailAddress
		result["created"] = issue.Fields.Created
		result["url"] = issueURL
		result["key"] = issue.Key
		result["label"] = label
		result["name"] = name
		result["raw_label"] = raw_label
		result["source"] = name

		resultsChan <- result
	}
}

var Searcher searchPlugin
