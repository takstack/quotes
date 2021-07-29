package netiex

import (
	"encoding/json"
	//LOG "github.com/takstack/logger"
	//"os"
	"sync"
	"time"

	"github.com/takstack/quotesformat"
)

func iexrtunmarshal(res []byte) format.IEXrtquote {
	//LOG.GL.Println(string(res))
	//writedump(res)
	var e format.IEXrtquote
	err := json.Unmarshal(res, &e)
	he(err)
	//LOG.GL.Println("unmarshaled eod")
	return e
}
func iexrtsender(st format.IEXrtquote, c1 chan<- format.IEXrtquote, wg *sync.WaitGroup) {
	wg.Add(1)
	st.Importdate = time.Now()
	c1 <- st

}
