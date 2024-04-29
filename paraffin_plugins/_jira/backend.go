package main

// backend maps backend definition from configuration
// there can be more than one check_mk server to query
type backend struct {
	Endpoint string `json:"endpoint" yaml:"Endpoint"`
	Username string `json:"username" yaml:"Username"`
	Password string `json:"password" yaml:"Password"`
}
