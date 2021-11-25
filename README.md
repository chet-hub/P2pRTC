# p2pRTC

### Golang usage

```golang

	client, _ := NewClient(false,"","")
	server, _ := NewServer(false,"","")

	server.OnSignal = func(message string) {
		fmt.Printf("Signal for Client: [%s] \n", message)
		client.Signal(message)
	}
	server.OnOpen = func(connection *P2pSocket){
		connection.Send([]byte("\tsend from server"))
	}
	server.OnMessage= func(connection *P2pSocket, msg []byte){
		fmt.Printf("\tserver get message: [%s] \n", msg)
	}
	server.OnError= func(connection *P2pSocket,err []byte){
		fmt.Printf("\tserver OnError [%s] \n", err)
	}
	server.OnClose= func(connection *P2pSocket){
		fmt.Printf("\tserver OnClose\n")
	}


	client.OnSignal = func(message string) {
		fmt.Printf("Signal for Server: [%s] \n", message)
		server.Signal(message)
	}
	client.OnOpen = func(connection *P2pSocket){
		connection.Send([]byte("I am client"))
	}
	client.OnMessage= func(connection *P2pSocket,msg []byte){
		connection.Send([]byte("send from client"))
		fmt.Printf("client get message: [%s] \n", msg)
	}
	client.OnError= func(connection *P2pSocket,err []byte){
		fmt.Printf("client OnError [%s]", err)
	}
	client.OnClose= func(connection *P2pSocket){
		fmt.Printf("client OnClose\n")
	}

	client.Connect()

	select {}
  
```


### Javascript usage

```js

 let client = P2pRTC.NewClient()
    let server = P2pRTC.NewServer()

    server.OnSignal = function (message) {
        console.log("server OnSignal: " + message)
        client.Signal(message)
    }
    server.OnOpen = function (connection) {
        connection.Send("I am server")
    }
    server.OnMessage = function (connection, msg) {
        console.log("server get message: " + msg)
    }
    server.OnError = function (err) {
        console.log(err)
    }
    server.OnClose = function (connection) {
        console.log("server closed")
    }

    client.OnSignal = function (message) {
        console.log("client OnSignal: " + message)
        server.Signal(message)
    }
    client.OnOpen = function (connection) {
        connection.Send("I am client")
    }
    client.OnMessage = function (connection, msg) {
        console.log("client get message: " + msg)
        connection.Send("from client echo: "+ msg)
    }
    client.OnError = function (connection,err) {
        console.log(err)
    }
    client.OnClose = function (connection) {
        console.log()
    }


### clang version - by compile using cgo

```shell
cd examples/go2c
go build -o webrtc.so -buildmode=c-shared bridge.go
```


