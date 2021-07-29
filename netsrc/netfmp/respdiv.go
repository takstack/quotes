package netfmp

import (
	"encoding/json"
	//LOG "github.com/takstack/logger"
	//"os"
	"sync"
	"time"

	"github.com/takstack/quoteformat"
)

func fmpdivunmarshal(res []byte) format.FMPhistdiv {
	//LOG.GL.Println(string(res))
	//writedump(res)
	var e format.FMPhistdiv
	err := json.Unmarshal(res, &e)
	he(err)
	//LOG.GL.Println("unmarshaled div")
	return e
}
func fmpdivsender(st format.FMPhistdiv, c1 chan<- format.FMPdiv, wg *sync.WaitGroup) {
	//for elem := range st.Historicalstocklist {
	for j := range st.Alldiv {

		wg.Add(1)
		st.Alldiv[j].Sym = st.Symbol
		st.Alldiv[j].Source = "FMP"
		st.Alldiv[j].Importdate = time.Now()
		//LOG.GL.Println("fmpdivsender sending: ", st.Historicalstocklist[elem].Alldiv[j])
		c1 <- st.Alldiv[j]
	}
	//}
}
