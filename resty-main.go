package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {
	// Load CA certificate (used to validate the proxy's server certificate)
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		log.Fatalf("Failed to read CA certificate: %v", err)
	}
	caCertPool, err := x509.SystemCertPool()
	if err != nil {
		log.Fatalf("Failed to load system CA pool: %v", err)
	}
	if caCertPool == nil {
		caCertPool = x509.NewCertPool()
	}
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatalf("Failed to append CA certificate")
	}

	// Load client certificate and key (used for client authentication)
	clientCert, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		log.Fatalf("Failed to load client certificate and key: %v", err)
	}

	// Create a Resty client
	client := resty.New()

	// Set up the proxy
	client.SetProxy("https://localhost:8080") // Use the same hostname as in ServerName

	// Configure the transport to use custom TLS settings for the proxy
	client.SetTransport(&http.Transport{
		Proxy: http.ProxyURL(mustParseURL("https://localhost:8080")),
		TLSClientConfig: &tls.Config{
			RootCAs:      caCertPool,
			Certificates: []tls.Certificate{clientCert},
		},
	})

	// Set timeout for requests
	client.SetTimeout(10 * time.Second)

	// Make HTTPS request through the proxy
	resp, err := client.R().Get("https://ipv4.icanhazip.com")
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}

	// Print the response
	fmt.Println(string(resp.Body()))
}

// mustParseURL parses a URL or panics if it fails.
func mustParseURL(rawURL string) *url.URL {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("Invalid URL: %v", err)
	}
	return parsedURL
}
