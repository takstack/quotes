package netfmp

import (
	"github.com/takstack/key"
	LOG "github.com/takstack/logger"
)

var token string
var addr string

func init() {
	a := key.Getkeys("FMP")
	token = a[1]
	addr = a[3]
}

func he(err error) {
	if err != nil {
		LOG.EL.Fatalln(err)
	}
}
