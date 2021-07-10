package Deece

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/bbalet/stopwords"
	"github.com/gen2brain/go-fitz"
	"github.com/otiai10/gosseract"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	getMessage    = "GET,x"
	StopCharacter = "\r\n\r\n"
)

//old functions v1
//-------------------------------------------------------------------------------------

////function for converting simple pdf's to text
//func extractPdfData(filename string) ([]string, error) {
//	filename = filename + ".pdf"
//	file, reader, err := pdf.Open(filename)
//	defer func() {
//		_ = file.Close()
//	}()
//	if err != nil {
//		log.Println(err)
//		return nil, err
//	}
//
//	var buf bytes.Buffer
//	b, err := reader.GetPlainText()
//	if err != nil {
//		log.Println(err)
//		return nil, err
//	}
//	_, err = buf.ReadFrom(b)
//	if err != nil {
//		log.Println(err)
//		return nil, err
//	}
//	v := stopwords.CleanString(buf.String(), "en", false)
//	s := strings.Split(v, " ")
//
//	return s, nil
//
//}
//
//
////function to create a NR entry with name registry metadata added
//func createIndexEntryNR(data []string, domain string, cid string) {
//
//	//index all words in data
//	for _, s := range data {
//
//		//remove empty or 1 letter entries
//		if s == "" || len(s) == 1 {
//			continue
//		}
//		exist, err := checkExists(s)
//		if err != nil {
//			log.Println(err)
//			continue
//		}
//
//		if exist {
//
//			f, err := os.Open("./test_index/" + s + ".csv")
//			if err != nil {
//				log.Println(err)
//			}
//
//			//	check if already on the page
//			reader := csv.NewReader(f)
//			records, _ := reader.ReadAll()
//			entryExist := false
//
//			for i := 0; i < len(records); i++ {
//				if records[i][0] == cid {
//					entryExist = true
//				}
//			}
//			err = f.Close()
//			if err != nil {
//				log.Println(err)
//			}
//
//			if entryExist {
//				continue
//
//			} else {
//				//fmt.Println("entry DOES NOT exists in index")
//
//				cont := "false"
//				if strings.Contains(domain, s) {
//					cont = "true"
//				}
//				var entry = []string{cid, cont}
//
//				file, err := os.OpenFile("./test_index/"+s+".csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//				if err != nil {
//					log.Println(err)
//					continue
//				}
//				writer := csv.NewWriter(file)
//				err = writer.Write(entry)
//				if err != nil {
//					log.Println(err)
//					continue
//				}
//				writer.Flush()
//				err = file.Close()
//				if err != nil {
//					log.Println(err)
//				}
//			}
//
//		} else {
//			//fmt.Println("no exists")
//			f, err := os.Create("./test_index/" + s + ".csv")
//			if err != nil {
//				log.Println(err)
//				continue
//			}
//			//here need to add count of times in entry
//			cont := "false"
//			if strings.Contains(domain, s) {
//				cont = "true"
//			}
//			var entry = [][]string{{cid, cont}}
//
//			writer := csv.NewWriter(f)
//
//			err = writer.WriteAll(entry)
//			if err != nil {
//				log.Println(err)
//				continue
//			}
//			//writer.Flush()
//			err = f.Close()
//			if err != nil {
//				log.Println(err)
//			}
//		}
//
//	}
//
//	//	check if exist in database (local for now, later using server or NR)
//
//	//	if exist add to record, otherwise create record and update database
//
//}
//
//
//func checkExists(name string) (bool, error) {
//	if _, err := os.Stat("./test_index/" + name + ".csv"); err == nil {
//		return true, nil
//	} else if os.IsNotExist(err) {
//		return false, nil
//	} else {
//		return false, &existCheckFail{name}
//	}
//}

//old functions v2
//-------------------------------------------------------------------------------------

