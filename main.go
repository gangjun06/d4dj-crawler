package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/gangjun06/d4dj-crawler/awsutil"
	"github.com/gangjun06/d4dj-crawler/crawler"
	"github.com/gangjun06/d4dj-crawler/parser"
	"github.com/gangjun06/d4dj-crawler/parser/crypto"
)

func main() {
	crawl := flag.Bool("crawl", false, "crawl assets from server")
	flag.Parse()

	if *crawl {
		awsutil.InitAWS()
		crawler.Start()
	} else {
		args := flag.Args()
		if len(args) < 1 {
			log.Fatalln("argument is empty")
		}
		info, err := os.Stat(args[0])
		if os.IsNotExist(err) {
			log.Fatalln("file not exists")
		}
		if info.IsDir() {
			files, err := ioutil.ReadDir(args[0])
			if err != nil {
				log.Fatal(err)
			}
			wg := new(sync.WaitGroup)

			for _, f := range files {
				if f.IsDir() {
					continue
				}
				wg.Add(1)
				go func(f fs.FileInfo) {
					if err := Parse(path.Join(args[0], f.Name())); err != nil {
						fmt.Println("Error ", info.Name(), ":", err.Error())
					}
					wg.Done()
				}(f)
			}
			wg.Wait()
		} else {
			if err := Parse(args[0]); err != nil {
				fmt.Println("Error ", info.Name(), " : ", err.Error())
				return
			}
		}
	}
}

func Parse(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("file not exists")
	}

	if strings.HasSuffix(filePath, ".enc") {
		decrypt, err := crypto.New().Decrypt(data)
		if err != nil {
			return fmt.Errorf("error while decrypt")
		}
		data = decrypt
	}
	if err := parser.Parse(filePath, data); err != nil {
		return err
	}
	fmt.Println("Success ", filePath)
	return nil
}
