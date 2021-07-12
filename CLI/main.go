package main

import (
	"Deece"
	"context"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"log"
	"os"
)

type Configuration struct {
	ServerIP   string
	ServerPort int
	EthGateway string
	TLI        string
	PassW      string
	ServerAddr string
}

var (
	configuration Configuration
	app1          = cli.NewApp()
)

//load parameters from config file
//requires config file to be present
func setConfig() error {
	file, err := os.Open("./config1.json")
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		return err
	}
	return nil
}

func info() {
	app1.Name = "Deece"
	app1.Usage = "Decentralised Search for IPFS"
	app1.Authors = []*cli.Author{
		{Name: "Navin V. Keizer", Email: "navin.keizer.15@ucl.ac.uk"},
		{Name: "Puneet K. Bindlish", Email: "p.k.bindlish@vu.nl"},
	}
	app1.Version = "0.0.1"

}

func commands() {
	ty := ""
	app1.Commands = []*cli.Command{
		{
			Name:        "search",
			Usage:       "Performs a decentralised search on IPFS",
			Description: "Retrieves the index to find which pages contain the keywords",

			Action: func(c *cli.Context) error {
				if c.Args().Len() < 1 {
					return &Deece.IncorrrectInput{}
				}
				//ensure swarm is connected to gateway peer
				err := Deece.Shell.SwarmConnect(context.Background(), configuration.ServerAddr)
				if err != nil {
					log.Println(err)
				}
				fmt.Println("searching...")
				err = Deece.DoSearchClient(c.Args().Slice())
				if err != nil {
					return err
				}
				return nil
			},
		},

		{
			Name:        "crawl",
			Usage:       "Crawls a page to add to the decentralised index",
			Description: "Crawls the page, extracts keywords using OCR, and adds to the index stored on IPFS",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "type",
					Aliases:     []string{"t"},
					Usage:       "Input domain type (default CID)",
					Destination: &ty,
				},
			},

			Before: func(c *cli.Context) error {
				fmt.Println("Start crawl...")
				//ensure swarm is connected to gateway peer
				err := Deece.Shell.SwarmConnect(context.Background(), configuration.ServerAddr)
				if err != nil {
					log.Println(err)
				}
				return nil
			},

			Action: func(c *cli.Context) error {
				if c.Args().Len() != 1 {
					return &Deece.IncorrrectInput{}
				}
				id := c.Args().Get(0)
				if ty == "" || ty == "CID" {
					err := Deece.DoCrawlClient(id, "CID")
					if err != nil {
						return err
					}
					return nil
				} else if ty == "ENS" {
					err := Deece.DoCrawlClient(id, "ENS")
					if err != nil {
						return err
					}
					return nil
				} else if ty == "DNS" {
					err := Deece.DoCrawlClient(id, "DNS")
					if err != nil {
						return err
					}
					return nil
				} else if ty == "IPNS" {
					err := Deece.DoCrawlClient(id, "IPNS")
					if err != nil {
						return err
					}
					return nil
				}
				return &Deece.IncorrrectInput{}
			},
		},
	}
}

//main CLI program
func main() {
	err := setConfig()
	if err != nil {
		log.Fatal(err)
	}
	Deece.Shell, Deece.Client = Deece.ConnectClient(configuration.EthGateway,
		configuration.TLI, configuration.ServerIP, configuration.ServerPort, configuration.PassW)
	info()
	commands()

	err = app1.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
