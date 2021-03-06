#!/usr/bin/env bash

: "${local_proxySearchEngine_PORT?Need to set local_proxySearchEngine_PORT what is the container port}"
: "${local_proxySearchEngine_HOST?Need to set local_proxySearchEngine_HOST what is the domain of the proxy}"
: "${local_proxySearchEngine_LETSENCRYPT_EMAIL?Need to set local_proxySearchEngine_LETSENCRYPT_EMAIL what is the email for lets encrypt}"

PORT="$local_proxySearchEngine_PORT" \
HOST="$local_proxySearchEngine_HOST" \
VIRTUAL_HOST="$HOST" \
LETSENCRYPT_HOST="$HOST" \
LETSENCRYPT_EMAIL="$local_proxySearchEngine_LETSENCRYPT_EMAIL" \
docker-compose up --build -d
