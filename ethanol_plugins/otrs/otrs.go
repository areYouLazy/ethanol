package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	name        = "otrs_search_plugin"
	provider    = "otrs_6_username_password"
	description = "This plugins get results from an otrs 6 installation through username/password authentication"
	label       = "O.T.R.S."
	raw_label   = "otrs"
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
	err := viper.UnmarshalKey("plugins.otrs", &backends)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error in otrs unmarshal")
		return
	}

	// log
	logrus.WithFields(logrus.Fields{
		"configurations": backends,
	}).Debug("otrs configurations")

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
	var tsr TicketSearchResponse

	query = fmt.Sprintf("*%s*", query)

	rawURL, _ := url.JoinPath(backend.TicketSearchEndpoint, url.QueryEscape(query))
	ticketSearchURL, err := url.Parse(rawURL)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
			"url":   rawURL,
		}).Error("error in otrs url parse")
		return
	}

	// prepare request body for authentication
	body, err := json.Marshal(AuthenticationBody{
		UserLogin: backend.Username,
		Password:  backend.Password,
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error in otrs request body construction")

		return
	}

	// get a new ethanol HTTP post request
	request := utils.NewEthanolHTTPClientPOSTRequest()

	// add content-type header
	request.Header.Del("Content-Type")
	request.Header.Add("Content-Type", "application/json")

	request.URL = ticketSearchURL

	request.Body = io.NopCloser(bytes.NewReader(body))

	// o.t.r.s. does not like chunked requests
	request.ContentLength = int64(len(body))

	// get a new ethanol HTTP client
	client := utils.NewEthanolHTTPClient(backend.InsecureSkipSSLVerify)

	// dump request
	utils.DumpHTTPRequest(request, "otrs request for ticketSearch")

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
	utils.DumpHTTPResponse(res, "otrs response for ticketSearch")

	// iterate results to get ticket for every returned id
	json.NewDecoder(res.Body).Decode(&tsr)

	for _, idx := range tsr.TicketID {
		var tr TicketResponse

		rawURL, _ := url.JoinPath(backend.TicketEndpoint, idx)
		ticketURL, err := url.Parse(rawURL)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
				"url":   backend.TicketEndpoint,
			}).Error("error in otrs url parse")
			return
		}

		ticketRequest := utils.NewEthanolHTTPClientPOSTRequest()
		ticketRequest.Header.Del("Content-Type")
		ticketRequest.Header.Add("Content-Type", "application/json")

		ticketRequest.URL = ticketURL
		ticketRequest.Body = io.NopCloser(bytes.NewReader(body))

		// o.t.r.s. does not like chunked requests
		ticketRequest.ContentLength = int64(len(body))

		// dump request
		utils.DumpHTTPRequest(ticketRequest, "otrs request for ticketGet")

		// do request
		res, err := client.Do(ticketRequest)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("error in otrs ticketGet")
			return
		}
		defer res.Body.Close()

		// dump response
		utils.DumpHTTPResponse(res, "otrs response for ticketGet")

		json.NewDecoder(res.Body).Decode(&tr)

		result := make(map[string]interface{}, 0)
		result["id"] = uuid.New().String()
		result["ticket_id"] = tr.Ticket[0].TicketID
		result["title"] = tr.Ticket[0].Title
		result["number"] = tr.Ticket[0].TicketNumber
		result["state"] = tr.Ticket[0].State
		result["owner"] = tr.Ticket[0].Owner
		result["label"] = label
		result["raw_label"] = raw_label
		result["source"] = name
		result["url"] = fmt.Sprintf("%s%s%s", backend.Endpoint, "?Action=AgentTicketZoom;TicketID=", tr.Ticket[0].TicketID)

		resultsChan <- result
	}
}

// expose structure as symbol
var Searcher searchPlugin
