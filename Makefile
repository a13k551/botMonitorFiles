fmt:
	$(GOPATH)\bin\goimports $(CURDIR)
.PHONY:fmt

lint: fmt
	$(GOPATH)\bin\golangci-lint run ./...
.PHONY:lint