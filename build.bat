set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -o ./bin/api ./cmd/api/main.go
go build -o ./bin/worker ./cmd/worker/main.go
REM docker build -t trust/watchmarket:api-local --build-arg SERVICE=api .
REM #docker build -t trust/watchmarket:worker-local --build-arg SERVICE=worker .
docker build -t trust/watchmarket:seed-local -f seed/Dockerfile seed/
docker build -t trust/watchmarket:proxy-local -f nginx/Dockerfile nginx/
docker build -t trust/watchmarket:pg-health-local -f scripts/pg-check/Dockerfile scripts/pg-check/