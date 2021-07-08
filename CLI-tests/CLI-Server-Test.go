package main

import (
	"Deece"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
)

var app = cli.NewApp()


func info() {
	app.Name = "Deece Search"
	app.Usage = "Decentralised search for IPFS"
	app.Authors = []*cli.Author{
		{Name:  "Navin V. Keizer", Email: "navin.keizer.15@ucl.ac.uk",},
		{Name:  "Puneet K. Bindlish", Email: "p.k.bindlish@vu.nl",},
	}
	app.Version = "0.0.1"

}

func commands(){
	ty := ""
	app.Commands = []*cli.Command{
		{
			Name:        "search",
			//Aliases:     []string{"s"},
			Usage:       "Performs a decentralised search on IPFS",
			Description: "Retrieves the index to find which pages contain the keywords",

			Action: func(c *cli.Context) error {
				if c.Args().Len() < 1{
					return &Deece.IncorrrectInput{}
				}
				Deece.DoSearch(c.Args().Slice())
				return nil
			},
		},

		{
			Name:        "crawl",
			//Aliases:     []string{"c"},
			Usage:       "Crawls a page to add to the decentralised index",
			Description: "Crawls the page, extracts keywords using OCR, and adds to the index stored on IPFS",

			Flags: []cli.Flag{
				&cli.StringFlag{
					Name: "type",
					Aliases: []string{"t"},
					Usage: "Input domain type (default CID)",
					Destination: &ty,
				},
			},

			Before: func(c *cli.Context) error {
				fmt.Println("Start crawl...")
				return nil
			},

			Action: func(c *cli.Context) error {
				if c.Args().Len() != 1{
					return &Deece.IncorrrectInput{}
				}
				id := c.Args().Get(0)
				if ty == "" || ty == "CID"{
					Deece.DoCrawlServer(id,"CID")
					return nil
				} else if ty == "ENS"{
					Deece.DoCrawlServer(id,"ENS")
					return nil
				}else if ty == "DNS"{
					Deece.DoCrawlServer(id,"DNS")
					return nil
				}else if ty == "IPNS"{
					Deece.DoCrawlServer(id,"IPNS")
					return nil
				}
				return &Deece.IncorrrectInput{}
			},
		},
	}
}

func main1() {
	Deece.Shell, Deece.Client, Deece.TLI = Deece.ConnectServer()
	info()
	commands()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}