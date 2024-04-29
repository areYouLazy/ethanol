package main

//
// Should we use https://github.com/andygrunwald/go-jira ?
//

import (
	"net/http"
	"sync"

	"github.com/areYouLazy/ethanol/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// define plugin constants
const (
	name    = "jira_search_plugin"
	label   = "Jira"
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

	err := viper.UnmarshalKey("Plugins.JIRA", &backends)
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

			search(getNewHTTPClient(), getNewHTTPGetRequest(), query, bck, resultsChan)
		}(b)
	}

	backendWG.Wait()
}

func search(client *http.Client, request *http.Request, query string, backend backend, resultsChan chan<- types.SearchResult) {
}
