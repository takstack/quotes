package netiex

import (
	//"fmt"
	//"io/ioutil"
	key "github.com/takstack/keys"
	LOG "github.com/takstack/logger"
)

var token string
var addr string

func init() {
	a := key.Getkeys("IEX")
	token = a[2]
	addr = a[3]
}

func he(err error) {
	if err != nil {
		LOG.EL.Fatalln(err)
	}
}
