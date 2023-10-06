package scheduler

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/go-rod/rod"
	"github.com/taintlesscupcake/khu_cafeteria/webcrawl"
)

func AutoCrawler(browser *rod.Browser, db *badger.DB) {
	fmt.Println("AutoCrawler is running...")

	ticker := time.NewTicker(30 * time.Minute)

	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("Crawling...")
		webcrawl.Crawl(db, browser)
	}
}
