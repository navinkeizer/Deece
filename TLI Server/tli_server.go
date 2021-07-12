package main

import (
	"bufio"
	"encoding/csv"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var shell *ipfsapi.Shell
var tli string

const (
	TLIcid        = "k51qzi5uqu5dm51kzdsr3pu33tkrzca5pse1kt3i9a7c5m1ai0kremf5c0ooe9"
	StopCharacter = "\r\n\r\n"
	passWord      = "FsXEzxp1EVmJjSNAZh"
)

func SetTLIDirectory() error {
	if _, err := os.Stat("./TLI"); os.IsNotExist(err) {
		err := os.Mkdir("./TLI", 0700)
		if err != nil {
			return err
		}
	}
	return nil
}

func pinAll(newcid string) {

	cat, err := shell.Cat(newcid)
	if err != nil {
		log.Println(err)
	}

	result, err := ioutil.ReadAll(cat)
	if err != nil {
		log.Println(err)
	}

	err = cat.Close()
	if err != nil {
		log.Println(err)
	}

	err = ioutil.WriteFile("./TLI/temp.csv", result, 0644)
	if err != nil {
		log.Println(err)
	}

	f, err := os.Open("./TLI/temp.csv")
	if err != nil {
		log.Println(err)
	}

	reader := csv.NewReader(f)
	records, _ := reader.ReadAll()

	for i := 0; i < len(records); i++ {
		err := shell.Pin(records[i][1])
		if err != nil {
			log.Println(err)
		}
	}
	_ = os.Truncate("./TLI/temp.csv", 0)
}

func setup() error {
	shell.SetTimeout(time.Duration(1000000000000))
	t, err := shell.Resolve(TLIcid)
	tli = t
	err = shell.Pin(tli)
	if err != nil {
		return err
	}
	return nil
}

func isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, StopCharacter)
	return
}

func updateIPNS(newcid string) error {

	err := shell.Publish("", "/ipfs/"+newcid)
	if err != nil {
		return err
	}

	err = shell.Pin(newcid)
	err = shell.Unpin(tli)
	if err != nil {
		return err
	}
	go pinAll(newcid)
	tli = newcid
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
		if request[0] == "SET" {
			if request[1] == passWord {
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
			} else {
				log.Println("RESPONSE: [No password supplied]")
				_, err = w.Write([]byte("Please add the password in the config file, or request one the authors"))
				err = w.Flush()
				if err != nil {
					log.Println(err)
				}
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

	err := SetTLIDirectory()
	if err != nil {
		log.Println(err)
	}
	shell = ipfsapi.NewShell("localhost:5001")

	err = setup()
	if err != nil {
		log.Println(err)
	}
	SocketServer(3333)
}
