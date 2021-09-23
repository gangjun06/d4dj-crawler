package main

import (
	"flag"
	"time"

	"github.com/gangjun06/d4dj-crawler/conf"
	"github.com/gangjun06/d4dj-crawler/utils/crawler"
)

func main() {
	mode := flag.String("mode", "parse", "Run Mode.\nparse -> decrypt, convert, extract unity asset\ncrawl -> crawl assets from server")
	flag.Parse()
	if *mode == "parse" {

	} else if *mode == "crawl" {
		crawler.Start()
		if conf.Get().CrawlerTimer > 0 {
			ticker := time.NewTicker(time.Minute * 15)
			for range ticker.C {
				crawler.Start()
			}
		}
	}
}
