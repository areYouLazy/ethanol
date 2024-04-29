package plugins

import (
	"net/http"
	"net/http/cookiejar"
	"os"
	"path"
	"plugin"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/areYouLazy/ethanol/proxy"
	"github.com/areYouLazy/ethanol/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	pluginsManager *PluginsManager
)

type PluginsManager struct {
	plugins []types.SearchPlugin
}

func Init() {
	pluginsManager = &PluginsManager{
		plugins: make([]types.SearchPlugin, 0),
	}

	// load plugins
	plugins, err := os.ReadDir(viper.GetString("Ethanol.Server.PluginsFolder"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err.Error(),
			"folder": viper.GetString("Ethanol.Server.PluginsFolder"),
		}).Error("error reading plugin folder")

		return
	}

	// iterate files
	for _, v := range plugins {
		// check if it ends with .so and is not a fodler
		if strings.HasSuffix(v.Name(), ".so") && !v.IsDir() {
			// construct absolute path
			pluginName := path.Join(
				viper.GetString("Ethanol.Server.PluginsFolder"),
				v.Name(),
			)

			// open plugin
			plug, err := plugin.Open(pluginName)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error":           err.Error(),
					"plugin_filename": v.Name(),
				}).Error("error loading search plugin")
				continue
			} else {
				logrus.WithFields(logrus.Fields{
					"plugin_filename": v.Name(),
				}).Debug("loading plugin")
			}

			// search for Searcher symbol
			searchSymbol, err := plug.Lookup("Searcher")
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error":           err.Error(),
					"plugin_filename": v.Name(),
					"symbol":          "Searcher",
				}).Error("error getting symbol from plugin")
				continue
			} else {
				logrus.WithFields(logrus.Fields{
					"plugin_filename": v.Name(),
					"symbol":          searchSymbol,
					"typeOf":          reflect.TypeOf(searchSymbol),
				}).Debug("found symbol for plugin")
			}

			// typecast plugin
			res := searchSymbol.(types.SearchPlugin)

			// add new plugin to plugins manager
			pluginsManager.plugins = append(pluginsManager.plugins, res)
		}
	}

	logrus.WithFields(logrus.Fields{
		"number_of_plugins": len(pluginsManager.plugins),
	}).Debug("plugins loaded")
}

func BulkSearch(query string) ([]types.SearchResult, error) {
	// define a results Channel
	resultsChan := make(chan types.SearchResult)
	results := make([]types.SearchResult, 0)

	// create a waitgroup to wait for all plugins to execute
	var wg sync.WaitGroup

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// create mutex to handle concurrent writes
	// var mu sync.Mutex

	// iterate plugins
	for _, plugin := range pluginsManager.plugins {
		// increment waitgroup counter
		wg.Add(1)

		// goroutine for every plugin
		go func(plugin types.SearchPlugin) {
			// decrement waitgroup on exit
			defer wg.Done()

			// log
			logrus.WithFields(logrus.Fields{
				"name":    plugin.Name(),
				"version": plugin.Version(),
			}).Debug("collecting search results from plugin")

			// fire search, results are sent to plugin channel
			plugin.Search(getNewHTTPClient, GetNewHTTPGETRequest, GetNewHTTPPOSTRequest, query, resultsChan)
		}(plugin)
	}

	// collect results from results channel
	for result := range resultsChan {
		results = append(results, result)
	}

	// return collected results
	return results, nil
}

func getNewHTTPClient() *http.Client {
	// setup cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil
	}

	// setup transport with proxy (if required)
	transport := http.Transport{
		Proxy: proxy.GetHTTPProxyURL(),
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

func GetNewHTTPClient() *http.Client {
	return getNewHTTPClient()
}

func getNewHTTPRequests() (*http.Request, *http.Request) {
	// setup GET http request
	getRequest, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		return nil, nil
	}
	// setup user-agent header
	getRequest.Header.Add("User-Agent", viper.GetString("Ethanol.Client.UserAgent"))

	// setup post HTTP request
	postRequest, err := http.NewRequest(http.MethodPost, "", nil)
	if err != nil {
		return nil, nil
	}
	// setup user-agent header
	postRequest.Header.Add("User-Agent", viper.GetString("Ethanol.Client.UserAgent"))

	// return requests
	return getRequest, postRequest
}

func GetNewHTTPGETRequest() *http.Request {
	// setup GET http request
	getRequest, err := http.NewRequest(http.MethodGet, "", nil)
	if err != nil {
		return nil
	}

	// setup user-agent header
	getRequest.Header.Add("User-Agent", viper.GetString("Ethanol.Client.UserAgent"))

	return getRequest
}

func GetNewHTTPPOSTRequest() *http.Request {
	// setup post HTTP request
	postRequest, err := http.NewRequest(http.MethodPost, "", nil)
	if err != nil {
		return nil
	}
	// setup user-agent header
	postRequest.Header.Add("User-Agent", viper.GetString("Ethanol.Client.UserAgent"))

	// return requests
	return postRequest
}
