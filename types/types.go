package types

import (
	"net/http"
)

type SearchResult map[string]interface{}

// SearchPlugin plugins must satisfy this interface
type SearchPlugin interface {
	Name() string
	Version() string
	Search(func() *http.Client, func() *http.Request, func() *http.Request, string, chan<- SearchResult)
}
