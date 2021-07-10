package main

import (
	"bufio"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

var shell *ipfsapi.Shell

const (
	TLIcid        = "k51qzi5uqu5dm51kzdsr3pu33tkrzca5pse1kt3i9a7c5m1ai0kremf5c0ooe9"
	StopCharacter = "\r\n\r\n"
	passWord      = "FsXEzxp1EVmJjSNAZh"
)

func isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, StopCharacter)
	return
}

func updateIPNS(newcid string) error {
	shell.SetTimeout(time.Duration(1000000000000))

	err := shell.Publish("", "/ipfs/"+newcid)
	if err != nil {
		log.Println(err)
		return err
	}

	shell.SetTimeout(time.Duration(10000000000))
	return nil
}

func SocketServer(port int) {

	listen, err := net.Listen("tcp4", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("Socket listen port %d failed,%s", port, err)
	}
	defer listen.Close()
	log.Printf("Begin listen port: %d", port)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handler(conn)
	}

}

func handler(conn net.Conn) {

	log.Println("Connected to: " + conn.RemoteAddr().String())
	defer conn.Close()

	var (
		buf = make([]byte, 1024)
		r   = bufio.NewReader(conn)
		w   = bufio.NewWriter(conn)
	)

	n, err := r.Read(buf)
	if err != nil {
		log.Println(err)
		return
	}

	data := string(buf[:n])
	log.Println("Request: [" + strings.Trim(data, StopCharacter) + "]")

	if isTransportOver(data) {
		request := strings.Split(data, ",")
		if request[0] == "SET" && request[1] == passWord {
			log.Println("RESPONSE: [Adding " + strings.Trim(request[2], StopCharacter) + " to IPNS]")
			_, err = w.Write([]byte("Adding " + strings.Trim(request[2], StopCharacter) + " as TLI"))
			err = w.Flush()
			if err != nil {
				log.Println(err)
			}
			err = updateIPNS(strings.Trim(request[2], StopCharacter))
			if err != nil {
				log.Println(err)
				return
			}

		} else if request[0] == "GET" {
			log.Println("RESPONSE: [" + TLIcid + "]")
			_, err = w.Write([]byte(TLIcid))
			err = w.Flush()
			if err != nil {
				log.Println(err)
			}
		} else {
			log.Println("RESPONSE: [Unknown request]")
			_, err = w.Write([]byte("Unknown request"))
			err = w.Flush()
			if err != nil {
				log.Println(err)
			}
		}
	}

	log.Println("Closed connection: " + conn.RemoteAddr().String())

}

func main() {
	shell = ipfsapi.NewShell("localhost:5001")
	port := 3333
	SocketServer(port)

}
