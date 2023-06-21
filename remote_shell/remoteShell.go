package main

/**
Create a remote shell access option for your device.

On "target" or host computer. Run the app.
go run remoteShel.go

From a remote computer connect. ie.
telnet localhost 20089
**/

import (
	"io"
	"log"
	"net"
	"os/exec"
)

func main() {
	listener, err := net.Listen("tcp", ":20089")
	if err != nil {
		log.Fatalln(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	// Create a command object
	cmd := exec.Command("/bin/bash", "-i") // windows "cmd.exe"

	// Create a pipe object (as in unix) (Writer process [wp], Reader process [wp)]
	rp, wp := io.Pipe()
	// Stdin (conn) Stdout (pipe's wp)
	cmd.Stdin = conn
	cmd.Stdout = wp
	// wp to conn
	go io.Copy(conn, rp)

	cmd.Run()
	conn.Close()
}
