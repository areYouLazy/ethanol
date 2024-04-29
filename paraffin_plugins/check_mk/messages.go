package main

// rawResponse maps response from check_mk backend
type rawResponse struct {
	Links []struct {
		DomainType string `json:"domainType"`
		Rel        string `json:"rel"`
		Href       string `json:"href"`
		Method     string `json:"method"`
		Type       string `json:"type"`
	}
	Id         string `json:"id"`
	DomainType string `json:"domainType"`
	Value      []struct {
		Links []struct {
			DomainType string `json:"domainType"`
			Rel        string `json:"rel"`
			Href       string `json:"href"`
			Method     string `json:"method"`
			Type       string `json:"type"`
		}
		DomainType string `json:"domainType"`
		Id         string `json:"id"`
		Title      string `json:"title"`
		Members    struct{}
		Extensions struct {
			Name    string `json:"name"`
			Address string `json:"address"`
		}
	}
	Extensions struct{}
}
