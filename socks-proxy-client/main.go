package main

import (
	"fmt"
	"net"
	"os"

	"github.com/armon/go-socks5"
	"github.com/hashicorp/yamux"
)

func main() {
	if len(os.Args) != 2 {
		usage()
	}
	conn, err := net.Dial("tcp", os.Args[1])
	for err != nil {
		println("Dial Error: " + err.Error())
		conn, err = net.Dial("tcp", os.Args[1])
	}
	println("Server connection OK")
	cfg := yamux.DefaultConfig()
	cfg.EnableKeepAlive = false
	ses, _ := yamux.Server(conn, cfg)
	println("mux session opened")
	socksProxy(ses)
}

func socksProxy(socksSession *yamux.Session) {
	println("starting socks proxy")
	conf := &socks5.Config{}
	server, _ := socks5.New(conf)
	println("socks start ok")
	for {
		println("waiting for socks session conn")
		c, err := socksSession.Accept()
		println("session conn established")
		if err != nil {
			break
		}
		go server.ServeConn(c)
	}
}

func usage() {
	fmt.Printf("%s [server]:[port]\n", os.Args[0])
	os.Exit(1)
}
