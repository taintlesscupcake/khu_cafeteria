package webcrawl

import (
	"fmt"
	// "os"
	"regexp"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/utils"
)

func formatDateString(input string) (string, string, string) {
	// 정규 표현식을 사용하여 숫자만 추출
	re := regexp.MustCompile(`\d+`)
	matches := re.FindAllString(input, -1)

	if len(matches) < 2 {
		return "", "", ""
	}

	// 숫자가 6자리인 경우 앞의 두 자리를 제거
	start := matches[0]
	end := matches[1]

	if len(start) == 6 {
		start = start[2:]
	}

	return start, end, start + "-" + end
}

func Crawl(db *badger.DB, browser *rod.Browser) {

	page := browser.MustPage("https://www.khu.ac.kr/kor/user/bbs/BMSR00040/list.do?menuNo=200283&catId=137").MustWindowFullscreen()

	page.MustWaitLoad().MustWaitStable()

	list := page.MustElements(".tal")

	print(len(list))

	for i, item := range list {

		if item.MustHas("a") {
			id := item.MustElement("a").MustText()
			id = strings.TrimSpace(id)
			print(id)
			start, end, id := formatDateString(id)

			err := db.View(func(txn *badger.Txn) error {
				_, err := txn.Get([]byte(id))
				return err
			})

			if err == nil {
				fmt.Println("Already saved", id)
				continue
			}

			newPage := browser.MustPage("https://www.khu.ac.kr/kor/user/bbs/BMSR00040/list.do?menuNo=200283&catId=137").MustWaitLoad().MustWaitStable()

			newPage.MustElements(".tal")[i].MustElement("a").MustClick()

			newPage.MustWaitLoad().MustWaitStable()

			el, err := newPage.MustElement(".row.contents.clearfix").Element("img")
			if err != nil {
				fmt.Println("Error getting image:", err)
				continue
			}

			_ = utils.OutputFile("img/"+id+".jpg", el.MustResource())

			err = db.Update(func(txn *badger.Txn) error {
				err := txn.Set([]byte(id), []byte(el.MustResource()))

				return err
			})
			if err != nil {
				fmt.Println("Error saving to db:", err)
			}

			// start와 end를 time.Time으로 파싱
			startDate, err := time.Parse("0102", start)
			if err != nil {
				fmt.Println("Error parsing start date:", err)
				return
			}
			endDate, err := time.Parse("0102", end)
			if err != nil {
				fmt.Println("Error parsing end date:", err)
				return
			}

			// start와 end 사이의 모든 날짜에 대해 DB에 저장
			for date := startDate; date.Before(endDate.AddDate(0, 0, 1)); date = date.AddDate(0, 0, 1) {
				dateKey := date.Format("0102")
				fmt.Println("Saving:", dateKey, id)

				err := db.Update(func(txn *badger.Txn) error {
					return txn.Set([]byte(dateKey), []byte(id))
				})
				if err != nil {
					fmt.Println("Error saving to db:", err)
				}
			}

			newPage.MustClose()

		}

	}

	page.MustClose()
	browser.MustClose()
}
