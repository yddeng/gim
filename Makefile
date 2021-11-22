make proto:
	cd internal/protocol; protoc --go_out=. *.proto; cd ../../