//
//func updateTLIServer(entries [][]string) {
//
//	//fmt.Println("starting update")
//
//	k, err := os.Open("./TLI/TLI.csv")
//	if err != nil {
//	}
//	reader := csv.NewReader(k)
//	records, err := reader.ReadAll()
//
//	for i := 0; i < len(entries); i++ {
//		records = ipnsEntryAdd(entries[i], records)
//		if err != nil {
//		}
//	}
//
//	err = k.Close()
//	if err != nil {
//		log.Println(err)
//	}
//
//	_ = os.Truncate("./TLI/TLI.csv", 0)
//
//	f, err := os.OpenFile("./TLI/TLI.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//	writer := csv.NewWriter(f)
//	_ = writer.WriteAll(records)
//	writer.Flush()
//
//	err = f.Close()
//	if err != nil {
//		log.Println(err)
//	}
//
//	z, err := os.OpenFile("./TLI/TLI.csv", os.O_RDONLY, 0644)
//	cid, err := Shell.Add(z)
//	if err != nil {
//		log.Println(err)
//	}
//
//	Shell.SetTimeout(time.Duration(1000000000000))
//	err = Shell.Publish("", "/ipfs/"+cid)
//	if err != nil {
//		log.Println(err)
//	}
//	Shell.SetTimeout(time.Duration(10000000000))
//
//}
//
//func updateTLIClient(entries [][]string) {
//
//	//fmt.Println("starting update")
//
//	k, err := os.Open("./TLI/TLI.csv")
//	if err != nil {
//	}
//	reader := csv.NewReader(k)
//	records, err := reader.ReadAll()
//
//	for i := 0; i < len(entries); i++ {
//		records = ipnsEntryAdd(entries[i], records)
//		if err != nil {
//		}
//	}
//
//	err = k.Close()
//	if err != nil {
//		log.Println(err)
//	}
//
//	_ = os.Truncate("./TLI/TLI.csv", 0)
//
//	f, err := os.OpenFile("./TLI/TLI.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//	writer := csv.NewWriter(f)
//	_ = writer.WriteAll(records)
//	writer.Flush()
//
//	err = f.Close()
//	if err != nil {
//		log.Println(err)
//	}
//
//	z, err := os.OpenFile("./TLI/TLI.csv", os.O_RDONLY, 0644)
//	cid, err := Shell.Add(z)
//	if err != nil {
//		log.Println(err)
//	}
//
//	//fmt.Println(cid)
//	//err = z.Close()
//	//if err != nil {
//	//	panic(err)
//	//}
//
//	//from here will run server side
//	err = serverTLIUpdate(cid)
//	if err != nil {
//		log.Println(err)
//	}
//}
//
//func checkTli(entry string) (bool, string, error) {
//
//	latestTLI, err := Shell.Resolve(TLI)
//	if err != nil {
//		log.Println(err)
//		return false, "", &noipns{TLI}
//	}
//	cidTLI := strings.Split(latestTLI, "s/")[1]
//
//	cat, err := Shell.Cat(cidTLI)
//	if err != nil {
//		log.Println(err)
//		return false, "", &cIDmissing{cidTLI}
//	}
//
//	result, err := ioutil.ReadAll(cat)
//	if err != nil {
//		return false, "", err
//	}
//
//	err = cat.Close()
//	if err != nil {
//		return false, "", err
//	}
//
//	err = ioutil.WriteFile("./TLI/TLI.csv", result, 0644)
//	if err != nil {
//		return false, "", err
//	}
//
//	f, err := os.Open("./TLI/TLI.csv")
//	if err != nil {
//		return false, "", err
//	}
//
//	csvr := csv.NewReader(f)
//	records, _ := csvr.ReadAll()
//	entryExist := false
//	IndexFileCid := ""
//
//	for i := 0; i < len(records); i++ {
//		if records[i][0] == entry {
//			entryExist = true
//			IndexFileCid = records[i][1]
//		}
//	}
//
//	err = f.Close()
//	if err != nil {
//		return false, "", err
//	}
//
//	return entryExist, IndexFileCid, nil
//}
//
//func ipnsEntryAdd(entry []string, records [][]string) [][]string {
//	//fmt.Println(entry)
//	//fmt.Println(records)
//
//	if records == nil {
//		//fmt.Println("empty")
//		records = append(records, entry)
//		return records
//	}
//
//	for i := 0; i < len(records); i++ {
//		//fmt.Println("comparing")
//		//fmt.Println(entry[0])
//		//fmt.Println(records[i][0])
//		//fmt.Println()
//
//		if entry[0] == records[i][0] {
//			records[i][1] = entry[1]
//			//fmt.Println("return" ,records)
//			return records
//		}
//	}
//
//	//fmt.Println("not in TLI")
//	records = append(records, entry)
//
//	//fmt.Println("return",records)
//	return records
//
//}
//
//func createIndexEntryServer(data []string, cid string) {
//
//	var ipnsEntries [][]string
//
//	for _, s := range data {
//		//remove empty or 1 letter entries
//		if s == "" || len(s) == 1 {
//			continue
//		}
//
//		//check if there is an index file available
//		exist, indexCID, err := checkTli(s)
//		if err != nil {
//			log.Println(err)
//			continue
//		}
//
//		// if index file is available add to it
//		// otherwise create one
//		if exist {
//			cat, err := Shell.Cat(indexCID)
//			if err != nil {
//
//			}
//			result, err := ioutil.ReadAll(cat)
//			if err != nil {
//
//			}
//
//			err = cat.Close()
//			if err != nil {
//
//			}
//
//			err = ioutil.WriteFile("./test_index/"+s+".csv", result, 0644)
//			if err != nil {
//				log.Println(err)
//			}
//
//			f, err := os.Open("./test_index/" + s + ".csv")
//			if err != nil {
//				log.Println(err)
//			}
//
//			//	check if already on the page
//			reader := csv.NewReader(f)
//			records, _ := reader.ReadAll()
//			entryExist := false
//
//			for i := 0; i < len(records); i++ {
//				if records[i][0] == cid {
//					entryExist = true
//				}
//			}
//
//			err = f.Close()
//			if err != nil {
//				log.Println(err)
//			}
//
//			if entryExist {
//				//fmt.Println("entry exists in the index")
//				//continue
//
//			} else {
//				//fmt.Println(" entry DOES NOT exists in index")
//				var entry = []string{
//					cid, "false",
//				}
//
//				file, err := os.OpenFile("./test_index/"+s+".csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//				if err != nil {
//					log.Println(err)
//					continue
//				}
//				writer := csv.NewWriter(file)
//
//				err = writer.Write(entry)
//				if err != nil {
//					log.Println(err)
//					continue
//				}
//				writer.Flush()
//				err = file.Close()
//				if err != nil {
//					log.Println(err)
//				}
//
//				//add file to ipfs
//				k, err := os.OpenFile("./test_index/"+s+".csv", os.O_RDONLY, 0644)
//				if err != nil {
//				}
//
//				id, err := Shell.Add(k)
//				if err != nil {
//				}
//
//				//fmt.Println("added to ipfs")
//
//				e := []string{s, id}
//				ipnsEntries = ipnsEntryAdd(e, ipnsEntries)
//
//			}
//
//		} else {
//			//fmt.Println(s + " does not exist in TLI.")
//
//			f, err := os.Create("./test_index/" + s + ".csv")
//			if err != nil {
//				log.Println(err)
//				continue
//			}
//
//			//here need to add count of times in entry
//			var entry = [][]string{
//				{cid, "false"},
//			}
//
//			writer := csv.NewWriter(f)
//			err = writer.WriteAll(entry)
//			if err != nil {
//				log.Println(err)
//				continue
//			}
//			err = f.Close()
//			if err != nil {
//				log.Println(err)
//			}
//
//			//add file to ipfs
//			k, err := os.OpenFile("./test_index/"+s+".csv", os.O_RDONLY, 0644)
//			if err != nil {
//			}
//
//			id, err := Shell.Add(k)
//			if err != nil {
//			}
//
//			//fmt.Println("added to ipfs")
//
//			e := []string{s, id}
//			ipnsEntries = ipnsEntryAdd(e, ipnsEntries)
//
//		}
//
//	}
//
//	//add ipns name entries to TLI
//
//	//fmt.Println(ipnsEntries)
//	if ipnsEntries != nil {
//		//fmt.Println("adding to TLI")
//
//		updateTLIServer(ipnsEntries)
//	}
//
//}
//
//func createIndexEntryClient(data []string, cid string) {
//
//	var ipnsEntries [][]string
//
//	//index all words in data
//	for _, s := range data {
//		//fmt.Println(s)
//		//remove empty or 1 letter entries
//		fmt.Println(s, time.Second)
//		if s == "" || len(s) == 1 {
//			continue
//		}
//
//		//check if there is an index file available
//		exist, indexCID, err := checkTli(s)
//		if err != nil {
//			log.Println(err)
//			continue
//		}
//
//		//fmt.Println(indexCID)
//
//		// if index file is available add to it
//		// otherwise create one
//		if exist {
//
//			//fmt.Println(s + " exists in TLI at " + indexCID)
//
//			cat, err := Shell.Cat(indexCID)
//			if err != nil {
//
//			}
//			result, err := ioutil.ReadAll(cat)
//			if err != nil {
//
//			}
//
//			err = cat.Close()
//			if err != nil {
//
//			}
//
//			err = ioutil.WriteFile("./test_index/"+s+".csv", result, 0644)
//			if err != nil {
//				log.Println(err)
//			}
//
//			f, err := os.Open("./test_index/" + s + ".csv")
//			if err != nil {
//				log.Println(err)
//			}
//
//			//	check if already on the page
//			reader := csv.NewReader(f)
//			records, _ := reader.ReadAll()
//			entryExist := false
//
//			for i := 0; i < len(records); i++ {
//				if records[i][0] == cid {
//					entryExist = true
//				}
//			}
//
//			err = f.Close()
//			if err != nil {
//				log.Println(err)
//			}
//
//			if entryExist {
//				//fmt.Println("entry exists in the index")
//				//continue
//
//			} else {
//				//fmt.Println(" entry DOES NOT exists in index")
//				var entry = []string{
//					cid, "false",
//				}
//
//				file, err := os.OpenFile("./test_index/"+s+".csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//				if err != nil {
//					log.Println(err)
//					continue
//				}
//				writer := csv.NewWriter(file)
//
//				err = writer.Write(entry)
//				if err != nil {
//					log.Println(err)
//					continue
//				}
//				writer.Flush()
//				err = file.Close()
//				if err != nil {
//					log.Println(err)
//				}
//
//				//add file to ipfs
//				k, err := os.OpenFile("./test_index/"+s+".csv", os.O_RDONLY, 0644)
//				if err != nil {
//				}
//
//				id, err := Shell.Add(k)
//				if err != nil {
//				}
//
//				//fmt.Println("added to ipfs")
//
//				e := []string{s, id}
//				ipnsEntries = ipnsEntryAdd(e, ipnsEntries)
//
//			}
//
//		} else {
//			//fmt.Println(s + " does not exist in TLI.")
//
//			f, err := os.Create("./test_index/" + s + ".csv")
//			if err != nil {
//				log.Println(err)
//				continue
//			}
//
//			//here need to add count of times in entry
//			var entry = [][]string{
//				{cid, "false"},
//			}
//
//			writer := csv.NewWriter(f)
//			err = writer.WriteAll(entry)
//			if err != nil {
//				log.Println(err)
//				continue
//			}
//			err = f.Close()
//			if err != nil {
//				log.Println(err)
//			}
//
//			//add file to ipfs
//			k, err := os.OpenFile("./test_index/"+s+".csv", os.O_RDONLY, 0644)
//			if err != nil {
//			}
//
//			id, err := Shell.Add(k)
//			if err != nil {
//			}
//
//			//fmt.Println("added to ipfs")
//
//			e := []string{s, id}
//			ipnsEntries = ipnsEntryAdd(e, ipnsEntries)
//
//		}
//
//	}
//
//	//add ipns name entries to TLI
//
//	//fmt.Println(ipnsEntries)
//	if ipnsEntries != nil {
//		//fmt.Println("adding to TLI")
//		fmt.Println("check3", time.Second)
//
//		updateTLIClient(ipnsEntries)
//	}
//
//}

