package v1beta1

import "crypto/tls"

type Config struct {
	insecure           bool
	url                string
	httpUrl            string
	enableOIDCAuth     bool
	clientId           string
	clientSecret       string
	authServerTokenUrl string
	tlsConfig          *tls.Config
}

type ClientOptions func(*Config)

func WithAuthEnabled(clientId string, clientSecret string, authServerTokenUrl string) ClientOptions {
	return func(c *Config) {
		c.enableOIDCAuth = true
		c.clientId = clientId
		c.clientSecret = clientSecret
		c.authServerTokenUrl = authServerTokenUrl
	}
}

func WithgRPCUrl(url string) ClientOptions {
	return func(c *Config) {
		c.url = url
	}
}

func WithHTTPUrl(url string) ClientOptions {
	return func(c *Config) {
		c.httpUrl = url
	}
}

func WithTLSInsecure(insecure bool) ClientOptions {
	return func(c *Config) {
		c.insecure = insecure
	}
}

func WithHTTPTLSConfig(tlsConfig *tls.Config) ClientOptions {
	return func(c *Config) {
		c.insecure = false
		c.tlsConfig = tlsConfig
	}
}

func NewConfig(options ...func(*Config)) *Config {
	svr := &Config{}
	for _, o := range options {
		o(svr)
	}
	return svr
}
