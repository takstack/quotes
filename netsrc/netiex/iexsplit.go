package netiex

import (
	"fmt"
	//"io/ioutil"
	//"keys"
	LOG "github.com/takstack/logger"
	//"net/http"
	"github.com/takstack/qdb"
	format "github.com/takstack/quoteformat"
	"github.com/takstack/quotes/netsrc/netutil"
	"github.com/takstack/quotes/parser"

	//"strings"
	"sync"
	"time"
)

//IEXgetsplits main handler to update missing quote date ranges from FMP, symsource =all or file
func IEXgetsplits(symsource string, term string) {
	LOG.GL.Println("Starting process: getsplits")
	r1, c1, c2, done, q1, lim, wg, netwg, thruch, thruwg := createchannelssplit()
	go spiniexsplitgophers(r1, c1, lim, wg, netwg, 20, thruch, "split", thruwg, term)
	go netutil.Limiter(lim, q1)
	go parser.Spiniexparsesplitgophers(c1, c2, wg, 3)
	go qdb.DBinsertsplit(c2, done, wg, "split")
	go netutil.Receivethru(thruch, thruwg)

	var symlist2d [][]string
	switch symsource {
	case "all":
		symlist2d = qdb.Restick(qdb.Performquery(qdb.Qryselalltickers()))
	case "file":
		symlist2d = netutil.Listindvsyms() //use this to request only syms in file
	}
	l := len(symlist2d)
	for i := range symlist2d {
		maxdate := maxsplitdateafter(symlist2d[i][0])
		//maxdate = false
		if !maxdate { //false means max date is not after
			//if symlist2d[i][0] == "BRK.B" || symlist2d[i][0] == "BF.B" {
			//	LOG.GL.Println("skipped net req:", symlist2d[i][0])
			//	continue
			//} else {
			netwg.Add(1)
			rsplitsender(symlist2d[i][0], r1, i, l)
			//}
			//} else {
			//LOG.GL.Println("maxdateafter is not after so sent:", symlist2d[i][0])
		}

	}

	close(r1)
	netwg.Wait()
	wg.Wait()
	//LOG.GL.Println("at end of main, after wait")
	close(c1)
	close(q1)
	close(c2)
	close(thruch)
	LOG.GL.Println("at end of splits, before done exit")
	thruwg.Wait()
	<-done
}
func rsplitsender(s string, r1 chan<- string, x int, i int) {
	//LOG.GL.Println("r1sender str:", s, x, ":", i)
	r1 <- s
}

//worker receive symbols(and date ranges) and sends to iexreq to get split quotes
func makesplitreq(r1 <-chan string, c1 chan format.IEXsplit, lim chan interface{}, wg *sync.WaitGroup, netwg *sync.WaitGroup, wk int, thruch chan format.Updthru, mode string, thruwg *sync.WaitGroup, term string) { //maybe add quit ch to kill goroutine
	for str := range r1 {
		now := time.Now()
		addr := iexaddrsplit(str, term)
		respData, RS := netutil.Netreq(addr, str, lim)
		if RS == 200 {
			thruwg.Add(1)
			thruch <- format.Updthru{Sym: str, Upddate: now, Mode: mode}
			LOG.R.Println("resp body:", str, string(respData))
			//LOG.GL.Println("in main respdata:", string(respData))
			res := iexsplitunmarshal(respData)
			//LOG.GL.Println("res unmarshaled:", res)
			iexsplitsender(res, str, c1, wg)

			LOG.GL.Printf("Net request for %s done in:%v by worker %d", str, time.Since(now), wk)
			netwg.Done()
		}
	}
}
func createchannelssplit() (chan string, chan format.IEXsplit, chan format.Split, chan interface{}, chan interface{}, chan interface{}, *sync.WaitGroup, *sync.WaitGroup, chan format.Updthru, *sync.WaitGroup) {
	r1 := make(chan string)          //request channel for net requests
	c1 := make(chan format.IEXsplit) //send channel
	c2 := make(chan format.Split)    //completed channel
	done := make(chan interface{})   //req completed
	q1 := make(chan interface{})     //stop lim
	lim := make(chan interface{})
	var wg sync.WaitGroup
	var netwg sync.WaitGroup
	var thruwg sync.WaitGroup
	thruwg.Add(1) //initial add for final batch completion
	thruch := make(chan format.Updthru)
	return r1, c1, c2, q1, done, lim, &wg, &netwg, thruch, &thruwg
}
func maxsplitdateafter(sym string) bool {
	dbdate := qdb.Resdate(qdb.Performquery(qdb.Qryselsplitthrudate(sym)))
	//start, _ := time.Parse("2006-01-02", chkdate)
	chkdate := time.Now().AddDate(0, 0, -2)
	//LOG.GL.Println("start in max date after: ", start, "dbdate:", dbdate)

	if dbdate.After(chkdate) {
		LOG.GL.Println(sym, " is current through ", dbdate)
		return true
	}
	LOG.GL.Println(sym, " is NOT current. ", dbdate)
	return false

}
func spiniexsplitgophers(r1 chan string, c1 chan format.IEXsplit, lim chan interface{}, wg *sync.WaitGroup, netwg *sync.WaitGroup, numworkers int, thruch chan format.Updthru, mode string, thruwg *sync.WaitGroup, term string) {
	for wk := 0; wk < numworkers; wk++ {
		go makesplitreq(r1, c1, lim, wg, netwg, wk, thruch, mode, thruwg, term)
	}
	LOG.GL.Println("IEX Split netrequest gophers spun up,", numworkers)
}

//iexaddrsplit returns web address for split quotes
func iexaddrsplit(sym string, term string) string {
	//fmt.Println("token:", token)
	addr := fmt.Sprintf("%s/stock/%s/splits/%s?token=%s", addr, sym, term, token)
	//fmt.Println("addr:", addr)
	return addr
}
