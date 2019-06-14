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
	var mode = flag.Int("mode", 0, "Mode 0:client+server 1:server only 2:client only")

	flag.Parse()
	log.Println("Detecting host port = ", *host, ":", *port, " mode:", *mode)

	if *mode == 0 || *mode == 1 {

		go pingPongServer(*port)
	}

	if *mode == 0 || *mode == 2 {
		go CheckPortReachableRoutine(*host, *port)
	}
	for {
		log.Println("Ping Pong Result: ", peerPortReachable)
		time.Sleep(5 * time.Second)
	}

}

func pingPongServer(port string) {
	log.Println("----Running Server----")
	var listenHost string
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	log.Printf("--- %s listen on %s --- \n", listenHost, port)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting: ", err.Error())
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
		log.Println("Error reading:", err.Error())
	}
	// Send a response back to person contacting us.
	conn.Write([]byte("Pong"))
	// Close the connection when you're done with it.
	conn.Close()
}

//CheckPortReachableRoutine is a Connection Check Routine
func CheckPortReachableRoutine(host, port string) {
	log.Println("----Running Client Routine---")
	status := make(chan bool)
	log.Printf("CheckPortReachableRoutine:%s:%s\n", host, port)
	for {
		go func(updateStatus chan<- bool) {
			connected := true
			for retry := 0; retry < retryTimes; retry++ {
				connected = connToPort(host, port)
				if !connected {
					log.Printf("NOT able to connect %s:%s\n", host, port)
					time.Sleep(retryDelay)
				} else {
					//log.Printf("Connect to %s:%s\n", host, port)
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
		log.Println("Ping")
		resp := make([]byte, 1024)
		lenResp, err := conn.Read(resp)
		if err != nil {
			log.Println("connToPort Error:", err)
		}
		log.Println(string(resp[:lenResp]))
		conn.Close()
		return true
	}
}
