CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/linux/api ./cmd/api/main.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/linux/worker ./cmd/worker/main.go