package netfmp

import (
	"encoding/json"
	//LOG "github.com/takstack/logger"
	//"os"
	"sync"
	"time"

	"github.com/takstack/quotesformat"
)

func fmpsplitunmarshal(res []byte) format.FMPhistsplit {
	//LOG.GL.Println(string(res))
	//writedump(res)
	var e format.FMPhistsplit
	err := json.Unmarshal(res, &e)
	he(err)
	//LOG.GL.Println("unmarshaled split")
	return e
}
func fmpsplitsender(st format.FMPhistsplit, c1 chan<- format.FMPsplit, wg *sync.WaitGroup) {
	//for elem := range st.Historicalstocklist {
	for j := range st.Allsplit {

		wg.Add(1)
		st.Allsplit[j].Sym = st.Symbol
		st.Allsplit[j].Source = "FMP"
		st.Allsplit[j].Importdate = time.Now()
		//LOG.GL.Println("fmpsplitsender sending: ", st.Historicalstocklist[elem].Allsplit[j])
		c1 <- st.Allsplit[j]
	}
	//}
}
