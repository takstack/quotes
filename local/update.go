package local

import (
	LOG "github.com/takstack/logger"
	"github.com/takstack/qdb"
	"github.com/takstack/qrw"
	"io"
)

//Updgics update gics sector/ind from file
func Updgics() {
	sl := getsects()
	qdb.Sendupdsectbatch(sl)

}

func getsects() [][]string {
	CSVr := qrw.StartCSVreader(qrw.Getreadfile("quotes/file/gics_sp500.csv", 0))
	var sl [][]string
	for {
		row, err := CSVr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			LOG.EL.Fatal(err)
		}

		sl = append(sl, []string{row[0], row[1], row[2]})
	}
	LOG.GL.Println("Updtgtalloc completed")
	return sl
}
