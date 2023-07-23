fmt:
	$(GOPATH)\bin\goimports
.PHONY:fmt

lint: fmt
	$(GOPATH)\bin\golangci-lint run ./...
.PHONY:lint