//latest functions
//-------------------------------------------------------------------------------------

//todo: standardise error handling and defer close files

//update the ipns record through the TLI server
func serverTLIUpdate(newcid string) error {

	addr := strings.Join([]string{serverIP, strconv.Itoa(serverPort)}, ":")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write([]byte(setMessage(newcid)))
	_, err = conn.Write([]byte(StopCharacter))
	if err != nil {
		return err
	}
	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		return err
	}
	log.Println(string(buff[:n]) + "...")

	return nil
}

//set message for server TLI update
func setMessage(cid string) string {
	return "SET," + cid
}

//function for extracting keywords from pdf with tesseract OCR
func ExtractPdfDataOCR(name string) ([]string, error) {
	var keywords string
	filename := name + ".pdf"
	doc, err := fitz.New(filename)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer doc.Close()

	client := gosseract.NewClient()
	defer client.Close()

	// Extract pages as images
	for n := 0; n < doc.NumPage(); n++ {

		img, err := doc.Image(n)
		if err != nil {
			log.Println(err)
			continue
		}

		buf := new(bytes.Buffer)
		err = jpeg.Encode(buf, img, &jpeg.Options{400})
		if err != nil {
			log.Println(err)
			continue
		}

		imbyte := buf.Bytes()

		err = client.SetImageFromBytes(imbyte)
		if err != nil {
			log.Println(err)
			continue
		}
		text, _ := client.Text()
		keywords = keywords + " " + text

	}

	v := stopwords.CleanString(keywords, "en", false)
	s := strings.Split(v, " ")

	return s, nil

}

