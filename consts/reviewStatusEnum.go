package consts

/**
 * Created with GoLand 2022.2.3.
 * @author: 炸薯条
 * Date: 2023/6/11
 * Time: 18:24
 * Description: 审核状态
 */

type ReviewStatusType int

const (
	REVIEWING ReviewStatusType = iota
	PASS
	REJECT
)

var ReviewStatusEnumToInt = map[ReviewStatusType]int{
	REVIEWING: 0,
	PASS:      1,
	REJECT:    2,
}

var ReviewStatusIntToEnum = map[int]ReviewStatusType{
	0: REVIEWING,
	1: PASS,
	2: REJECT,
}

func GetReviewStatus(status int) bool {
	_, ok := ReviewStatusIntToEnum[status]
	return ok
}
