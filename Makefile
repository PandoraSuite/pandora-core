VERSION ?= v1
OUT_BASE := internal/adapters/grpc
PROTO_DIR := proto

PROTO_FILES := $(shell find $(PROTO_DIR) -type f -path "*/$(VERSION)/*.proto")

PROTOC_GEN_VALIDATE := $(shell which protoc-gen-validate)

.PHONY: all proto clean

all:
	@echo "ℹ️  Run 'make proto VERSION=v1' to generate Go code from .proto files"

proto:
	@if [ -z "$(PROTO_FILES)" ]; then \
		echo "❌ No .proto files found for version '$(VERSION)' in $(PROTO_DIR)"; \
		exit 1; \
	fi
	@echo "🔧 Generating Go code for version: $(VERSION)..."
	@for file in $(PROTO_FILES); do \
		entity=$$(echo $$file | sed -E 's|^$(PROTO_DIR)/([^/]+)/$(VERSION)/.*|\1|'); \
		tmp_out_dir=$(OUT_BASE)/$$entity/pb/tmp; \
		final_out_dir=$(OUT_BASE)/$$entity/pb; \
		mkdir -p $$tmp_out_dir; \
		protoc \
			-I $(PROTOC_GEN_VALIDATE) \
			--proto_path=$(PROTO_DIR) \
			--go_out=$$tmp_out_dir \
			--go_opt=paths=source_relative \
			--go-grpc_out=$$tmp_out_dir \
			--go-grpc_opt=paths=source_relative \
			--validate_out="lang=go:$$tmp_out_dir" \
			$$file; \
		find $$tmp_out_dir -type f -name '*.go' -exec mv {} $$final_out_dir \; ; \
		rm -rf $$tmp_out_dir; \
		echo "✅ Processed: $$file → $$final_out_dir"; \
	done
	@echo "🎉 Code generation completed successfully."

clean:
	@echo "🧹 Removing generated files..."
	@find $(OUT_BASE) -type f \( -name '*.pb.go' -o -name '*.pb.validate.go' \) -delete
	@echo "✔️  Cleanup complete."
