FROM node:22-alpine AS builder

WORKDIR /web
COPY apps/portal-web/package.json ./
RUN npm install

COPY apps/portal-web ./
RUN npm run build

FROM nginx:1.27-alpine

ENV PORTAL_API_BASE_URL=http://localhost:8080/api
ENV PORTAL_APP_TITLE=Portal

COPY apps/portal-web/nginx.conf /etc/nginx/conf.d/default.conf
COPY apps/portal-web/docker-entrypoint.sh /docker-entrypoint.d/40-portal-config.sh
COPY --from=builder /web/dist /usr/share/nginx/html
COPY apps/portal-web/public/portal-config.template.js /usr/share/nginx/html/portal-config.template.js

RUN chmod +x /docker-entrypoint.d/40-portal-config.sh

EXPOSE 80
