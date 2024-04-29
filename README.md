# Ethanol

Ethanol aims to be a search aggregator, capable of querying multiple applications and return aggregate results.
It is like a [MetaSearch Engine](https://en.wikipedia.org/wiki/Metasearch_engine) but for applications like [Jira](https://www.atlassian.com/software/jira), [MediaWiki](https://www.mediawiki.org/wiki/MediaWiki), [Check_MK](https://checkmk.com/) and others...

## Table of Content

- [Table of Content](#table-of-content)
- [Background](#background)
- [Installation](#installation)
- [Usage](#usage)
- [Notes](#notes)
- [Help](#help)
- [List of supported plugins](#list-of-supported-plugins)
- [Build Core](#build-core)
- [Build Plugins](#build-plugins)
- [Write Plugins](#write-plugins)
- [Utils](#utils)
	- [Dump HTTP Requests and Responses](#dump-http-requests-and-responses)
	- [Pass plugins data to results](#pass-plugins-data-to-results)

## Background

This project started as a personal playground to improve my `golang` skills, so you'll find test code, dead code, prepared code, any kind of code!

Also I'm not a professional developer, so some decision may sound stranges or completely insane, any help is welcome!

## Installation

Once the 1st release is out you'll be able to `apt install ethanol`, until that moment you need to go the raw way

```console
root@localhost$ # Clone and access repository
root@localhost$ git clone https://github.com/areYouLazy/ethanol
root@localhost$ cd ethanol
root@localhost$
root@localhost$ # start the build.sh script
root@localhost$ ./build.sh
root@localhost$
root@localhost$ # generate a valid configuration file (you can copy the minimal example)
root@localhost$ cp config.minimal.yml config.yml
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
  -config-file string
    	use a custom configuration file
  -config-json
    	read configuration file as json
  -debug
    	print debug log messages (default false)
  -json
    	print log messages in json format (default false)
```

## List of supported plugins

- [Check_MK](https://checkmk.com/)
- [SysPass](https://syspass.org/)

## Build Core

The `build.sh` script is designed to build the ethanol core, or you can compile it by hand

```console
root@localhost$ go build
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
	Version() string
	Search(func() *http.Client, func() *http.Request, func() *http.Request, string, chan<- SearchResult)
}
```

You should use `GetNewHTTPClient()` to get a preconfigured HTTP Client. This will help Ethanol to be consist across HTTP calls 
when it acts as an HTTP Client

In the same way, `GetNewHTTPGETRequest()` will provides both a GET and a POST requests. This will help Ethanol to be consist across HTTP calls 
when it acts as an HTTP Client

The Plugin `export` is usually done at the end of the file

```go
// example_plugin.go
// compile with # go build -buildmode=plugin -o example_plugin.so example_plugin.go
package main

const (
    name = "example_search_plugin"
    version = "0.1"
)

type searchPlugin interface{}

func (s *searchPlugin) Name() string {
    return name
}

func (s *searchPlugin) Version() string {
    return version
}

func (s *searchPlugin) Search(getNewHTTPClient func() *http.Client, getNewHTTPGetRequest func() *http.Request, getNewHTTPPostRequest func() *http.Request, query string, resultsChan chan<- types.SearchResult) {

	var backendWG sync.WaitGroup

    for _, b := range backends {
		backendWG.Add(1)

		client := getNewHTTPClient()
		req := getNewHTTPGetRequest()

		go func() {
			defer backendWG.Done()
			search(client, req, query, b, resultsChan)
		}(b)
    }

	backendWG.Wait()
}

// [...]

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

You should append plugins data to every result, like this.

Those info's can be used in the UI to better diplay results

```go
var (
	name = "example_search_plugin"
	label = "Example"
	version = "0.1"
)

// [...]

result["source_name"] = name
result["source_label"] = label
result["source_version"] = version

resultsChan <- result
```
