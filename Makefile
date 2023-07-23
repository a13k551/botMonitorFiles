fmt:
	go fmt ./...
.PHONY:fmt

lint: fmt
	$(GOPATH)\bin\golangci-lint run ./...
.PHONY:lint