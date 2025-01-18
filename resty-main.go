package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"net"
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
	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		log.Fatalf("Failed to append CA certificate")
	}

	// Load client certificate and key (used for client authentication)
	clientCert, err := tls.LoadX509KeyPair("client.crt", "client.key")
	if err != nil {
		log.Fatalf("Failed to load client certificate and key: %v", err)
	}

	// Configure TLS settings for the proxy
	proxyTLSConfig := &tls.Config{
		Certificates:       []tls.Certificate{clientCert}, // Client authentication
		RootCAs:            caCertPool,                    // Trust the proxy's CA
		ServerName:         "localhost",                   // ServerName for SNI (match the proxy's certificate's common name or SAN)
		InsecureSkipVerify: false,                         // Ensure proper certificate validation
	}

	// Create a Resty client
	client := resty.New()

	// Set up the proxy
	client.SetProxy("https://localhost:8080") // Use the same hostname as in ServerName

	// Configure the transport to use custom TLS settings for the proxy
	client.SetTransport(&http.Transport{
		Proxy: http.ProxyURL(mustParseURL("https://localhost:8080")),
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{Timeout: 10 * time.Second}
			conn, err := dialer.DialContext(ctx, network, addr)
			if err != nil {
				return nil, err
			}
			// Upgrade to a TLS connection using the custom proxy TLS config
			tlsConn := tls.Client(conn, proxyTLSConfig)
			if err := tlsConn.Handshake(); err != nil {
				return nil, err
			}
			return tlsConn, nil
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
