package handleimg

import (
	"time"

	"github.com/dgraph-io/badger/v4"
)

func getMondayNFriday(date string) (string, string, error) {
	// 문자열을 time.Time으로 변환
	t, err := time.Parse("20060102", "2023"+date)
	if err != nil {
		return "", "", err
	}

	// 월요일과의 차이 계산
	offsetMonday := int(time.Monday - t.Weekday())
	if offsetMonday > 0 {
		offsetMonday -= 7
	}

	// 금요일과의 차이 계산
	offsetFriday := int(time.Friday - t.Weekday())

	// 월요일과 금요일의 날짜 계산
	monday := t.AddDate(0, 0, offsetMonday).Format("0102")
	friday := t.AddDate(0, 0, offsetFriday).Format("0102")

	return monday, friday, nil
}

func HandleImg(db *badger.DB, id string) string {
	mon, fri, err := getMondayNFriday(id)
	if err != nil {
		return ""
	}

	return mon + "-" + fri

}
