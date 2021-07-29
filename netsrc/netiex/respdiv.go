package netiex

import (
	"encoding/json"
	//LOG "github.com/takstack/logger"
	//"os"
	"sync"

	"github.com/takstack/quoteformat"
)

func iexdivunmarshal(res []byte) format.IEXdivarr {
	//LOG.GL.Println(string(res))
	//writedump(res)
	var e format.IEXdivarr
	err := json.Unmarshal(res, &e.Alldiv)
	he(err)
	//LOG.GL.Println("unmarshaled div")
	return e
}
func iexdivsender(st format.IEXdivarr, sym string, c1 chan<- format.IEXdiv, wg *sync.WaitGroup) {
	for elem := range st.Alldiv {
		wg.Add(1)
		st.Alldiv[elem].Symbol = sym
		st.Alldiv[elem].Source = "IEX"
		//LOG.GL.Printf("%+v", st)
		//LOG.GL.Println("iexdivsender sending: ", st.Alldiv[elem])
		c1 <- st.Alldiv[elem]
	}
}
