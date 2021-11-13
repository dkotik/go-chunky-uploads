default:
	go test ./... --tags "sqlite_stat4 sqlite_secure_delete"
example:
	go get .
	go run examples/index.go
