package consts

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/18
 * Time: 16:12
 * Description: 举报状态枚举
 */

type ReportStatus int

const (
	DEFAULT ReportStatus = iota
	HANDLED
)

var ReportStatusEnumToInt = map[ReportStatus]int{
	DEFAULT: 0,
	HANDLED: 1,
}

var ReportStatusIntToEnum = map[int]ReportStatus{
	0: DEFAULT,
	1: HANDLED,
}

func GetReportStatus(status int) bool {
	_, ok := ReportStatusIntToEnum[status]
	return ok
}
