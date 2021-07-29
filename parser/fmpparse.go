package parser

import (
	"sync"
	"time"

	LOG "github.com/takstack/logger"

	"github.com/takstack/quoteformat"
)

func fmpeodparse(c1 <-chan format.FMPEOD, c2 chan<- format.EOD, wg *sync.WaitGroup, wk int) {
	var res format.EOD //not updated
	for st := range c1 {
		hdate, err := time.Parse("2006-01-02", st.Thisdate)
		if err != nil {
			LOG.EL.Println("time.Time parse error:", err)
		}
		st.Hdate = hdate
		c2 <- res
	}
}
func fmpdivparse(c1 <-chan format.FMPdiv, c2 chan<- format.Div, wg *sync.WaitGroup, wk int) {
	var res format.Div //not updated
	for st := range c1 {
		hdate, err := time.Parse("2006-01-02", st.Thisdate)
		if err != nil {
			LOG.EL.Println("time.Time parse error:", err)
		}
		st.Hdate = hdate
		//LOG.GL.Println("Parsed ", st)
		c2 <- res
	}
}
func fmpsplitparse(c1 <-chan format.FMPsplit, c2 chan<- format.Split, wg *sync.WaitGroup, wk int) {
	var res format.Split //not updated
	for st := range c1 {
		hdate, err := time.Parse("2006-01-02", st.Thisdate)
		if err != nil {
			LOG.EL.Println("time.Time parse error:", err)
		}
		st.Hdate = hdate
		c2 <- res
	}
}
func fmpprofparse(c1 <-chan format.FMPprofile, c2 chan<- format.FMPprofile, wg *sync.WaitGroup, wk int) {
	for st := range c1 {

		c2 <- st

	}
}

//Spinfmpparseeodgophers parse FMP profiles
func Spinfmpparseeodgophers(c1 chan format.FMPEOD, c2 chan format.EOD, wg *sync.WaitGroup, numworkers int) {
	for wk := 0; wk < numworkers; wk++ {
		go fmpeodparse(c1, c2, wg, wk)
	}
	LOG.GL.Println("FMP EOD parse gophers spun up,", numworkers)

}

//Spinfmpparsedivgophers parse FMP profiles
func Spinfmpparsedivgophers(c1 chan format.FMPdiv, c2 chan format.Div, wg *sync.WaitGroup, numworkers int) {
	for wk := 0; wk < numworkers; wk++ {
		go fmpdivparse(c1, c2, wg, wk)
	}
	LOG.GL.Println("FMP div parse gophers spun up,", numworkers)

}

//Spinfmpparsesplitgophers parse FMP profiles
func Spinfmpparsesplitgophers(c1 chan format.FMPsplit, c2 chan format.Split, wg *sync.WaitGroup, numworkers int) {
	for wk := 0; wk < numworkers; wk++ {
		go fmpsplitparse(c1, c2, wg, wk)
	}
	LOG.GL.Println("FMP split parse gophers spun up,", numworkers)

}

//Spinfmpparsegophers parse FMP profiles
func Spinfmpparsegophers(c1 chan format.FMPprofile, c2 chan format.FMPprofile, wg *sync.WaitGroup, numworkers int) {
	for wk := 0; wk < numworkers; wk++ {
		go fmpprofparse(c1, c2, wg, wk)
	}
	LOG.GL.Println("FMP profile parse gophers spun up,", numworkers)

}
