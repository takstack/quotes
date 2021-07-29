package netiex

import (
	"encoding/json"
	//LOG "github.com/takstack/logger"
	//"os"
	"sync"

	"github.com/takstack/quotesformat"
)

func iexsplitunmarshal(res []byte) format.IEXsplitarr {
	//LOG.GL.Println(string(res))
	//writedump(res)
	var e format.IEXsplitarr
	err := json.Unmarshal(res, &e.Allsplit)
	he(err)
	//LOG.GL.Println("unmarshaled split")
	return e
}
func iexsplitsender(st format.IEXsplitarr, sym string, c1 chan<- format.IEXsplit, wg *sync.WaitGroup) {
	for elem := range st.Allsplit {
		wg.Add(1)
		st.Allsplit[elem].Symbol = sym
		st.Allsplit[elem].Source = "IEX"
		//LOG.GL.Printf("%+v", st)
		//LOG.GL.Println("iexdivsender sending: ", st.Allsplit[elem])
		c1 <- st.Allsplit[elem]
	}
}
