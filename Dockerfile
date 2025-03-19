FROM alpine:3.10 as proxy
RUN set -xe && \
    apk add --no-cache tinyproxy
COPY tinyproxy.conf /etc/tinyproxy/

FROM alpine:3.10 as stunnel
COPY ca.crt server.pem stunnel.conf /etc/stunnel/
RUN set -xe \
 && apk add --update --no-cache stunnel \
 && chmod 600 /etc/stunnel/server.pem \
 && mkdir -p /client-certs
COPY ca.crt client.crt client.key client.pem /client-certs/