package plugins

import (
	"os"
	"path"
	"plugin"
	"reflect"
	"strings"
	"sync"

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
	plugins, err := os.ReadDir(viper.GetString("ethanol.server.pluginsfolder"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":  err.Error(),
			"folder": viper.GetString("ethanol.server.pluginsfolder"),
		}).Error("error reading plugin folder")

		return
	}

	// iterate files
	for _, v := range plugins {
		// check if it ends with .so and is not a fodler
		if strings.HasSuffix(v.Name(), ".so") && !v.IsDir() {
			// construct absolute path
			pluginName := path.Join(
				viper.GetString("ethanol.server.pluginsfolder"),
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
	}).Info("plugins loaded")
}

// BulkSearch runs a search for every loaded plugin
func BulkSearch(query string) ([]types.SearchResult, error) {
	// define a results Channel
	resultsChan := make(chan types.SearchResult)
	results := make([]types.SearchResult, 0)

	// create a waitgroup to wait for all plugins to execute
	var wg sync.WaitGroup

	// TODO(areYouLazy) : is this even correct?
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// iterate plugins
	for _, plugin := range pluginsManager.plugins {
		// increment waitgroup counter
		wg.Add(1)

		// start a goroutine for every plugin
		go func(plugin types.SearchPlugin) {
			// decrement waitgroup on exit
			defer wg.Done()

			// log
			logrus.WithFields(logrus.Fields{
				"name":     plugin.Name(),
				"version":  plugin.Version(),
				"provider": plugin.Provider(),
			}).Debug("collecting search results from plugin")

			// fire search, results are sent to plugin channel
			plugin.Search(query, resultsChan)
		}(plugin)
	}

	// collect results from results channel
	for result := range resultsChan {
		results = append(results, result)
	}

	logrus.WithFields(logrus.Fields{
		"query": query,
	}).Info("search completed")

	// return collected results
	return results, nil
}
