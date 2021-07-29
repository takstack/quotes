package parser

import (
	LOG "github.com/takstack/logger"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/takstack/quotesformat"
)

func iexprofparse(c1 <-chan format.IEXprof, c2 chan<- format.Prof, wg *sync.WaitGroup, wk int) {
	var res format.Prof
	for st := range c1 {

		//LOG.GL.Println("Dates to parse: ", st.Thisopendate, st.Thisclosedate, st.Thislastupddate)
		//LOG.GL.Println("Parsed dates: ", hdateo, hdatec, hdateu)

		res.Importdate = time.Now()

		res.Companyname = st.Companyname
		res.Description = st.Description
		res.Sp500memb = false //change for each
		res.Spweight = 0
		res.Djiamemb = false
		res.Memberof = "Sector"
		res.Exch = st.Exch
		res.Vtype = "ETF"
		res.Sym = st.Sym
		res.Sector = st.Sector
		res.Industry = st.Industry

		res.Sym = st.Sym

		c2 <- res
	}
}
func iexeodparse(c1 <-chan format.IEXEOD, c2 chan<- format.EOD, wg *sync.WaitGroup, wk int) {
	var res format.EOD
	for st := range c1 {
		hdate := Parsedate(st.Thisdate)
		//LOG.GL.Println("Dates to parse: ", st.Thisopendate, st.Thisclosedate, st.Thislastupddate)
		//LOG.GL.Println("Parsed dates: ", hdateo, hdatec, hdateu)

		res.Importdate = time.Now()
		res.Hdate = hdate

		res.Open = st.Uopen
		res.High = st.Uhigh
		res.Low = st.Ulow
		res.Close = st.Uclose
		res.Volume = st.Uvolume
		res.Adjopen = st.Open
		res.Adjhigh = st.High
		res.Adjlow = st.Low
		res.Adjclose = st.Close
		res.Adjvolume = st.Volume
		res.Change = st.Change
		res.Changepercent = st.Changepercent
		res.Sym = st.Symbol
		res.Source = st.Source

		c2 <- res
	}
}
func iexrtparse(c1 <-chan format.IEXrtquote, c2 chan<- format.Rtquote, wg *sync.WaitGroup, wk int) {
	var res format.Rtquote
	for st := range c1 {
		hdateo := time.Unix(0, st.Thisopendate*1000000)
		hdatec := time.Unix(0, st.Thisclosedate*1000000)
		hdateu := time.Unix(0, st.Thislastupddate*1000000)
		//LOG.GL.Println("Dates to parse: ", st.Thisopendate, st.Thisclosedate, st.Thislastupddate)
		//LOG.GL.Println("Parsed dates: ", hdateo, hdatec, hdateu)

		res.Opentime = hdateo
		res.Closetime = hdatec
		res.Lastupdated = hdateu
		res.Importdate = time.Now()

		res.Latestsource = st.Latestsource
		res.Change = st.Change
		res.Changepercent = st.Changepercent
		res.Volume = st.Volume
		res.Hopen = st.Hopen
		res.Hclose = st.Hclose
		res.Previousclose = st.Previousclose
		res.Previousvolume = st.Previousvolume
		res.High = st.High
		res.Low = st.Low
		res.Extendedprice = st.Extendedprice
		res.Extendedchange = st.Extendedchange
		res.Extendedchangepercent = st.Extendedchangepercent
		res.Ytdchange = st.Ytdchange
		res.Symbol = st.Symbol
		res.Primaryexchange = st.Primaryexchange
		res.Peratio = st.Peratio
		res.Isusmarketopen = st.Isusmarketopen
		res.Importdate = st.Importdate

		if st.Iexrealtimeprice != 0 {
			res.Realtimeprice = st.Iexrealtimeprice
		} else {
			res.Realtimeprice = st.Latestprice
		}

		c2 <- res
	}
}
func iexdivparse(c1 <-chan format.IEXdiv, c2 chan<- format.Div, wg *sync.WaitGroup, wk int) {
	var res format.Div
	for st := range c1 {
		exdate := Parsedate(st.Exdate)
		paydate := Parsedate(st.Paymentdate)
		recdate := Parsedate(st.Recorddate)
		decdate := Parsedate(st.Declareddate)
		//LOG.GL.Println("Dates to parse: ", st.Thisopendate, st.Thisclosedate, st.Thislastupddate)
		//LOG.GL.Println("Parsed dates: ", hdateo, hdatec, hdateu)

		flt := convAtof(st.Amount)
		res.Amount = flt

		res.Importdate = time.Now()
		res.Exdate = exdate
		res.Paydate = paydate
		res.Recdate = recdate
		res.Decdate = decdate

		res.Sym = st.Symbol
		res.Flag = st.Flag
		res.Currency = st.Currency
		res.Descrip = st.Description
		res.Frequency = st.Frequency
		res.Source = st.Source
		//LOG.GL.Println("IEX parser sending res to c2:", res)
		c2 <- res
	}
}

