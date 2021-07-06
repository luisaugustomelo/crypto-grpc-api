#example: make migration state=down or state=up
migration:
	@migrate -database mongodb://127.0.0.1/klever -path migrations/ $(state)

#generate protobuf
gen:
	@protoc --go_out=plugins=grpc:upvote ./upvote/upvotesystem.proto

mocks:
	@mockery --all --keeptree

coverage:
	@go test -coverprofile ./../coverage/fmtcoverage.html fmt