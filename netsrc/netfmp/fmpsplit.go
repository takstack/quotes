package netfmp

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

//FMPgetsplit main handler to update missing quote date ranges from FMP, symsource =all or file
func FMPgetsplit(symsource string) {
	r1, c1, c2, done, q1, lim, wg, netwg := createchannelssplit()
	go spinfmpsplitgophers(r1, c1, lim, wg, netwg, 20)
	go netutil.Limiter(lim, q1)
	go parser.Spinfmpparsesplitgophers(c1, c2, wg, 3)
	go qdb.DBinsertsplit(c2, done, wg, "split")

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
		rsplitsender(symlist2d[i][0], r1, i, l)
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
func rsplitsender(s string, r1 chan<- string, x int, i int) {
	//LOG.GL.Println("r1sender str:", s, x, ":", i)
	r1 <- s
}

//worker receive symbols(and date ranges) and sends to fmpreq to get split quotes
func makesplitreq(r1 <-chan string, c1 chan format.FMPsplit, lim chan interface{}, wg *sync.WaitGroup, netwg *sync.WaitGroup, wk int) { //maybe add quit ch to kill goroutine
	for str := range r1 {
		now := time.Now()
		addr := fmpaddrsplit(str)
		respData, _ := netutil.Netreq(addr, str, lim)
		//LOG.GL.Println("in main respdata:", string(respData))
		res := fmpsplitunmarshal(respData)
		//LOG.GL.Println("res unmarshaled:", res)
		fmpsplitsender(res, c1, wg)
		LOG.GL.Printf("Net request for %s done in:%v by worker %d", str, time.Since(now), wk)
		netwg.Done()
	}
}
func createchannelssplit() (chan string, chan format.FMPsplit, chan format.Split, chan interface{}, chan interface{}, chan interface{}, *sync.WaitGroup, *sync.WaitGroup) {
	r1 := make(chan string)          //request channel for net requests
	c1 := make(chan format.FMPsplit) //send channel
	c2 := make(chan format.Split)    //completed channel
	done := make(chan interface{})   //req completed
	q1 := make(chan interface{})     //stop lim
	lim := make(chan interface{})
	var wg sync.WaitGroup
	var netwg sync.WaitGroup
	return r1, c1, c2, q1, done, lim, &wg, &netwg
}
func spinfmpsplitgophers(r1 chan string, c1 chan format.FMPsplit, lim chan interface{}, wg *sync.WaitGroup, netwg *sync.WaitGroup, numworkers int) {
	for wk := 0; wk < numworkers; wk++ {
		go makesplitreq(r1, c1, lim, wg, netwg, wk)
	}
	LOG.GL.Println("split netrequest gophers spun up,", numworkers)
}

//fmpaddrsplit returns web address for split quotes
func fmpaddrsplit(sym string) string {
	addr := fmt.Sprintf("%s/historical-price-full/stock_split/%s?apikey=%s", addr, sym, token)
	return addr
}