func iexsplitparse(c1 <-chan format.IEXsplit, c2 chan<- format.Split, wg *sync.WaitGroup, wk int) {
	var res format.Split
	for st := range c1 {
		exdate := Parsedate(st.Exdate)
		decdate := Parsedate(st.Decdate)
		//LOG.GL.Println("Dates to parse: ", st.Thisopendate, st.Thisclosedate, st.Thislastupddate)
		//LOG.GL.Println("Parsed dates: ", hdateo, hdatec, hdateu)

		res.Importdate = time.Now()
		res.Exdate = exdate
		res.Decdate = decdate

		res.Sym = st.Symbol
		res.Ratio = st.Ratio
		res.Tofactor = st.Tofactor
		res.Fromfactor = st.Fromfactor
		res.Descrip = st.Descrip
		res.Source = st.Source
		//LOG.GL.Println("IEX parser sending res to c2:", res)
		c2 <- res
	}
}

//Spiniexparseprofgophers parse IEX eod
func Spiniexparseprofgophers(c1 chan format.IEXprof, c2 chan format.Prof, wg *sync.WaitGroup, numworkers int) {
	for wk := 0; wk < numworkers; wk++ {
		go iexprofparse(c1, c2, wg, wk)
	}
	LOG.GL.Println("IEX Prof parse gophers spun up,", numworkers)

}

//Spiniexparseeodgophers parse IEX eod
func Spiniexparseeodgophers(c1 chan format.IEXEOD, c2 chan format.EOD, wg *sync.WaitGroup, numworkers int) {
	for wk := 0; wk < numworkers; wk++ {
		go iexeodparse(c1, c2, wg, wk)
	}
	LOG.GL.Println("IEX EOD parse gophers spun up,", numworkers)

}

//Spiniexparsertgophers parse IEX realtime
func Spiniexparsertgophers(c1 chan format.IEXrtquote, c2 chan format.Rtquote, wg *sync.WaitGroup, numworkers int) {
	for wk := 0; wk < numworkers; wk++ {
		go iexrtparse(c1, c2, wg, wk)
	}
	LOG.GL.Println("IEX Rt parse gophers spun up,", numworkers)

}

//Spiniexparsedivgophers parse IEX div
func Spiniexparsedivgophers(c1 chan format.IEXdiv, c2 chan format.Div, wg *sync.WaitGroup, numworkers int) {
	for wk := 0; wk < numworkers; wk++ {
		go iexdivparse(c1, c2, wg, wk)
	}
	LOG.GL.Println("IEX Div parse gophers spun up,", numworkers)

}

//Spiniexparsesplitgophers parse IEX split
func Spiniexparsesplitgophers(c1 chan format.IEXsplit, c2 chan format.Split, wg *sync.WaitGroup, numworkers int) {
	for wk := 0; wk < numworkers; wk++ {
		go iexsplitparse(c1, c2, wg, wk)
	}
	LOG.GL.Println("IEX Split parse gophers spun up,", numworkers)

}

//Parsedate checks for empty string and parses format 2006-01-02
func Parsedate(sd string) time.Time {
	var t time.Time
	if sd == "" {
		t = time.Time{}
		return t
	}
	var e error
	t, e = time.Parse("2006-01-02", sd)
	if e != nil {
		LOG.EL.Println("IEXparse cannot parse time,", sd)
	}
	return t
}
func convAtof(s string) float64 {
	var r float64
	if s == "" {
		return 0.0
	}
	c, err := strconv.ParseFloat(s, 64)
	//LOG.GL.Println("Parsefloat:", r)
	if err != nil {
		LOG.EL.Println("parsetocks: Error converting atof", s)
	}
	r = math.Round(c*10000) / 10000 //10,000=4 decimal places
	//LOG.GL.Println("Parsefloat raw and rounded:", c, r)
	return r
}
func convAtoi(s string) bool {
	res, _ := strconv.Atoi(s)
	if res == 1 {
		return true
	} else if res == 0 {
		return false
	} else {
		LOG.EL.Fatalln("parsetocks: Error converting atoi", s)
	}
	return false
}
