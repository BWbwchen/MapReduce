all: clean grpc build_plugin
	-rm output/*
	-rm mr-out-*
	go run -race cmd/main.go -i 'txt/*' -p 'cmd/wc.so' -r 4 -w 8
	go run -race cmd/main.go -i 'mr-out-*' -p 'cmd/merge.so' -r 1 -w 4
	diff -q mr-out-0.txt ~/work/6.824/src/main/mr-out-0
	rm /dev/shm/imd-*

clean :
	-rm output/*
	-rm cmd/*.so

grpc :
	protoc --go-grpc_out=rpc --go_out=rpc rpc/*.proto

build_plugin:
	go build -race -buildmode=plugin -o cmd/wc.so ./mrapps/wc.go
	go build -race -buildmode=plugin -o cmd/merge.so ./mrapps/merge.go

build_crash:
	go build -race -buildmode=plugin -o cmd/crash.so ./mrapps/crash.go

test_multi_node: clean grpc build_plugin
	-rm output/*
	-rm mr-out-*
	go run -race cmd/master.go -i 'txt/*' -p 'cmd/wc.so' -r 4 -w 8 &
	sleep 1
	go run -race cmd/worker.go -i 'txt/*' -p 'cmd/wc.so' -r 4 -w 1 &
	go run -race cmd/worker.go -i 'txt/*' -p 'cmd/wc.so' -r 4 -w 2 &
	go run -race cmd/worker.go -i 'txt/*' -p 'cmd/wc.so' -r 4 -w 3 &
	go run -race cmd/worker.go -i 'txt/*' -p 'cmd/wc.so' -r 4 -w 4 &
	go run -race cmd/worker.go -i 'txt/*' -p 'cmd/wc.so' -r 4 -w 5 &
	go run -race cmd/worker.go -i 'txt/*' -p 'cmd/wc.so' -r 4 -w 6 &
	go run -race cmd/worker.go -i 'txt/*' -p 'cmd/wc.so' -r 4 -w 7 &
	go run -race cmd/worker.go -i 'txt/*' -p 'cmd/wc.so' -r 4 -w 8 
	go run -race cmd/master.go -i 'mr-out-*' -p 'cmd/merge.so' -r 1 -w 4 &
	sleep 1
	go run -race cmd/worker.go -i 'mr-out-*' -p 'cmd/merge.so' -r 1 -w 1 &
	go run -race cmd/worker.go -i 'mr-out-*' -p 'cmd/merge.so' -r 1 -w 2 &
	go run -race cmd/worker.go -i 'mr-out-*' -p 'cmd/merge.so' -r 1 -w 3 &
	go run -race cmd/worker.go -i 'mr-out-*' -p 'cmd/merge.so' -r 1 -w 4
	diff -q mr-out-0.txt ~/work/6.824/src/main/mr-out-0
	rm /dev/shm/imd-*

test_crash: clean grpc build_crash clean_port
	-rm output/*
	-rm mr-out-*
	go run -race cmd/master.go -i 'txt/*' -p 'cmd/crash.so' -r 1 -w 8 &
	sleep 2
	go run -race cmd/worker.go -i 'txt/*' -p 'cmd/crash.so' -r 1 -w 1 &
	go run -race cmd/worker.go -i 'txt/*' -p 'cmd/crash.so' -r 1 -w 2 &
	go run -race cmd/worker.go -i 'txt/*' -p 'cmd/crash.so' -r 1 -w 3 &
	go run -race cmd/worker.go -i 'txt/*' -p 'cmd/crash.so' -r 1 -w 4 &
	go run -race cmd/worker.go -i 'txt/*' -p 'cmd/crash.so' -r 1 -w 5 &
	go run -race cmd/worker.go -i 'txt/*' -p 'cmd/crash.so' -r 1 -w 6 &
	go run -race cmd/worker.go -i 'txt/*' -p 'cmd/crash.so' -r 1 -w 7 &
	go run -race cmd/worker.go -i 'txt/*' -p 'cmd/crash.so' -r 1 -w 8
	# diff mr-out-0.txt ~/work/6.824/src/main/mr-out-0
	rm /dev/shm/imd-*

clean_port :
	-./clean_port.sh
