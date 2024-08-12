package main

type backend struct {
	Endpoint              string `json:"endpoint" yaml:"endpoint"`
	APIToken              string `json:"api_token" yamk:"apitoken"`
	InsecureSkipSSLVerify bool   `json:"insecure_skip_ssl_verify" yaml:"insecureskipsslverify"`
}
