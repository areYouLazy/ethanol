package main

type backend struct {
	Endpoint              string `json:"endpoint" yaml:"endpoint"`
	TicketSearchEndpoint  string `json:"ticket_search_endpoint" yaml:"ticketsearchendpoint"`
	TicketEndpoint        string `json:"ticket_endpoint" yaml:"ticketendpoint"`
	Username              string `json:"username" yaml:"username"`
	Password              string `json:"password" yaml:"password"`
	InsecureSkipSSLVerify bool   `json:"insecure_skip_ssl_verify" yaml:"insecureskipsslverify"`
}
