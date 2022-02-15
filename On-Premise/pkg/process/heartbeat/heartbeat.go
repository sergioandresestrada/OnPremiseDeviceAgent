package heartbeat

import (
	"fmt"
	"net"

	"On-Premise/pkg/types"
)

type Message = types.Message

func ProcessHeartbeat(msg Message) {
	fmt.Println("Processing Heartbeat Job")

	sendToClient(msg.Message)
}

func sendToClient(message string) {
	host := "localhost"
	port := "9999"
	conType := "tcp"

	fmt.Printf("Connecting to %s on port %s.\n", host, port)

	conn, err := net.Dial(conType, host+":"+port)

	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		panic(err)
	}

	fmt.Println("Connection established correctly")

	_, err = conn.Write([]byte(message))

	if err != nil {
		fmt.Println("Error sending message:", err.Error())
		panic(err)
	}

}
