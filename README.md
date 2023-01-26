<img src="./images/DEECElogo.png" width="500" >


# Deece Search

[![cli passing](https://img.shields.io/badge/cli-passing-green)](CLI)
[![library passing](https://img.shields.io/badge/library-passing-blue)](Deece)
[![readme deece](https://img.shields.io/badge/readme-deece-red)](README.md)
[![license deece](https://img.shields.io/badge/license-Apache%202-orange)](LICENSE)

Deece Search is an open, collaborative, and decentralised search mechanism for IPFS. Any node running the client is able to crawl content on IPFS and add this to the index, which itself is stored in a decentralised manner on IPFS. This allows for decentralised search on decentralised content. 


***The current implementation is still highly experimental. We are working on the second version without central gateway and automatic crawling, so for the time being our server is down.***



## Table of Contents

- [Overview](#overview)
  - [Motivation](#motivation)
  - [State-of-the-Art](#state-of-the-art)
- [Architecture Design](#architecture-design)
  - [Search](#search)
  - [Crawl](#crawl)
  - [Performance Considerations](#performance-considerations)
- [Components](#components)
  - [`Client`](#client) 
  - [`Gateway`](#gateway) 
  - [`Library`](#library)
- [Install](#install)
  - [Client CLI](#client-cli)
  - [Go Library](#go-library)  
- [Project Status](#project-status)
- [Maintainers](#maintainers)



## Overview

Deece Search allows for decentralised search on IPFS data. This is achieved by a network of IPFS nodes who participate in crawling and indexing of the data on the network. The index is stored on IPFS, and split up into a two-layer hierarchy, the first being the *Top Level Index* (TLI) and the second being the *Keyword Specific Indexes* (KSI). The TLI contains the identifiers (CID) for the KSI for each keyword, and is constantly updated when a node submits a crawl. When crawling, the nodes add to the current KSI a list of the identifiers of files that contain that keyword. 


<img src="./images/indexes.jpg" width="550" >


Deece Search allows for two specific actions: `search` and `crawl`. Search queries the latest TLI to find the KSI for each keyword in the user query, and then fetches the results from these, which are displayed to the user. Currently ranking of results is ordered based on CID, but more sophisticated mechanisms should be developed. We allow for combined results for up to two keywords, which will be extended in the future. 

There are currently three ways to access Deece Search. First, there is the client software which uses a command line interface. Second, we have implemented a gateway service(www.deece.nl/web/), which runs an instance of our client node and allows "light clients" to interact with the search without installing other software. Finally, we have released our code used by the CLI and gateway in the form of a Go library.  

*The initial version of Deece Search relies on a trusted node (the same node as our gateway) to update the IPNS record pointing to the latest version of the TLI. When clients crawl, the final step involves them sending an update request to this server. At the moment, clients will need to specify a password in their config file, which can be obtained from the maintainers, as security measures will be implemented later.*


### Motivation

Currently, Web users have few alternatives to *centralised* search engines. These engines maintain centralised control, policy and trust, which may lead to issues in censorship, privacy protection, and transparency. 

Furthermore, these engines generally focus their efforts on traditional Web content (hosted at Web servers, accessed through the DNS). However, in a Web3 paradigm, where content is expected to be stored at decentralised storage networks (e.g. IPFS) and name resolution to take place through blockchain solutions (e.g. ENS), an alternative search engine is required. 

In short, a search mechanism is needed which searches decentralised data, and does so in a decentralised manner. 


### State-of-the-Art

There are a number of comparable projects, which have attempted to solve the problem of centralisation in Web search. First of all, there are implementations and proposals from research for distributed / decentralised search mechanisms for the current Web data. Early projects include Yacy, Faroo, and Seeks. More recently, Presearch aims to create a collaborative search engine using blockchain rewards for incentives. 

Similarly, a number of works aimed to provide distributed search for P2P storage networks. More recently, The Graph has built a decentralised indexing protocol for blockchain data using cryptocurrency incentives. 

However, none of the above projects entirely captures our specific use case of decentralised search for decentralised Web3 data. 



## Architecture Design

Our architecture relies on a number of client nodes, which collectively maintain and add to the index, and are able to perform searches. We have taken the approach of finishing a working protype of our architecture first, and adding features incrementally. Therefore, our current version relies on a trusted node (gateway) to update the TLI IPNS record. As there is no added security or incentivisation implemented, we have used a simple password to allow new client nodes to add to the index. While security may be insufficient in the future, we assume an altruistic model for our early stage release. 

In the future we envision there to be added security and incentives in place, which align nodes to be honest when updating the index. These may be in the form of cryptocurreny rewards, slashing, reputation, etc. One way to fund rewards to honest nodes could be by incorporating advertisement into the protocol and allow advertisement fees to be delegated to the nodes maintaining the network. 

Our current version only supports PDF files on IPFS to be added to the index. In the future, we would like to expand this to more file types and directories, and support different decentralised storage networks. Finally, we aim to incorporate blockchain based data such as smart contracts into search. 

We now present an overview of the two main operations in our mechanism.



### Search

Search starts with a query by the client containing a number of search terms. The client then fetches the latest TLI by resolving the IPNS name set by the gateway to the corresponding CID. This TLI is then fetched and traversed to check if the keywords have KSI's. If this is the case, the relevant KSI's are queries, to return the content that contains the keywords. The client can then retrieve these files from the network. 

<img src="./images/search.jpg" width="550" >

One important aspect in search engines is the ranking mechanism. This generally happens in a centralised manner, without much influence from the clients. While we have not implemented sophisticated ranking mechanisms, we envision there to be ranking at the clients of the results, which gives them greater power and transparency. This allows clients to be in control of ranking functions and to personalise these based on specific needs. At present, our mechanism returns ordered results based on CID's. When two search terms are entered, the pages where these both occur are returned first, after which the pages are returned which contain only one of the terms. 


### Crawl

An important aspect of any search engine is the addition of entries to the index. This process involves a number of steps, which we describe below. 

The first decision to be made is what content will be added to the index, which we call *curation*. In traditional engines this includes all public Web content. Although this achieves high performance, it may add too much overhead when executed in a decentralised network. Another approach may be curation based on network consensus of important content. For our current system, we allow anyone who believes content to be important to add this to the network. Content can be addressed by CID, DNSLink, ENS, or IPNS identifiers. 

Next, crawling happens, which involves fetching and analysing files to extract important keywords. As mentioned above, our system crawls when someone decides content should be added, and thus manually submits the CID to be crawled. In the future, we envision this to happen automatically when content is uploaded or visited on the network. Besides extracting keywords, other metadata may be added. We currently use the file type (PDF) and timestamp when crawled, but in the future intent to add title, count, size etc.

<img src="./images/crawl.jpg" width="550" >

After extracting the keywords (and producing the RWI), the index needs to be stored. For storage we use IPFS, as this allows for decentralised collaborative storage. We have decided to maintain a two-level hierarchy. Each keyword will have an associated index file (KSI) where nodes can find which content contains those keywords. A separate index is kept (TLI) to point to the identifiers of the KSI's, and this is published to an IPNS name from our gateway server. When a node updates the KSI's after crawling a file, they update the pointer in the TLI to these files, and requests the gateway to update the pointer which the IPNS record resolves to. This way the IPNS record points to the latest version of the TLI, which in turn points to the latest versions of the KSI's.

At present, client nodes can change the TLI if they possess a password, which can be obtained from the maintainers of this project. This way, potential malicious entries are less likely. Nodes with the passwords can be seen as 'authorities' in the network. 



### Performance Considerations

During development and testing we have made a number of observations with regards to performance. As our solution relies heavily on IPFS, so does our performance. We found that significant delays may occur when nodes have not added the gateway peer in their swarm of peers. While we have added this to our CLI, connection still drops occasionally. While this does not break the system, it does add delays. 

Furthermore, updating our IPNS entry from the gateway can be very slow, and may become a performance bottleneck when crawl traffic increases. We have started looking into alternatives, but leave implementation to future releases. One option is to use the DNS to store a pointer to the latest TLI record, but this brings a number of additional challenges inherent to the DNS. A blockchain based name registry such as ENS may also be used, although frequent updates to the resolver contract may become a large expense.       


## Components

There are a number of ways to access Deece Search: 

### `Client`

The client software can be used by any node running IPFS, and provides a simple command line interface. 

```shell

NAME:
   Deece - Decentralised search for IPFS

USAGE:
   Deece [global options] command [command options] [arguments...]

VERSION:
   0.0.1

AUTHORS:
   Navin V. Keizer <navin.keizer.15@ucl.ac.uk>
   Puneet K. Bindlish <p.k.bindlish@vu.nl>

COMMANDS:
   search   Performs a decentralised search on IPFS
   crawl    Crawls a page to add to the decentralised index
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)

```


### `Gateway`

For easy and lightweight access we have implemented a gateway for our search clients. This can be found at: www.deece.nl/web/, and allows for search and crawls on the network based on identifiers (CID's).

<img src="./images/webinterface.png" width="700" >

***Note: the gateway is currently susspended while upgrading to version 2.***


### `Library`

Both the CLI and gateway run using our Deece Search package for Go. We have released this, as this can be used for easy integrations and extensions. 



## Install

***Further installation instructions will be added once tested across different platforms. For now we have provided instructions based on our installation on Linux.***

For Deece Search to work there are a number of requirements and dependencies. To run as a client, a local IPFS daemon needs to be running, and to speed up results it helps to add the gateway maintaining the TLI in the peer swarm. To submit changes to the TLI as a client a password is required. Finally, a config file needs to be present in the same directory as the executable to load results. An incomplete config file can be found in this repository.  

### Client CLI

To run the client, first [IPFS](https://docs.ipfs.io/install/command-line/#system-requirements), [Go](https://golang.org/doc/install) (tested for version 1.13.7, newer versions should work with minor modifications), and [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) must be installed. 

Next we need to install from source:

```shell
git clone github.com/navinkeizer/Deece
```
Next [tesseract-ocr](https://tesseract-ocr.github.io/tessdoc/Installation.html) needs to be installed, as well as other dependencies. For Linux this may look like this: 
```shell
sudo apt-get install g++ 
sudo apt-get install autoconf automake libtool
sudo apt-get install autoconf-archive
sudo apt-get install pkg-config
sudo apt-get install libpng-dev
sudo apt-get install libjpeg8-dev
sudo apt-get install libtiff5-dev
sudo apt-get install zlib1g-dev
wget http://www.leptonica.org/source/leptonica-1.81.1.tar.gz
sudo tar xf leptonica-1.81.1.tar.gz
cd leptonica-1.81.1 &&\
sudo ./configure &&\
sudo apt install make
sudo make &&\
sudo make install
sudo apt-get install tesseract-ocr # or sudo apt install tesseract-ocr
sudo apt install libtesseract-dev
```

other relevant Go packages may then be installed:

```
$ go get -t github.com/otiai10/gosseract
$ go get github.com/navinkeizer/Deece
$ go get github.com/ipfs/go-ipfs-api 
$ go get github.com/wealdtech/go-ens/v3 
$ go get github.com/otiai10/gosseract/v2 
```

and the CLI built:

```
$ sudo go build Deece/CLI/.
```

and run:
```
$ ./CLI [command] [arguments]
```



### Go Library

The package can also be used as a library.

```
go get github.com/navinkeizer/Deece
```


## Project Status

The current implementation of Deece Search is still experimental, and therefore may experience instabilities. As described in this document, we have made simplifying assumptions (altruism) and focused on limited functionality (PDF only). Furthermore, the gateway presents a centralised aspect, which in the future should be replaced by decentralised network consensus, and the protocol should be secured by incentives. 

Our implementation takes a first principle approach. We aimed to build from the ground up, rather than relying on existing approaches and solutions for system components. We believe this is necessary as existing solutions may not be optimal for decentralised Web3 content. In other words, there is much work to be done.

***Currently we occasionally experience issues in the crawling process due to IPNS updates timing out. We are working on resolving this with alternative solutions.***

## Maintainers

- [@navinkeizer](https://github.com/navinkeizer/)
- [@igpkb](https://github.com/igpkb/)



