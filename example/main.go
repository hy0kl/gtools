package main

import (
	"log"

	"github.com/hy0kl/gtools"
)

func main() {
	strAmount := `12.00`
	amount, err := gtools.DecimalMoneyMul100(strAmount)
	log.Printf(`strAmount: %s, amount: %d, err: %v`, strAmount, amount, err)

	timeNow := gtools.GetUnixMillis()
	timeStr := gtools.UnixMsec2Date(timeNow, `Y-m-d H:i:s`)
	log.Println("timeNow:", timeNow, ", timeStr:", timeStr)

	timeOffset := gtools.NaturalDay(-2)
	offsetDate := gtools.UnixMsec2Date(timeOffset, `Y-m-d H:i:s`)
	timeDiff := timeNow - timeOffset
	humanTime := gtools.HumanUnixMillis(timeDiff)
	log.Println(`timeOffset:`, timeOffset, `, offsetDate:`, offsetDate, `, humanTime:`, humanTime)

	log.Println("Unix Timestamp 0 is:", gtools.UnixMsec2Date(0, `Y-m-d H:i:s`))

	url := `https://cn.bing.com/`
	code, respHeader, respBody, err := gtools.SimpleHttpClient(gtools.HttpMethodGet, url, gtools.DefaultFormReqHeaders(), "", gtools.DefaultHttpTimeout())
	log.Printf(`code: %d, header: %v, body: %s, err: %v`, code, respHeader, gtools.SubString(string(respBody), 0, 128), err)
}
