#!/bin/bash

# Set variables
CA_SUBJECT="/CN=TestCA"
SERVER_SUBJECT="/CN=localhost"
CLIENT_SUBJECT="/CN=client"
SAN="subjectAltName=DNS:localhost,DNS:stunnel.default,IP:127.0.0.1"

# Clean up old files
rm -f ca.* server.* client.*

# Generate CA private key and certificate
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -sha256 -days 3650 -out ca.crt -subj "$CA_SUBJECT"

# Generate server private key and CSR
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -subj "$SERVER_SUBJECT"

# Create server certificate with SAN
cat > server.ext <<EOF
basicConstraints=CA:FALSE
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
$SAN
EOF
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
    -out server.crt -days 365 -sha256 -extfile server.ext

# Generate client private key and CSR
openssl genrsa -out client.key 2048
openssl req -new -key client.key -out client.csr -subj "$CLIENT_SUBJECT"

# Create client certificate with SAN
cat > client.ext <<EOF
basicConstraints=CA:FALSE
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth
$SAN
EOF
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial \
    -out client.crt -days 365 -sha256 -extfile client.ext

# Combine client certificate, CA certificate, and private key into PEM
cat client.crt ca.crt client.key > client.pem

# Combine server certificate, CA certificate, and private key into PEM
cat server.crt ca.crt server.key > server.pem

echo "Certificates generated successfully."
