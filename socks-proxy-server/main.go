package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"

	"github.com/hashicorp/yamux"
)

func main() {
	if len(os.Args) != 3 {
		usage()
	}
	startAddr, err := strconv.Atoi(os.Args[2])
	if err != nil {
		usage()
	}
	l, err := net.Listen("tcp", ":"+os.Args[1])
	if err != nil {
		usage()
	}
	for {
		c, _ := l.Accept()
		go handleSocks(strconv.Itoa(startAddr), c)
		startAddr++
	}

}

func handleSocks(lPort string, conn net.Conn) {
	println("Entered socks handling")
	l, _ := net.Listen("tcp", "127.0.0.1:"+lPort)
	println("Local listener opened at port " + lPort)
	cfg := yamux.DefaultConfig()
	cfg.EnableKeepAlive = false
	socksChannel, _ := yamux.Client(conn, cfg)
	println("Socks channel open OK, entering loop")
	for {
		println("waiting for a local connection to server")
		localConn, err := l.Accept()
		if err != nil {
			println("Accept error: " + err.Error())
			return
		}
		println("local connection established")
		nestedStream, err := socksChannel.Open()
		if err != nil {
			println("MUX open error: " + err.Error())
			return
		}
		println("channel opened")
		go io.Copy(localConn, nestedStream)
		go io.Copy(nestedStream, localConn)
		println("i/o OK")
	}
}

func usage() {
	fmt.Printf("%s [listen port] [first socks port]\n", os.Args[0])
	os.Exit(1)
}
