package netfmp

import (
	"encoding/json"
	//LOG "github.com/takstack/logger"
	//"os"
	"sync"
	"time"

	"github.com/takstack/quoteformat"
)

func fmpprofunmarshal(res []byte) format.FMPprofreq {
	//LOG.GL.Println("response body", string(res))
	//writedump(res)
	var prof format.FMPprofreq
	err := json.Unmarshal(res, &prof.Profiles)
	he(err)
	return prof
}
func fmpprofsender(st format.FMPprofreq, c1 chan<- format.FMPprofile, wg *sync.WaitGroup) {
	for elem := range st.Profiles {
		wg.Add(1)
		st.Profiles[elem].Source = "FMP"
		st.Profiles[elem].Importdate = time.Now()
		c1 <- st.Profiles[elem]
	}
}
