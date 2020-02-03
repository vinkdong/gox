/*
This is a small go library for some conditions
This Main function just a example func to use all functions
*/
package main

import (
	"github.com/vinkdong/gox/encrypt"
	"github.com/vinkdong/gox/http/server"
	"github.com/vinkdong/gox/log"
	"github.com/vinkdong/gox/random"
	"github.com/vinkdong/gox/slice"
	"github.com/vinkdong/gox/vtime"
	"time"
)

func main() {

	// log
	log.Info("start gox example")
	log.Debug("--------------------")

	// encrypt
	log.Infof("base64 encode gox, get %s", encrypt.Base64StringEncode("gox"))
	d, _ := encrypt.Base64StringDecode("YWJjZDEyMzQ1Ng==")
	log.Infof("base64 decode Z294, get %s ", d)

	// http
	s := server.NewServer(":30000")
	go s.Start()

	// random
	random.Seed(time.Now().UnixNano())
	log.Infof("get random int (1-100) %d", random.RangeInt(1, 100))
	log.Infof("get random string (a-Z) %s", random.String(1, 100))

	// slice
	sa := []string{"a", "b", "c"}
	sb := []string{"b", "c", "e"}
	log.Infof("slice %v with %v difference %v", sa, sb, slice.Difference(sa, sb))
	log.Infof("slice %v with %v difference %v", sb, sa, slice.Difference(sb, sa))

	// vtime
	log.Successf("time of 1562920273000 is %s", vtime.ParserTimestampMs(1562920273000).String())
	log.Successf("time of 1562920273000110000 is %s", vtime.ParserTimestampNs(1562920273000110000))

	t := vtime.Time{
		Format: "2006-01-02 15:04:05",
		TZ:     "Asia/Shanghai",
	}
	t.FromRelativeTime("now-1h")
	log.Infof("now -1h is %s", t.Time.String())
}
