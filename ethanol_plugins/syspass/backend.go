package main

// backend maps backend definition from configuration
// there can be more than one syspass server to query
type backend struct {
	Endpoint              string `json:"endpoint" yaml:"Endpoint"`
	APIKey                string `json:"api_key" yaml:"APIKey"`
	APIKeyPassPhrase      string `json:"api_key_passphrase" yaml:"APIKeyPassPhrase"`
	Count                 int    `json:"count" yaml:"Count"`
	InsecureSkipSSLVerify bool   `json:"insecure_skip_ssl_verify" yaml:"insecureskipsslverify"`
}
