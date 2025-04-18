all: server client

server: cmd/server/main.go
	go build -o _output/ github.com/LiangNing7/go-example/go-tcp/cmd/server
client: cmd/client/main.go 
	go build -o _output/ github.com/LiangNing7/go-example/go-tcp/cmd/client

clean:
	rm -rf ./_output
