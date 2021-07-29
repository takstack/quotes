package netiex

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

//IEXgetrtq main handler to update missing quote date ranges from FMP, symsource =all or file
func IEXgetrtq(symsource string) {
	r1, c1, c2, done, q1, lim, wg, netwg := createchannelsrt()
	go spiniexrtgophers(r1, c1, lim, wg, netwg, 20)
	go netutil.Limiter(lim, q1)
	go parser.Spiniexparsertgophers(c1, c2, wg, 3)
	go qdb.DBinsertrt(c2, done, wg, "rt")

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
		rrtsender(symlist2d[i][0], r1, i, l)
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
func rrtsender(s string, r1 chan<- string, x int, i int) {
	//LOG.GL.Println("r1sender str:", s, x, ":", i)
	r1 <- s
}

//worker receive symbols(and date ranges) and sends to iexreq to get rt quotes
func makertreq(r1 <-chan string, c1 chan format.IEXrtquote, lim chan interface{}, wg *sync.WaitGroup, netwg *sync.WaitGroup, wk int) { //maybe add quit ch to kill goroutine
	for str := range r1 {
		now := time.Now()
		addr := iexaddrrt(str)
		respData, RS := netutil.Netreq(addr, str, lim)
		if RS == 200 {
			//LOG.GL.Println("in main respdata:", string(respData))
			res := iexrtunmarshal(respData)
			//LOG.GL.Println("res unmarshaled:", res)
			iexrtsender(res, c1, wg)
			LOG.GL.Printf("Net request for %s done in:%v by worker %d", str, time.Since(now), wk)
		}
		netwg.Done()
	}
}
func createchannelsrt() (chan string, chan format.IEXrtquote, chan format.Rtquote, chan interface{}, chan interface{}, chan interface{}, *sync.WaitGroup, *sync.WaitGroup) {
	r1 := make(chan string)            //request channel for net requests
	c1 := make(chan format.IEXrtquote) //send channel
	c2 := make(chan format.Rtquote)    //completed channel
	done := make(chan interface{})     //req completed
	q1 := make(chan interface{})       //stop lim
	lim := make(chan interface{})
	var wg sync.WaitGroup
	var netwg sync.WaitGroup
	return r1, c1, c2, q1, done, lim, &wg, &netwg
}
func spiniexrtgophers(r1 chan string, c1 chan format.IEXrtquote, lim chan interface{}, wg *sync.WaitGroup, netwg *sync.WaitGroup, numworkers int) {
	for wk := 0; wk < numworkers; wk++ {
		go makertreq(r1, c1, lim, wg, netwg, wk)
	}
	LOG.GL.Println("Rt netrequest gophers spun up,", numworkers)
}

//iexaddrrt returns web address for rt quotes
func iexaddrrt(sym string) string {
	//fmt.Println("token:", token)
	addr := fmt.Sprintf("%s/stock/%s/quote?token=%s", addr, sym, token)
	//fmt.Println("addr:", addr)
	return addr
}
