package router

import (
	"fmt"
	"os"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/go-rod/rod"
	"github.com/gofiber/fiber/v2"
	"github.com/taintlesscupcake/khu_cafeteria/webcrawl"
)

type GetRequest struct {
	Id fiber.Map `json:"userRequest"`
}

func fixId(id string) string {
	currentTime := time.Now()

	switch id {
	case "오늘":
		return currentTime.Format("0102")
	case "내일":
		tomorrow := currentTime.AddDate(0, 0, 1)
		return tomorrow.Format("0102")
	case "모레":
		dayAfterTomorrow := currentTime.AddDate(0, 0, 2)
		return dayAfterTomorrow.Format("0102")
	default:
		return id
	}
}

func MsgHandler(Msg string) fiber.Map {
	return fiber.Map{
		"contents": []fiber.Map{
			{
				"type": "text",
				"text": Msg,
			},
		},
	}
}

func Router(app *fiber.App, browser *rod.Browser, db *badger.DB) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Post("/crawl", func(c *fiber.Ctx) error {
		webcrawl.Crawl(db, browser)
		return c.JSON(MsgHandler("크롤링이 완료되었습니다."))
	})

	app.Post("/get", func(c *fiber.Ctx) error {
		// Print the request body.
		fmt.Println("Request Body:", string(c.Body()))

		// Parse the request body into a struct.
		var request GetRequest
		err := c.BodyParser(&request)
		if err != nil {
			fmt.Println(err)
			return c.JSON(MsgHandler("입력이 잘못되었습니다."))
		}

		// Print the struct.
		id := request.Id["utterance"].(string)
		id = fixId(id)

		fmt.Println("Request Id:", id)

		if id == "발화 내용" {
			return c.JSON(fiber.Map{
				"contents": []fiber.Map{
					{
						"type": "text",
						"text": "hi",
					},
				},
			})
		}
		if id == "테스트" {
			return c.JSON(MsgHandler("테스트 답장"))
		}

		if id == "크롤링" {
			webcrawl.Crawl(db, browser)
			return c.JSON(MsgHandler("크롤링이 완료되었습니다."))
		}

		err = db.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte(id))

			if err != nil {
				return c.JSON(MsgHandler("해당하는 날짜의 식단이 없습니다."))
			}

			value, _ := item.ValueCopy(nil)

			fmt.Println("Value:", string(value))

			return c.JSON(fiber.Map{
				"contents": []fiber.Map{
					{
						"type": "image",
						"image": fiber.Map{
							"url": "https://cafe.sungjin.dev/img/" + string(value) + ".jpg",
						},
					},
				},
			})
		})
		if err != nil {
			fmt.Println(err)
			return c.JSON(MsgHandler("해당하는 날짜의 식단이 없습니다."))
		}
		return nil
	})

	app.Static("/img", "./img")

	app.Listen(":" + os.Getenv("PORT"))
}
