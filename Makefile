all: build_plugin
	-rm output/*
	go run -race cmd/main.go -i 'txt/pg-grimm.txt,txt/pg-being_ernest.txt' -p 'cmd/wc.so' -r 2 -w 4
	go run -race cmd/main.go -i 'output/mr-out-*' -p 'cmd/merge.so' -r 1 -w 4
	diff -q output/mr-out-0.txt ~/work/6.824/src/main/mr-out-0

clean :
	-rm output/*

grpc :
	protoc --go-grpc_out=rpc --go_out=rpc rpc/*.proto

build_plugin:
	go build -race -buildmode=plugin -o cmd/wc.so ./mrapps/wc.go
	go build -race -buildmode=plugin -o cmd/merge.so ./mrapps/merge.go

