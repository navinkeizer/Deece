package Deece

import (
	"encoding/csv"
	"github.com/ethereum/go-ethereum/ethclient"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type IncorrrectInput struct{}

func (zz *IncorrrectInput) Error() string {
	return "Input type is not recognised."
}

func ConnectServer() (*ipfsapi.Shell, *ethclient.Client, string) {
	sh := ipfsapi.NewShell("localhost:5001")
	sh.SetTimeout(time.Duration(10000000000))
	cli, err := ethclient.Dial(infura)
	if err != nil {
		panic(err)
	}
	tli := "k2k4r8oxynrlparmnoh62lhk0ozhsdw8lizrcxjxs2w3jlllrqzi2bm8"
	return sh, cli, tli
}

func ConnectClient() (*ipfsapi.Shell, *ethclient.Client, string) {
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

//TODO: return special structure output

func DoCrawlServer(name string, t string) {

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

func DoCrawlClient(name string, t string) {

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

func DoSearchClient(searchTerms []string) {}

func DoSearchServer(searchTerms []string) {
	latestTLI, err := Shell.Resolve(TLI)
	if err != nil {
		log.Println(err)
	}
	cidTLI := strings.Split(latestTLI, "s/")[1]

	cat, err := Shell.Cat(cidTLI)
	if err != nil {
		log.Println(err)
	}

	result, err := ioutil.ReadAll(cat)
	if err != nil {
	}

	err = cat.Close()
	if err != nil {
	}

	err = ioutil.WriteFile("./TLI/TLI.csv", result, 0644)
	if err != nil {
	}

	f, err := os.Open("./TLI/TLI.csv")
	if err != nil {
	}

	csvr := csv.NewReader(f)
	records, _ := csvr.ReadAll()

	var indexLocations []string
	// need to change to better performance search mechanism
	for i := 0; i < len(searchTerms); i++ {
		indexLocations = append(indexLocations, fetchIndex(searchTerms[i], records))
	}

	err = f.Close()
	if err != nil {

	}

	err = printResultsWordServer(searchTerms, indexLocations)
	if err != nil {

	}
}