func checkTli1(entry string, records [][]string) (bool, string, int, error) {

	i := sort.Search(len(records), func(i int) bool { return entry <= records[i][0] })
	if i < len(records) && records[i][0] == entry {
		//fmt.Println("found")
		return true, records[i][1], i, nil

	} else {
		//fmt.Println("not found")
		return false, "", i, nil
	}
}

func setTLI() ([][]string, error) {
	latestTLI, err := Shell.Resolve(TLI)
	if err != nil {
		log.Println(err)
		return nil, &noipns{TLI}
	}
	cidTLI := strings.Split(latestTLI, "s/")[1]

	cat, err := Shell.Cat(cidTLI)
	if err != nil {
		log.Println(err)
		return nil, &cIDmissing{cidTLI}
	}

	result, err := ioutil.ReadAll(cat)
	if err != nil {
		return nil, err
	}

	err = cat.Close()
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile("./TLI/TLI.csv", result, 0644)
	if err != nil {
		return nil, err
	}

	f, err := os.Open("./TLI/TLI.csv")
	if err != nil {
		return nil, err
	}

	csvr := csv.NewReader(f)
	records, _ := csvr.ReadAll()

	err = f.Close()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func CreateIndexEntryServer1(data []string, cid string) {

	TopLevelIndex, err := setTLI()
	change := false
	if err != nil {
		log.Println(err)
	}

	//index all words in data
	for _, s := range data {

		//remove empty or 1 letter entries
		if s == "" || len(s) == 1 {
			continue
		}

		//check if there is an index file available
		exist, indexCID, TLIposition, err := checkTli1(s, TopLevelIndex)
		if err != nil {
			log.Println(err)
			continue
		}

		// if index file is available add to it
		// otherwise create one
		if exist {

			cat, err := Shell.Cat(indexCID)
			if err != nil {

			}
			result, err := ioutil.ReadAll(cat)
			if err != nil {

			}

			err = cat.Close()
			if err != nil {

			}

			err = ioutil.WriteFile("./test_index/"+s+".csv", result, 0644)
			if err != nil {
				log.Println(err)
			}

			f, err := os.Open("./test_index/" + s + ".csv")
			if err != nil {
				log.Println(err)
			}

			reader := csv.NewReader(f)
			records, _ := reader.ReadAll()

			err = f.Close()
			if err != nil {
				log.Println(err)
			}

			//check if the entry for the CID is already in the index
			//if so, continue, otherwise, add to the sorted index
			i := sort.Search(len(records), func(i int) bool { return cid <= records[i][0] })
			if i < len(records) && records[i][0] == cid {

			} else {
				var entry = []string{
					cid, "false",
				}
				records = append(records, []string{""})
				copy(records[i+1:], records[i:])
				records[i] = entry

				//write new records to the file
				_ = os.Truncate("./test_index/"+s+".csv", 0)

				file, err := os.OpenFile("./test_index/"+s+".csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Println(err)
					continue
				}

				writer := csv.NewWriter(file)
				_ = writer.WriteAll(records)
				writer.Flush()

				err = file.Close()
				if err != nil {
					log.Println(err)
				}

				//add new sub-index file to ipfs
				k, err := os.OpenFile("./test_index/"+s+".csv", os.O_RDONLY, 0644)
				if err != nil {

				}

				id, err := Shell.Add(k)
				if err != nil {

				}
				fmt.Println("changes TLI")
				TopLevelIndex[TLIposition][1] = id
				change = true
			}

		} else {
			f, err := os.Create("./test_index/" + s + ".csv")
			if err != nil {
				log.Println(err)
				continue
			}

			var entry = [][]string{
				{cid, "false"},
			}

			writer := csv.NewWriter(f)
			err = writer.WriteAll(entry)
			if err != nil {
				log.Println(err)
				continue
			}
			err = f.Close()
			if err != nil {
				log.Println(err)
			}

			//add file to ipfs
			k, err := os.OpenFile("./test_index/"+s+".csv", os.O_RDONLY, 0644)
			if err != nil {
			}

			id, err := Shell.Add(k)
			if err != nil {
				log.Println(err)
			}

			e := []string{s, id}

			TopLevelIndex = append(TopLevelIndex, []string{""})
			copy(TopLevelIndex[TLIposition+1:], TopLevelIndex[TLIposition:])
			TopLevelIndex[TLIposition] = e
			fmt.Println("changes TLI")
			change = true
		}

	}

	if change != false {
		log.Println("Updating the TLI...")

		//safe new TLI and add to ipfs
		_ = os.Truncate("./TLI/TLI.csv", 0)
		f, err := os.OpenFile("./TLI/TLI.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		writer := csv.NewWriter(f)
		err = writer.WriteAll(TopLevelIndex)
		writer.Flush()
		err = f.Close()
		if err != nil {
			log.Println(err)
		}

		z, err := os.OpenFile("./TLI/TLI.csv", os.O_RDONLY, 0644)
		contentid, err := Shell.Add(z)
		if err != nil {
			log.Println(err)
		}

		Shell.SetTimeout(time.Duration(1000000000000))
		err = Shell.Publish("", "/ipfs/"+contentid)
		if err != nil {
			log.Println(err)
		}
		Shell.SetTimeout(time.Duration(10000000000))

	}
}

//todo: test client function
func CreateIndexEntryClient1(data []string, cid string) {

	TopLevelIndex, err := setTLI()
	change := false
	if err != nil {
		log.Println(err)
	}

	//index all words in data
	for _, s := range data {
		fmt.Println(TopLevelIndex)

		fmt.Println(s)
		//remove empty or 1 letter entries
		if s == "" || len(s) == 1 {
			continue
		}

		//check if there is an index file available
		exist, indexCID, TLIposition, err := checkTli1(s, TopLevelIndex)
		if err != nil {
			log.Println(err)
			continue
		}

		// if index file is available add to it
		// otherwise create one
		if exist {

			fmt.Println("found2")
			cat, err := Shell.Cat(indexCID)
			if err != nil {

			}
			result, err := ioutil.ReadAll(cat)
			if err != nil {

			}

			err = cat.Close()
			if err != nil {

			}

			err = ioutil.WriteFile("./test_index/"+s+".csv", result, 0644)
			if err != nil {
				log.Println(err)
			}

			f, err := os.Open("./test_index/" + s + ".csv")
			if err != nil {
				log.Println(err)
			}

			reader := csv.NewReader(f)
			records, _ := reader.ReadAll()

			err = f.Close()
			if err != nil {
				log.Println(err)
			}

			//check if the entry for the CID is already in the index
			//if so, continue, otherwise, add to the sorted index
			i := sort.Search(len(records), func(i int) bool { return cid <= records[i][0] })
			if i < len(records) && records[i][0] == cid {
				fmt.Println("found in sub-index")

			} else {
				fmt.Println("not found in sub-index")
				var entry = []string{
					cid, "false",
				}
				records = append(records, []string{""})
				copy(records[i+1:], records[i:])
				records[i] = entry

				//write new records to the file
				_ = os.Truncate("./test_index/"+s+".csv", 0)

				file, err := os.OpenFile("./test_index/"+s+".csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					log.Println(err)
					continue
				}
				writer := csv.NewWriter(file)
				_ = writer.WriteAll(records)
				writer.Flush()

				err = file.Close()
				if err != nil {
					log.Println(err)
				}

				//add new sub-index file to ipfs
				k, err := os.OpenFile("./test_index/"+s+".csv", os.O_RDONLY, 0644)
				if err != nil {

				}

				id, err := Shell.Add(k)
				if err != nil {

				}

				TopLevelIndex[TLIposition][1] = id
				change = true

			}

		} else {

			f, err := os.Create("./test_index/" + s + ".csv")
			if err != nil {
				log.Println(err)
				continue
			}

			var entry = [][]string{
				{cid, "false"},
			}

			writer := csv.NewWriter(f)
			err = writer.WriteAll(entry)
			if err != nil {
				log.Println(err)
				continue
			}
			err = f.Close()
			if err != nil {
				log.Println(err)
			}

			//add file to ipfs
			k, err := os.OpenFile("./test_index/"+s+".csv", os.O_RDONLY, 0644)
			if err != nil {
			}

			id, err := Shell.Add(k)
			if err != nil {
				log.Println(err)
			}

			e := []string{s, id}

			TopLevelIndex = append(TopLevelIndex, []string{""})
			copy(TopLevelIndex[TLIposition+1:], TopLevelIndex[TLIposition:])
			TopLevelIndex[TLIposition] = e
			change = true
		}

	}

	if change != false {
		//safe new TLI and add to ipfs
		_ = os.Truncate("./TLI/TLI.csv", 0)
		f, err := os.OpenFile("./TLI/TLI.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		writer := csv.NewWriter(f)
		err = writer.WriteAll(TopLevelIndex)
		writer.Flush()
		err = f.Close()
		if err != nil {
			log.Println(err)
		}

		z, err := os.OpenFile("./TLI/TLI.csv", os.O_RDONLY, 0644)
		contentid, err := Shell.Add(z)
		if err != nil {
			log.Println(err)
		}

		//from here will run server side
		err = serverTLIUpdate(contentid)
		if err != nil {
			log.Println(err)
		}

	}
}
