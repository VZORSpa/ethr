package ethr

import (
	"log"
	"net"
	"testing"
)

func TestNWTool(t *testing.T) {

	data := ethrNetStat{}
	getNetDevStats(&data)

	log.Printf("%+v", data)

	for c := range data.netDevStats {
		log.Printf("%+v %v %s", data.netDevStats[c], data.netDevStats[c].flags&net.FlagPointToPoint, data.netDevStats[c].hwAddr)

	}
}
