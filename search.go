package Deece

import (
	"encoding/csv"
	"sort"

	//shell "github.com/ipfs/go-ipfs-api"
	"log"
)

//type resultKeyword struct {
//	CID      string
//	metadata string
//}

//func fetchIndex(searchTerm string, records [][]string) string {
//	for j := 0; j < len(records); j++ {
//		if searchTerm == records[j][0] {
//			return records[j][1]
//		}
//	}
//	return "-"
//}

//func printPerTermClient(items []string, cids []string) error {
//	for i := 0; i < len(items); i++ {
//		fmt.Println(items[i] + ": ")
//		if cids[i] == "-" {
//			continue
//		}
//
//		cat, err := Shell.Cat(cids[i])
//		if err != nil {
//			log.Println(err)
//		}
//
//		csvr := csv.NewReader(cat)
//		records, err := csvr.ReadAll()
//		if err != nil {
//		}
//
//		for j := 0; j < len(records); j++ {
//			fmt.Println(records[j][0])
//		}
//		fmt.Println()
//
//		err = cat.Close()
//		if err != nil {
//		}
//
//	}
//	return nil
//}
//
//func printPerTermServer(items []string, cids []string) ([]string, [][]resultKeyword) {
//
//	var searchTerms []string
//	var searchResult = make([][]resultKeyword, int(len(items)))
//
//	for i := 0; i < len(items); i++ {
//		searchTerms = append(searchTerms, "'"+items[i]+"'")
//		//fmt.Println(items[i] + ": ")
//		if cids[i] == "-" {
//			continue
//		}
//
//		cat, err := Shell.Cat(cids[i])
//		if err != nil {
//			log.Println(err)
//		}
//
//		csvr := csv.NewReader(cat)
//		records, err := csvr.ReadAll()
//		if err != nil {
//		}
//
//		for j := 0; j < len(records); j++ {
//			searchResult[i] = append(searchResult[i], resultKeyword{records[j][0], records[j][1]})
//			//fmt.Println(records[j][0])
//		}
//		//fmt.Println()
//
//		err = cat.Close()
//		if err != nil {
//		}
//
//	}
//
//	fmt.Println(searchTerms)
//	fmt.Println(searchResult)
//
//	return searchTerms, searchResult
//}

//func printResultsWordClient(items []string, cids []string) error {
//
//	if len(items) == 2 {
//
//		if cids[0] == "-" || cids[1] == "-" {
//			err := printPerTermClient(items, cids)
//			if err != nil {
//				return err
//			}
//			return nil
//		}
//
//		cat1, err := Shell.Cat(cids[0])
//		if err != nil {
//			return err
//		}
//		csvr1 := csv.NewReader(cat1)
//		records1, err := csvr1.ReadAll()
//
//		cat2, err := Shell.Cat(cids[1])
//		if err != nil {
//			return err
//		}
//		csvr2 := csv.NewReader(cat2)
//		records2, err := csvr2.ReadAll()
//
//		//fmt.Println(len(records1))
//		//fmt.Println(len(records2))
//
//		var combinedresult [][]string
//
//		if len(records1) > len(records2) {
//			for z := 0; z < len(records1); z++ {
//				for b := 0; b < len(records2); b++ {
//					if records1[z][0] == records2[b][0] {
//						combinedresult = append(combinedresult, records1[z])
//					}
//				}
//			}
//
//			fmt.Println(items[0] + " " + items[1] + " : ")
//			for j := 0; j < len(combinedresult); j++ {
//				fmt.Println(combinedresult[j][0])
//			}
//
//		} else {
//			for z := 0; z < len(records2); z++ {
//				for b := 0; b < len(records1); b++ {
//					if records2[z][0] == records1[b][0] {
//						combinedresult = append(combinedresult, records2[z])
//					}
//				}
//			}
//
//			fmt.Println(items[0] + " " + items[1] + " : ")
//			for j := 0; j < len(combinedresult); j++ {
//				fmt.Println(combinedresult[j][0])
//			}
//
//		}
//		fmt.Println()
//		fmt.Println(items[0] + " : ")
//		for p := 0; p < len(records1); p++ {
//			fmt.Println(records1[p][0])
//		}
//		fmt.Println()
//		fmt.Println(items[1] + " : ")
//		for q := 0; q < len(records2); q++ {
//			fmt.Println(records2[q][0])
//		}
//
//	} else {
//		err := printPerTermClient(items, cids)
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}
//
//func printResultsWordServer(items []string, cids []string) ([]string, [][]resultKeyword, error) {
//
//	if len(items) == 2 {
//
//		if cids[0] == "-" || cids[1] == "-" {
//			var searchTerms []string
//			var searchResult = make([][]resultKeyword, len(items))
//			searchTerms, searchResult = printPerTermServer(items, cids)
//			return searchTerms, searchResult, nil
//		}
//
//		cat1, err := Shell.Cat(cids[0])
//		if err != nil {
//			return nil, nil, err
//		}
//		csvr1 := csv.NewReader(cat1)
//		records1, err := csvr1.ReadAll()
//
//		cat2, err := Shell.Cat(cids[1])
//		if err != nil {
//			return nil, nil, err
//		}
//		csvr2 := csv.NewReader(cat2)
//		records2, err := csvr2.ReadAll()
//
//		var searchTerms []string
//		var searchResult = make([][]resultKeyword, len(items)+1)
//
//		var combinedresult [][]string
//		searchTerms = append(searchTerms, "'"+items[0]+" "+items[1]+"'")
//
//		if len(records1) > len(records2) {
//			for z := 0; z < len(records1); z++ {
//				for b := 0; b < len(records2); b++ {
//					if records1[z][0] == records2[b][0] {
//						combinedresult = append(combinedresult, records1[z])
//					}
//				}
//			}
//		} else {
//			for z := 0; z < len(records2); z++ {
//				for b := 0; b < len(records1); b++ {
//					if records2[z][0] == records1[b][0] {
//						combinedresult = append(combinedresult, records2[z])
//					}
//				}
//			}
//		}
//
//		for j := 0; j < len(combinedresult); j++ {
//			searchResult[0] = append(searchResult[0], resultKeyword{combinedresult[j][0], combinedresult[j][1]})
//		}
//
//		searchTerms = append(searchTerms, "'"+items[0]+"'")
//		for p := 0; p < len(records1); p++ {
//			searchResult[1] = append(searchResult[1], resultKeyword{records1[p][0], records1[p][1]})
//		}
//
//		searchTerms = append(searchTerms, "'"+items[1]+"'")
//		for q := 0; q < len(records2); q++ {
//			searchResult[2] = append(searchResult[2], resultKeyword{records2[q][0], records2[q][1]})
//		}
//
//		return searchTerms, searchResult, nil
//	} else {
//		var searchTerms []string
//		var searchResult = make([][]resultKeyword, len(items))
//		searchTerms, searchResult = printPerTermServer(items, cids)
//		return searchTerms, searchResult, nil
//	}
//
//}
//

