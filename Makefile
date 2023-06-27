@PHONY: run dev

run :
	@go run .
dev :
	@RUN_MODE=dev go run .