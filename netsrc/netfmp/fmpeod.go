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
	"github.com/takstack/quotesformat"

	//"strings"
	"sync"
	"time"
)

//FMPgetEOD main handler to update missing quote date ranges from FMP, symsource =all or file
func FMPgetEOD(symsource string) {
	r1, c1, c2, done, q1, lim, wg, netwg := createchannelseod()
	go spinfmpeodgophers(r1, c1, lim, wg, netwg, 20)
	go netutil.Limiter(lim, q1)
	go parser.Spinfmpparseeodgophers(c1, c2, wg, 3)
	go qdb.DBinserteod(c2, done, wg, "eod")

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
		reodsender(symlist2d[i][0], r1, i, l)
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
func reodsender(s string, r1 chan<- string, x int, i int) {
	//LOG.GL.Println("r1sender str:", s, x, ":", i)
	r1 <- s
}

/*
//sendsymbatch sends symbols for netword request in batches
func sendsymbatch() {

		symlist2d := checkindvsyms()
		x := 0
		l := len(symlist2d)
		for i := 3; i < l; i = i + 3 {
			//if i > 11 {
			//	break
			//}
			netwg.Add(1)
			symlist1d := qdb.Convert2to1d(symlist2d[x:i])
			r1sender(strings.Join(symlist1d, ","), r1, x, i)
			x = i
		}
		symlist1d := qdb.Convert2to1d(symlist2d[x:l])
		r1sender(strings.Join(symlist1d, ","), r1, x, l)

}
*/

//worker receive symbols(and date ranges) and sends to fmpreq to get eod quotes
func makeeodreq(r1 <-chan string, c1 chan format.FMPEOD, lim chan interface{}, wg *sync.WaitGroup, netwg *sync.WaitGroup, wk int) { //maybe add quit ch to kill goroutine
	for str := range r1 {
		now := time.Now()
		addr := fmpaddreod(str)
		respData, _ := netutil.Netreq(addr, str, lim)
		//LOG.GL.Println("in main respdata:", string(respData))
		res := fmpeodunmarshal(respData)
		//LOG.GL.Println("res unmarshaled:", res)
		fmpeodsender(res, c1, wg)
		LOG.GL.Printf("Net request for %s done in:%v by worker %d", str, time.Since(now), wk)
		netwg.Done()
	}
}
func createchannelseod() (chan string, chan format.FMPEOD, chan format.EOD, chan interface{}, chan interface{}, chan interface{}, *sync.WaitGroup, *sync.WaitGroup) {
	r1 := make(chan string)        //request channel for net requests
	c1 := make(chan format.FMPEOD) //send channel
	c2 := make(chan format.EOD)    //completed channel
	done := make(chan interface{}) //req completed
	q1 := make(chan interface{})   //stop lim
	lim := make(chan interface{})
	var wg sync.WaitGroup
	var netwg sync.WaitGroup
	return r1, c1, c2, q1, done, lim, &wg, &netwg
}
func spinfmpeodgophers(r1 chan string, c1 chan format.FMPEOD, lim chan interface{}, wg *sync.WaitGroup, netwg *sync.WaitGroup, numworkers int) {
	for wk := 0; wk < numworkers; wk++ {
		go makeeodreq(r1, c1, lim, wg, netwg, wk)
	}
	LOG.GL.Println("EOD netrequest gophers spun up,", numworkers)
}

//fmpaddreod returns web address for eod quotes
func fmpaddreod(sym string) string {
	addr := fmt.Sprintf("%s/historical-price-full/%s?from=1999-12-31&to=2020-06-12&apikey=%s", addr, sym, token)
	return addr
}
