default: testacc

.PHONY: testacc sweep

testacc:
	TF_ACC=1 go test ./... -v -timeout 60m $(TESTARGS)

sweep:
	TF_ACC=1 go test ./internal/provider -args -sweep='no-region'