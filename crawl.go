package Deece

import (
	"github.com/ipfs/go-dnslink"
	//shell "github.com/ipfs/go-ipfs-api"
	"github.com/wealdtech/go-ens/v3"
	"io/ioutil"
	"log"
	"strings"
)

//function to retrieve a file using ens name
//uses the gateway provider specified in a config1.json file
func getFromEns(ENSdomain string) ([]byte, string, error) {
	resol, err := ens.NewResolver(Client, ENSdomain)
	if err != nil {
		log.Println(err)
		return nil, "", err
	}

	q, err := resol.Contenthash()
	if err != nil {
		log.Println(err)
		return nil, "", err
	}

	//convert to ipfs readable format
	link, err := b32Cid(q)
	if err != nil {
		log.Println(err)
		return nil, "", err
	}

	cid := strings.Split(link, "//")[1]

	file, err := getFromCid(cid)
	if err != nil {
		log.Println(err)
		return nil, cid, err
	}

	return file, cid, nil
}

//function to retrieve a file using dnslink name
// it seems like most dnslink names are not in use
func getFromDns(DNSdomain string) ([]byte, string, error) {
	link, err := dnslink.Resolve(DNSdomain)
	if err != nil {
		log.Println(err)
		return nil, "", err
	}

	id := strings.Split(link, "s/")[1]

	file, err := getFromCid(id)
	if err != nil {
		log.Println(err)
		return nil, id, err
	}

	return file, id, nil
}

//function to retrieve a file using ipns name
func getFromIpns(ipns string) ([]byte, string, error) {
	link, err := Shell.Resolve(ipns)
	if err != nil {
		log.Println(err)
		return nil, "", err
	}

	cid := strings.Split(link, "s/")[1]

	file, err := getFromCid(cid)
	if err != nil {
		log.Println(err)
		return nil, cid, err
	}

	return file, cid, nil

}

//function to retrieve a file using ipfs CID
func getFromCid(CID string) ([]byte, error) {
	cat, err := Shell.Cat(CID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	result, err := ioutil.ReadAll(cat)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = cat.Close()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if !isValidPdf(result) {
		err = &nopdf{CID}
		log.Println(err)
		return nil, err
	}

	return result, nil

}

//function to start a crawl. Takes in the name of a file and the type
func crawlInput(domaintype string, filename string) ([]byte, string, error) {

	switch domaintype {
	case "CID":
		dat, err := getFromCid(filename)
		if err != nil {
			return nil, "", &cIDmissing{filename}
		}
		return dat, filename, nil

	case "ENS":
		dat, id, err := getFromEns(filename)
		if err != nil && id == "" {
			return nil, "", &noENS{filename}
		}
		if err != nil && id != "" {
			return nil, "", &noNSresolve{id, filename}
		}
		return dat, id, nil

	case "DNS":
		dat, id, err := getFromDns(filename)
		if err != nil && id == "" {
			return nil, "", &noDNS{}
		}
		if err != nil && id != "" {
			return nil, "", &noNSresolve{id, filename}
		}
		return dat, id, nil

	case "IPNS":
		dat, id, err := getFromIpns(filename)
		if err != nil && id == "" {
			return nil, "", &noipns{filename}
		}
		if err != nil && id != "" {
			return nil, "", &noNSresolve{id, filename}
		}
		return dat, id, nil
	}
	return nil, "", &noNSType{}
}
