package main

import (
	"fmt"
	"log"

	"github.com/hy0kl/gtools"
)

func main() {
	strAmount := `12.00`
	amount, err := gtools.DecimalMoneyMul100(strAmount)
	log.Println(fmt.Sprintf(`strAmount: %s, amount: %d, err: %v`, strAmount, amount, err))

	timeNow := gtools.GetUnixMillis()
	timeStr := gtools.UnixMsec2Date(timeNow, `Y-m-d H:i:s`)
	log.Println("timeNow:", timeNow, ", timeStr:", timeStr)

	timeOffset := gtools.NaturalDay(-2)
	offsetDate := gtools.UnixMsec2Date(timeOffset, `Y-m-d H:i:s`)
	timeDiff := timeNow - timeOffset
	humanTime := gtools.HumanUnixMillis(timeDiff)
	log.Println(`timeOffset:`, timeOffset, `, offsetDate:`, offsetDate, `, humanTime:`, humanTime)

	log.Println("Unix Timestamp 0 is:", gtools.UnixMsec2Date(0, `Y-m-d H:i:s`))
}
