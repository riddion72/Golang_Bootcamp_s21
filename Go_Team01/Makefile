
PORT ?= 1234
PBD ?= 8765
server:
	cd bd && go build cmd/server.go 

serverP:
	cd bd &&./server

serverF:
	cd bd && ./server -HBD 127.0.0.1 -P $(PORT) -PBD $(PBD)

client:
	cd client && go run cmd/client.go

clean:
	cd bd && rm -rf server
	cd client && rm -rf client

# make serverF PORT=1235  