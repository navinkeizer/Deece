package main
import (
	"Deece"
	"fmt"
	//"github.com/ethereum/go-ethereum/ethclient"
	//ipfsapi "github.com/ipfs/go-ipfs-api"
	"github.com/urfave/cli"
	"log"
	"os"
)

var app1 = cli.NewApp()


func info1() {
	app1.Name = "Deece Search"
	app1.Usage = "Decentralised search for IPFS"
	app1.Authors = []*cli.Author{
		{Name:  "Navin V. Keizer", Email: "navin.keizer.15@ucl.ac.uk",},
		{Name:  "Puneet K. Bindlish", Email: "p.k.bindlish@vu.nl",},
	}
	app1.Version = "0.0.1"

}

func commands1(){
	ty := ""
	app1.Commands = []*cli.Command{
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
					Deece.DoCrawlClient(id,"CID")
					return nil
				} else if ty == "ENS"{
					Deece.DoCrawlClient(id,"ENS")
					return nil
				}else if ty == "DNS"{
					Deece.DoCrawlClient(id,"DNS")
					return nil
				}else if ty == "IPNS"{
					Deece.DoCrawlClient(id,"IPNS")
					return nil
				}
				return &Deece.IncorrrectInput{}
			},
		},
	}
}

func main() {
	Deece.Shell, Deece.Client, Deece.TLI = Deece.ConnectClient()
	info1()
	commands1()
	err := app1.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
