package main

import (
	"log"
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/igpkb/navsearch/Deece"
	/*"strings"*/
)

func getSearchResultHandler(w http.ResponseWriter, r *http.Request) {
	// the `ParseForm` method of the request, parses HTML form data
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var searchResults []Deece.QueryResult
	log.Println(r.Form.Get("q"))
	searchResults,err = Deece.DoSearch1(r.Form.Get("q"))//,r.Form.Get("ctype"))
	log.Println(searchResults)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	searchResultListBytes, err := json.Marshal(searchResults)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(searchResultListBytes)
}

type ResponseMessage struct {
	Text     string `json:"text"`	
}

type CrawlTarget struct {
	ContentID     string `json:"CID"`
	Type     string `json:"type"`
	Description string `json:"description"`
}

var crawlTargets []CrawlTarget
func getCrawlTargetHandler(w http.ResponseWriter, r *http.Request) {
	crawlTargetListBytes, err := json.Marshal(crawlTargets)
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(crawlTargetListBytes)
}

func createCrawlTargetHandler(w http.ResponseWriter, r *http.Request) {
	crawlTarget := CrawlTarget{}
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	crawlTarget.ContentID = r.Form.Get("CID")
	crawlTarget.Type = r.Form.Get("type")
	crawlTargets = append(crawlTargets, crawlTarget)
	log.Println("CID=" + r.Form.Get("CID") + "  type=" + r.Form.Get("type"))
	go Deece.DoCrawlServer(r.Form.Get("CID"),r.Form.Get("type"))
	var responseMessages []ResponseMessage
	responseMessages = append(responseMessages,ResponseMessage{"Crawl submitted successfully."})
	respMsgBytes, err := json.Marshal(responseMessages)
	w.Write(respMsgBytes)
	//http.Redirect(w, r, "/web/", http.StatusFound)
}
