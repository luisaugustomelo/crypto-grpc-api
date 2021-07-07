#example: make migration state=down or state=up
migration:
	@migrate -database mongodb://127.0.0.1/klever -path ./databases/migrations/ $(state)

#generate protobuf
gen:
	@protoc --go_out=plugins=grpc:proto ./proto/upvotesystem.proto

mocks:
	@mockery --all --keeptree

cover:
	@go test ./server -cover -coverprofile ./../coverage/fmtcoverage.html fmt