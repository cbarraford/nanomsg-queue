package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/rep"
	"github.com/go-mangos/mangos/protocol/req"
	"github.com/go-mangos/mangos/transport/tcp"
)

type Work struct {
	Work int `json:"work"`
}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func pop(url string) {
	var sock mangos.Socket
	var err error
	var msg []byte
	if sock, err = rep.NewSocket(); err != nil {
		die("can't get new rep socket: %s", err)
	}
	sock.AddTransport(tcp.NewTransport())
	if err = sock.Listen(url); err != nil {
		die("can't listen on rep socket: %s", err.Error())
	}

	// Could also use sock.RecvMsg to get header
	msg, err = sock.Recv()
	fmt.Printf("AGENT: RECEIVED REQUEST %s\n", msg)
	fmt.Printf("AGENT: Sending Acknowledgement\n")
	err = sock.Send([]byte(url))
	if err != nil {
		die("can't send reply: %s", err.Error())
	}
	fmt.Println("Doing the work...")
	work := Work{}
	json.Unmarshal(msg, &work)
	fmt.Printf("Sleeping for %d seconds\n", work.Work)
	time.Sleep(time.Duration(work.Work) * time.Second)
	fmt.Println("Done.")
}

func todo(agents []string, msg []byte) {
	var sock mangos.Socket
	var err error

	fmt.Println("Writing todo to disk for persistence")

	if sock, err = req.NewSocket(); err != nil {
		die("can't get new req socket: %s", err.Error())
	}
	sock.AddTransport(tcp.NewTransport())
	sock.SetOption("RETRY-TIME", 5*time.Second)
	for _, agent := range agents {
		if err = sock.Dial(agent); err != nil {
			die("can't dial on req socket: %s", err.Error())
		}
	}
	fmt.Printf("Waiting for REQUEST ACK (%s)\n", msg)
	if err = sock.Send([]byte(msg)); err != nil {
		die("can't send message on push socket: %s", err.Error())
	}
	if msg, err = sock.Recv(); err != nil {
		die("can't receive msg: %s", err.Error())
	}
	fmt.Printf("RECEIVED ACK FROM %s\n", string(msg))
	fmt.Println("Removed todo from persistent disk")
	sock.Close()
}

func main() {
	agents := []string{"tcp://localhost:40897", "tcp://localhost:40898", "tcp://localhost:40899"}
	if os.Args[1] == "todo" {
		todo(agents, []byte(os.Args[2]))
	}
	if os.Args[1] == "pop" {
		pop(os.Args[2])
	}
}
