package netfmp

import (
	"fmt"
	//"io/ioutil"
	//"keys"
	LOG "github.com/takstack/logger"
	//"net/http"

	"github.com/takstack/qdb"
	"github.com/takstack/quotes/netsrc/netutil"
	"github.com/takstack/quotes/parser"
	format "github.com/takstack/quotesformat"

	//"strings"

	"sync"
	"time"
)

//FMPgetdiv main handler to update missing quote date ranges from FMP, symsource =all or file
func FMPgetdiv(symsource string) {
	r1, c1, c2, done, q1, lim, wg, netwg := createchannelsdiv()
	go spinfmpdivgophers(r1, c1, lim, wg, netwg, 20)
	go netutil.Limiter(lim, q1)
	go parser.Spinfmpparsedivgophers(c1, c2, wg, 3)
	go qdb.DBinsertdiv(c2, done, wg, "div")

	var symlist2d [][]string
	switch symsource {
	case "all":
		symlist2d = qdb.Restick(qdb.Performquery(qdb.Qryselalltickers()))
	case "file":
		symlist2d = netutil.Listindvsyms() //use this to request only syms in file
	}
	l := len(symlist2d)
	for i := range symlist2d {
		//if i > 11 {
		//	break
		//}
		netwg.Add(1)
		rdivsender(symlist2d[i][0], r1, i, l)
	}

	close(r1)
	netwg.Wait()
	wg.Wait()
	//LOG.GL.Println("at end of main, after wait")
	close(c1)
	close(q1)
	close(c2)
	LOG.GL.Println("at end of main, before done exit")
	<-done
}
func rdivsender(s string, r1 chan<- string, x int, i int) {
	//LOG.GL.Println("r1sender str:", s, x, ":", i)
	r1 <- s
}

//worker receive symbols(and date ranges) and sends to fmpreq to get div quotes
func makedivreq(r1 <-chan string, c1 chan format.FMPdiv, lim chan interface{}, wg *sync.WaitGroup, netwg *sync.WaitGroup, wk int) { //maybe add quit ch to kill goroutine
	for str := range r1 {
		now := time.Now()
		addr := fmpaddrdiv(str)
		respData, _ := netutil.Netreq(addr, str, lim)
		//LOG.GL.Println("in main respdata:", string(respData))
		res := fmpdivunmarshal(respData)
		//LOG.GL.Println("res unmarshaled:", res)
		fmpdivsender(res, c1, wg)
		LOG.GL.Printf("Net request for %s done in:%v by worker %d", str, time.Since(now), wk)
		netwg.Done()
	}
}
func createchannelsdiv() (chan string, chan format.FMPdiv, chan format.Div, chan interface{}, chan interface{}, chan interface{}, *sync.WaitGroup, *sync.WaitGroup) {
	r1 := make(chan string)        //request channel for net requests
	c1 := make(chan format.FMPdiv) //send channel
	c2 := make(chan format.Div)    //completed channel
	done := make(chan interface{}) //req completed
	q1 := make(chan interface{})   //stop lim
	lim := make(chan interface{})
	var wg sync.WaitGroup
	var netwg sync.WaitGroup
	return r1, c1, c2, q1, done, lim, &wg, &netwg
}
func spinfmpdivgophers(r1 chan string, c1 chan format.FMPdiv, lim chan interface{}, wg *sync.WaitGroup, netwg *sync.WaitGroup, numworkers int) {
	for wk := 0; wk < numworkers; wk++ {
		go makedivreq(r1, c1, lim, wg, netwg, wk)
	}
	LOG.GL.Println("div netrequest gophers spun up,", numworkers)
}

//fmpaddrdiv returns web address for div quotes
func fmpaddrdiv(sym string) string {
	addr := fmt.Sprintf("%s/historical-price-full/stock_dividend/%s?apikey=%s", addr, sym, token)
	return addr
}
