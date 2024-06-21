// -----------------------------------------------------------------------------
// Copyright (C) Microsoft. All rights reserved.
// Licensed under the MIT license.
// See LICENSE.txt file in the project root for full license information.
// -----------------------------------------------------------------------------
package ethr

import (
	"net"
)

type ethrNetStat struct {
	netDevStats []ethrNetDevStat
	tcpStats    ethrTCPStat
}

type ethrNetDevStat struct {
	interfaceName string
	hwAddr        net.HardwareAddr
	rxBytes       uint64
	txBytes       uint64
	rxPkts        uint64
	txPkts        uint64
	txErrPkts     uint64
	rxErrPkts     uint64
	txDrops       uint64
	rxDrops       uint64
	flags         net.Flags
}

type ethrTCPStat struct {
	segRetrans uint64
}

func getNetDevStatDiff(curStats ethrNetDevStat, prevNetStats ethrNetStat, seconds uint64) ethrNetDevStat {
	for _, prevStats := range prevNetStats.netDevStats {
		if prevStats.interfaceName != curStats.interfaceName {
			continue
		}

		if curStats.rxBytes >= prevStats.rxBytes {
			curStats.rxBytes -= prevStats.rxBytes
		} else {
			curStats.rxBytes += (^uint64(0) - prevStats.rxBytes)
		}

		if curStats.txBytes >= prevStats.txBytes {
			curStats.txBytes -= prevStats.txBytes
		} else {
			curStats.txBytes += (^uint64(0) - prevStats.txBytes)
		}

		if curStats.rxPkts >= prevStats.rxPkts {
			curStats.rxPkts -= prevStats.rxPkts
		} else {
			curStats.rxPkts += (^uint64(0) - prevStats.rxPkts)
		}

		if curStats.txPkts >= prevStats.txPkts {
			curStats.txPkts -= prevStats.txPkts
		} else {
			curStats.txPkts += (^uint64(0) - prevStats.txPkts)
		}

		break
	}
	curStats.rxBytes /= seconds
	curStats.txBytes /= seconds
	curStats.rxPkts /= seconds
	curStats.txPkts /= seconds
	return curStats
}
