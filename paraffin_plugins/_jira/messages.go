package main

type rawResponse struct {
	Expand string `json:"expand"`
	ID     string `json:"id"`
	Self   string `json:"self"`
	Key    string `json:"key"`
	Fields struct {
		Description string `json:"description"`
		Creator     struct {
			Self         string `json:"self"`
			Name         string `json:"name"`
			Key          string `json:"key"`
			EmailAddress string `json:"emailAddress"`
			DisplayName  string `json:"displayName"`
			TimeZone     string `json:"timeZone"`
		}
	}
}
