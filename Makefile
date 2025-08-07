build:
	@mkdir -p bin
	go build -o bin/llm-context ./cmd/llm-context

.PHONY: build
