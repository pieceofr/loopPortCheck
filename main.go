package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

/* ./nodePortCheck -host=118.163.120.180 -port=2136 */
var peerPortReachable = false

const (
	retryDelay   = time.Duration(200 * time.Millisecond)
	retryTimes   = 3
	checkInterMs = 1000
	dialTimeout  = 2 * time.Second
	tcpType6     = "tcp6"
	tcpType4     = "tcp4"
)

func main() {
	var host = flag.String("host", "127.0.0.1", "enter the host ip for examine the node")
	var port = flag.String("port", "2136", "enter the host port for examine the node")

	flag.Parse()
	log.Println("Detecting host port = ", *host, ":", *port)
	go pingPongServer(*port)

	go CheckPortReachableRoutine(*host, *port)

	for {
		fmt.Println("Ping Pong Result: ", peerPortReachable)
		time.Sleep(5 * time.Second)
	}

}

func pingPongServer(port string) {
	var listenHost string

	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Printf("--- %s listen on %s --- \n", listenHost, port)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	// Handle connections in a new goroutine.
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	// Send a response back to person contacting us.
	conn.Write([]byte("Pong"))
	// Close the connection when you're done with it.
	conn.Close()
}

//CheckPortReachableRoutine is a Connection Check Routine
func CheckPortReachableRoutine(host, port string) {
	status := make(chan bool)
	fmt.Printf("CheckPortReachableRoutine:%s:%s\n", host, port)
	for {
		go func(updateStatus chan<- bool) {
			connected := true
			for retry := 0; retry < retryTimes; retry++ {
				connected = connToPort(host, port)
				if !connected {
					fmt.Printf("NOT able to connect %s:%s\n", host, port)
					time.Sleep(retryDelay)
				} else {
					//fmt.Printf("Connect to %s:%s\n", host, port)
					retry = retryTimes + 1
				}
			}
			updateStatus <- connected

		}(status)

		peerPortReachable = <-status
		time.Sleep(time.Duration(checkInterMs) * time.Millisecond)
	}
}

func connToPort(host, port string) bool {
	if host == "" {
		return false
	}
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host, port), dialTimeout)
	if err != nil {
		return false
	} else {
		conn.Write([]byte("Ping"))
		fmt.Println("Ping")
		resp := make([]byte, 1024)
		lenResp, err := conn.Read(resp)
		if err != nil {
			fmt.Println("connToPort Error:", err)
		}
		fmt.Println(string(resp[:lenResp]))
		conn.Close()
		return true
	}
}
