set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -o ./bin/linux/api ./cmd/api/main.go
go build -o ./bin/linux/worker ./cmd/worker/main.go