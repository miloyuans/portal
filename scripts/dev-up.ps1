$ErrorActionPreference = "Stop"

if (-not (Test-Path ".env")) {
  Copy-Item ".env.example" ".env"
}

docker compose -f deployments/docker-compose/docker-compose.yml up --build -d
