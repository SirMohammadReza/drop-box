PROTO_DIR=proto
OUT_DIR=internal

gen_user:
	protoc -I . \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/user/user.proto

gen_token:
	protoc -I . \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/token/token.proto

gen_uploader:
	protoc -I . \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/uploader/uploader.proto

.PHONY: gen_user gen_token gen_uploader