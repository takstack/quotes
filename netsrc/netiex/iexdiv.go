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

//IEXgetdivs main handler to update missing quote date ranges from FMP, symsource =all or file
func IEXgetdivs(symsource string, term string) {
	LOG.GL.Println("Starting process: getdivs")
	r1, c1, c2, done, q1, lim, wg, netwg, thruch, thruwg := createchannelsdiv()
	go spiniexdivgophers(r1, c1, lim, wg, netwg, 20, thruch, "div", thruwg, term)
	go netutil.Limiter(lim, q1)
	go parser.Spiniexparsedivgophers(c1, c2, wg, 3)
	go qdb.DBinsertdiv(c2, done, wg, "div")
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
		maxdate := maxdivdateafter(symlist2d[i][0])
		//maxdate = false
		if !maxdate { //false means max date is not after
			//if symlist2d[i][0] == "BRK.B" || symlist2d[i][0] == "BF.B" {
			//	LOG.GL.Println("skipped net req:", symlist2d[i][0])
			//	continue
			//} else {
			netwg.Add(1)
			rdivsender(symlist2d[i][0], r1, i, l)

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
	LOG.GL.Println("at end of divs, before done exit")
	thruwg.Wait()
	<-done
}
func rdivsender(s string, r1 chan<- string, x int, i int) {
	//LOG.GL.Println("r1sender str:", s, x, ":", i)
	r1 <- s
}

//worker receive symbols(and date ranges) and sends to iexreq to get div quotes
func makedivreq(r1 <-chan string, c1 chan format.IEXdiv, lim chan interface{}, wg *sync.WaitGroup, netwg *sync.WaitGroup, wk int, thruch chan format.Updthru, mode string, thruwg *sync.WaitGroup, term string) { //maybe add quit ch to kill goroutine
	for str := range r1 {
		now := time.Now()
		addr := iexaddrdiv(str, term)
		respData, RS := netutil.Netreq(addr, str, lim)
		if RS == 200 {
			thruwg.Add(1)
			thruch <- format.Updthru{Sym: str, Upddate: now, Mode: mode}
			LOG.R.Println("resp body:", str, string(respData))
			//LOG.GL.Println("in main respdata:", string(respData))
			res := iexdivunmarshal(respData)
			//LOG.GL.Println("res unmarshaled:", res)
			iexdivsender(res, str, c1, wg)

			LOG.GL.Printf("Net request for %s done in:%v by worker %d", str, time.Since(now), wk)
			netwg.Done()
		}
	}
}
func createchannelsdiv() (chan string, chan format.IEXdiv, chan format.Div, chan interface{}, chan interface{}, chan interface{}, *sync.WaitGroup, *sync.WaitGroup, chan format.Updthru, *sync.WaitGroup) {
	r1 := make(chan string)        //request channel for net requests
	c1 := make(chan format.IEXdiv) //send channel
	c2 := make(chan format.Div)    //completed channel
	done := make(chan interface{}) //req completed
	q1 := make(chan interface{})   //stop lim
	lim := make(chan interface{})
	var wg sync.WaitGroup
	var netwg sync.WaitGroup
	var thruwg sync.WaitGroup
	thruwg.Add(1) //initial add for final batch completion
	thruch := make(chan format.Updthru)
	return r1, c1, c2, q1, done, lim, &wg, &netwg, thruch, &thruwg
}
func maxdivdateafter(sym string) bool {
	dbdate := qdb.Resdate(qdb.Performquery(qdb.Qryseldivthrudate(sym)))
	//start := time.Parse("2006-01-02", chkdate)
	chkdate := time.Now().AddDate(0, 0, -2)
	//LOG.GL.Println("start in max date after: ", start, "dbdate:", dbdate)

	if dbdate.After(chkdate) {
		LOG.GL.Println(sym, " is current through ", dbdate)
		return true
	}
	LOG.GL.Println(sym, " is NOT current. ", dbdate)
	return false

}
func spiniexdivgophers(r1 chan string, c1 chan format.IEXdiv, lim chan interface{}, wg *sync.WaitGroup, netwg *sync.WaitGroup, numworkers int, thruch chan format.Updthru, mode string, thruwg *sync.WaitGroup, term string) {
	for wk := 0; wk < numworkers; wk++ {
		go makedivreq(r1, c1, lim, wg, netwg, wk, thruch, mode, thruwg, term)
	}
	LOG.GL.Println("IEX Div netrequest gophers spun up,", numworkers)
}

//iexaddrdiv returns web address for div quotes
func iexaddrdiv(sym string, term string) string {
	//fmt.Println("token:", token)
	addr := fmt.Sprintf("%s/stock/%s/dividends/%s?token=%s", addr, sym, term, token)
	//fmt.Println("addr:", addr)
	return addr
}
