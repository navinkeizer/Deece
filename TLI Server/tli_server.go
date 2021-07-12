package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var shell *ipfsapi.Shell
var tli string
var configuration Configuration

const StopCharacter = "\r\n\r\n"

type Configuration struct {
	PassWord string
	TLIcid   string
}

func setConfig() error {
	file, err := os.Open("./config1.json")
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		return err
	}
	return nil
}

//function to check if TLI directory exists, otherwise create one
func SetTLIDirectory() error {
	if _, err := os.Stat("./TLI"); os.IsNotExist(err) {
		err := os.Mkdir("./TLI", 0700)
		if err != nil {
			return err
		}
	}
	return nil
}

//constantly adds externally added files to pins on server
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
	log.Println("Finished pinning files")
}

//sets up the tli and starts pinning
func setup() error {
	//shell.SetTimeout(time.Duration(1000000000000))
	t, err := shell.Resolve(configuration.TLIcid)
	tli = t
	err = shell.Pin(tli)
	go pinAll(tli)
	if err != nil {
		return err
	}
	return nil
}

//function to check if traffic from client is over
func isTransportOver(data string) (over bool) {
	over = strings.HasSuffix(data, StopCharacter)
	return
}

//updates the ipns record and starts pinning new file
func updateIPNS(newcid string) error {

	err := shell.Publish("", "/ipfs/"+newcid)
	if err != nil {
		return err
	}
	log.Println("Success adding to IPNS")
	err = shell.Pin(newcid)
	err = shell.Unpin(tli)
	if err != nil {
		return err
	}
	go pinAll(newcid)
	tli = newcid
	return nil
}

//sets up socket
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

//connects to client and returns results
//depending on the message
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
			if request[1] == configuration.PassWord {
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
			log.Println("RESPONSE: [" + configuration.TLIcid + "]")
			_, err = w.Write([]byte(configuration.TLIcid))
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
	err := setConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = SetTLIDirectory()
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
