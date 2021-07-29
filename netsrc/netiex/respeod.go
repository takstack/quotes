package netiex

import (
	"encoding/json"
	//LOG "github.com/takstack/logger"
	//"os"
	"sync"

	"github.com/takstack/quotesformat"
)

func iexeodunmarshal(res []byte) format.IEXeodarr {
	//LOG.GL.Println(string(res))
	//writedump(res)
	var e format.IEXeodarr
	err := json.Unmarshal(res, &e.Alleod)
	he(err)
	//LOG.GL.Println("unmarshaled eod")
	return e
}
func iexeodsender(st format.IEXeodarr, sym string, c1 chan<- format.IEXEOD, wg *sync.WaitGroup) {
	for elem := range st.Alleod {
		wg.Add(1)
		st.Alleod[elem].Symbol = sym
		st.Alleod[elem].Source = "IEX"
		//LOG.GL.Printf("%+v", st)
		//LOG.GL.Println("fmpeodsender sending: ", st.Historicalstocklist[elem].AllEOD[j])
		c1 <- st.Alleod[elem]
	}
}
