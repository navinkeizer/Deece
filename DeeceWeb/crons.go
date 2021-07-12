package main

import (
	"log"
	"os"
	"path/filepath"
	"os/exec"
	"bufio"
	"time"
	"github.com/igpkb/navsearch/Deece"
	"strings"
	"io"
	"io/ioutil"
	"strconv"
)

func RunCrons (options string,corpusDir string,crawlfile string) {
if(options == "all"){
	log.Println("Runcrons(): Running loadCorpus for directory:"+corpusDir+" and try to put in crawlCIDfile:"+crawlfile)
  //go loadCorpus(corpusDir, crawlfile)
  go crawlCron(crawlfile);
}
}


func loadCorpus(corpusDir string,crawlRequestFile string) error {
if _, err := os.Stat(corpusDir); os.IsNotExist(err) {return err} //if _, err := os.Stat(corpusDir); !os.IsNotExist(err) {path/to/whatever exists}
if _, err := os.Stat(crawlRequestFile); os.IsNotExist(err) {if _, err := os.Create(crawlRequestFile); err !=nil {return err}}
// unzip all before creating a files
unzipAll (corpusDir, true)
var files []string
var output []byte
	files,err := fileList (corpusDir); if (err !=nil) {return err }
for _, file := range files {
	log.Println("loadCorpus(): Adding file to ipfs:"+file)	
	output, err = (exec.Command( "ipfs", "add", file )).Output(); if(err != nil) {return err} 
		log.Println(output); 
		outputstr :=string(output)
		if(strings.Contains(outputstr, "added ")) {InsertStringToFile(crawlRequestFile,strings.Split(outputstr," ")[1],-1)}
			   }
	return nil
}

func crawlCron(crawlRequestFile string) error {
	src := crawlRequestFile
	file, err := os.Open(src)
    if err != nil {
	    log.Println("Crawl File couldn't be accessed: "+src+ "with error: "+err.Error())  //log.Fatal(err)
	    return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    // optionally, resize scanner's capacity for lines over 64K, see next example
    for scanner.Scan() {
	    CID := scanner.Text()
        log.Println("Submitting crawl for " + CID)
    	go Deece.DoCrawlServer(CID,"CID")
    }
    if err := scanner.Err(); err != nil {
         log.Println("Crawl process failed for: "+src+ "with error: "+err.Error())  //log.Fatal(err)
	    return err
    } 
    //rename crawlCID.txt to crawlCID.txt.datetimenow
	dst := src+"."+(time.Now().Format("01-02-2006-15-04-05"))
	if  err := os.Rename(src, dst); err != nil {
		 log.Println("File renaming failed for "+src + "with error: "+err.Error()) //log.Fatal(err)
		 return err
   		 }    
	return nil
}

func File2lines(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return LinesFromReader(f)
}

func unzipAll (root string, recursive bool) error {
	log.Println("unzipAll(): Running for directory:"+root+" with recursive option:"+strconv.FormatBool(recursive))
	var zipfiles []string
	err := filepath.Walk (root, func(path string, info os.FileInfo, err error) error {
						if(!info.IsDir() && filepath.Ext(path) == ".zip")  {
								log.Println(path);
								zipfiles = append(zipfiles, path) }
        					return nil })
	if(err !=nil) {
		return err
	}
	for _, file := range zipfiles {
		log.Println("unzipAll(): Unzipping:"+string(file))
		output, err := (exec.Command( "tar", "-xzvf",file )).Output(); if(err != nil) {return err} 
					log.Println(output); 
					} 
		return nil
}

func fileList(root string) (files []string, err error) {
	err = filepath.Walk ( root, func(path string, info os.FileInfo, err error) error {
						if(info.IsDir()) {
								childfiles,err := fileList(path)
								if err != nil {
									return err
								}
								for _, file := range childfiles {
												log.Println("fileList(): Adding child file:"+path)
												log.Println(path);
												files = append(files, file)
												}
								 }
						//info.Name() carries the file name only without the path
						if(!info.IsDir() && (filepath.Ext(path) != ".zip"))  { 
							log.Println("fileList(): Adding file:"+path)	
							files = append(files, path) }
        					return nil })
	if err != nil {
		return nil,err }
	
	return files, nil
}

func LinesFromReader(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
func InsertStringToFile(path, str string, index int) error {
	lines, err := File2lines(path)
	if err != nil {
		return err
	}
	//-1 index corresponds to append to the end of file
	if (index>-1){
			index = len(lines)
		    }
	fileContent := ""
	for i, line := range lines {
		if i == index {
			fileContent += str
		}
		fileContent += line
		fileContent += "\n"
	}
	return ioutil.WriteFile(path, []byte(fileContent), 0644)
}
