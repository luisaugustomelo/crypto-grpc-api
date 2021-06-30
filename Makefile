#example: make migration state=down or state=up
migration:
	@migrate -database mongodb://127.0.0.1/test -path migrations/ $(state)