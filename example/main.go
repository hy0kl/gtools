package main

import (
	"context"
	"log"
	"sync"
	"time"

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

	exited := make(chan struct{})
	ctx, cancelFunc := context.WithCancel(context.Background())
	gtools.ClearOnSignal(func() {
		log.Print(`exit signal received`)
		cancelFunc()
		close(exited)
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context, wg *sync.WaitGroup) {
		defer wg.Done()

		var i int64
		for {
			select {
			case <-ctx.Done():
				log.Print(`ctx done, i will exit`)
				return

			case <-exited:
				log.Print(`exited ...`)
				return

			default:
				log.Print(`current idx:`, i)
				i++
				time.Sleep(time.Second)
			}
		}
	}(ctx, &wg)
	wg.Wait()

	<-exited

	log.Print(`exit normally`)
}
