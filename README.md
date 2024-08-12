# Ethanol

Ethanol aims to be a search aggregator, capable of querying multiple applications and return aggregate results.
It is like a [MetaSearch Engine](https://en.wikipedia.org/wiki/Metasearch_engine) but for applications like [Jira](https://www.atlassian.com/software/jira), [MediaWiki](https://www.mediawiki.org/wiki/MediaWiki), [Check_MK](https://checkmk.com/) and others...

## Table of Content

- [Ethanol](#ethanol)
	- [Table of Content](#table-of-content)
	- [Supported Environments](#supported-environments)
	- [Installation](#installation)
	- [Usage](#usage)
	- [Notes](#notes)
	- [Help](#help)
	- [Build Core](#build-core)
	- [Build Plugins](#build-plugins)
	- [Write Plugins](#write-plugins)
	- [Utils](#utils)
		- [Dump HTTP Requests and Responses](#dump-http-requests-and-responses)
		- [Pass plugins data to results](#pass-plugins-data-to-results)

## Supported Environments

Ethanol is actually capable of querying the following systems:

- [Check_MK](https://checkmk.com/)
- [SysPass](https://syspass.org/en)
- [O.T.R.S.](https://otrs.com/)
- [Jira](https://www.atlassian.com/software/jira)

## Installation

```console
root@localhost$ # Clone and access repository
root@localhost$ git clone https://github.com/areYouLazy/ethanol
root@localhost$ cd ethanol
root@localhost$
root@localhost$ # start the build.sh script
root@localhost$ ./build.sh
root@localhost$
root@localhost$ # generate a valid configuration file
root@localhost$ cp config.template.yml config.yml
root@localhost$
root@localhost$ # edit the configuration file according to your needs, than you can start ethanol
root@localhost$ ./ethanol
```

## Usage

Once `ethanol` is started, you can connect to the WebUI. The default address is `http://<IP-address>:8888/ui/index`

From there you can use the Search Bar to input your searches.

## Notes

Keep in mind that, while we love `RegEx`, some applications may not be able to correctly handle special kind of search, so the same query (let's say: `SRV*`) can produce a wide range of results for one plugin while returning no results in another one because there's no object literaly called `SRV*`

This is not a thing `ethanol` can handle (as far as I know!)

## Help

```console
root@localhost$
Usage of ./ethanol:
  -caller
        print log messages caller (default false)
  -config string
        use a custom configuration file (default "config.yml")
  -debug
        print debug log messages (default false)
  -json
        print log messages in json format (default false)
```

## Build Core

The `build.sh` script is designed to build the ethanol core, along with all available plugins

```console
root@localhost$ ./build.sh
```

## Build Plugins

The `build.sh` script is designed to build every plugin it find in the `ethanol_plugins` folder,
ignoring folders whose name begins with an underscore _

```console
root@localhost$ cd ./ethanol_plugins/$f
root@localhost$ go build -buildmode=plugin -o ../../search_providers/"$f".so ./*.go
root@localhost$ cd ..

root@localhost$ # example
root@localhost$ cd ./ethanol_plugins/check_mk
root@localhost$ go build -buildmode=plugins -o ../../search_providers/check_mk.so ./*.go
root@localhost$ cd ..
```

## Write Plugins

A plugin must export the `Searcher` symbol as a structure that satisfies the `types.SearchPlugin` interface.

```go
type SearchResult map[string]interface{}

// SearchPlugin plugins must satisfy this interface
type SearchPlugin interface {
	Name() string
	Provider() string
	Description() string
	Version() string
	Search(string, chan<- SearchResult)
}
```

You should use `utils.NewEthanolHTTPClient()` to get a preconfigured HTTP Client. This will help Ethanol to be consist across HTTP calls 
when it acts as an HTTP Client

In the same way, `utils.NewEthanolHTTPClientGETRequest()` and `utils.NewEthanolHTTPClientPOSTRequest()` will provides a GET and a POST requests. This will help Ethanol to be consist across HTTP calls when it acts as an HTTP Client

The Plugin `export` is usually done at the end of the file

```go
// example_plugin.go
// compile with # go build -buildmode=plugin -o example_plugin.so example_plugin.go
package main

const (
    name = "example_search_plugin"
	provider = "example_search_1.0_username_password"
	description = "get results from an example_search installation through username/password authentication"
    version = "0.1"
	label = "Example Search"
	raw_label = "example_search"
)

type searchPlugin interface{}

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

// Search prepare environment for queries
func (s *searchPlugin) Search(query string, resultsChan chan<- types.SearchResult) {

	var backendWG sync.WaitGroup

    for _, b := range backends {
		backendWG.Add(1)

		go func() {
			defer backendWG.Done()
			search(query, b, resultsChan)
		}(b)
    }

	backendWG.Wait()
}

// search performs the actual search
func search(query string, backend backend, results chan<- types.SearchResult) {
	client := utils.NewEthanolHTTPClient()

	request := utils.NewEthanolHTTPClientGETRequest()

	// do stuffs and send results to the channel
	results <- res
}

// export symbol
var Searcher searchPlugin
```

## Utils

### Dump HTTP Requests and Responses

Use `utils.DumpHTTPRequest` and `utils.DumpHTTPResponse` in your plugins (this will help in troubleshooting http related problems) like this:

```go
// REQUESTS: just before client.Do()
// RESPONSES: just after client.Do()

// dump request
utils.DumpHTTPRequest(request, "request to plugin <plugin_name>")

// do request
response, err := client.Do(request)
if err != nil {
	logrus.WithFields(logrus.Fields{
		"error": err.Error(),
	}).Error("error in plugin <plugin_name> request")
	return nil, err
}
defer response.Body.Close()

// dump response
utils.DumpHTTPResponse(response, "response from plugin <plugin_name>")
```

### Pass plugins data to results

You should append plugins data to every result before send it.

Those info's can be used in the UI to better diplay results

`raw_label` is actualy used to select the correct template to render the result's card

```go
// [...]

result["name"] = name
result["label"] = label
result["raw_label"] = raw_label
result["description"] = description
result["provider"] = provider
result["version"] = version

resultsChan <- result
```
