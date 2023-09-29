package router

import (
	"github.com/dgraph-io/badger/v4"
	"github.com/go-rod/rod"
	"github.com/gofiber/fiber/v2"
	"github.com/taintlesscupcake/khu_cafeteria/webcrawl"
)

func Router(app *fiber.App, browser *rod.Browser, db *badger.DB) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/crawl", func(c *fiber.Ctx) error {
		webcrawl.Crawl(db, browser)
		return c.SendString("Crawled")
	})

	app.Post("/get", func(c *fiber.Ctx) error {
		id := c.FormValue("id")

		err := db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(id))

			if err != nil {
				return err
			}

			value, _ := item.ValueCopy(nil)
			img, err := txn.Get(value)

			if err != nil {
				return err
			}

			val, _ := img.ValueCopy(nil)

			c.Type("jpeg")
			return c.Send(val)
		})
		if err != nil {
			return c.SendString(err.Error())
		}
		return nil
	})

	app.Listen(":4000")
}