type QueryResult struct {
	SearchTerm string `json:"searchTerm"`
	CID        string `json:"CID"`
	Metadata   string `json:"metadata"`
}

func perTermServer1(terms []string, locations []string) ([]QueryResult, error) {

	var searchResult []QueryResult

	for i := 0; i < len(terms); i++ {
		if locations[i] == "-" {
			continue
		}
		cat, err := Shell.Cat(locations[i])
		if err != nil {
			log.Println(err)
		}

		csvr := csv.NewReader(cat)
		records, err := csvr.ReadAll()
		if err != nil {

		}

		err = cat.Close()
		if err != nil {
		}

		for j := 0; j < len(records); j++ {
			r := QueryResult{
				SearchTerm: terms[i],
				CID:        records[j][0],
				Metadata:   records[j][1],
			}

			searchResult = append(searchResult, r)
		}
	}

	return searchResult, nil
}

func twoTerm(terms []string, locations []string) ([]QueryResult, error) {
	var combinedsearchResult []QueryResult

	cat1, err := Shell.Cat(locations[0])
	if err != nil {
		return nil, err
	}
	csvr1 := csv.NewReader(cat1)
	records1, err := csvr1.ReadAll()

	cat2, err := Shell.Cat(locations[1])
	if err != nil {
		return nil, err
	}
	csvr2 := csv.NewReader(cat2)
	records2, err := csvr2.ReadAll()

	//find overlapping records
	if len(records1) > len(records2) {
		for z := 0; z < len(records1); z++ {

			i := sort.Search(len(records2), func(i int) bool { return records1[z][0] <= records2[i][0] })
			if i < len(records2) && records2[i][0] == records1[z][0] {
				r := QueryResult{
					SearchTerm: terms[0] + " " + terms[1],
					CID:        records2[i][0],
					Metadata:   records2[i][1],
				}
				combinedsearchResult = append(combinedsearchResult, r)
			}
		}
		return combinedsearchResult, nil
	} else {
		for z := 0; z < len(records2); z++ {

			i := sort.Search(len(records1), func(i int) bool { return records2[z][0] <= records1[i][0] })
			if i < len(records1) && records1[i][0] == records2[z][0] {
				r := QueryResult{
					SearchTerm: terms[0] + " " + terms[1],
					CID:        records1[i][0],
					Metadata:   records1[i][1],
				}
				combinedsearchResult = append(combinedsearchResult, r)
			}
		}
		return combinedsearchResult, nil

	}
}

func resultsWordServer1(searchterms []string, indexlocation []string) ([]QueryResult, error) {

	if len(searchterms) == 2 {
		if indexlocation[0] != "-" && indexlocation[1] != "-" {

			comres, err := twoTerm(searchterms, indexlocation)
			if err != nil {
				return nil, err
			}

			singres, err := perTermServer1(searchterms, indexlocation)
			if err != nil {
				return nil, err
			}

			comres = append(comres, singres...)

			return comres, nil
		} else {
			r, err := perTermServer1(searchterms, indexlocation)
			if err != nil {
				return nil, err
			}
			return r, nil
		}
	} else {
		r, err := perTermServer1(searchterms, indexlocation)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
}
