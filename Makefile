DESTDIR ?= ./bin

$(DESTDIR)/spelunk: go.mod main.go spelunk/walk.go
	goimports -w .
	gofmt -s -w .
	go vet ./...
	go test ./...
	go build -o $@
