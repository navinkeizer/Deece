package Deece

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ipfs/go-cid"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	"github.com/multiformats/go-multibase"
	"github.com/pkg/errors"
	"github.com/wealdtech/go-multicodec"
	"net"
	"strconv"
	"strings"
)

//defining the global variables used
var (
	Shell      *ipfsapi.Shell
	Client     *ethclient.Client
	TLI        string
	serverPort int
	serverIP   string
	passWord   string
)

//function to convert to base32 encoding
func b32Cid(bytes []byte) (string, error) {
	data, codec, err := multicodec.RemoveCodec(bytes)
	if err != nil {
		return "", err
	}
	codecName, err := multicodec.Name(codec)
	if err != nil {
		return "", err
	}

	if codecName == "ipfs-ns" {
		thisCID, err := cid.Parse(data)
		if err != nil {
			return "", errors.Wrap(err, "failed to parse CID")
		}
		str, err := thisCID.StringOfBase(multibase.Base32)
		if err != nil {
			return "", errors.Wrap(err, "failed to obtain base36 representation")
		}
		return fmt.Sprintf("ipfs://%s", str), nil
	}

	return "", fmt.Errorf("unknown codec name %s", codecName)

}

//simple way of ensuring we have pdf's
//this will become more sophisticated in future releases
func isValidPdf(stream []byte) bool {
	//fmt.Println(stream)

	l := len(stream)
	var isHeaderValid = stream[0] == 0x25 && stream[1] == 0x50 && stream[2] == 0x44 && stream[3] == 0x46 //%PDF

	//(.%%EOF)
	var isTrailerValid1 = stream[l-6] == 0xa && stream[l-5] == 0x25 && stream[l-4] == 0x25 &&
		stream[l-3] == 0x45 && stream[l-2] == 0x4f && stream[l-1] == 0x46
	if isHeaderValid && isTrailerValid1 {
		return true
	}

	//(.%%EOF.)
	var isTrailerValid2 = stream[l-7] == 0xa && stream[l-6] == 0x25 && stream[l-5] == 0x25 && stream[l-4] == 0x45 &&
		stream[l-3] == 0x4f && stream[l-2] == 0x46 && stream[l-1] == 0xa
	if isHeaderValid && isTrailerValid2 {
		return true
	}

	//(.%%EOF.)
	var isTrailerValid4 = stream[l-7] == 0xd && stream[l-6] == 0x25 && stream[l-5] == 0x25 && stream[l-4] == 0x45 &&
		stream[l-3] == 0x4f && stream[l-2] == 0x46 && stream[l-1] == 0xd
	if isHeaderValid && isTrailerValid4 {
		return true
	}

	//(..%%EOF..)
	var isTrailerValid3 = stream[l-8] == 0xd && stream[l-7] == 0x25 && stream[l-6] == 0x25 && stream[l-5] == 0x45 && stream[l-4] == 0x4f &&
		stream[l-3] == 0x46 && stream[l-2] == 0xd && stream[l-1] == 0xa
	if isHeaderValid && isTrailerValid3 {
		return true
	}

	return false
}

//function to get the latest ipns record from the server for the TLI
func getTLI() (string, error) {
	addr := strings.Join([]string{serverIP, strconv.Itoa(serverPort)}, ":")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(getMessage))
	_, err = conn.Write([]byte(StopCharacter))
	if err != nil {
		return "", err
	}
	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	data := string(buff[:n])
	tli := strings.Trim(data, "\r\n\r\n")
	//log.Printf("Update TLI: %s",tli)
	return tli, nil
}
