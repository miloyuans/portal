SHELL := /bin/bash

.PHONY: dev-start dev-stop keycloak-bootstrap mongo-indexes test api web

dev-start:
	./scripts/dev-start.sh

dev-stop:
	docker compose down -v

keycloak-bootstrap:
	./scripts/keycloak-bootstrap.sh

mongo-indexes:
	./scripts/init-mongo-indexes.sh

test:
	go test ./...

api:
	go run ./apps/portal-api

web:
	cd apps/portal-web && npm install && npm run dev
