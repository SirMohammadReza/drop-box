PROTO_DIR=proto

gen_user:
	protoc -I . \
		--go_out=./internal/user --go_opt=paths=source_relative \
		--go-grpc_out=./internal/user --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/user/user.proto

gen_token:
	protoc -I . \
		--go_out=./internal/token --go_opt=paths=source_relative \
		--go-grpc_out=./internal/token --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/token/token.proto

.PHONY: gen_user gen_token