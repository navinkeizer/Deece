package Deece

import (
	"encoding/csv"
	"sort"

	//shell "github.com/ipfs/go-ipfs-api"
	"log"
)

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
