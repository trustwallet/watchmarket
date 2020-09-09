set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64
go build -o ./bin/linux/api ./cmd/api/main.go || exit /b 1
go build -o ./bin/linux/worker ./cmd/worker/main.go || exit /b 1