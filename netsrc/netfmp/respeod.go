package netfmp

import (
	"encoding/json"
	//LOG "github.com/takstack/logger"
	//"os"
	"sync"
	"time"

	"github.com/takstack/quotesformat"
)

func fmpeodunmarshal(res []byte) format.FMPhistEOD {
	//LOG.GL.Println(string(res))
	//writedump(res)
	var e format.FMPhistEOD
	err := json.Unmarshal(res, &e)
	he(err)
	//LOG.GL.Println("unmarshaled eod")
	return e
}
func fmpeodsender(st format.FMPhistEOD, c1 chan<- format.FMPEOD, wg *sync.WaitGroup) {
	//for elem := range st.Historicalstocklist {
	for j := range st.AllEOD {

		wg.Add(1)
		st.AllEOD[j].Sym = st.Symbol
		st.AllEOD[j].Source = "FMP"
		st.AllEOD[j].Importdate = time.Now()
		//LOG.GL.Println("fmpeodsender sending: ", st.Historicalstocklist[elem].AllEOD[j])
		c1 <- st.AllEOD[j]
	}
	//}
}
