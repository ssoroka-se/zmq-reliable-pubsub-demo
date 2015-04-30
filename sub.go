package main

// Based on: Reliable pub-sub
// http://zguide.zeromq.org/page:all#toc119

import (
	"log"
	"strings"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	// start api to listen for subscribe commands
	// expire subscribe commands after so long?
	// periodically rerank auth servers (or just recheck rankings)
	// connect to nearest auth server
	// tell it what devices we want to subscribe to

	// listen for updates

	client_public, client_secret, err := zmq.NewCurveKeypair()
	if err != nil {
		log.Println(err)
		return
	}

	//  Tell authenticator to use this public client key
	// zmq.AuthCurveAdd("domain1", client_public)

	//  Create and connect client socket
	client, err := zmq.NewSocket(zmq.DEALER)
	if err != nil {
		log.Println(err)
		return
	}
	defer client.Close()
	client.SetIdentity("Client1" + client_public)
	server_public := "um%8LQT60u7I=RA)I8M4kaq4c]xJdfK)0YiEzB!!"
	client.ClientAuthCurve(server_public, client_public, client_secret)
	err = client.Connect("tcp://127.0.0.1:9000")
	if err != nil {
		log.Println(err)
		return
	}

	//  Send a message from client to server
	_, err = client.SendMessage("subscribe", "1,2,3,4")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("sent...")
	var parts []string
	// loop forever receiving updates..
	for i := 0; i < 10; i++ {
		parts, err = client.RecvMessage(0)
		log.Println(strings.Join(parts, ", "))
	}
}
