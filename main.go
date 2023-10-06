package main

import (
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/go-rod/rod"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/taintlesscupcake/khu_cafeteria/router"
	"github.com/taintlesscupcake/khu_cafeteria/scheduler"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db, err := badger.Open(badger.DefaultOptions(os.Getenv("DB_PATH")))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	app := fiber.New()
	browser := rod.New().MustConnect().NoDefaultDevice()

	go scheduler.AutoCrawler(browser, db)

	router.Router(app, browser, db)

}
