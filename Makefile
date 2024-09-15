run-books:
	@cd services/book-service && go run .

run-categories:
	@cd services/book-category-service && go run .

run-users:
	@cd services/user-service && go run .

gen-api:
	@protoc \
    --proto_path=protobuf "protobuf/api.proto" \
    --go_out=services/common/grpc-api --go_opt=paths=source_relative \
    --go-grpc_out=services/common/grpc-api --go-grpc_opt=paths=source_relative