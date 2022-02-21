package main

import (
	"On-Premise/pkg/types"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
)

const (
	JOB_LISTEN_PORT       = "55555"
	HEARTBEAT_LISTEN_PORT = "44444"
)

func main() {

	JobListener, err := net.Listen("tcp", "localhost:"+JOB_LISTEN_PORT)
	if err != nil {
		fmt.Println("Error while creating the socket")
	}

	go receiveJob(JobListener)

	HeartbeatListener, err := net.Listen("tcp", "localhost:"+HEARTBEAT_LISTEN_PORT)
	if err != nil {
		fmt.Println("Error while creating the socket")
	}

	receiveHB(HeartbeatListener)

}

func receiveJob(listener net.Listener) {
	for {
		var Job types.JobClient
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error while accepting the connection")
			return
		}

		buffer := make([]byte, 256)
		n, err := conn.Read(buffer)
		buffer = bytes.Trim(buffer, "\x00")

		fmt.Printf("Received %v bytes\n", n)
		if err != nil {
			fmt.Println("Error while reading from the socket")
			return
		}
		json.Unmarshal(buffer, &Job)
		fmt.Printf("Job: %v\n", Job)

		fd, err := os.Create("clientFiles/" + Job.FileName)
		if err != nil {
			fmt.Println("Error while creating the file")
			return
		}
		defer fd.Close()

		// send a single byte buffer for sync
		buffer = make([]byte, 1)
		conn.Write(buffer)

		nBytes, err := io.Copy(fd, conn)
		if err != nil {
			fmt.Println("Error while receiving the file")
			return
		}
		fmt.Printf("Received file %s, length: %v bytes\n", Job.FileName, nBytes)
		conn.Close()
	}
}

func receiveHB(listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error while accepting the connection")
			return
		}

		buffer, err := ioutil.ReadAll(conn)
		receivedMessage := string(buffer)
		fmt.Printf("Received Heartbeat: %s \n", receivedMessage)
		conn.Close()
	}

}
