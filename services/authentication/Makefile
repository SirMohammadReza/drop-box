gen_user:
	protoc --go_out=./internal/user --go_opt=paths=source_relative --go-grpc_out=./internal/user --go-grpc_opt=paths=source_relative proto/user/user.proto

gen_token:
	protoc --go_out=./internal/token --go_opt=paths=source_relative --go-grpc_out=./internal/token --go-grpc_opt=paths=source_relative proto/token/token.proto

.PHONY: gen_user, gen_token