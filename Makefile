.PHONY: dev

dev:
	templ generate --watch --cmd="go run ."
dev-fake:
	templ generate --watch --cmd="go run ./internal/cmd/uifake"
