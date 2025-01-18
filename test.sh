#! /usr/bin/env bash

URL=https://ipv4.icanhazip.com

echo -e "\n******\n"

echo "This should fail - no client cert spec'd -- https"
echo "curl ${URL}"
curl \
    --proxy https://localhost:8080 \
    --proxy-cacert ca.crt \
    ${URL}
echo -e "\n******\n"

echo "This should fail; no client private key spec'd -- https"
echo "curl ${URL}"
curl \
    --proxy https://localhost:8080 \
    --proxy-cacert ca.crt \
    --proxy-cert client.crt \
    ${URL}
echo -e "\n******\n"

echo "This should pass; client cert and key in separate files -- https"
echo "curl ${URL}"
curl \
    --proxy https://localhost:8080 \
    --proxy-cacert ca.crt \
    --proxy-cert client.crt \
    --proxy-key client.key \
    ${URL}
echo -e "\n******\n"

echo "This should pass; combined client cert and key -- https"
echo "curl ${URL}"
curl \
    --proxy https://localhost:8080 \
    --proxy-cacert ca.crt \
    --proxy-cert client.pem \
    ${URL}
echo -e "\n******\n"

URL=http://ipv4.icanhazip.com
echo "This should pass; combined client cert and key -- http"
echo "curl ${URL}"
curl \
    --proxy https://localhost:8080 \
    --proxy-cacert ca.crt \
    --proxy-cert client.pem \
    ${URL}
echo -e "\n******\n"
