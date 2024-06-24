package ethr

// -----------------------------------------------------------------------------
// Copyright (C) VZOR Spa .
// Licensed under the MIT license.
// See LICENSE.txt file in the project root for full license information.
// -----------------------------------------------------------------------------
import (
	"log"
	"net"
	"testing"
)

func TestNWTool(t *testing.T) {

	data := EthrNetStat{}
	err := getNetDevStats(&data)

	if err != nil {
		log.Println(err)
		t.FailNow()
	}

	for c := range data.NetDevStats {
		log.Printf("%+v %v %s", data.NetDevStats[c], data.NetDevStats[c].flags&net.FlagPointToPoint, data.NetDevStats[c].hwAddr)

	}
}
