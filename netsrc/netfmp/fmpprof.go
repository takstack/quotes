package netfmp

import (
	"fmt"
	//"io/ioutil"
	//"keys"
	LOG "github.com/takstack/logger"
	//"net/http"

	"strings"
	"sync"

	"github.com/takstack/qdb"
	"github.com/takstack/quotes/netsrc/netutil"
	"github.com/takstack/quotes/parser"
	format "github.com/takstack/quotesformat"
)

//FMPgetprofiles main handler to update profiles from FMP
func FMPgetprofiles() {
	c1, c2, done, q1, lim, wg := createchannelsprof()
	go netutil.Limiter(lim, q1)
	go parser.Spinfmpparsegophers(c1, c2, wg, 3)
	go qdb.DBinsertprofiles(c2, done, wg, "profile")

	symlist := qdb.Restick(qdb.Performquery(qdb.Qryselalltickers()))
	//for elem := range symlist {
	addr := fmpaddrprofile(strings.Join(symlist[0], ","))
	respData, _ := netutil.Netreq(addr, "", lim)
	//LOG.GL.Println("in main respdata:", string(respData))
	res := fmpprofunmarshal(respData)
	fmpprofsender(res, c1, wg)
	//}
	wg.Wait()
	LOG.GL.Println("at end of main, after wait")
	close(c1)
	close(q1)
	close(c2)
	LOG.GL.Println("at end of main, before done exit")
	<-done
}
func createchannelsprof() (chan format.FMPprofile, chan format.FMPprofile, chan interface{}, chan interface{}, chan interface{}, *sync.WaitGroup) {
	c1 := make(chan format.FMPprofile) //send channel
	c2 := make(chan format.FMPprofile) //completed channel
	done := make(chan interface{})     //req completed
	q1 := make(chan interface{})       //stop lim
	lim := make(chan interface{})
	var wg sync.WaitGroup
	return c1, c2, q1, done, lim, &wg
}
func spinfmpprofilegophers() {

}

//addrprofile get req address for profiles
func fmpaddrprofile(sym string) string {
	addr := fmt.Sprintf("%s/profile/%s?apikey=%s", addr, sym, token)
	return addr
}
