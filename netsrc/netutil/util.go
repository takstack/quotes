package netutil

import (
	"sync"

	LOG "github.com/takstack/logger"
	"github.com/takstack/qdb"
	"github.com/takstack/qrw"

	"github.com/takstack/quoteformat"
)

//Listindvsyms gets list of syms from file
func Listindvsyms() [][]string {
	f := qrw.Getreadfile("/quotes/file/targetsyms.txt", 0)
	scanner := qrw.Startbufscanner(f)

	var sl [][]string
	for scanner.Scan() {
		var sltmp []string
		sltmp = append(sltmp, scanner.Text())
		sl = append(sl, sltmp)
	}

	LOG.GL.Println(sl)
	return sl
}

//determines if sym is in string already
func slcontains(sl []format.Updthru, e string) bool {
	for _, a := range sl {
		if a.Sym == e {
			return true
		}
	}
	return false
}

//Receivethru receives syms and dates of stocks updated and enters lastupdated in symbols
func Receivethru(thruch chan format.Updthru, thruwg *sync.WaitGroup) {
	var thruSL []format.Updthru
	for i := range thruch {
		if slcontains(thruSL, i.Sym) {
			continue
		} else {
			//LOG.GL.Println("receivethru received on thruch channel: ", i)
			thruSL = append(thruSL, i)
		}
		if len(thruSL) >= 5 {
			qdb.Sendupddatebatch(thruSL)
			thruSL = nil
		}
		thruwg.Done()
	}

	qdb.Sendupddatebatch(thruSL)

	thruSL = nil
	thruwg.Done()
}
