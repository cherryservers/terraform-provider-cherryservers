default: testacc

# Run acceptance tests
.PHONY: testacc

testacc:
	TF_ACC=1 go test ./... -v -timeout 60m $(TESTARGS)
