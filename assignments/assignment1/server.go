/*****************************************************************************
 * server.go
 * Name:
 * NetId:
 *****************************************************************************/

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

const RECV_BUFFER_SIZE = 2048

/* TODO: server()
 * Open socket and wait for client to connect
 * Print received message to stdout
 */
func server(server_port string) {
	listener, err := net.Listen("tcp", "127.0.0.1:"+server_port)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err.Error())
		}
		go func(c net.Conn) {
			for {
				var buff []byte = make([]byte, RECV_BUFFER_SIZE)
				_, err := conn.Read(buff)
				if err != nil {
					if err != io.EOF {
						fmt.Println("error:", err)
					}
					break
				}
				fmt.Println("From server: ")
				fmt.Println(string(buff))
			}
			c.Close()
		}(conn)
	}
}

// Main parses command-line arguments and calls server function
func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: ./server [server port]")
	}
	server_port := os.Args[1]
	server(server_port)
}
