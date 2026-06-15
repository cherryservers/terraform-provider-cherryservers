default: testacc

# Run acceptance tests
.PHONY: testacc

testacc:
	TF_AC=1 go test ./... -v -timeout 60m $(TESTARGS)
