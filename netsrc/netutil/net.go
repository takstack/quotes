package netutil

import (
	"log"
	//"io/ioutil"
	//"keys"
	LOG "github.com/takstack/logger"
	//"net/http"
	//"github.com/takstack/qdb"
	//"github.com/takstack/quoteformat"
	//"github.com/takstack/quotes/parser"
	//"strings"
	//"sync"
	"io/ioutil"
	//"keys"
	"net/http"
	//"os"
	//"github.com/takstack/qrw"
)

//Netreq requests data, sends response to parsedata
func Netreq(addr string, sym string, lim chan<- interface{}) ([]byte, int) {
	httpClient := new(http.Client)
	//LOG.GL.Println("req address:", addr)
	req, err := http.NewRequest(http.MethodGet, addr, nil)
	he(err)
	lim <- 1
	resp, err := httpClient.Do(req)
	he(err)

	RS := resp.StatusCode
	LOG.GL.Println("HTTP Response Status:", RS, http.StatusText(RS), sym)
	if RS != 200 {
		LOG.EL.Println("HTTP Response Status:", RS, http.StatusText(RS), sym)
		log.Fatalln("HTTP Response Status:", RS, http.StatusText(RS), sym)
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	//LOG.R.Println(addr)
	//LOG.R.Println("resp body:", sym, string(respData))
	he(err)
	return respData, RS
}
func he(err error) {
	if err != nil {
		LOG.EL.Fatalln(err)
	}
}
