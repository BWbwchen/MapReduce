all: build_plugin
	rm output/*
	go run -race cmd/main.go 
	diff -q output/mr-out.txt ~/work/6.824/src/main/mr-out-0

grpc :
	protoc --go-grpc_out=rpc --go_out=rpc rpc/*.proto

build_plugin:
	go build -race -buildmode=plugin -o cmd/wc.so ./mrapps/wc.go

