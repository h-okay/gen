.PHONY: gen
gen:
	@oapi-codegen -config server.cfg.yaml openapi.yaml
	@sqlc generate