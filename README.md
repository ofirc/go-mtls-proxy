# go-mtls-proxy
A simple mTLS proxy (tinyproxy + stunnel, no TLS inspection) and a Go client program that uses the proxy.

It demonstrates how to use a client certificate to authenticate to a proxy server that requires it and curl a public URL.

This is based on the following repo: [docker-mtls-https-proxy](
https://github.com/fancy-owl/docker-mtls-https-proxy/).

It differs from the above repo in the following ways:
- It generates certs with SAN to make it comply with Go (needed for Go crypto/tls)
- It adds a Go client program that uses the proxy

The project is meant for purely demonstration purposes, do not use it in production.

# Related projects
* [HTTP CONNECT forward proxy](https://github.com/ofirc/k8s-sniff-https) - useful for a simple HTTP(s) proxy that does not apply deep packet inspection

* [Man-in-the-middle TLS inspection proxy](https://github.com/ofirc/k8s-sniff-https) - useful for deep packet inspection and reverse engineering HTTPS encrypted traffic

## Docker Compose
Generate the certificates and build the container images:
```bash
./scripts/generate-certs.sh
docker compose build
# docker compose push
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

## Kubernetes
Create a local kind cluster:
```bash
kind create cluster --name test-proxy
```

Deploy the resources to the cluster:
```bash
kubectl apply -f deploy
```

Copy the certificates from the Pod:
```bash
POD_NAME=$(kubectl get pod -oname -lapp=stunnel | cut -d'/' -f2)
kubectl cp $POD_NAME:/client-certs/ca.crt ca.crt
kubectl cp $POD_NAME:/client-certs/client.crt client.crt
kubectl cp $POD_NAME:/client-certs/client.key client.key
kubectl cp $POD_NAME:/client-certs/client.pem client.pem
```

Port forward the stunnel:
```bash
kubectl port-forward svc/stunnel 8080
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

curl \
  --proxy https://localhost:8080 \
  --proxy-cacert ca.crt \
  --proxy-cert client.pem \
  https://ipv4.icanhazip.com
```

For example:
```bash
$ curl \
  --proxy https://localhost:8080 \
  --proxy-cacert ca.crt \
  --proxy-cert client.crt \
  --proxy-key client.key \
  https://ipv4.icanhazip.com
84.228.242.243
$
```
