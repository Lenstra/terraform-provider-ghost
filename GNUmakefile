.PHONY: default
default: fmt lint

.PHONY: fmt
fmt:
	gofmt -s -w -e .

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test-client
test-client:
	go test -cover ./client/...

.PHONY: test
test:
	go test -v -cover -timeout=120s -parallel=10 ./internal/...

.PHONY: testacc
testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./internal/...

.PHONY: testacc-refresh
testacc-refresh:
	TFTEST_REFRESH_STATE=1 TF_ACC=1 go test -v -cover -timeout 120m ./internal/...
