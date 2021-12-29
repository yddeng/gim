make proto:
	protoc -I $(shell pwd) --go_out=paths=source_relative:. protocol/*.proto;cd im/protocol;rm *.go;mv ../../protocol/*.go ./;cd ../../
	#cd protocol/; protoc --go_out=. *.proto; cd ../im/pb;rm *.go;mv ../../protocol/*.go ./;cd ../../