#!/bin/sh
set -eu

envsubst '${PORTAL_API_BASE_URL} ${PORTAL_APP_TITLE}' \
  < /usr/share/nginx/html/portal-config.template.js \
  > /usr/share/nginx/html/portal-config.js
