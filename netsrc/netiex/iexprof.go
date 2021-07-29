package netiex

import (
	"fmt"
	//"io/ioutil"
	//"keys"
	LOG "github.com/takstack/logger"
	//"net/http"
	"github.com/takstack/qdb"
	"github.com/takstack/quoteformat"
	"github.com/takstack/quotes/netsrc/netutil"
	"github.com/takstack/quotes/parser"

	//"strings"
	"sync"
	"time"
)

//IEXgetprofs main handler to update missing quote date ranges from FMP, symsource =all or file
func IEXgetprofs(symsource string) {
	LOG.GL.Println("Starting process: getprofs")
	r1, c1, c2, done, q1, lim, wg, netwg, thruch, thruwg := createchannelsprof()
	go spiniexprofgophers(r1, c1, lim, wg, netwg, 20, thruch, "prof", thruwg)
	go netutil.Limiter(lim, q1)
	go parser.Spiniexparseprofgophers(c1, c2, wg, 3)
	go qdb.DBinsertprof(c2, done, wg, "prof")
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
		netwg.Add(1)
		rprofsender(symlist2d[i][0], r1, i, l)
	}

	close(r1)
	netwg.Wait()

	wg.Wait()
	//LOG.GL.Println("at end of main, after wait")
	close(c1)
	close(q1)
	close(c2)
	close(thruch)
	LOG.GL.Println("at end of profs, before done exit")
	thruwg.Wait()
	<-done
}
func rprofsender(s string, r1 chan<- string, x int, i int) {
	//LOG.GL.Println("r1sender str:", s, x, ":", i)
	r1 <- s
}

//worker receive symbols(and date ranges) and sends to iexreq to get prof quotes
func makeprofreq(r1 <-chan string, c1 chan format.IEXprof, lim chan interface{}, wg *sync.WaitGroup, netwg *sync.WaitGroup, wk int, thruch chan format.Updthru, mode string, thruwg *sync.WaitGroup) { //maybe add quit ch to kill goroutine
	for str := range r1 {
		now := time.Now()
		addr := iexaddrprof(str)
		respData, RS := netutil.Netreq(addr, str, lim)
		if RS == 200 {
			thruwg.Add(1)
			thruch <- format.Updthru{Sym: str, Upddate: now, Mode: mode}
			LOG.R.Println("resp body:", str, string(respData))
			//LOG.GL.Println("in main respdata:", string(respData))
			res := iexprofunmarshal(respData)
			//LOG.GL.Println("res unmarshaled:", res)
			iexprofsender(res, str, c1, wg)

			LOG.GL.Printf("Net request for %s done in:%v by worker %d", str, time.Since(now), wk)
			netwg.Done()
		}
	}
}
func createchannelsprof() (chan string, chan format.IEXprof, chan format.Prof, chan interface{}, chan interface{}, chan interface{}, *sync.WaitGroup, *sync.WaitGroup, chan format.Updthru, *sync.WaitGroup) {
	r1 := make(chan string)         //request channel for net requests
	c1 := make(chan format.IEXprof) //send channel
	c2 := make(chan format.Prof)    //completed channel
	done := make(chan interface{})  //req completed
	q1 := make(chan interface{})    //stop lim
	lim := make(chan interface{})
	var wg sync.WaitGroup
	var netwg sync.WaitGroup
	var thruwg sync.WaitGroup
	thruwg.Add(1) //initial add for final batch completion
	thruch := make(chan format.Updthru)
	return r1, c1, c2, q1, done, lim, &wg, &netwg, thruch, &thruwg
}
func spiniexprofgophers(r1 chan string, c1 chan format.IEXprof, lim chan interface{}, wg *sync.WaitGroup, netwg *sync.WaitGroup, numworkers int, thruch chan format.Updthru, mode string, thruwg *sync.WaitGroup) {
	for wk := 0; wk < numworkers; wk++ {
		go makeprofreq(r1, c1, lim, wg, netwg, wk, thruch, mode, thruwg)
	}
	LOG.GL.Println("IEX Prof netrequest gophers spun up,", numworkers)
}

//iexaddrprof returns web address for prof quotes
func iexaddrprof(sym string) string {
	//fmt.Println("token:", token)
	addr := fmt.Sprintf("%s/stock/%s/company?token=%s", addr, sym, token)
	//fmt.Println("addr:", addr)
	return addr
}
