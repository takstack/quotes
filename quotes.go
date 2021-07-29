package main

import (
	//"fmt"
	LOG "github.com/takstack/logger"
	"github.com/takstack/qdb"
	"github.com/takstack/quotes/local"
	"github.com/takstack/quotes/netsrc/netfmp"
	"github.com/takstack/quotes/netsrc/netiex"
)

func main() {
	scope := "all"
	term := "1m"

	updategics() //one-time gics sect update
	//getnewestdates()
	//updateprofiles()

	//getprofs("file")
	getmissingiexeod(scope, term)
	//getquotes("all")
	getdivs(scope, term)
	getsplits(scope, term)

	//getmissingiexeod("file", "2020-06-15") //adj address date range and maxdate check before calling
}
func onetimenetreq() {

}
func onetimedbcall() {

}

//get profiles from iex
func getprofs(s string) {
	netiex.IEXgetprofs(s)
}

//get realtime quotes from iex
func getquotes(s string) {
	netiex.IEXgetrtq(s)
}
func getsplits(s string, term string) {
	netiex.IEXgetsplits(s, term)
}
func getdivs(s string, term string) {
	netiex.IEXgetdivs(s, term)
}

//update company profiles
func updateprofiles() {
	netfmp.FMPgetprofiles()
}

//get quotes for date ranges not added already to db
func getmissingiexeod(s string, term string) {
	netiex.IEXgeteodq(s, term)
}

//get quotes for date ranges not added already to db
func getfmpmissingeodquotes(s string) {
	netfmp.FMPgetEOD(s)
}
func getfmpdivs(s string) {
	netfmp.FMPgetdiv(s)
}
func getfmpsplits(s string) {
	netfmp.FMPgetsplit(s)
}

//find out what dates eod is updated through
func getnewestdates() {
	sl := qdb.Restickdate(qdb.Performquery(qdb.Qryrecdateall()))
	for sym := range sl {
		LOG.GL.Println(sl[sym], sym)
	}
}

//update gics sector from file
func updategics() {
	local.Updgics()
}
