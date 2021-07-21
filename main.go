package main

import (
	"flag"

	"github.com/gangjun06/d4dj-info-server/routes"
	"github.com/gangjun06/d4dj-info-server/utils/crawler"
)

func main() {
	crawl := flag.Bool("crawl", false, "")
	flag.Parse()
	if *crawl {
		crawler.Start()
		return
	}
	routes.InitServer()
}
