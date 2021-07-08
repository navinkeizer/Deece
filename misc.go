package Deece

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ipfs/go-cid"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	"github.com/multiformats/go-multibase"
	"github.com/pkg/errors"
	"github.com/wealdtech/go-multicodec"
	"log"
	"time"
)

//TODO: fetch this from config file, as well as tli
var infura string = "https://mainnet.infura.io/v3/5c04e573d61b4e5a8fc0f3312becfdbc"

var (
	shell  *ipfsapi.Shell
	client *ethclient.Client
	TLI    string
)

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
//need more sophisticated way in future, e.g. directly getting the type of the file from ipfs
func IsValidPdf(stream []byte) bool {
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

func connectServer() (*ipfsapi.Shell, *ethclient.Client, string) {
	sh := ipfsapi.NewShell("localhost:5001")
	sh.SetTimeout(time.Duration(10000000000))
	cli, err := ethclient.Dial(infura)
	if err != nil {
		panic(err)
	}
	tli := "k2k4r8oxynrlparmnoh62lhk0ozhsdw8lizrcxjxs2w3jlllrqzi2bm8"
	return sh, cli, tli
}

func connectClient() (*ipfsapi.Shell, *ethclient.Client, string) {
	sh := ipfsapi.NewShell("localhost:5001")
	sh.SetTimeout(time.Duration(10000000000))
	cli, err := ethclient.Dial(infura)
	if err != nil {
		panic(err)
	}
	tli, err := getTLI()
	if err != nil || tli == "" {
		log.Println(err)
		tli = "k2k4r8oxynrlparmnoh62lhk0ozhsdw8lizrcxjxs2w3jlllrqzi2bm8"
	}
	return sh, cli, tli
}
