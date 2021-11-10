package P2pRTC

import (
	"fmt"
	"testing"
)

/**
Js method for test

function PostData(url, str) {
    return new Promise((ok, error) => {
        const xhr = new XMLHttpRequest();
        xhr.onload = () => {
            if (xhr.readyState === 4 && xhr.status === 200) {
                ok(xhr.responseText);
            } else {
                error(xhr.status + ":" + xhr.statusText)
            }
        }
        xhr.open('POST', url, true);
        xhr.setRequestHeader('Content-Type', 'text/plain');
        xhr.send(str);
    })
}

 */


func Test1_go_connect_go(t *testing.T) {

	//client, _ := newP2pSocket(false,false,"","")
	//server, _ := newP2pSocket(true,false,"","")

	client, _ := NewClient(false,"","")
	server, _ := NewServer(false,"","")

	server.OnSignal = func(message string) {
		fmt.Printf("Signal for Client: [%s] \n", message)
		client.Signal(message)
	}
	server.OnOpen = func(connection *DataConnection){
		connection.SendText("\tsend from server")
	}
	server.OnMessage= func(connection *DataConnection, msg []byte){
		fmt.Printf("\tserver get message: [%s] \n", msg)
	}
	server.OnError= func(connection *DataConnection,err []byte){
		fmt.Printf("\tserver OnError [%s] \n", err)
	}
	server.OnClose= func(connection *DataConnection){
		fmt.Printf("\tserver OnClose\n")
	}


	client.OnSignal = func(message string) {
		fmt.Printf("Signal for Server: [%s] \n", message)
		server.Signal(message)
	}
	client.OnOpen = func(connection *DataConnection){
		connection.SendText("I am client")
	}
	client.OnMessage= func(connection *DataConnection,msg []byte){
		connection.SendText("send from client")
		fmt.Printf("client get message: [%s] \n", msg)
	}
	client.OnError= func(connection *DataConnection,err []byte){
		fmt.Printf("client OnError [%s]", err)
	}
	client.OnClose= func(connection *DataConnection){
		fmt.Printf("client OnClose\n")
	}

	client.Connect()

	select {}

}
