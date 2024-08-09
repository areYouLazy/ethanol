package main

// backend maps backend definition from configuration
// there can be more than one jira server to query
type backend struct {
	Endpoint              string `json:"endpoint" yaml:"Endpoint"`
	Username              string `json:"username" yaml:"Username"`
	UserEmail             string `json:"useremail" yaml:"useremail"`
	Password              string `json:"password" yaml:"Password"`
	APIToken              string `json:"api_token" yaml:"apitoken"`
	InsecureSkipSSLVerify bool   `json:"insecure_skip_ssl_verify" yaml:"insecureskipsslverify"`
}
