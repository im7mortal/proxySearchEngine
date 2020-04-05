#!/usr/bin/env bash

: "${local_proxySearchEngine_PORT?Need to set local_proxySearchEngine_PORT what is the container port}"
: "${local_proxySearchEngine_HOST?Need to set local_proxySearchEngine_HOST what is the domain of the proxy}"
: "${local_proxySearchEngine_LETSENCRYPT_EMAIL?Need to set local_proxySearchEngine_LETSENCRYPT_EMAIL what is the email for lets encrypt}"

export PORT="$local_proxySearchEngine_PORT"
export HOST="$local_proxySearchEngine_HOST"
export VIRTUAL_HOST="$HOST"
export LETSENCRYPT_HOST="$HOST"
export LETSENCRYPT_EMAIL="$local_proxySearchEngine_LETSENCRYPT_EMAIL"

docker-compose up --build -d
