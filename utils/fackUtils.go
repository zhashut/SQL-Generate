package utils

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
	"math/rand"
	. "sql_generate/consts"
	"time"
)

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/7
 * Time: 11:49
 * Description: 随机数生成工具
 */

func GetRandomValue(randomTypeEnum MockParamsRandomTypeEnum) string {
	defaultValues := []string{"lily", "lele", "john", "jack", "ani"}
	defaultValue := defaultValues[rand.Intn(len(defaultValues))]
	if randomTypeEnum == "" {
		return defaultValue
	}

	layout := "2006-01-02 15:04:05"

	switch randomTypeEnum {
	case NAME:
		return gofakeit.Name()
	case CITY:
		return gofakeit.City()
	case EMAIL:
		return gofakeit.Email()
	case URL:
		return gofakeit.URL()
	case IP:
		return gofakeit.IPv4Address()
	case INTEGER:
		return fmt.Sprintf("%d", gofakeit.Number(0, 100000))
	case DECIMALS:
		return fmt.Sprintf("%.2f", gofakeit.Float64Range(0, 100000))
	case UNIVERSITY:
		return gofakeit.Company()
	case DATES:
		startDate, _ := time.Parse(time.RFC3339, "2021-07-14T09:15:00Z")
		endDate, _ := time.Parse(time.RFC3339, "2021-07-14T09:15:00Z")
		randomDate := gofakeit.DateRange(startDate, endDate)
		dateStr := randomDate.Format(layout)
		return fmt.Sprintf("%s", dateStr)
	case TIMESTAMPS:
		startTime, _ := time.Parse(time.RFC3339, "2022-01-01 00:00:00")
		endTime, _ := time.Parse(time.RFC3339, "2023-01-01 00:00:00")
		randomTime := gofakeit.DateRange(startTime, endTime)
		return fmt.Sprintf("%d", randomTime.Unix())
	case PHONE:
		return gofakeit.Phone()
	default:
		return defaultValue
	}
}
