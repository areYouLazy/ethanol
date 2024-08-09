package types

type SearchResult map[string]interface{}

// SearchPlugin plugins must satisfy this interface
type SearchPlugin interface {
	Name() string
	Provider() string
	Description() string
	Version() string

	Search(string, chan<- SearchResult)
}
