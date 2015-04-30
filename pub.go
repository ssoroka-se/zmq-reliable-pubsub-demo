package main

// Based on: Reliable pub-sub
// http://zguide.zeromq.org/page:all#toc119
import (
	"fmt"
	"log"
	"strings"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	zmq.AuthSetVerbose(true)

	//  Start authentication engine
	err := zmq.AuthStart()
	if err != nil {
		log.Println(err)
		return
	}
	defer zmq.AuthStop()

	zmq.AuthSetMetadataHandler(func(version, request_id, domain, address, identity, mechanism string, credentials ...string) (metadata map[string]string) {
		return map[string]string{
			"identity":    identity,
			"request_id":  request_id,
			"domain":      domain,
			"address":     address,
			"mechanism":   mechanism,
			"credentials": strings.Join(credentials, ","),
		}
	})

	zmq.AuthAllow("domain1", "127.0.0.1")
	// zmq.AuthCurveAdd("domain1", client_public)
	zmq.AuthCurveAdd("domain1", zmq.CURVE_ALLOW_ANY)

	//  We need two certificates, one for the client and one for
	//  the server. The client must know the server's public key
	//  to make a CURVE connection.
	// server_public, server_secret, err := zmq.NewCurveKeypair()
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// fmt.Println(server_public)
	// fmt.Println(server_secret)
	// server_public := "um%8LQT60u7I=RA)I8M4kaq4c]xJdfK)0YiEzB!!"
	server_secret := "YNqHnkTTi]Ygh08X5%X0fJs!0C*BJ9huqTiF0J^$"

	//  Create and bind server socket
	server, err := zmq.NewSocket(zmq.ROUTER)
	if err != nil {
		log.Println(err)
		return
	}
	defer server.Close()
	server.SetIdentity("AuthServer1")
	server.SetRouterHandover(true)
	server.SetRouterMandatory(1)
	server.ServerAuthCurve("domain1", server_secret)
	err = server.Bind("tcp://*:9000")
	if err != nil {
		log.Println(err)
		return
	}

	for {
		parts, err := server.RecvMessage(0)
		if err != nil {
			log.Println(err)
			return
		}
		client := parts[0]
		message := parts[1:]
		fmt.Printf("%v\n", message)
		for i := 0; i < 9; i++ {
			// we can route messages to specific clients interested in them
			server.SendMessageDontwait(client, fmt.Sprintf("Reply %d!", i))
		}
	}
}
