package Deece

import (
	"github.com/ipfs/go-dnslink"
	"github.com/wealdtech/go-ens"
	"io/ioutil"
	"log"
	"strings"
)

//use go-ens to get CID's
func getFromEns(ENSdomain string) ([]byte, string, error) {

	resol, err := ens.NewResolver(client, ENSdomain)
	if err != nil {
		log.Println(err)
		return nil, "", err
	}

	q, err := resol.Contenthash()
	if err != nil {
		log.Println(err)
		return nil, "", err
	}

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

// it seems like most dnslinks do not work
//example which does work: "originprotocol.com"
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

func getFromIpns(ipns string) ([]byte, string, error) {
	link, err := shell.Resolve(ipns)
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

func getFromCid(CID string) ([]byte, error) {
	cat, err := shell.Cat(CID)
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

	if !IsValidPdf(result) {
		err = &nopdf{CID}
		log.Println(err)
		return nil, err
	}

	return result, nil

}

func crawlInput(domaintype string, filename string) ([]byte, string, error) {

	switch domaintype {
	case "CID":
		dat, err := getFromCid(filename)
		if err != nil {
			return nil, "", &CIDmissing{filename}
		}
		return dat, filename, nil

	case "ENS":
		dat, id, err := getFromEns(filename)
		if err != nil && id == "" {
			return nil, "", &NoENS{filename}
		}
		if err != nil && id != "" {
			return nil, "", &NoNSresolve{id, filename}
		}
		return dat, id, nil

	case "DNS":
		dat, id, err := getFromDns(filename)
		if err != nil && id == "" {
			return nil, "", &NoDNS{}
		}
		if err != nil && id != "" {
			return nil, "", &NoNSresolve{id, filename}
		}
		return dat, id, nil

	case "IPNS":
		dat, id, err := getFromIpns(filename)
		if err != nil && id == "" {
			return nil, "", &Noipns{filename}
		}
		if err != nil && id != "" {
			return nil, "", &NoNSresolve{id, filename}
		}
		return dat, id, nil
	}
	return nil, "", &NoNSType{}
}

func doCrawlServer(name string, t string) {

	//CRAWLER
	d, id, err := crawlInput(t, name)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Success retrieving file...")

	err = ioutil.WriteFile("./retrieved_files/"+id+".pdf", d, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Success saving file locally...")

	//ADD TO INDEX
	content, err := extractPdfDataOCR("./retrieved_files/" + id)
	if err != nil {
		log.Fatal(&pdfreadfail{"./retrieved_files/" + id + ".pdf"})
	}
	createIndexEntryServer(content, id)

	log.Println("Successful indexing.")
}

func doCrawlClient(name string, t string) {

	//CRAWLER
	d, id, err := crawlInput(t, name)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Success retrieving file...")

	err = ioutil.WriteFile("./retrieved_files/"+id+".pdf", d, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Success saving file locally...")

	//ADD TO INDEX
	content, err := extractPdfDataOCR("./retrieved_files/" + id)
	if err != nil {
		log.Fatal(&pdfreadfail{"./retrieved_files/" + id + ".pdf"})
	}
	createIndexEntryClient(content, id)

	log.Println("Successful indexing.")
}
