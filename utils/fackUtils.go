package utils

import (
	"fmt"
	"github.com/brianvoe/gofakeit/v6"
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
	defaultValue := MockParamsRandomTypeEnumToString[STRING]
	if randomTypeEnum == "" {
		return defaultValue
	}

	layout := "2006-01-02"

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
		startDate, _ := time.Parse(layout, "2022-01-01 00:00:00")
		endDate, _ := time.Parse(layout, "2023-01-01 00:00:00")
		randomDate := gofakeit.DateRange(startDate, endDate)
		return randomDate.Format(layout)
	case TIMESTAMPS:
		startTime, _ := time.Parse(layout, "2022-01-01 00:00:00")
		endTime, _ := time.Parse(layout, "2023-01-01 00:00:00")
		randomTime := gofakeit.DateRange(startTime, endTime)
		return fmt.Sprintf("%d", randomTime.Unix())
	case PHONE:
		return gofakeit.Phone()
	default:
		return defaultValue
	}
}
