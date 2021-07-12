package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/igpkb/navsearch/Deece"
	/*"strings"*/
)

func main() {
	Deece.Shell, Deece.Client = Deece.ConnectServer("https://mainnet.infura.io/v3/5c04e573d61b4e5a8fc0f3312becfdbc","k51qzi5uqu5dm51kzdsr3pu33tkrzca5pse1kt3i9a7c5m1ai0kremf5c0ooe9") //, "127.0.0.1",3333)
	RunCrons("all","corpus","crawlCIDs.txt")
	r := newRouter()
	err := http.ListenAndServe(":80", r)
	if err != nil {
		panic(err.Error())
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hello", handler).Methods("GET")
	staticFileDirectory := http.Dir("./web/")
	staticFileHandler := http.StripPrefix("/web/", http.FileServer(staticFileDirectory))
	r.PathPrefix("/web/").Handler(staticFileHandler).Methods("GET")
	r.HandleFunc("/search", getSearchResultHandler).Methods("GET")
	r.HandleFunc("/getcrawlTargets", getCrawlTargetHandler).Methods("GET")
	r.HandleFunc("/crawlTarget", createCrawlTargetHandler).Methods("GET")
	return r
}
