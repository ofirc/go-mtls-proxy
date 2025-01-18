package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
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
		ServerName:         "localhost",                   // ServerName for SNI (match the proxy hostname)
		InsecureSkipVerify: false,                         // Ensure proper certificate validation
	}

	// Create HTTPS client with proxy configuration
	client := &http.Client{
		Transport: &http.Transport{
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
		},
	}

	// Make HTTPS request through the proxy
	resp, err := client.Get("https://ipv4.icanhazip.com")
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Print response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	fmt.Println(string(body))
}

// mustParseURL parses a URL or panics if it fails.
func mustParseURL(rawURL string) *url.URL {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		log.Fatalf("Invalid URL: %v", err)
	}
	return parsedURL
}
