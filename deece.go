package Deece

import (
	"encoding/csv"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

//function to setup the local connections to ipfs, eth gateway etc.
//to be used at gateway server running the web interface
func ConnectServer(Infura string, tli string, ip string, port int) (*ipfsapi.Shell, *ethclient.Client) {

	sh := ipfsapi.NewShell("localhost:5001")
	cli, err := ethclient.Dial(Infura)
	if err != nil {
		log.Println(err)
	}
	TLI = tli
	return sh, cli
}

//function to setup the local connections to ipfs, eth gateway, gateway server address etc.
//to be used by clients using the CLI application
func ConnectClient(Infura string, tli string, ip string, port int) (*ipfsapi.Shell, *ethclient.Client) {
	sh := ipfsapi.NewShell("localhost:5001")
	cli, err := ethclient.Dial(Infura)
	if err != nil {
		log.Println(err)
	}
	serverPort = port
	serverIP = ip
	TLI, err = getTLI()
	if err != nil || TLI == "" {
		log.Println(err)
		//fallback to config file if server does not respond
		TLI = tli
	}
	return sh, cli
}

//initiates a crawl for a name and content type t
//to be used on server (and therefore performs name publishing step locally)
//files are stored locally. Future releases will use tmpdir
func DoCrawlServer(name string, t string) {

	//start crawling process
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

	//starts adding to the index by extracting keywords using OCR
	content, err := ExtractPdfDataOCR("./retrieved_files/" + id)
	if err != nil {
		log.Fatal(&pdfreadfail{"./retrieved_files/" + id + ".pdf"})
	}

	log.Println("Success extracting keywords...")

	CreateIndexEntryServer1(content, id)

	log.Println("Successful indexing.")
}

//initiates a crawl for a name and content type t
//to be used at client and uses the server to update the TLI entry on IPNS
//files are stored locally. Future releases will use tmpdir
func DoCrawlClient(name string, t string) {

	//start crawling process
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

	//starts adding to the index by extracting keywords using OCR
	content, err := ExtractPdfDataOCR("./retrieved_files/" + id)
	if err != nil {
		log.Fatal(&pdfreadfail{"./retrieved_files/" + id + ".pdf"})
	}

	log.Println("Success extracting keywords...")

	CreateIndexEntryClient1(content, id)

	log.Println("Successful indexing.")
}

////starts the search process for an array of search terms
////to be run on CLI
//func DoSearchClient(searchTerms []string) {
//	//get the latest TLI file
//	latestTLI, err := Shell.Resolve(TLI)
//	if err != nil {
//		log.Println(err)
//	}
//	cidTLI := strings.Split(latestTLI, "s/")[1]
//	//retrieve the TLI file
//	cat, err := Shell.Cat(cidTLI)
//	if err != nil {
//		log.Println(err)
//	}
//
//	result, err := ioutil.ReadAll(cat)
//	if err != nil {
//	}
//
//	err = cat.Close()
//	if err != nil {
//	}
//
//	err = ioutil.WriteFile("./TLI/TLI.csv", result, 0644)
//	if err != nil {
//	}
//
//	f, err := os.Open("./TLI/TLI.csv")
//	if err != nil {
//	}
//
//	csvr := csv.NewReader(f)
//	records, _ := csvr.ReadAll()
//
//	//retrieve the index location for each keyword
//	var indexLocations []string
//
//	// TODO:need to change to better performance search mechanism
//	for i := 0; i < len(searchTerms); i++ {
//		indexLocations = append(indexLocations, fetchIndex(searchTerms[i], records))
//	}
//
//	err = f.Close()
//	if err != nil {
//
//	}
//	//print the search results to the terminal
//	err = printResultsWordClient(searchTerms, indexLocations)
//	if err != nil {
//
//	}
//}
//
////starts the search process for an array of search terms
////to be run on gateway server
//func DoSearchServer(searchTerms []string) ([]string, [][]resultKeyword, error) {
//
//	//get the latest TLI file
//	latestTLI, err := Shell.Resolve(TLI)
//	if err != nil {
//		log.Println(err)
//	}
//	cidTLI := strings.Split(latestTLI, "s/")[1]
//	//retrieve the TLI file
//	cat, err := Shell.Cat(cidTLI)
//	if err != nil {
//		log.Println(err)
//	}
//
//	result, err := ioutil.ReadAll(cat)
//	if err != nil {
//	}
//
//	err = cat.Close()
//	if err != nil {
//	}
//
//	err = ioutil.WriteFile("./TLI/TLI.csv", result, 0644)
//	if err != nil {
//	}
//
//	f, err := os.Open("./TLI/TLI.csv")
//	if err != nil {
//	}
//
//	csvr := csv.NewReader(f)
//	records, _ := csvr.ReadAll()
//
//	//retrieve the index location for each keyword
//	var indexLocations []string
//	// need to change to better performance search mechanism
//	for i := 0; i < len(searchTerms); i++ {
//		indexLocations = append(indexLocations, fetchIndex(searchTerms[i], records))
//	}
//
//	err = f.Close()
//	if err != nil {
//
//	}
//
//	//return the results in structure to be used by web interface
//	ST, SR, err := printResultsWordServer(searchTerms, indexLocations)
//	if err != nil {
//
//	}
//	return ST, SR, nil
//}
//

func DoSearch1(query string) ([]QueryResult, error) {

	searchTerms := strings.Split(query, " ")

	//get the latest TLI file
	latestTLI, err := Shell.Resolve(TLI)
	if err != nil {
		log.Println(err)
	}

	cidTLI := strings.Split(latestTLI, "s/")[1]

	//retrieve the TLI file
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

	//retrieve the index location for each keyword
	var indexLocations []string

	for j := 0; j < len(searchTerms); j++ {
		i := sort.Search(len(records), func(i int) bool { return searchTerms[j] <= records[i][0] })
		if i < len(records) && records[i][0] == searchTerms[j] {
			indexLocations = append(indexLocations, records[i][1])
		} else {
			indexLocations = append(indexLocations, "-")
		}

	}
	err = f.Close()
	if err != nil {
		log.Println(err)
	}

	qResult, err := resultsWordServer1(searchTerms, indexLocations)
	if err != nil {

	}

	fmt.Println(qResult)

	return qResult, nil
}
