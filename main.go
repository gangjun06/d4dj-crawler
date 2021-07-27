package main

import (
	"flag"
	"time"

	"github.com/gangjun06/d4dj-info-server/routes"
	"github.com/gangjun06/d4dj-info-server/utils/crawler"
)

func main() {
	crawl := flag.Bool("crawl", false, "")
	flag.Parse()
	if *crawl {
		crawler.Start()
		crawler.Start()
		return
	}
	ticker := time.NewTicker(time.Minute * 15)
	go func() {
		crawler.Start()
		for range ticker.C {
			crawler.Start()
		}
	}()
	routes.InitServer()
}
