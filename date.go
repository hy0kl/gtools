package gtools

import (
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	SecondAHour      int64 = 3600
	MillsSecondAHour       = SecondAHour * 1000
	SecondADay       int64 = 86400
	MillsSecondADay        = SecondADay * 1000
	MillsSecondAYear       = MillsSecondADay * 365
)

func RFC3339TimeTransfer(datetime string) int64 {
	timeLayout := "2006-01-02T15:04:05Z" // 转化所需模板
	loc, _ := time.LoadLocation("Local") // 获取时区

	tmp, _ := time.ParseInLocation(timeLayout, datetime, loc)
	timestamp := tmp.Unix() * 1000 //转化为时间戳 类型是int64

	return timestamp
}

// GetUnixMillis 取当前系统时间的毫秒
func GetUnixMillis() int64 {
	t := time.Now()
	return t.UnixNano() / 1000000
}

// NaturalDay 自然日的0点
func NaturalDay(offset int64) (um int64) {
	layout := `2006-01-02`
	t := time.Now()
	tm := time.Unix(t.Unix(), 0)
	local, _ := time.LoadLocation("Local")
	date := tm.In(local).Format(layout)
	parse, _ := time.ParseInLocation(layout, date, local)
	baseUm := parse.Unix() * 1000

	offsetUm := MillsSecondADay * offset
	um = baseUm + offsetUm

	return
}

func HumanUnixMillis(t int64) (display string) {
	t = t / 1000

	var second int64 = 1
	var minute = 60 * second
	var oneHour = minute * 60
	var oneDay = oneHour * 24
	var oneWeek = oneDay * 7
	var oneMonth = oneDay * 30
	var oneYear = oneDay * 365

	var box []string
	if t >= oneYear {
		y := t / oneYear
		box = append(box, fmt.Sprintf(`%d year(s)`, y))
		t -= y * oneYear
	}
	if t >= oneMonth {
		m := t / oneMonth
		box = append(box, fmt.Sprintf(`%d month(s)`, m))
		t -= m * oneMonth
	}
	if t >= oneWeek {
		w := t / oneWeek
		box = append(box, fmt.Sprintf(`%d week(s)`, w))
		t -= w * oneWeek
	}
	if t >= oneHour {
		h := t / oneHour
		box = append(box, fmt.Sprintf(`%d hour(s)`, h))
		t -= h * oneHour
	}
	if t >= minute {
		m := t / minute
		box = append(box, fmt.Sprintf(`%d minute(s)`, m))
		t -= m * minute
	}

	if t > 0 {
		box = append(box, fmt.Sprintf(`%d second(s)`, t))
	}

	if len(box) > 0 {
		display = strings.Join(box, ", ")
	}

	return
}

func HumanUnixMillisV2(t int64) (display string) {
	t = t / 1000

	var second int64 = 1
	var minute = 60 * second
	var oneHour = minute * 60
	var oneDay = oneHour * 24
	var oneWeek = oneDay * 7
	var oneMonth = oneDay * 30
	var oneYear = oneDay * 365

	var box []string
	if t >= oneYear {
		y := t / oneYear
		box = append(box, fmt.Sprintf(`%d year(s)`, y))
		t -= y * oneYear
	}
	if t >= oneMonth {
		m := t / oneMonth
		box = append(box, fmt.Sprintf(`%d month(s)`, m))
		t -= m * oneMonth
	}
	if t >= oneWeek {
		w := t / oneWeek
		box = append(box, fmt.Sprintf(`%d week(s)`, w))
		t -= w * oneWeek
	}
	if t >= oneHour {
		h := t / oneHour
		box = append(box, fmt.Sprintf(`%02d`, h))
		t -= h * oneHour
	} else {
		box = append(box, "00")
	}
	if t >= minute {
		m := t / minute
		box = append(box, fmt.Sprintf(`%02d`, m))
		t -= m * minute
	} else {
		box = append(box, "00")
	}

	if t > 0 {
		box = append(box, fmt.Sprintf(`%02d`, t))
	} else {
		box = append(box, `00`)
	}

	if len(box) > 0 {
		display = strings.Join(box, ":")
	}

	return
}

func CalculateAgeByBirthday(birthday string) int {
	exp := strings.Split(birthday, "-")
	if len(exp) < 1 {
		return 0
	}

	year, _ := Str2Int(exp[0])
	age := time.Now().Year() - year
	if age < 0 {
		age = 0
	}
	return age
}

// 针对 golang 的时间函数库难记难用,封装以下两个函数,采用共识标识符来简化原始库的使用 {{{
// millisecond <-> msec
// see: https://www.php.net/manual/zh/function.date.php
// 采用类 linux 时间格式
// 仅取以下值:
// 日: d, D, l, j
// 月: m, M, n
// 年:  Y, y
// 时间: a, H, i, s
// 时区: e
var (
	find = []string{
		`a`, `M`, `n`, // 需要优先替换,否则出现误替换
		`d`, `D`, `l`, `j`,
		`m`,
		`Y`, `y`,
		`H`, `i`, `s`,
		`e`,
	}

	replace = []string{
		`3:04PM`, `Jan`, `1`,
		`02`, `Mon`, `Monday`, `2`,
		`01`,
		`2006`, `06`,
		`15`, `04`, `05`,
		`MST`,
	}
)

func UnixMsec2Date(um int64, layout string) string {
	timestamp := um / 1000
	if timestamp <= 0 {
		return `-`
	}

	tm := time.Unix(timestamp, 0)
	local, _ := time.LoadLocation("Local")

	for i, f := range find {
		layout = strings.Replace(layout, f, replace[i], -1)
	}

	//log.Println(fmt.Sprintf("[UnixMsec2Date] layout: %s", layout))
	return tm.In(local).Format(layout)
}

func Date2UnixMsec(dateStr, layout string) int64 {
	if "" == dateStr {
		return 0
	}

	for i, f := range find {
		layout = strings.Replace(layout, f, replace[i], -1)
	}

	loc, _ := time.LoadLocation("Local")
	parse, err := time.ParseInLocation(layout, dateStr, loc)
	if err != nil {
		log.Println(fmt.Sprintf("[Date2UnixMsec] parse layout get exception, layout: %s, err: %v", layout, err))
		return 0
	}

	return parse.UnixNano() / 1000000
}
