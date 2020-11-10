default: testacc

.PHONY: testacc docs

# Run acceptance tests
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
