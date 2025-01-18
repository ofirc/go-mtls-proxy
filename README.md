# go-mtls-proxy
A simple mTLS proxy (tinyproxy + stunnel, no TLS inspection) and a Go client program that uses the proxy.

It demonstrates hwo to use a client certificate to authenticate to a proxy server that requires it and curl a public URL.

This is based on the following repo: [docker-mtls-https-proxy](
https://github.com/fancy-owl/docker-mtls-https-proxy/).

It differs from the project in the following ways:
- It generates certs with SAN to make it comply with Go (needed for Go crypto/tls)
- It adds a Go client program that uses the proxy

The project is meant for purely demonstration purposes, do not use it in production.

To run:
```bash
./scripts/generate-certs.sh
docker compose build
docker compose up
```

And then on a separate shell:
```bash
go run main.go
curl \
  --proxy https://localhost:8080 \
  --proxy-cacert ca.crt \
  --proxy-cert client.crt \
  --proxy-key client.key \
  https://ipv4.icanhazip.com
```

For example:
```bash
$ curl --proxy https://localhost:8080 --proxy-cacert ca.crt --proxy-cert client.crt --proxy-key client.key https://ipv4.icanhazip.com
84.228.242.243
$ go run main.go                                                                                                                     
84.228.242.243
```