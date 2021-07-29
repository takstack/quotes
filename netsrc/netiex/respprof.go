package netiex

import (
	"encoding/json"
	//LOG "github.com/takstack/logger"
	//"os"
	"sync"

	"github.com/takstack/quoteformat"
)

func iexprofunmarshal(res []byte) format.IEXprof {
	//LOG.GL.Println(string(res))
	//writedump(res)
	var e format.IEXprof
	err := json.Unmarshal(res, &e)
	he(err)
	//LOG.GL.Println("unmarshaled div")
	return e
}
func iexprofsender(st format.IEXprof, sym string, c1 chan<- format.IEXprof, wg *sync.WaitGroup) {
	wg.Add(1)
	//LOG.GL.Printf("%+v", st)
	//LOG.GL.Println("iexdivsender sending: ", st[elem])
	c1 <- st
}
